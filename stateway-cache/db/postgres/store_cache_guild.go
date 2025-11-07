package postgres

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/merlinfuchs/stateway/stateway-cache/db/postgres/pgmodel"
	"github.com/merlinfuchs/stateway/stateway-cache/store"
)

func (c *Client) UpsertGuilds(ctx context.Context, guilds ...store.UpsertGuildParams) error {
	if len(guilds) == 0 {
		return nil
	}

	params := make([]pgmodel.UpsertGuildsParams, len(guilds))
	for i, guild := range guilds {
		params[i] = pgmodel.UpsertGuildsParams{
			GroupID:  guild.GroupID,
			ClientID: int64(guild.ClientID),
			GuildID:  int64(guild.GuildID),
			Data:     guild.Data,
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
		GroupID:  params.GroupID,
		ClientID: int64(params.ClientID),
		GuildID:  int64(params.GuildID),
	})
}

func (c *Client) DeleteGuild(ctx context.Context, params store.GuildIdentifier) error {
	return c.Q.DeleteGuild(ctx, pgmodel.DeleteGuildParams{
		GroupID:  params.GroupID,
		ClientID: int64(params.ClientID),
		GuildID:  int64(params.GuildID),
	})
}
