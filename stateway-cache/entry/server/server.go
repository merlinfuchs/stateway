package server

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"time"

	"github.com/disgoorg/disgo/gateway"
	"github.com/disgoorg/snowflake/v2"
	"github.com/merlinfuchs/stateway/stateway-cache/db/postgres"
	"github.com/merlinfuchs/stateway/stateway-cache/store"
	"github.com/merlinfuchs/stateway/stateway-lib/broker"
	"github.com/merlinfuchs/stateway/stateway-lib/config"
	"github.com/merlinfuchs/stateway/stateway-lib/event"
)

func Run(ctx context.Context, pg *postgres.Client, cfg *config.RootCacheConfig) error {
	slog.Info(
		"Starting cache server and listening to gateway events",
		slog.Any("gateway_ids", cfg.Cache.GatewayIDs),
	)

	// Discord some times sends unquoted snowflake IDs, so we need to allow them
	snowflake.AllowUnquoted = true

	br, err := broker.NewNATSBroker(cfg.Broker.NATS.URL)
	if err != nil {
		return fmt.Errorf("failed to create NATS broker: %w", err)
	}

	if len(cfg.Cache.GatewayIDs) == 0 {
		slog.Info("Listening to events from all gateways")
		err = broker.Listen(ctx, br, &CacheListener{
			cacheStore: pg,
		})
		if err != nil {
			return fmt.Errorf("failed to listen to gateway events: %w", err)
		}
	} else {
		for _, gatewayID := range cfg.Cache.GatewayIDs {
			slog.Info("Listening to events from gateway", slog.Int("gateway_id", gatewayID))
			err = broker.Listen(ctx, br, &CacheListener{
				cacheStore: pg,
				gatewayIDs: []int{gatewayID},
			})
			if err != nil {
				return fmt.Errorf("failed to listen to gateway events: %w", err)
			}
		}
	}

	cacheService := broker.NewCacheService(NewCaches(pg))
	err = broker.Provide(ctx, br, cacheService)
	if err != nil {
		return fmt.Errorf("failed to provide cache service: %w", err)
	}

	<-ctx.Done()
	return nil
}

type CacheListener struct {
	cacheStore store.CacheStore
	gatewayIDs []int
}

func (l *CacheListener) BalanceKey() string {
	return "cache"
}

