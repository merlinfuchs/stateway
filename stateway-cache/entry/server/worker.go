package server

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/disgoorg/disgo/gateway"
	"github.com/merlinfuchs/stateway/stateway-cache/store"
	"github.com/merlinfuchs/stateway/stateway-lib/broker"
	"github.com/merlinfuchs/stateway/stateway-lib/event"
)

type CacheWorker struct {
	cacheStore store.CacheStore
	gatewayIDs []int
}

func (l *CacheWorker) BalanceKey() string {
	key := "cache"
	for _, gatewayID := range l.gatewayIDs {
		key += fmt.Sprintf("_%d", gatewayID)
	}
	if len(l.gatewayIDs) == 0 {
		key += "_all"
	}
	return key
}

func (l *CacheWorker) EventFilter() broker.EventFilter {
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

func (l *CacheWorker) HandleEvent(ctx context.Context, event *event.GatewayEvent) error {
	slog.Debug("Received event:", slog.String("type", event.Type))

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
			roles[i] = store.UpsertRoleParams{
				AppID:     event.AppID,
				GuildID:   e.ID,
				RoleID:    role.ID,
				Data:      role,
				CreatedAt: time.Now().UTC(),
				UpdatedAt: time.Now().UTC(),
			}
		}

		channels := make([]store.UpsertChannelParams, 0, len(e.Channels)+len(e.Threads))
		for _, channel := range e.Channels {
			channel := ensureChannelGuildID(channel, e.ID)
			channels = append(channels, store.UpsertChannelParams{
				AppID:     event.AppID,
				GuildID:   e.ID,
				ChannelID: channel.ID(),
				Data:      channel,
				CreatedAt: time.Now().UTC(),
				UpdatedAt: time.Now().UTC(),
			})
		}
		for _, thread := range e.Threads {
			channel := ensureChannelGuildID(thread, e.ID)
			channels = append(channels, store.UpsertChannelParams{
				AppID:     event.AppID,
				GuildID:   e.ID,
				ChannelID: thread.ID(),
				Data:      channel,
				CreatedAt: time.Now().UTC(),
				UpdatedAt: time.Now().UTC(),
			})
		}

		emojis := make([]store.UpsertEmojiParams, len(e.Emojis))
		for _, emoji := range e.Emojis {
			emojis = append(emojis, store.UpsertEmojiParams{
				AppID:     event.AppID,
				GuildID:   e.ID,
				EmojiID:   emoji.ID,
				Data:      emoji,
				CreatedAt: time.Now().UTC(),
				UpdatedAt: time.Now().UTC(),
			})
		}

		stickers := make([]store.UpsertStickerParams, len(e.Stickers))
		for i, sticker := range e.Stickers {
			stickers[i] = store.UpsertStickerParams{
				AppID:     event.AppID,
				GuildID:   e.ID,
				StickerID: sticker.ID,
				Data:      sticker,
				CreatedAt: time.Now().UTC(),
				UpdatedAt: time.Now().UTC(),
			}
		}

		err = l.cacheStore.MassUpsertEntities(ctx, store.MassUpsertEntitiesParams{
			AppID: event.AppID,
			Guilds: []store.UpsertGuildParams{
				{
					AppID:     event.AppID,
					GuildID:   e.Guild.ID,
					Data:      e.Guild,
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
			return fmt.Errorf("failed to mass upsert entities for guild %s: %w", e.ID, err)
		}

		return nil
	case gateway.EventGuildUpdate:
		err = l.cacheStore.UpsertGuilds(ctx, store.UpsertGuildParams{
			AppID:     event.AppID,
			GuildID:   e.Guild.ID,
			Data:      e.Guild,
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
		err = l.cacheStore.UpsertRoles(ctx, store.UpsertRoleParams{
			AppID:     event.AppID,
			GuildID:   e.GuildID,
			RoleID:    e.Role.ID,
			Data:      e.Role,
			CreatedAt: time.Now().UTC(),
			UpdatedAt: time.Now().UTC(),
		})
		if err != nil {
			return fmt.Errorf("failed to upsert role: %w", err)
		}
	case gateway.EventGuildRoleUpdate:
		err = l.cacheStore.UpsertRoles(ctx, store.UpsertRoleParams{
			AppID:     event.AppID,
			GuildID:   e.GuildID,
			RoleID:    e.Role.ID,
			Data:      e.Role,
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
	case gateway.EventChannelCreate:
		err = l.cacheStore.UpsertChannels(ctx, store.UpsertChannelParams{
			AppID:     event.AppID,
			GuildID:   e.GuildID(),
			ChannelID: e.ID(),
			Data:      e.GuildChannel,
			CreatedAt: time.Now().UTC(),
			UpdatedAt: time.Now().UTC(),
		})
		if err != nil {
			return fmt.Errorf("failed to upsert channel: %w", err)
		}
	case gateway.EventChannelUpdate:

		err = l.cacheStore.UpsertChannels(ctx, store.UpsertChannelParams{
			AppID:     event.AppID,
			GuildID:   e.GuildID(),
			ChannelID: e.ID(),
			Data:      e.GuildChannel,
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
		err = l.cacheStore.UpsertChannels(ctx, store.UpsertChannelParams{
			AppID:     event.AppID,
			GuildID:   e.GuildID(),
			ChannelID: e.ID(),
			Data:      e.GuildThread,
			CreatedAt: time.Now().UTC(),
			UpdatedAt: time.Now().UTC(),
		})
		if err != nil {
			return fmt.Errorf("failed to upsert thread: %w", err)
		}
	case gateway.EventThreadUpdate:
		err = l.cacheStore.UpsertChannels(ctx, store.UpsertChannelParams{
			AppID:     event.AppID,
			GuildID:   e.GuildID(),
			ChannelID: e.ID(),
			Data:      e.GuildThread,
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
			emojis[i] = store.UpsertEmojiParams{
				AppID:     event.AppID,
				GuildID:   e.GuildID,
				EmojiID:   emoji.ID,
				Data:      emoji,
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
			stickers[i] = store.UpsertStickerParams{
				AppID:     event.AppID,
				GuildID:   e.GuildID,
				StickerID: sticker.ID,
				Data:      sticker,
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
