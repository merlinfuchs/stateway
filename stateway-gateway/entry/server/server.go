package server

import (
	"context"
	"fmt"

	"github.com/merlinfuchs/stateway/stateway-gateway/app"
	"github.com/merlinfuchs/stateway/stateway-gateway/db/postgres"
	"github.com/merlinfuchs/stateway/stateway-lib/broker"
	"github.com/merlinfuchs/stateway/stateway-lib/config"
)

func Run(ctx context.Context, pg *postgres.Client, cfg *config.RootGatewayConfig) error {
	broker, err := broker.NewNATSBroker(cfg.Broker.NATS.URL)
	if err != nil {
		return fmt.Errorf("failed to create NATS broker: %w", err)
	}

	err = broker.CreateGatewayStream(ctx)
	if err != nil {
		return fmt.Errorf("failed to create gateway stream: %w", err)
	}

	eventHandler := NewEventHandler(broker)
	go eventHandler.Run(ctx)

	appManager := app.NewAppManager(
		app.AppManagerConfig{
			GatewayCount: cfg.Gateway.GatewayCount,
			GatewayID:    cfg.Gateway.GatewayID,
		},
		pg,
		pg,
		eventHandler,
	)

	appManager.Run(ctx)
	return nil
}