func (l *CacheListener) EventFilter() broker.EventFilter {
	return broker.EventFilter{
		GatewayIDs: l.gatewayIDs,
		EventTypes: []string{
			"ready",
			"guild.>",
			"channel.>",
			"thread.>",
		},
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
			AppID:      event.AppID,
			ShardCount: e.Shard[1],
			ShardID:    e.Shard[0],
		})
		if err != nil {
			return fmt.Errorf("failed to mark shard guilds as tainted: %w", err)
		}
	case gateway.EventGuildCreate:
		roles := make([]store.UpsertRoleParams, len(e.Roles))
		for i, role := range e.Roles {
			data, err := json.Marshal(role)
			if err != nil {
				return fmt.Errorf("failed to marshal role data: %w", err)
			}

			roles[i] = store.UpsertRoleParams{
				AppID:     event.AppID,
				GuildID:   e.ID,
				RoleID:    role.ID,
				Data:      data,
				CreatedAt: time.Now().UTC(),
				UpdatedAt: time.Now().UTC(),
			}
		}

		channels := make([]store.UpsertChannelParams, len(e.Channels)+len(e.Threads))
		for i, channel := range e.Channels {
			data, err := json.Marshal(channel)
			if err != nil {
				return fmt.Errorf("failed to marshal channel data: %w", err)
			}
			channels[i] = store.UpsertChannelParams{
				AppID:     event.AppID,
				GuildID:   e.ID,
				ChannelID: channel.ID(),
				Data:      data,
				CreatedAt: time.Now().UTC(),
				UpdatedAt: time.Now().UTC(),
			}
		}
		for i, thread := range e.Threads {
			data, err := json.Marshal(thread)
			if err != nil {
				return fmt.Errorf("failed to marshal thread data: %w", err)
			}
			channels[i+len(e.Channels)] = store.UpsertChannelParams{
				AppID:     event.AppID,
				GuildID:   e.ID,
				ChannelID: thread.ID(),
				Data:      data,
				CreatedAt: time.Now().UTC(),
				UpdatedAt: time.Now().UTC(),
			}
		}

		emojis := make([]store.UpsertEmojiParams, len(e.Emojis))
		for i, emoji := range e.Emojis {
			data, err := json.Marshal(emoji)
			if err != nil {
				return fmt.Errorf("failed to marshal emoji data: %w", err)
			}
			emojis[i] = store.UpsertEmojiParams{
				AppID:     event.AppID,
				GuildID:   e.ID,
				EmojiID:   emoji.ID,
				Data:      data,
				CreatedAt: time.Now().UTC(),
				UpdatedAt: time.Now().UTC(),
			}
		}

		stickers := make([]store.UpsertStickerParams, len(e.Stickers))
		for i, sticker := range e.Stickers {
			data, err := json.Marshal(sticker)
			if err != nil {
				return fmt.Errorf("failed to marshal sticker data: %w", err)
			}
			stickers[i] = store.UpsertStickerParams{
				AppID:     event.AppID,
				GuildID:   e.ID,
				StickerID: sticker.ID,
				Data:      data,
				CreatedAt: time.Now().UTC(),
				UpdatedAt: time.Now().UTC(),
			}
		}

		guildData, err := json.Marshal(e.Guild)
		if err != nil {
			return fmt.Errorf("failed to marshal guild data: %w", err)
		}

		err = l.cacheStore.MassUpsertEntities(ctx, store.MassUpsertEntitiesParams{
			Guilds: []store.UpsertGuildParams{
				{
					AppID:     event.AppID,
					GuildID:   e.Guild.ID,
					Data:      guildData,
					CreatedAt: time.Now().UTC(),
					UpdatedAt: time.Now().UTC(),
				},
			},
			Roles:    roles,
			Channels: channels,
			Emojis:   emojis,
			Stickers: stickers,
		})
		if err != nil {
			return fmt.Errorf("failed to mass upsert entities: %w", err)
		}

		return nil
	case gateway.EventGuildUpdate:
		data, err := json.Marshal(e.Guild)
		if err != nil {
			return fmt.Errorf("failed to marshal guild data: %w", err)
		}

		err = l.cacheStore.UpsertGuilds(ctx, store.UpsertGuildParams{
			AppID:     event.AppID,
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
			err = l.cacheStore.MarkGuildUnavailable(ctx, event.AppID, e.ID)
			if err != nil {
				return fmt.Errorf("failed to mark guild as unavailable: %w", err)
			}
		} else {
			err = l.cacheStore.DeleteGuild(ctx, event.AppID, e.ID)
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
			AppID:     event.AppID,
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
			AppID:     event.AppID,
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
		err = l.cacheStore.DeleteRole(ctx, event.AppID, e.GuildID, e.RoleID)
		if err != nil {
			return fmt.Errorf("failed to delete role: %w", err)
		}
		if err != nil {
			return fmt.Errorf("failed to delete role: %w", err)
		}
	case gateway.EventChannelCreate:
		data, err := json.Marshal(e.GuildChannel)
		if err != nil {
			return fmt.Errorf("failed to marshal channel data: %w", err)
		}

		err = l.cacheStore.UpsertChannels(ctx, store.UpsertChannelParams{
			AppID:     event.AppID,
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
			AppID:     event.AppID,
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
		err = l.cacheStore.DeleteChannel(ctx, event.AppID, e.GuildID(), e.ID())
		if err != nil {
			return fmt.Errorf("failed to delete channel: %w", err)
		}
	case gateway.EventThreadCreate:
		data, err := json.Marshal(e.GuildThread)
		if err != nil {
			return fmt.Errorf("failed to marshal thread data: %w", err)
		}

		err = l.cacheStore.UpsertChannels(ctx, store.UpsertChannelParams{
			AppID:     event.AppID,
			GuildID:   e.GuildID(),
			ChannelID: e.ID(),
			Data:      data,
			CreatedAt: time.Now().UTC(),
			UpdatedAt: time.Now().UTC(),
		})
		if err != nil {
			return fmt.Errorf("failed to upsert thread: %w", err)
		}
	case gateway.EventThreadUpdate:
		data, err := json.Marshal(e.GuildThread)
		if err != nil {
			return fmt.Errorf("failed to marshal thread data: %w", err)
		}

		err = l.cacheStore.UpsertChannels(ctx, store.UpsertChannelParams{
			AppID:     event.AppID,
			GuildID:   e.GuildID(),
			ChannelID: e.ID(),
			Data:      data,
			CreatedAt: time.Now().UTC(),
			UpdatedAt: time.Now().UTC(),
		})
		if err != nil {
			return fmt.Errorf("failed to upsert thread: %w", err)
		}
	case gateway.EventThreadDelete:
		err = l.cacheStore.DeleteChannel(ctx, event.AppID, e.GuildID, e.ID)
		if err != nil {
			return fmt.Errorf("failed to delete thread: %w", err)
		}
	case gateway.EventGuildEmojisUpdate:
		emojis := make([]store.UpsertEmojiParams, len(e.Emojis))
		for i, emoji := range e.Emojis {
			data, err := json.Marshal(emoji)
			if err != nil {
				return fmt.Errorf("failed to marshal emoji data: %w", err)
			}
			emojis[i] = store.UpsertEmojiParams{
				AppID:     event.AppID,
				GuildID:   e.GuildID,
				EmojiID:   emoji.ID,
				Data:      data,
				CreatedAt: time.Now().UTC(),
				UpdatedAt: time.Now().UTC(),
			}
		}

		err = l.cacheStore.UpsertEmojis(ctx, emojis...)
		if err != nil {
			return fmt.Errorf("failed to upsert emojis: %w", err)
		}
	case gateway.EventGuildStickersUpdate:
		stickers := make([]store.UpsertStickerParams, len(e.Stickers))
		for i, sticker := range e.Stickers {
			data, err := json.Marshal(sticker)
			if err != nil {
				return fmt.Errorf("failed to marshal sticker data: %w", err)
			}
			stickers[i] = store.UpsertStickerParams{
				AppID:     event.AppID,
				GuildID:   e.GuildID,
				StickerID: sticker.ID,
				Data:      data,
				CreatedAt: time.Now().UTC(),
				UpdatedAt: time.Now().UTC(),
			}
		}

		err = l.cacheStore.UpsertStickers(ctx, stickers...)
		if err != nil {
			return fmt.Errorf("failed to upsert stickers: %w", err)
		}
	}

	return nil
}
