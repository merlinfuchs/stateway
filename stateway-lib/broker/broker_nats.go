package broker

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/merlinfuchs/stateway/stateway-lib/event"
	"github.com/merlinfuchs/stateway/stateway-lib/service"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
)

type NATSBroker struct {
	nc *nats.Conn
	js jetstream.JetStream
}

const (
	GatewayStreamName    = "GATEWAY"
	GatewayStreamSubject = "gateway.>"
)

func streamFromService(st service.ServiceType) (string, error) {
	switch st {
	case service.ServiceTypeGateway:
		return GatewayStreamName, nil
	default:
		return "", fmt.Errorf("unknown service type: %s", st)
	}
}

func gatewayEventSubject(e *event.GatewayEvent) string {
	eventType := strings.ToLower(strings.ReplaceAll(e.EventType(), "_", "."))

	return fmt.Sprintf("gateway.%d.%s.%d.%s", e.GatewayID, e.GroupID, e.AppID, eventType)
}

func NewNATSBroker(url string) (*NATSBroker, error) {
	if url == "" {
		url = nats.DefaultURL
	}

	nc, err := nats.Connect(url)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to NATS: %w", err)
	}

	js, err := jetstream.New(nc, jetstream.WithPublishAsyncMaxPending(5000))
	if err != nil {
		return nil, fmt.Errorf("failed to create JetStream context: %w", err)
	}

	// Verify JetStream is available by checking account info
	_, err = js.AccountInfo(context.Background())
	if err != nil {
		return nil, fmt.Errorf("JetStream is not available: %w (ensure NATS server is started with -js flag)", err)
	}

	broker := &NATSBroker{nc: nc, js: js}

	return broker, nil
}

func (b *NATSBroker) CreateGatewayStream(ctx context.Context) error {
	_, err := b.js.CreateOrUpdateStream(ctx, jetstream.StreamConfig{
		Name:      GatewayStreamName,
		Subjects:  []string{GatewayStreamSubject},
		Retention: jetstream.InterestPolicy,
		MaxAge:    1 * time.Hour,
		MaxBytes:  32 * 1024 * 1024 * 1024, // 32GB
		MaxMsgs:   -1,
		Discard:   jetstream.DiscardOld,
		Storage:   jetstream.FileStorage,
		Replicas:  1,
	})
	if err != nil && !errors.Is(err, jetstream.ErrStreamNotFound) {
		return fmt.Errorf("failed to create or update stream: %w", err)
	}

	return nil
}

func (b *NATSBroker) Publish(ctx context.Context, evt event.Event) error {
	switch e := evt.(type) {
	case *event.GatewayEvent:
		rawEvent, err := event.MarshalEvent(e)
		if err != nil {
			return fmt.Errorf("failed to marshal event: %w", err)
		}

		subject := gatewayEventSubject(e)

		_, err = b.js.PublishAsync(subject, rawEvent)
		if err != nil {
			// Check if error is due to stream not existing
			if errors.Is(err, nats.ErrNoStreamResponse) {
				return fmt.Errorf("stream %s does not exist or JetStream is not properly configured: %w", GatewayStreamName, err)
			}
			return fmt.Errorf("failed to publish event to %s: %w", subject, err)
		}
		return nil
	default:
		return fmt.Errorf("unsupported event type: %T", e)
	}
}

func (b *NATSBroker) PublishComplete(ctx context.Context) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-b.js.PublishAsyncComplete():
		return nil
	}
}

