package server

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"time"

	"github.com/disgoorg/disgo/gateway"
	"github.com/merlinfuchs/stateway/stateway-cache/db/postgres"
	"github.com/merlinfuchs/stateway/stateway-cache/store"
	"github.com/merlinfuchs/stateway/stateway-lib/broker"
	"github.com/merlinfuchs/stateway/stateway-lib/config"
	"github.com/merlinfuchs/stateway/stateway-lib/event"
)

func Run(ctx context.Context, pg *postgres.Client, cfg *config.CacheConfig) error {
	br, err := broker.NewNATSBroker(cfg.Broker.NATS.URL)
	if err != nil {
		return fmt.Errorf("failed to create NATS broker: %w", err)
	}

	err = broker.Listen(ctx, br, &CacheListener{
		guildCacheStore: pg,
	})
	if err != nil {
		return fmt.Errorf("failed to listen to gateway events: %w", err)
	}

	<-ctx.Done()
	return nil
}

type CacheListener struct {
	guildCacheStore store.CacheGuildStore
}

func (l *CacheListener) BalanceKey() string {
	return "cache"
}

func (l *CacheListener) EventFilters() []string {
	return []string{
		"ready",
		"guild.>",
	}
}

func (l *CacheListener) HandleEvent(ctx context.Context, event *event.GatewayEvent) error {
	slog.Info("Received event:", slog.String("type", event.Type))

	e, err := gateway.UnmarshalEventData(event.Data, gateway.EventType(event.Type))
	if err != nil {
		return fmt.Errorf("failed to unmarshal event data: %w", err)
	}

	switch e := e.(type) {
	case gateway.EventGuildCreate:
		data, err := json.Marshal(e.Guild)
		if err != nil {
			return fmt.Errorf("failed to marshal guild data: %w", err)
		}

		err = l.guildCacheStore.UpsertGuild(ctx, store.UpsertGuildParams{
			ID:        e.Guild.ID,
			AppID:     event.AppID,
			Data:      data,
			CreatedAt: time.Now().UTC(),
			UpdatedAt: time.Now().UTC(),
		})
		if err != nil {
			return fmt.Errorf("failed to upsert guild: %w", err)
		}
	}

	return nil
}
