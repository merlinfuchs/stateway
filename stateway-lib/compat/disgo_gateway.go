package compat

import (
	"context"
	"fmt"
	"time"

	"github.com/disgoorg/disgo/gateway"
	"github.com/disgoorg/snowflake/v2"
	"github.com/merlinfuchs/stateway/stateway-lib/broker"
	"github.com/merlinfuchs/stateway/stateway-lib/event"
	"github.com/nats-io/nats.go/jetstream"
)

type DisgoGatewayConfig struct {
	GatewayCount int
	GroupIDs     []string
	AppIDs       []snowflake.ID
	EventTypes   []string
}

type DisgoGateway struct {
	broker broker.Broker
	config DisgoGatewayConfig
	close  context.CancelFunc

	EventHandlerFunc gateway.EventHandlerFunc
}

func NewDisgoGateway(broker broker.Broker, cfg DisgoGatewayConfig) *DisgoGateway {
	return &DisgoGateway{
		broker: broker,
		config: cfg,
	}
}

func (g *DisgoGateway) ShardCount() int {
	return 0
}

func (g *DisgoGateway) Close(ctx context.Context) {
	if g.close != nil {
		g.close()
	}
}

func (g *DisgoGateway) CloseWithCode(ctx context.Context, code int, message string) {
	g.Close(ctx)
}

func (g *DisgoGateway) Status() gateway.Status {
	return gateway.StatusReady
}

func (g *DisgoGateway) Send(ctx context.Context, op gateway.Opcode, data gateway.MessageData) error {
	return fmt.Errorf("Gateway.Send is not supported")
}

func (g *DisgoGateway) Latency() time.Duration {
	return 0
}

func (g *DisgoGateway) Presence() *gateway.MessageDataPresenceUpdate {
	return nil
}

func (g *DisgoGateway) Intents() gateway.Intents {
	return gateway.IntentsAll
}

func (g *DisgoGateway) ResumeURL() *string {
	return nil
}

func (g *DisgoGateway) LastSequenceReceived() *int {
	return nil
}

func (g *DisgoGateway) SessionID() *string {
	return nil
}

func (g *DisgoGateway) ShardID() int {
	return 0
}

func (g *DisgoGateway) Open(ctx context.Context) error {
	ctx, cancel := context.WithCancel(ctx)
	g.close = cancel

	for i := 0; i < g.config.GatewayCount; i++ {
		listener := gatewayListener{
			eventFilter: broker.EventFilter{
				GatewayIDs: []int{i},
				GroupIDs:   g.config.GroupIDs,
				AppIDs:     g.config.AppIDs,
				EventTypes: g.config.EventTypes,
			},
			gateway:     g,
			handlerFunc: g.EventHandlerFunc,
		}

		err := broker.Listen(ctx, g.broker, listener)
		if err != nil {
			return fmt.Errorf("failed to listen to gateway events: %w", err)
		}
	}

	return nil
}

type gatewayListener struct {
	eventFilter broker.EventFilter

	gateway     gateway.Gateway
	handlerFunc gateway.EventHandlerFunc
}

func (l gatewayListener) BalanceKey() string {
	key := "gateway"
	for _, gatewayID := range l.eventFilter.GatewayIDs {
		key += fmt.Sprintf("_%d", gatewayID)
	}
	return key
}

func (l gatewayListener) EventFilter() broker.EventFilter {
	return l.eventFilter
}

func (l gatewayListener) ConsumerConfig() broker.ConsumerConfig {
	return broker.ConsumerConfig{
		AckPolicy: jetstream.AckNonePolicy,
	}
}

func (l gatewayListener) HandleEvent(ctx context.Context, event *event.GatewayEvent) (bool, error) {
	e, err := gateway.UnmarshalEventData(event.Data, gateway.EventType(event.Type))
	if err != nil {
		return false, fmt.Errorf("failed to unmarshal event data: %w", err)
	}

	l.handlerFunc(l.gateway, gateway.EventType(event.Type), 0, e)
	return true, nil
}
