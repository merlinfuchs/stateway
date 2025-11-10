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

func (c *Client) GetGuild(ctx context.Context, appID snowflake.ID, guildID snowflake.ID) (*model.Guild, error) {
	row, err := c.Q.GetGuild(ctx, pgmodel.GetGuildParams{
		AppID:   int64(appID),
		GuildID: int64(guildID),
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, store.ErrNotFound
		}
		return nil, err
	}
	return rowToGuild(row), nil
}

func (c *Client) GetGuilds(ctx context.Context, appID snowflake.ID, limit int, offset int) ([]*model.Guild, error) {
	rows, err := c.Q.GetGuilds(ctx, pgmodel.GetGuildsParams{
		AppID:  int64(appID),
		Limit:  int32(limit),
		Offset: int32(offset),
	})
	if err != nil {
		return nil, err
	}

	guilds := make([]*model.Guild, len(rows))
	for i, row := range rows {
		guilds[i] = rowToGuild(row)
	}
	return guilds, nil
}

func (c *Client) SearchGuilds(ctx context.Context, params store.SearchGuildsParams) ([]*model.Guild, error) {
	rows, err := c.Q.SearchGuilds(ctx, pgmodel.SearchGuildsParams{
		AppID:  int64(params.AppID),
		Data:   params.Data,
		Limit:  int32(params.Limit),
		Offset: int32(params.Offset),
	})
	if err != nil {
		return nil, err
	}

	guilds := make([]*model.Guild, len(rows))
	for i, row := range rows {
		guilds[i] = rowToGuild(row)
	}
	return guilds, nil
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

func (c *Client) MarkGuildUnavailable(ctx context.Context, appID snowflake.ID, guildID snowflake.ID) error {
	return c.Q.MarkGuildUnavailable(ctx, pgmodel.MarkGuildUnavailableParams{
		AppID:   int64(appID),
		GuildID: int64(guildID),
	})
}

func (c *Client) DeleteGuild(ctx context.Context, appID snowflake.ID, guildID snowflake.ID) error {
	return c.Q.DeleteGuild(ctx, pgmodel.DeleteGuildParams{
		AppID:   int64(appID),
		GuildID: int64(guildID),
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
