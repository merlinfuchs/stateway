package broker

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/merlinfuchs/stateway/stateway-lib/event"
	"github.com/merlinfuchs/stateway/stateway-lib/service"
	"github.com/nats-io/nats.go"
)

type NATSBroker struct {
	nc *nats.Conn
	js nats.JetStreamContext
}

const (
	gatewayStreamName = "GATEWAY"
	gatewaySubject    = "gateway.>"
)

func NewNATSBroker(url string) (*NATSBroker, error) {
	if url == "" {
		url = nats.DefaultURL
	}

	nc, err := nats.Connect(url)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to NATS: %w", err)
	}

	js, err := nc.JetStream(nats.PublishAsyncMaxPending(256))
	if err != nil {
		return nil, fmt.Errorf("failed to create JetStream context: %w", err)
	}

	// Verify JetStream is available by checking account info
	_, err = js.AccountInfo()
	if err != nil {
		return nil, fmt.Errorf("JetStream is not available: %w (ensure NATS server is started with -js flag)", err)
	}

	broker := &NATSBroker{nc: nc, js: js}

	return broker, nil
}

func (b *NATSBroker) CreateGatewayStream() error {
	// Create the stream configuration
	streamConfig := &nats.StreamConfig{
		Name:      gatewayStreamName,
		Subjects:  []string{gatewaySubject},
		Retention: nats.InterestPolicy,
		MaxAge:    1 * time.Hour,
		MaxBytes:  4 * 1024 * 1024 * 1024, // 4GB
		Discard:   nats.DiscardOld,
		Storage:   nats.FileStorage,
		Replicas:  1,
	}

	// Check if stream already exists
	stream, err := b.js.StreamInfo(gatewayStreamName)
	if err != nil && !errors.Is(err, nats.ErrStreamNotFound) {
		return fmt.Errorf("failed to check stream info: %w", err)
	}

	// Stream already exists
	if stream != nil {
		slog.Info("JetStream stream already exists, updating it", slog.String("stream", gatewayStreamName))
		_, err = b.js.UpdateStream(streamConfig)
		if err != nil {
			return fmt.Errorf("failed to update stream: %w", err)
		}
		slog.Info("JetStream stream updated successfully", slog.String("stream", gatewayStreamName))
		return nil
	}

	slog.Info("Creating JetStream stream", slog.String("stream", gatewayStreamName), slog.String("subject", gatewaySubject))
	stream, err = b.js.AddStream(streamConfig)
	if err != nil {
		return fmt.Errorf("failed to create stream: %w", err)
	}

	// Verify stream was created successfully
	if stream == nil {
		return fmt.Errorf("stream creation returned nil")
	}

	slog.Info("JetStream stream created successfully", slog.String("stream", gatewayStreamName))

	// Verify the stream exists by querying it again
	stream, err = b.js.StreamInfo(gatewayStreamName)
	if err != nil {
		return fmt.Errorf("stream was created but cannot be verified: %w", err)
	}
	if stream == nil {
		return fmt.Errorf("stream was created but verification returned nil")
	}

	return nil
}

func (b *NATSBroker) Publish(ctx context.Context, evt event.Event) error {
	switch e := evt.(type) {
	case *event.GatewayEvent:
		rawEvent, err := json.Marshal(e)
		if err != nil {
			return fmt.Errorf("failed to marshal event: %w", err)
		}

		subject := fmt.Sprintf("gateway.%s", e.Type)

		_, err = b.js.Publish(subject, rawEvent, nats.Context(ctx))
		if err != nil {
			// Check if error is due to stream not existing
			if errors.Is(err, nats.ErrNoStreamResponse) {
				return fmt.Errorf("stream %s does not exist or JetStream is not properly configured: %w", gatewayStreamName, err)
			}
			return fmt.Errorf("failed to publish event to %s: %w", subject, err)
		}
		return nil
	default:
		return fmt.Errorf("unsupported event type: %T", e)
	}
}

func (b *NATSBroker) Listen(ctx context.Context, listener GenericListener) error {
	return nil
}

func (b *NATSBroker) Request(ctx context.Context, service service.ServiceType, method string, request any, opts ...RequestOption) (Response, error) {
	subject := fmt.Sprintf("%s.%s", service, method)

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
	subject := fmt.Sprintf("%s.>", service.ServiceType())

	sub, err := b.js.Subscribe(subject, func(msg *nats.Msg) {
		var request json.RawMessage
		err := json.Unmarshal(msg.Data, &request)
		if err != nil {
			return
		}

		data, err := service.HandleRequest(msg.Subject, request)

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
