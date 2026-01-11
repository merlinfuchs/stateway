package server

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/disgoorg/snowflake/v2"
	"github.com/merlinfuchs/stateway/stateway-cache/db/postgres"
	"github.com/merlinfuchs/stateway/stateway-cache/inmemory"
	"github.com/merlinfuchs/stateway/stateway-cache/store"
	"github.com/merlinfuchs/stateway/stateway-lib/broker"
	"github.com/merlinfuchs/stateway/stateway-lib/cache"
	"github.com/merlinfuchs/stateway/stateway-lib/config"
)

func Run(ctx context.Context, pg *postgres.Client, cfg *config.RootCacheConfig) error {
	slog.Info(
		"Starting cache server and listening to gateway events",
		slog.Any("gateway_ids", cfg.Cache.GatewayIDs),
	)

	var cacheStore store.CacheStore = pg
	if cfg.Cache.InMemory {
		slog.Info("Using in-memory cache store")
		cacheStore = inmemory.NewMapCacheStore()
	}

	// Discord some times sends unquoted snowflake IDs, so we need to allow them
	snowflake.AllowUnquoted = true

	br, err := broker.NewNATSBroker(cfg.Broker.NATS.URL)
	if err != nil {
		return fmt.Errorf("failed to create NATS broker: %w", err)
	}

	if len(cfg.Cache.GatewayIDs) == 0 {
		slog.Info("Listening to events from all gateways")
		err = broker.Listen(ctx, br, &CacheWorker{
			cacheStore: cacheStore,
		})
		if err != nil {
			return fmt.Errorf("failed to listen to gateway events: %w", err)
		}
	} else {
		for _, gatewayID := range cfg.Cache.GatewayIDs {
			slog.Info("Listening to events from gateway", slog.Int("gateway_id", gatewayID))
			err = broker.Listen(ctx, br, &CacheWorker{
				cacheStore: cacheStore,
				gatewayIDs: []int{gatewayID},
			})
			if err != nil {
				return fmt.Errorf("failed to listen to gateway events: %w", err)
			}
		}
	}

	cacheService := cache.NewCacheService(NewCaches(cacheStore))
	err = broker.Provide(ctx, br, cacheService)
	if err != nil {
		return fmt.Errorf("failed to provide cache service: %w", err)
	}

	<-ctx.Done()
	return nil
}
