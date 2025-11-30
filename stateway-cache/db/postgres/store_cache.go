package postgres

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"

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

	// Process each entity type, logging errors but continuing to commit as many operations as possible
	// If any error occurs, the upsert for the entity type will be rolled back because it's part of one batch INSERT
	// TODO?: To fix this we could process each entity type in chunks

	if len(params.Guilds) != 0 {
		guilds := make([]pgmodel.UpsertGuildsParams, 0, len(params.Guilds))
		for _, guild := range params.Guilds {
			data, err := json.Marshal(guild.Data)
			if err != nil {
				slog.Error("Failed to marshal guild data", slog.String("error", err.Error()), slog.Int64("guild_id", int64(guild.GuildID)))
				continue
			}
			guilds = append(guilds, pgmodel.UpsertGuildsParams{
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
			})
		}

		if len(guilds) > 0 {
			guildsRes := q.UpsertGuilds(ctx, guilds)
			guildsRes.Exec(func(idx int, err error) {
				if err != nil {
					guild := guilds[idx]
					slog.Error("Failed to upsert guild", slog.String("error", err.Error()), slog.Int64("app_id", guild.AppID), slog.Int64("guild_id", guild.GuildID))
				}
			})
		}
	}

	if len(params.Roles) != 0 {
		roles := make([]pgmodel.UpsertRolesParams, 0, len(params.Roles))
		for _, role := range params.Roles {
			data, err := json.Marshal(role.Data)
			if err != nil {
				slog.Error("Failed to marshal role data", slog.String("error", err.Error()), slog.Int64("role_id", int64(role.RoleID)))
				continue
			}
			roles = append(roles, pgmodel.UpsertRolesParams{
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
			})
		}

		if len(roles) > 0 {
			rolesRes := q.UpsertRoles(ctx, roles)
			rolesRes.Exec(func(idx int, err error) {
				if err != nil {
					role := roles[idx]
					slog.Error("Failed to upsert role", slog.String("error", err.Error()), slog.Int64("app_id", role.AppID), slog.Int64("guild_id", role.GuildID), slog.Int64("role_id", role.RoleID))
				}
			})
		}
	}

	if len(params.Channels) != 0 {
		channels := make([]pgmodel.UpsertChannelsParams, 0, len(params.Channels))
		for _, channel := range params.Channels {
			data, err := json.Marshal(channel.Data)
			if err != nil {
				slog.Error("Failed to marshal channel data", slog.String("error", err.Error()), slog.Int64("channel_id", int64(channel.ChannelID)))
				continue
			}
			channels = append(channels, pgmodel.UpsertChannelsParams{
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
			})
		}

		if len(channels) > 0 {
			channelsRes := q.UpsertChannels(ctx, channels)
			channelsRes.Exec(func(idx int, err error) {
				if err != nil {
					channel := channels[idx]
					slog.Error("Failed to upsert channel", slog.String("error", err.Error()), slog.Int64("app_id", channel.AppID), slog.Int64("guild_id", channel.GuildID), slog.Int64("channel_id", channel.ChannelID))
				}
			})
		}
	}

	if len(params.Emojis) != 0 {
		emojis := make([]pgmodel.UpsertEmojisParams, 0, len(params.Emojis))
		for _, emoji := range params.Emojis {
			data, err := json.Marshal(emoji.Data)
			if err != nil {
				slog.Error("Failed to marshal emoji data", slog.String("error", err.Error()), slog.Int64("emoji_id", int64(emoji.EmojiID)))
				continue
			}
			emojis = append(emojis, pgmodel.UpsertEmojisParams{
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
			})
		}

		if len(emojis) > 0 {
			emojisRes := q.UpsertEmojis(ctx, emojis)
			emojisRes.Exec(func(idx int, err error) {
				if err != nil {
					emoji := emojis[idx]
					slog.Error("Failed to upsert emoji", slog.String("error", err.Error()), slog.Int64("app_id", emoji.AppID), slog.Int64("guild_id", emoji.GuildID), slog.Int64("emoji_id", emoji.EmojiID))
				}
			})
		}
	}

	if len(params.Stickers) != 0 {
		stickers := make([]pgmodel.UpsertStickersParams, 0, len(params.Stickers))
		for _, sticker := range params.Stickers {
			data, err := json.Marshal(sticker.Data)
			if err != nil {
				slog.Error("Failed to marshal sticker data", slog.String("error", err.Error()), slog.Int64("sticker_id", int64(sticker.StickerID)))
				continue
			}
			stickers = append(stickers, pgmodel.UpsertStickersParams{
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
			})
		}

		if len(stickers) > 0 {
			stickersRes := q.UpsertStickers(ctx, stickers)
			stickersRes.Exec(func(idx int, err error) {
				if err != nil {
					sticker := stickers[idx]
					slog.Error("Failed to upsert sticker", slog.String("error", err.Error()), slog.Int64("app_id", sticker.AppID), slog.Int64("guild_id", sticker.GuildID), slog.Int64("sticker_id", sticker.StickerID))
				}
			})
		}
	}

	err = tx.Commit(ctx)
	if err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}
