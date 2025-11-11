package postgres

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/snowflake/v2"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/merlinfuchs/stateway/stateway-cache/db/postgres/pgmodel"
	"github.com/merlinfuchs/stateway/stateway-cache/model"
	"github.com/merlinfuchs/stateway/stateway-cache/store"
)

func (c *Client) GetChannel(ctx context.Context, appID snowflake.ID, guildID snowflake.ID, channelID snowflake.ID) (*model.Channel, error) {
	row, err := c.Q.GetChannel(ctx, pgmodel.GetChannelParams{
		AppID:     int64(appID),
		GuildID:   int64(guildID),
		ChannelID: int64(channelID),
	})
	if err != nil {
		return nil, err
	}
	return rowToChannel(row)
}

func (c *Client) GetChannels(ctx context.Context, appID snowflake.ID, guildID snowflake.ID, limit int, offset int) ([]*model.Channel, error) {
	rows, err := c.Q.GetChannels(ctx, pgmodel.GetChannelsParams{
		AppID:   int64(appID),
		GuildID: int64(guildID),
		Limit:   int32(limit),
		Offset:  int32(offset),
	})
	if err != nil {
		return nil, err
	}

	channels := make([]*model.Channel, len(rows))
	for i, row := range rows {
		channel, err := rowToChannel(row)
		if err != nil {
			return nil, err
		}
		channels[i] = channel
	}
	return channels, nil
}

func (c *Client) GetChannelsByType(ctx context.Context, appID snowflake.ID, guildID snowflake.ID, types []int, limit int, offset int) ([]*model.Channel, error) {
	types32 := make([]int32, len(types))
	for i, t := range types {
		types32[i] = int32(t)
	}

	rows, err := c.Q.GetChannelsByType(ctx, pgmodel.GetChannelsByTypeParams{
		AppID:   int64(appID),
		GuildID: int64(guildID),
		Types:   types32,
		Limit:   int32(limit),
		Offset:  int32(offset),
	})
	if err != nil {
		return nil, err
	}

	channels := make([]*model.Channel, len(rows))
	for i, row := range rows {
		channel, err := rowToChannel(row)
		if err != nil {
			return nil, err
		}
		channels[i] = channel
	}
	return channels, nil
}

func (c *Client) SearchChannels(ctx context.Context, params store.SearchChannelsParams) ([]*model.Channel, error) {
	rows, err := c.Q.SearchChannels(ctx, pgmodel.SearchChannelsParams{
		AppID:   int64(params.AppID),
		GuildID: int64(params.GuildID),
		Data:    params.Data,
		Limit:   int32(params.Limit),
		Offset:  int32(params.Offset),
	})
	if err != nil {
		return nil, err
	}

	channels := make([]*model.Channel, len(rows))
	for i, row := range rows {
		channel, err := rowToChannel(row)
		if err != nil {
			return nil, err
		}
		channels[i] = channel
	}
	return channels, nil
}

func (c *Client) UpsertChannels(ctx context.Context, channels ...store.UpsertChannelParams) error {
	if len(channels) == 0 {
		return nil
	}

	params := make([]pgmodel.UpsertChannelsParams, len(channels))
	for i, channel := range channels {
		data, err := json.Marshal(channel.Data)
		if err != nil {
			return fmt.Errorf("failed to marshal channel data: %w", err)
		}

		params[i] = pgmodel.UpsertChannelsParams{
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
	res := c.Q.UpsertChannels(ctx, params)
	return res.Close()
}

func (c *Client) DeleteChannel(ctx context.Context, appID snowflake.ID, guildID snowflake.ID, channelID snowflake.ID) error {
	return c.Q.DeleteChannel(ctx, pgmodel.DeleteChannelParams{
		AppID:     int64(appID),
		GuildID:   int64(guildID),
		ChannelID: int64(channelID),
	})
}

func rowToChannel(row pgmodel.CacheChannel) (*model.Channel, error) {
	var data discord.UnmarshalChannel
	err := json.Unmarshal(row.Data, &data)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal channel data: %w", err)
	}

	return &model.Channel{
		AppID:     snowflake.ID(row.AppID),
		GuildID:   snowflake.ID(row.GuildID),
		ChannelID: snowflake.ID(row.ChannelID),
		Data:      data.Channel,
		Tainted:   row.Tainted,
		CreatedAt: row.CreatedAt.Time,
		UpdatedAt: row.UpdatedAt.Time,
	}, nil
}
