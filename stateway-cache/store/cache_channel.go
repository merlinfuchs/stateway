package store

import (
	"context"
	"encoding/json"
	"time"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/snowflake/v2"
	"github.com/merlinfuchs/stateway/stateway-cache/model"
)

type UpsertChannelParams struct {
	AppID     snowflake.ID
	GuildID   snowflake.ID
	ChannelID snowflake.ID
	Data      discord.Channel
	CreatedAt time.Time
	UpdatedAt time.Time
}

type SearchChannelsParams struct {
	AppID   snowflake.ID
	GuildID snowflake.ID
	Limit   int
	Offset  int
	Data    json.RawMessage
}

type CacheChannelStore interface {
	GetChannel(ctx context.Context, appID snowflake.ID, guildID snowflake.ID, channelID snowflake.ID) (*model.Channel, error)
	GetChannels(ctx context.Context, appID snowflake.ID, guildID snowflake.ID, limit int, offset int) ([]*model.Channel, error)
	GetChannelsByType(ctx context.Context, appID snowflake.ID, guildID snowflake.ID, types []int, limit int, offset int) ([]*model.Channel, error)
	SearchChannels(ctx context.Context, params SearchChannelsParams) ([]*model.Channel, error)
	UpsertChannels(ctx context.Context, channels ...UpsertChannelParams) error
	DeleteChannel(ctx context.Context, appID snowflake.ID, guildID snowflake.ID, channelID snowflake.ID) error
}
