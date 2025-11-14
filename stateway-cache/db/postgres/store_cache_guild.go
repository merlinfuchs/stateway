package postgres

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/disgoorg/disgo/discord"
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
	return rowToGuild(row)
}

func (c *Client) GetGuilds(ctx context.Context, appID snowflake.ID, limit int, offset int) ([]*model.Guild, error) {
	rows, err := c.Q.GetGuilds(ctx, pgmodel.GetGuildsParams{
		AppID: int64(appID),
		Limit: pgtype.Int4{
			Int32: int32(limit),
			Valid: limit != 0,
		},
		Offset: pgtype.Int4{
			Int32: int32(offset),
			Valid: offset != 0,
		},
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, store.ErrNotFound
		}
		return nil, err
	}

	guilds := make([]*model.Guild, len(rows))
	for i, row := range rows {
		guild, err := rowToGuild(row)
		if err != nil {
			return nil, err
		}
		guilds[i] = guild
	}
	return guilds, nil
}

func (c *Client) SearchGuilds(ctx context.Context, params store.SearchGuildsParams) ([]*model.Guild, error) {
	rows, err := c.Q.SearchGuilds(ctx, pgmodel.SearchGuildsParams{
		AppID: int64(params.AppID),
		Data:  params.Data,
		Limit: pgtype.Int4{
			Int32: int32(params.Limit),
			Valid: params.Limit != 0,
		},
		Offset: pgtype.Int4{
			Int32: int32(params.Offset),
			Valid: params.Offset != 0,
		},
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, store.ErrNotFound
		}
		return nil, err
	}

	guilds := make([]*model.Guild, len(rows))
	for i, row := range rows {
		guild, err := rowToGuild(row)
		if err != nil {
			return nil, err
		}
		guilds[i] = guild
	}
	return guilds, nil
}

func (c *Client) UpsertGuilds(ctx context.Context, guilds ...store.UpsertGuildParams) error {
	if len(guilds) == 0 {
		return nil
	}

	params := make([]pgmodel.UpsertGuildsParams, len(guilds))
	for i, guild := range guilds {
		data, err := json.Marshal(guild.Data)
		if err != nil {
			return fmt.Errorf("failed to marshal guild data: %w", err)
		}

		params[i] = pgmodel.UpsertGuildsParams{
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

func rowToGuild(row pgmodel.CacheGuild) (*model.Guild, error) {
	var data discord.Guild
	err := json.Unmarshal(row.Data, &data)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal guild data: %w", err)
	}

	return &model.Guild{
		AppID:       snowflake.ID(row.AppID),
		GuildID:     snowflake.ID(row.GuildID),
		Data:        data,
		Unavailable: row.Unavailable,
		Tainted:     row.Tainted,
		CreatedAt:   row.CreatedAt.Time,
		UpdatedAt:   row.UpdatedAt.Time,
	}, nil
}
