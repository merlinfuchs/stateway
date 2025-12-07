package server

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/merlinfuchs/stateway/stateway-audit/batcher"
	"github.com/merlinfuchs/stateway/stateway-audit/db/clickhouse"
	"github.com/merlinfuchs/stateway/stateway-audit/db/postgres"
	"github.com/merlinfuchs/stateway/stateway-lib/broker"
	"github.com/merlinfuchs/stateway/stateway-lib/config"
)

func Run(ctx context.Context, pg *postgres.Client, ch *clickhouse.Client, cfg *config.RootAuditConfig) error {
	slog.Info(
		"Starting audit server and publishing events to NATS broker",
		slog.Any("gateway_ids", cfg.Audit.GatewayIDs),
	)

	br, err := broker.NewNATSBroker(cfg.Broker.NATS.URL)
	if err != nil {
		return fmt.Errorf("failed to create NATS broker: %w", err)
	}

	batcher := batcher.NewJetStreamBatcher(br.JetStream(), ch, batcher.JetStreamBatcherConfig{
		NamePrefix: cfg.Broker.NamePrefix,
	})
	err = batcher.CreateStream(ctx)
	if err != nil {
		return fmt.Errorf("failed to create JetStream stream: %w", err)
	}

	err = batcher.Start(ctx)
	if err != nil {
		return fmt.Errorf("failed to start batcher: %w", err)
	}

	err = broker.Listen(ctx, br, NewAuditWorker(
		pg,
		batcher,
		AuditWorkerConfig{
			GatewayIDs: cfg.Audit.GatewayIDs,
			NamePrefix: cfg.Broker.NamePrefix,
		},
	))
	if err != nil {
		return fmt.Errorf("failed to listen to audit events: %w", err)
	}

	<-ctx.Done()

	slog.Info("Shutting down audit server")

	closeCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = br.Close(closeCtx)
	if err != nil {
		return fmt.Errorf("failed to close NATS broker: %w", err)
	}

	return nil
}
