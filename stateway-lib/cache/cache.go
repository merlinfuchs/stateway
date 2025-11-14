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
	GetChannel(ctx context.Context, channelID snowflake.ID, opts ...CacheOption) (*Channel, error)
	GetChannels(ctx context.Context, opts ...CacheOption) ([]*Channel, error)
	CountChannels(ctx context.Context, opts ...CacheOption) (int, error)
	GetGuildChannel(ctx context.Context, guildID snowflake.ID, channelID snowflake.ID, opts ...CacheOption) (*Channel, error)
	GetGuildChannels(ctx context.Context, guildID snowflake.ID, opts ...CacheOption) ([]*Channel, error)
	CountGuildChannels(ctx context.Context, guildID snowflake.ID, opts ...CacheOption) (int, error)
	SearchChannels(ctx context.Context, data json.RawMessage, opts ...CacheOption) ([]*Channel, error)
	SearchGuildChannels(ctx context.Context, guildID snowflake.ID, data json.RawMessage, opts ...CacheOption) ([]*Channel, error)
}

type RoleCache interface {
	GetRole(ctx context.Context, roleID snowflake.ID, opts ...CacheOption) (*Role, error)
	GetRoles(ctx context.Context, opts ...CacheOption) ([]*Role, error)
	CountRoles(ctx context.Context, opts ...CacheOption) (int, error)
	GetGuildRole(ctx context.Context, guildID snowflake.ID, roleID snowflake.ID, opts ...CacheOption) (*Role, error)
	GetGuildRoles(ctx context.Context, guildID snowflake.ID, opts ...CacheOption) ([]*Role, error)
	CountGuildRoles(ctx context.Context, guildID snowflake.ID, opts ...CacheOption) (int, error)
	SearchRoles(ctx context.Context, data json.RawMessage, opts ...CacheOption) ([]*Role, error)
	SearchGuildRoles(ctx context.Context, guildID snowflake.ID, data json.RawMessage, opts ...CacheOption) ([]*Role, error)
}
