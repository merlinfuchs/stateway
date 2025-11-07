package store

import (
	"context"
	"encoding/json"
	"time"

	"github.com/disgoorg/snowflake/v2"
)

type UpsertChannelParams struct {
	AppID     snowflake.ID
	GuildID   snowflake.ID
	ChannelID snowflake.ID
	Data      json.RawMessage
	CreatedAt time.Time
	UpdatedAt time.Time
}

type ChannelIdentifier struct {
	AppID     snowflake.ID
	GuildID   snowflake.ID
	ChannelID snowflake.ID
}

type CacheChannelStore interface {
	UpsertChannels(ctx context.Context, channels ...UpsertChannelParams) error
	DeleteChannel(ctx context.Context, params ChannelIdentifier) error
}
