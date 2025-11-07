package server

import (
	"context"
	"fmt"

	"github.com/merlinfuchs/stateway/stateway-lib/broker"
	"github.com/merlinfuchs/stateway/stateway-lib/config"
	"github.com/merlinfuchs/stateway/stateway-lib/event"
)

func Run(ctx context.Context, cfg *config.CacheConfig) error {
	br, err := broker.NewNATSBroker(cfg.Broker.NATS.URL)
	if err != nil {
		return fmt.Errorf("failed to create NATS broker: %w", err)
	}

	err = broker.Listen(ctx, br, &CacheListener{})
	if err != nil {
		return fmt.Errorf("failed to listen to gateway events: %w", err)
	}

	<-ctx.Done()
	return nil
}

type CacheListener struct{}

func (l *CacheListener) EventFilters() []broker.EventFilter {
	return []broker.EventFilter{}
}

func (l *CacheListener) HandleEvent(ctx context.Context, event *event.GatewayEvent) error {
	return nil
}
