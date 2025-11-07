package postgres

import (
	"context"
	"errors"

	"github.com/disgoorg/snowflake/v2"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/merlinfuchs/stateway/stateway-cache/db/postgres/pgmodel"
	"github.com/merlinfuchs/stateway/stateway-cache/model"
	"github.com/merlinfuchs/stateway/stateway-cache/store"
)

func (c *Client) GetGuild(ctx context.Context, guild store.GuildIdentifier) (*model.Guild, error) {
	row, err := c.Q.GetGuild(ctx, pgmodel.GetGuildParams{
		AppID:   int64(guild.AppID),
		GuildID: int64(guild.GuildID),
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, store.ErrNotFound
		}
		return nil, err
	}
	return rowToGuild(row), nil
}

func (c *Client) UpsertGuilds(ctx context.Context, guilds ...store.UpsertGuildParams) error {
	if len(guilds) == 0 {
		return nil
	}

	params := make([]pgmodel.UpsertGuildsParams, len(guilds))
	for i, guild := range guilds {
		params[i] = pgmodel.UpsertGuildsParams{
			AppID:   int64(guild.AppID),
			GuildID: int64(guild.GuildID),
			Data:    guild.Data,
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
	res := c.Q.UpsertGuilds(ctx, params)
	return res.Close()
}

func (c *Client) MarkGuildUnavailable(ctx context.Context, params store.GuildIdentifier) error {
	return c.Q.MarkGuildUnavailable(ctx, pgmodel.MarkGuildUnavailableParams{
		AppID:   int64(params.AppID),
		GuildID: int64(params.GuildID),
	})
}

func (c *Client) DeleteGuild(ctx context.Context, params store.GuildIdentifier) error {
	return c.Q.DeleteGuild(ctx, pgmodel.DeleteGuildParams{
		AppID:   int64(params.AppID),
		GuildID: int64(params.GuildID),
	})
}

func rowToGuild(row pgmodel.CacheGuild) *model.Guild {
	return &model.Guild{
		AppID:       snowflake.ID(row.AppID),
		GuildID:     snowflake.ID(row.GuildID),
		Data:        row.Data,
		Unavailable: row.Unavailable,
		Tainted:     row.Tainted,
		CreatedAt:   row.CreatedAt.Time,
		UpdatedAt:   row.UpdatedAt.Time,
	}
}
