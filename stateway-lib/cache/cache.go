package cache

import (
	"context"
	"encoding/json"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/snowflake/v2"
)

type Cache interface {
	GuildCache
	ChannelCache
	RoleCache
	EmojiCache
	StickerCache
}

type GuildCache interface {
	GetGuild(ctx context.Context, id snowflake.ID, opts ...CacheOption) (*Guild, error)
	GetGuildWithPermissions(
		ctx context.Context,
		id snowflake.ID,
		userID snowflake.ID,
		roleIDs []snowflake.ID,
		abortAtPermissions discord.Permissions,
		opts ...CacheOption,
	) (*GuildWithPermissions, error)
	GetGuilds(ctx context.Context, opts ...CacheOption) ([]*Guild, error)
	CheckGuildsExist(ctx context.Context, guildIDs []snowflake.ID, opts ...CacheOption) ([]bool, error)
	SearchGuilds(ctx context.Context, data json.RawMessage, opts ...CacheOption) ([]*Guild, error)
	ComputeGuildPermissions(
		ctx context.Context,
		guildID snowflake.ID,
		userID snowflake.ID,
		roleIDs []snowflake.ID,
		opts ...CacheOption,
	) (discord.Permissions, error)
}

type ChannelCache interface {
	GetChannel(ctx context.Context, channelID snowflake.ID, opts ...CacheOption) (*Channel, error)
	GetChannels(ctx context.Context, opts ...CacheOption) ([]*Channel, error)
	CountChannels(ctx context.Context, opts ...CacheOption) (int, error)
	GetGuildChannel(ctx context.Context, guildID snowflake.ID, channelID snowflake.ID, opts ...CacheOption) (*Channel, error)
	GetGuildChannels(ctx context.Context, guildID snowflake.ID, opts ...CacheOption) ([]*Channel, error)
	GetGuildChannelsWithPermissions(
		ctx context.Context,
		guildID snowflake.ID,
		userID snowflake.ID,
		roleIDs []snowflake.ID,
		opts ...CacheOption,
	) ([]*ChannelWithPermissions, error)
	CountGuildChannels(ctx context.Context, guildID snowflake.ID, opts ...CacheOption) (int, error)
	SearchChannels(ctx context.Context, data json.RawMessage, opts ...CacheOption) ([]*Channel, error)
	SearchGuildChannels(ctx context.Context, guildID snowflake.ID, data json.RawMessage, opts ...CacheOption) ([]*Channel, error)
	ComputeChannelPermissions(
		ctx context.Context,
		channelID snowflake.ID,
		userID snowflake.ID,
		roleIDs []snowflake.ID,
		opts ...CacheOption,
	) (discord.Permissions, error)
	MassComputeChannelPermissions(
		ctx context.Context,
		guildID snowflake.ID,
		channelIDs []snowflake.ID,
		userID snowflake.ID,
		roleIDs []snowflake.ID,
		opts ...CacheOption,
	) ([]discord.Permissions, error)
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

type EmojiCache interface {
	GetGuildEmojis(ctx context.Context, guildID snowflake.ID, opts ...CacheOption) ([]*Emoji, error)
}

type StickerCache interface {
	GetGuildStickers(ctx context.Context, guildID snowflake.ID, opts ...CacheOption) ([]*Sticker, error)
}
