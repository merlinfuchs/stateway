package cache

import (
	"context"
	"encoding/json"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/snowflake/v2"
)

type GuildCache interface {
	GetGuild(ctx context.Context, id snowflake.ID, opts ...CacheOption) (*discord.Guild, error)
	GetGuilds(ctx context.Context, opts ...CacheOption) ([]*discord.Guild, error)
	SearchGuilds(ctx context.Context, data json.RawMessage, opts ...CacheOption) ([]*discord.Guild, error)
}

type ChannelCache interface {
	GetChannel(ctx context.Context, guildID snowflake.ID, channelID snowflake.ID, opts ...CacheOption) (*discord.Channel, error)
	GetChannels(ctx context.Context, guildID snowflake.ID, opts ...CacheOption) ([]*discord.Channel, error)
	SearchChannels(ctx context.Context, data json.RawMessage, opts ...CacheOption) ([]*discord.Channel, error)
}

type RoleCache interface {
	GetRole(ctx context.Context, guildID snowflake.ID, roleID snowflake.ID, opts ...CacheOption) (*discord.Role, error)
	GetRoles(ctx context.Context, guildID snowflake.ID, opts ...CacheOption) ([]*discord.Role, error)
	SearchRoles(ctx context.Context, data json.RawMessage, opts ...CacheOption) ([]*discord.Role, error)
}
