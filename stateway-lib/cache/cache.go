package cache

import (
	"context"
	"encoding/json"

	"github.com/disgoorg/snowflake/v2"
)

type Cache interface {
	GuildCache
	ChannelCache
	RoleCache
}

type GuildCache interface {
	GetGuild(ctx context.Context, id snowflake.ID, opts ...CacheOption) (*Guild, error)
	GetGuilds(ctx context.Context, opts ...CacheOption) ([]*Guild, error)
	SearchGuilds(ctx context.Context, data json.RawMessage, opts ...CacheOption) ([]*Guild, error)
}

type ChannelCache interface {
	GetChannel(ctx context.Context, guildID snowflake.ID, channelID snowflake.ID, opts ...CacheOption) (*Channel, error)
	GetChannels(ctx context.Context, guildID snowflake.ID, opts ...CacheOption) ([]*Channel, error)
	SearchChannels(ctx context.Context, data json.RawMessage, opts ...CacheOption) ([]*Channel, error)
}

type RoleCache interface {
	GetRole(ctx context.Context, guildID snowflake.ID, roleID snowflake.ID, opts ...CacheOption) (*Role, error)
	GetRoles(ctx context.Context, guildID snowflake.ID, opts ...CacheOption) ([]*Role, error)
	SearchRoles(ctx context.Context, data json.RawMessage, opts ...CacheOption) ([]*Role, error)
}
