package postgres

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/merlinfuchs/stateway/stateway-cache/db/postgres/pgmodel"
	"github.com/merlinfuchs/stateway/stateway-cache/store"
)

func (c *Client) UpsertChannels(ctx context.Context, channels ...store.UpsertChannelParams) error {
	if len(channels) == 0 {
		return nil
	}

	params := make([]pgmodel.UpsertChannelsParams, len(channels))
	for i, channel := range channels {
		params[i] = pgmodel.UpsertChannelsParams{
			AppID:     int64(channel.AppID),
			GuildID:   int64(channel.GuildID),
			ChannelID: int64(channel.ChannelID),
			Data:      channel.Data,
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
	res := c.Q.UpsertChannels(ctx, params)
	return res.Close()
}

func (c *Client) DeleteChannel(ctx context.Context, params store.ChannelIdentifier) error {
	return c.Q.DeleteChannel(ctx, pgmodel.DeleteChannelParams{
		AppID:     int64(params.AppID),
		GuildID:   int64(params.GuildID),
		ChannelID: int64(params.ChannelID),
	})
}
