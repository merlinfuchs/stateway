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
		cacheStore: pg,
	})
	if err != nil {
		return fmt.Errorf("failed to listen to gateway events: %w", err)
	}

	<-ctx.Done()
	return nil
}

type CacheListener struct {
	cacheStore store.CacheStore
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
	case gateway.EventReady:
		err = l.cacheStore.MarkShardEntitiesTainted(ctx, store.MarkShardEntitiesTaintedParams{
			GroupID:    event.GroupID,
			ClientID:   event.ClientID,
			ShardCount: e.Shard[1],
			ShardID:    e.Shard[0],
		})
		if err != nil {
			return fmt.Errorf("failed to mark shard guilds as tainted: %w", err)
		}
	case gateway.EventGuildCreate:
		data, err := json.Marshal(e.Guild)
		if err != nil {
			return fmt.Errorf("failed to marshal guild data: %w", err)
		}

		err = l.cacheStore.UpsertGuilds(ctx, store.UpsertGuildParams{
			GroupID:   event.GroupID,
			ClientID:  event.ClientID,
			GuildID:   e.Guild.ID,
			Data:      data,
			CreatedAt: time.Now().UTC(),
			UpdatedAt: time.Now().UTC(),
		})
		if err != nil {
			return fmt.Errorf("failed to upsert guild: %w", err)
		}

		roles := make([]store.UpsertRoleParams, len(e.Roles))
		for i, role := range e.Roles {
			data, err := json.Marshal(role)
			if err != nil {
				return fmt.Errorf("failed to marshal role data: %w", err)
			}

			roles[i] = store.UpsertRoleParams{
				GroupID:   event.GroupID,
				ClientID:  event.ClientID,
				GuildID:   e.ID,
				RoleID:    role.ID,
				Data:      data,
				CreatedAt: time.Now().UTC(),
				UpdatedAt: time.Now().UTC(),
			}
		}

		err = l.cacheStore.UpsertRoles(ctx, roles...)
		if err != nil {
			return fmt.Errorf("failed to upsert roles: %w", err)
		}

		channels := make([]store.UpsertChannelParams, len(e.Channels))
		for i, channel := range e.Channels {
			data, err := json.Marshal(channel)
			if err != nil {
				return fmt.Errorf("failed to marshal channel data: %w", err)
			}
			channels[i] = store.UpsertChannelParams{
				GroupID:   event.GroupID,
				ClientID:  event.ClientID,
				GuildID:   e.ID,
				ChannelID: channel.ID(),
				Data:      data,
				CreatedAt: time.Now().UTC(),
				UpdatedAt: time.Now().UTC(),
			}
		}

		err = l.cacheStore.UpsertChannels(ctx, channels...)
		if err != nil {
			return fmt.Errorf("failed to upsert channels: %w", err)
		}
	case gateway.EventGuildUpdate:
		data, err := json.Marshal(e.Guild)
		if err != nil {
			return fmt.Errorf("failed to marshal guild data: %w", err)
		}

		err = l.cacheStore.UpsertGuilds(ctx, store.UpsertGuildParams{
			GroupID:   event.GroupID,
			ClientID:  event.ClientID,
			GuildID:   e.Guild.ID,
			Data:      data,
			CreatedAt: time.Now().UTC(),
			UpdatedAt: time.Now().UTC(),
		})
		if err != nil {
			return fmt.Errorf("failed to upsert guild: %w", err)
		}
	case gateway.EventGuildDelete:
		if !e.Unavailable {
			err = l.cacheStore.MarkGuildUnavailable(ctx, store.GuildIdentifier{
				GroupID:  event.GroupID,
				ClientID: event.ClientID,
				GuildID:  e.ID,
			})
			if err != nil {
				return fmt.Errorf("failed to mark guild as unavailable: %w", err)
			}
		} else {
			err = l.cacheStore.DeleteGuild(ctx, store.GuildIdentifier{
				GroupID:  event.GroupID,
				ClientID: event.ClientID,
				GuildID:  e.ID,
			})
			if err != nil {
				return fmt.Errorf("failed to delete guild: %w", err)
			}
		}
	case gateway.EventGuildRoleCreate:
		data, err := json.Marshal(e.Role)
		if err != nil {
			return fmt.Errorf("failed to marshal role data: %w", err)
		}

		err = l.cacheStore.UpsertRoles(ctx, store.UpsertRoleParams{
			GroupID:   event.GroupID,
			ClientID:  event.ClientID,
			GuildID:   e.GuildID,
			RoleID:    e.Role.ID,
			Data:      data,
			CreatedAt: time.Now().UTC(),
			UpdatedAt: time.Now().UTC(),
		})
		if err != nil {
			return fmt.Errorf("failed to upsert role: %w", err)
		}
	case gateway.EventGuildRoleUpdate:
		data, err := json.Marshal(e.Role)
		if err != nil {
			return fmt.Errorf("failed to marshal role data: %w", err)
		}

		err = l.cacheStore.UpsertRoles(ctx, store.UpsertRoleParams{
			GroupID:   event.GroupID,
			ClientID:  event.ClientID,
			GuildID:   e.GuildID,
			RoleID:    e.Role.ID,
			Data:      data,
			CreatedAt: time.Now().UTC(),
			UpdatedAt: time.Now().UTC(),
		})
		if err != nil {
			return fmt.Errorf("failed to upsert role: %w", err)
		}
	case gateway.EventGuildRoleDelete:
		err = l.cacheStore.DeleteRole(ctx, store.RoleIdentifier{
			GroupID:  event.GroupID,
			ClientID: event.ClientID,
			GuildID:  e.GuildID,
			RoleID:   e.RoleID,
		})
		if err != nil {
			return fmt.Errorf("failed to delete role: %w", err)
		}
	case gateway.EventChannelCreate:
		data, err := json.Marshal(e.GuildChannel)
		if err != nil {
			return fmt.Errorf("failed to marshal channel data: %w", err)
		}

		err = l.cacheStore.UpsertChannels(ctx, store.UpsertChannelParams{
			GroupID:   event.GroupID,
			ClientID:  event.ClientID,
			GuildID:   e.GuildID(),
			ChannelID: e.ID(),
			Data:      data,
			CreatedAt: time.Now().UTC(),
			UpdatedAt: time.Now().UTC(),
		})
		if err != nil {
			return fmt.Errorf("failed to upsert channel: %w", err)
		}
	case gateway.EventChannelUpdate:
		data, err := json.Marshal(e.GuildChannel)
		if err != nil {
			return fmt.Errorf("failed to marshal channel data: %w", err)
		}

		err = l.cacheStore.UpsertChannels(ctx, store.UpsertChannelParams{
			GroupID:   event.GroupID,
			ClientID:  event.ClientID,
			GuildID:   e.GuildID(),
			ChannelID: e.ID(),
			Data:      data,
			CreatedAt: time.Now().UTC(),
			UpdatedAt: time.Now().UTC(),
		})
		if err != nil {
			return fmt.Errorf("failed to upsert channel: %w", err)
		}
	case gateway.EventChannelDelete:
		err = l.cacheStore.DeleteChannel(ctx, store.ChannelIdentifier{
			GroupID:   event.GroupID,
			ClientID:  event.ClientID,
			GuildID:   e.GuildID(),
			ChannelID: e.ID(),
		})
		if err != nil {
			return fmt.Errorf("failed to delete channel: %w", err)
		}
	}

	return nil
}
