package store

import (
	"context"
	"encoding/json"
	"time"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/snowflake/v2"
	"github.com/merlinfuchs/stateway/stateway-cache/model"
)

type UpsertGuildParams struct {
	AppID     snowflake.ID
	GuildID   snowflake.ID
	Data      discord.Guild
	CreatedAt time.Time
	UpdatedAt time.Time
}

type SearchGuildsParams struct {
	AppID  snowflake.ID
	Limit  int
	Offset int
	Data   json.RawMessage
}

type CacheGuildStore interface {
	GetGuild(ctx context.Context, appID snowflake.ID, guildID snowflake.ID) (*model.Guild, error)
	GetGuildOwnerID(ctx context.Context, appID snowflake.ID, guildID snowflake.ID) (snowflake.ID, error)
	GetGuilds(ctx context.Context, appID snowflake.ID, limit int, offset int) ([]*model.Guild, error)
	CheckGuildExist(ctx context.Context, appID snowflake.ID, guildID snowflake.ID) (bool, error)
	UpsertGuilds(ctx context.Context, guilds ...UpsertGuildParams) error
	MarkGuildUnavailable(ctx context.Context, appID snowflake.ID, guildID snowflake.ID) error
	DeleteGuild(ctx context.Context, appID snowflake.ID, guildID snowflake.ID) error
	SearchGuilds(ctx context.Context, params SearchGuildsParams) ([]*model.Guild, error)
}
