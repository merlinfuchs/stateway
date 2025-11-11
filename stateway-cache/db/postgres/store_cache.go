package postgres

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/merlinfuchs/stateway/stateway-cache/db/postgres/pgmodel"
	"github.com/merlinfuchs/stateway/stateway-cache/store"
)

func (c *Client) MarkShardEntitiesTainted(ctx context.Context, params store.MarkShardEntitiesTaintedParams) error {
	tx, err := c.DB.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	q := c.Q.WithTx(tx)
	err = q.MarkShardGuildsTainted(ctx, pgmodel.MarkShardGuildsTaintedParams{
		AppID:      int64(params.AppID),
		ShardCount: int64(params.ShardCount),
		ShardID:    int64(params.ShardID),
	})
	if err != nil {
		return fmt.Errorf("failed to mark shard guilds tainted: %w", err)
	}

	err = q.MarkShardRolesTainted(ctx, pgmodel.MarkShardRolesTaintedParams{
		AppID:      int64(params.AppID),
		ShardCount: int64(params.ShardCount),
		ShardID:    int64(params.ShardID),
	})
	if err != nil {
		return fmt.Errorf("failed to mark shard roles tainted: %w", err)
	}

	err = q.MarkShardChannelsTainted(ctx, pgmodel.MarkShardChannelsTaintedParams{
		AppID:      int64(params.AppID),
		ShardCount: int64(params.ShardCount),
		ShardID:    int64(params.ShardID),
	})
	if err != nil {
		return fmt.Errorf("failed to mark shard channels tainted: %w", err)
	}

	err = q.MarkShardEmojisTainted(ctx, pgmodel.MarkShardEmojisTaintedParams{
		AppID:      int64(params.AppID),
		ShardCount: int64(params.ShardCount),
		ShardID:    int64(params.ShardID),
	})
	if err != nil {
		return fmt.Errorf("failed to mark shard emojis tainted: %w", err)
	}

	err = q.MarkShardStickersTainted(ctx, pgmodel.MarkShardStickersTaintedParams{
		AppID:      int64(params.AppID),
		ShardCount: int64(params.ShardCount),
		ShardID:    int64(params.ShardID),
	})
	if err != nil {
		return fmt.Errorf("failed to mark shard stickers tainted: %w", err)
	}

	err = tx.Commit(ctx)
	if err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (c *Client) MassUpsertEntities(ctx context.Context, params store.MassUpsertEntitiesParams) error {
	tx, err := c.DB.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	q := c.Q.WithTx(tx)

	if len(params.Guilds) != 0 {
		guilds := make([]pgmodel.UpsertGuildsParams, len(params.Guilds))
		for i, guild := range params.Guilds {
			data, err := json.Marshal(guild.Data)
			if err != nil {
				return fmt.Errorf("failed to marshal guild data: %w", err)
			}
			guilds[i] = pgmodel.UpsertGuildsParams{
				AppID:   int64(guild.AppID),
				GuildID: int64(guild.GuildID),
				Data:    data,
				CreatedAt: pgtype.Timestamp{
					Time:  guild.CreatedAt,
					Valid: true,
				},
				UpdatedAt: pgtype.Timestamp{
					Time:  guild.UpdatedAt,
					Valid: true,
				},
			}
		}
		guildsRes := q.UpsertGuilds(ctx, guilds)
		if err := guildsRes.Close(); err != nil {
			return fmt.Errorf("failed to close upsert guilds results: %w", err)
		}
	}

	if len(params.Roles) != 0 {
		roles := make([]pgmodel.UpsertRolesParams, len(params.Roles))
		for i, role := range params.Roles {
			data, err := json.Marshal(role.Data)
			if err != nil {
				return fmt.Errorf("failed to marshal role data: %w", err)
			}
			roles[i] = pgmodel.UpsertRolesParams{
				AppID:   int64(role.AppID),
				GuildID: int64(role.GuildID),
				RoleID:  int64(role.RoleID),
				Data:    data,
				CreatedAt: pgtype.Timestamp{
					Time:  role.CreatedAt,
					Valid: true,
				},
				UpdatedAt: pgtype.Timestamp{
					Time:  role.UpdatedAt,
					Valid: true,
				},
			}
		}
		rolesRes := q.UpsertRoles(ctx, roles)
		if err := rolesRes.Close(); err != nil {
			return fmt.Errorf("failed to upsert roles: %w", err)
		}
	}

	if len(params.Channels) != 0 {
		channels := make([]pgmodel.UpsertChannelsParams, len(params.Channels))
		for i, channel := range params.Channels {
			data, err := json.Marshal(channel.Data)
			if err != nil {
				return fmt.Errorf("failed to marshal channel data: %w", err)
			}
			channels[i] = pgmodel.UpsertChannelsParams{
				AppID:     int64(channel.AppID),
				GuildID:   int64(channel.GuildID),
				ChannelID: int64(channel.ChannelID),
				Data:      data,
				CreatedAt: pgtype.Timestamp{
					Time:  channel.CreatedAt,
					Valid: true,
				},
				UpdatedAt: pgtype.Timestamp{
					Time:  channel.UpdatedAt,
					Valid: true,
				},
			}
		}
		channelsRes := q.UpsertChannels(ctx, channels)
		if err := channelsRes.Close(); err != nil {
			return fmt.Errorf("failed to upsert channels: %w", err)
		}
	}

	if len(params.Emojis) != 0 {
		emojis := make([]pgmodel.UpsertEmojisParams, len(params.Emojis))
		for i, emoji := range params.Emojis {
			data, err := json.Marshal(emoji.Data)
			if err != nil {
				return fmt.Errorf("failed to marshal emoji data: %w", err)
			}
			emojis[i] = pgmodel.UpsertEmojisParams{
				AppID:   int64(emoji.AppID),
				GuildID: int64(emoji.GuildID),
				EmojiID: int64(emoji.EmojiID),
				Data:    data,
				CreatedAt: pgtype.Timestamp{
					Time:  emoji.CreatedAt,
					Valid: true,
				},
				UpdatedAt: pgtype.Timestamp{
					Time:  emoji.UpdatedAt,
					Valid: true,
				},
			}
		}
		emojisRes := q.UpsertEmojis(ctx, emojis)
		if err := emojisRes.Close(); err != nil {
			return fmt.Errorf("failed to upsert emojis: %w", err)
		}
	}

	if len(params.Stickers) != 0 {
		stickers := make([]pgmodel.UpsertStickersParams, len(params.Stickers))
		for i, sticker := range params.Stickers {
			data, err := json.Marshal(sticker.Data)
			if err != nil {
				return fmt.Errorf("failed to marshal sticker data: %w", err)
			}
			stickers[i] = pgmodel.UpsertStickersParams{
				AppID:     int64(sticker.AppID),
				GuildID:   int64(sticker.GuildID),
				StickerID: int64(sticker.StickerID),
				Data:      data,
				CreatedAt: pgtype.Timestamp{
					Time:  sticker.CreatedAt,
					Valid: true,
				},
				UpdatedAt: pgtype.Timestamp{
					Time:  sticker.UpdatedAt,
					Valid: true,
				},
			}
		}
		stickersRes := q.UpsertStickers(ctx, stickers)
		if err := stickersRes.Close(); err != nil {
			return fmt.Errorf("failed to upsert stickers: %w", err)
		}
	}

	err = tx.Commit(ctx)
	if err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}
