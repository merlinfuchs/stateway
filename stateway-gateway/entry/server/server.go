package server

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/merlinfuchs/stateway/stateway-gateway/app"
	"github.com/merlinfuchs/stateway/stateway-gateway/db/postgres"
	"github.com/merlinfuchs/stateway/stateway-lib/broker"
	"github.com/merlinfuchs/stateway/stateway-lib/config"
	"github.com/merlinfuchs/stateway/stateway-lib/gateway"
)

func Run(ctx context.Context, pg *postgres.Client, cfg *config.RootGatewayConfig) error {
	slog.Info(
		"Starting gateway server and publishing events to NATS broker",
		slog.Int("gateway_count", cfg.Gateway.GatewayCount),
		slog.Int("gateway_id", cfg.Gateway.GatewayID),
	)

	br, err := broker.NewNATSBroker(cfg.Broker.NATS.URL)
	if err != nil {
		return fmt.Errorf("failed to create NATS broker: %w", err)
	}

	gatewayService := gateway.NewGatewayService(NewGateway(pg, pg))
	err = broker.Provide(ctx, br, gatewayService)
	if err != nil {
		return fmt.Errorf("failed to provide gateway service: %w", err)
	}

	err = br.CreateGatewayStream(ctx)
	if err != nil {
		return fmt.Errorf("failed to create gateway stream: %w", err)
	}

	eventHandler := NewEventHandler(br)
	go eventHandler.Run(ctx)

	appManager := app.NewAppManager(
		app.AppManagerConfig{
			GatewayCount: cfg.Gateway.GatewayCount,
			GatewayID:    cfg.Gateway.GatewayID,
			NoResume:     cfg.Gateway.NoResume,
		},
		pg,
		pg,
		pg,
		pg,
		eventHandler,
	)

	appManager.Run(ctx)
	return nil
}
