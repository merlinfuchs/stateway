package server

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/merlinfuchs/stateway/stateway-gateway/app"
	"github.com/merlinfuchs/stateway/stateway-gateway/db/postgres"
	"github.com/merlinfuchs/stateway/stateway-lib/broker"
	"github.com/merlinfuchs/stateway/stateway-lib/config"
	"github.com/merlinfuchs/stateway/stateway-lib/event"
)

type eventHandler struct {
	broker broker.Broker
}

func (h *eventHandler) HandleEvent(event event.Event) {
	err := h.broker.Publish(context.Background(), event)
	if err != nil {
		slog.Error(
			"Failed to publish event",
			slog.String("event_id", event.EventID().String()),
			slog.String("service_type", string(event.ServiceType())),
			slog.String("error", err.Error()),
		)
	}
}

func Run(ctx context.Context, pg *postgres.Client, cfg *config.GatewayConfig) error {
	broker, err := broker.NewNATSBroker(cfg.Broker.NATS.URL)
	if err != nil {
		return fmt.Errorf("failed to create NATS broker: %w", err)
	}

	err = broker.CreateGatewayStream(ctx)
	if err != nil {
		return fmt.Errorf("failed to create gateway stream: %w", err)
	}

	appManager := app.NewAppManager(pg, &eventHandler{broker: broker})

	appManager.Run(ctx)
	return nil
}