func (b *NATSBroker) Listen(ctx context.Context, listener GenericListener) error {
	stream, err := streamFromService(listener.ServiceType())
	if err != nil {
		return fmt.Errorf("failed to get stream from service: %w", err)
	}

	var filterSubjects []string
	for _, filter := range listener.EventFilter().Subjects() {
		filterSubjects = append(filterSubjects, fmt.Sprintf("%s.%s", listener.ServiceType(), filter))
	}

	consumer, err := b.js.CreateOrUpdateConsumer(ctx, stream, jetstream.ConsumerConfig{
		Name:              listener.BalanceKey(),
		Durable:           listener.BalanceKey(),
		FilterSubjects:    filterSubjects,
		AckPolicy:         jetstream.AckNonePolicy,
		InactiveThreshold: time.Minute * 15,
	})
	if err != nil {
		return fmt.Errorf("failed to create or update consumer: %w", err)
	}

	subject := fmt.Sprintf("%s.>", listener.ServiceType())

	cc, err := consumer.Consume(func(msg jetstream.Msg) {
		event, err := event.UnmarshalEvent(msg.Data())
		if err != nil {
			slog.Error(
				"Failed to unmarshal event",
				slog.String("subject", msg.Subject()),
				slog.String("error", err.Error()),
			)
			/* err = msg.NakWithDelay(time.Second)
			if err != nil {
				slog.Error(
					"Failed to nak message",
					slog.String("subject", msg.Subject()),
					slog.String("error", err.Error()),
				)
			} */
			return
		}

		err = listener.HandleEvent(ctx, event)
		if err != nil {
			slog.Error(
				"Failed to handle event",
				slog.String("subject", msg.Subject()),
				slog.String("error", err.Error()),
			)
			/* err = msg.NakWithDelay(time.Second)
			if err != nil {
				slog.Error(
					"Failed to nak message",
					slog.String("subject", msg.Subject()),
					slog.String("error", err.Error()),
				)
			} */
			return
		}

		/* err = msg.Ack()
		if err != nil {
			slog.Error(
				"Failed to ack message",
				slog.String("subject", msg.Subject()),
				slog.String("error", err.Error()),
			)
		} */
	}, jetstream.ConsumeErrHandler(func(consumeCtx jetstream.ConsumeContext, err error) {
		slog.Error(
			"Failed to consume message",
			slog.String("error", err.Error()),
		)
	}))
	if err != nil {
		return fmt.Errorf("failed to subscribe to %s: %w", subject, err)
	}

	go func() {
		<-ctx.Done()
		cc.Stop()
	}()

	return nil
}

func (b *NATSBroker) Request(ctx context.Context, service service.ServiceType, method string, request any, opts ...RequestOption) (Response, error) {
	subject := fmt.Sprintf("service.%s.%s", service, method)

	options := &RequestOptions{
		Timeout: 5 * time.Second,
	}
	for _, opt := range opts {
		opt(options)
	}

	rawRequest, err := json.Marshal(request)
	if err != nil {
		return Response{
			Success: false,
			Error:   &Error{Message: err.Error(), Code: "request_failed"},
			Data:    nil,
		}, err
	}

	response, err := b.nc.Request(subject, rawRequest, options.Timeout)
	if err != nil {
		return Response{
			Success: false,
			Error:   &Error{Message: err.Error(), Code: "request_failed"},
			Data:    nil,
		}, err
	}

	var resp Response
	err = json.Unmarshal(response.Data, &resp)
	if err != nil {
		return Response{
			Success: false,
			Error:   &Error{Message: err.Error(), Code: "response_failed"},
			Data:    nil,
		}, err
	}

	return resp, nil
}

func (b *NATSBroker) Provide(ctx context.Context, service GenericBrokerService) error {
	subject := fmt.Sprintf("service.%s.>", service.ServiceType())

	sub, err := b.nc.Subscribe(subject, func(msg *nats.Msg) {
		method := strings.SplitN(msg.Subject, ".", 3)[2]
		data, err := service.HandleRequest(ctx, method, msg.Data)

		var resp Response
		if err == nil {
			rawData, err := json.Marshal(data)
			if err != nil {
				slog.Error(
					"Failed to marshal response",
					slog.String("subject", msg.Subject),
					slog.String("error", err.Error()),
				)
				return
			}
			resp = Response{
				Success: true,
				Error:   nil,
				Data:    rawData,
			}
		} else {
			resp = Response{
				Success: false,
				Error:   &Error{Message: err.Error(), Code: "request_failed"},
				Data:    nil,
			}
		}

		rawResp, err := json.Marshal(resp)
		if err != nil {
			slog.Error(
				"Failed to marshal response",
				slog.String("subject", msg.Subject),
				slog.String("error", err.Error()),
			)
			return
		}

		err = msg.Respond(rawResp)
		if err != nil {
			slog.Error(
				"Failed to respond to %s: %w",
				slog.String("subject", msg.Subject),
				slog.String("error", err.Error()),
			)
		}
	})
	if err != nil {
		return fmt.Errorf("failed to subscribe to %s: %w", subject, err)
	}

	go func() {
		<-ctx.Done()
		err := sub.Unsubscribe()
		if err != nil {
			slog.Error("failed to unsubscribe from %s: %s", subject, err.Error())
		}
	}()

	return nil
}
