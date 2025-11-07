package store

import (
	"context"
	"encoding/json"
	"time"

	"github.com/disgoorg/snowflake/v2"
	"github.com/merlinfuchs/stateway/stateway-cache/model"
)

type UpsertGuildParams struct {
	AppID     snowflake.ID
	GuildID   snowflake.ID
	Data      json.RawMessage
	CreatedAt time.Time
	UpdatedAt time.Time
}

type GuildIdentifier struct {
	AppID   snowflake.ID
	GuildID snowflake.ID
}

type CacheGuildStore interface {
	GetGuild(ctx context.Context, guild GuildIdentifier) (*model.Guild, error)
	UpsertGuilds(ctx context.Context, guilds ...UpsertGuildParams) error
	MarkGuildUnavailable(ctx context.Context, params GuildIdentifier) error
	DeleteGuild(ctx context.Context, params GuildIdentifier) error
}
