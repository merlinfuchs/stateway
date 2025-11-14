package compat

import (
	"context"
	"iter"
	"log/slog"
	"time"

	discache "github.com/disgoorg/disgo/cache"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/snowflake/v2"
	"github.com/merlinfuchs/stateway/stateway-lib/cache"
	"github.com/merlinfuchs/stateway/stateway-lib/service"
)

var _ discache.Caches = &DisgoCaches{}

type DisgoCaches struct {
	discache.Caches
}

func NewDisgoCaches(ctx context.Context, cache cache.Cache) *DisgoCaches {
	return &DisgoCaches{Caches: discache.New(
		discache.WithGuildCache(&GuildCache{ctx: ctx, cache: cache}),
		discache.WithChannelCache(&ChannelCache{ctx: ctx, cache: cache}),
		discache.WithRoleCache(&RoleCache{ctx: ctx, cache: cache}),
	)}
}

type GuildCache struct {
	ctx   context.Context
	cache cache.GuildCache
}

func (c *GuildCache) GuildCache() discache.Cache[discord.Guild] {
	return &anyCache[discord.Guild]{
		getFunc: c.Guild,
		allFunc: c.Guilds,
		lenFunc: c.GuildsLen,
	}
}

func (c *GuildCache) IsGuildUnready(guildID snowflake.ID) bool {
	return false
}

func (c *GuildCache) SetGuildUnready(guildID snowflake.ID, unready bool) {
}

func (c *GuildCache) UnreadyGuildIDs() []snowflake.ID {
	return nil
}

func (c *GuildCache) IsGuildUnavailable(guildID snowflake.ID) bool {
	return false
}

func (c *GuildCache) SetGuildUnavailable(guildID snowflake.ID, unavailable bool) {
}

func (c *GuildCache) UnavailableGuildIDs() []snowflake.ID {
	return nil
}

func (c *GuildCache) Guild(guildID snowflake.ID) (discord.Guild, bool) {
	ctx, cancel := cacheCtx(c.ctx)
	defer cancel()

	guild, err := c.cache.GetGuild(ctx, guildID)
	if err != nil {
		if !service.IsErrorCode(err, service.ErrorCodeNotFound) {
			slog.Error(
				"Failed to get guild from cache",
				slog.String("guild_id", guildID.String()),
				slog.Any("error", err),
			)
		}
		return discord.Guild{}, false
	}

	return guild.Data, true
}

func (c *GuildCache) Guilds() iter.Seq[discord.Guild] {
	var offset int

	return func(fn func(discord.Guild) bool) {
		for {
			ctx, cancel := cacheCtx(c.ctx)

			guilds, err := c.cache.GetGuilds(ctx, cache.WithLimit(100), cache.WithOffset(offset))
			if err != nil {
				slog.Error(
					"Failed to get guilds from cache",
					slog.Any("error", err),
				)
			}

			cancel()

			for _, guild := range guilds {
				if fn(guild.Data) {
					return
				}
			}

			if len(guilds) < 100 {
				break
			}

			offset += 100
		}
	}
}

func (c *GuildCache) GuildsLen() int {
	return 0 // TODO: Implement cache method for this
}

func (c *GuildCache) AddGuild(guild discord.Guild) {
}

func (c *GuildCache) RemoveGuild(guildID snowflake.ID) (discord.Guild, bool) {
	return discord.Guild{}, false
}

type ChannelCache struct {
	ctx   context.Context
	cache cache.ChannelCache
}

func (c *ChannelCache) ChannelCache() discache.Cache[discord.GuildChannel] {
	return &anyCache[discord.GuildChannel]{
		getFunc: c.Channel,
		allFunc: c.Channels,
		lenFunc: c.ChannelsLen,
	}
}

func (c *ChannelCache) Channel(channelID snowflake.ID) (discord.GuildChannel, bool) {
	ctx, cancel := cacheCtx(c.ctx)
	defer cancel()

	channel, err := c.cache.GetChannel(ctx, channelID)
	if err != nil {
		if !service.IsErrorCode(err, service.ErrorCodeNotFound) {
			slog.Error(
				"Failed to get channel from cache",
				slog.String("channel_id", channelID.String()),
				slog.Any("error", err),
			)
		}

		return nil, false
	}

	guildChannel, ok := channel.Data.(discord.GuildChannel)
	if !ok {
		return nil, false
	}

	return guildChannel, true
}

func (c *ChannelCache) Channels() iter.Seq[discord.GuildChannel] {
	var offset int

	return func(fn func(discord.GuildChannel) bool) {
		for {
			ctx, cancel := cacheCtx(c.ctx)

			channels, err := c.cache.GetChannels(ctx, cache.WithLimit(100), cache.WithOffset(offset))
			if err != nil {
				slog.Error(
					"Failed to get channels from cache",
					slog.Any("error", err),
				)
			}

			cancel()

			for _, channel := range channels {
				guildChannel, ok := channel.Data.(discord.GuildChannel)
				if !ok {
					continue
				}

				if fn(guildChannel) {
					return
				}
			}

			if len(channels) < 100 {
				break
			}

			offset += 100
		}
	}
}

func (c *ChannelCache) ChannelsForGuild(guildID snowflake.ID) iter.Seq[discord.GuildChannel] {
	ctx, cancel := cacheCtx(c.ctx)
	defer cancel()

	channels, err := c.cache.GetGuildChannels(ctx, guildID)
	if err != nil {
		slog.Error(
			"Failed to get channels from cache",
			slog.String("guild_id", guildID.String()),
			slog.Any("error", err),
		)
	}

	return func(fn func(discord.GuildChannel) bool) {
		for _, channel := range channels {
			guildChannel, ok := channel.Data.(discord.GuildChannel)
			if !ok {
				continue
			}

			if fn(guildChannel) {
				return
			}
		}
	}
}

func (c *ChannelCache) ChannelsLen() int {
	ctx, cancel := cacheCtx(c.ctx)
	defer cancel()

	count, err := c.cache.CountChannels(ctx)
	if err != nil {
		slog.Error(
			"Failed to get channels count from cache",
			slog.Any("error", err),
		)
	}

	return count
}

func (c *ChannelCache) AddChannel(channel discord.GuildChannel) {
}

func (c *ChannelCache) RemoveChannel(channelID snowflake.ID) (discord.GuildChannel, bool) {
	return nil, false
}

func (c *ChannelCache) RemoveChannelsByGuildID(guildID snowflake.ID) {
}

type RoleCache struct {
	ctx   context.Context
	cache cache.RoleCache
}

func (c *RoleCache) RoleCache() discache.GroupedCache[discord.Role] {
	return &groupCache[discord.Role]{
		getFunc:      c.Role,
		allFunc:      func() iter.Seq2[snowflake.ID, discord.Role] { return nil },
		lenFunc:      c.RolesAllLen,
		groupAllFunc: c.Roles,
		groupLenFunc: c.RolesLen,
	}
}

func (c *RoleCache) Role(guildID snowflake.ID, roleID snowflake.ID) (discord.Role, bool) {
	ctx, cancel := cacheCtx(c.ctx)
	defer cancel()

	role, err := c.cache.GetGuildRole(ctx, guildID, roleID)
	if err != nil {
		if !service.IsErrorCode(err, service.ErrorCodeNotFound) {
			slog.Error(
				"Failed to get role from cache",
				slog.String("guild_id", guildID.String()),
				slog.String("role_id", roleID.String()),
				slog.Any("error", err),
			)
		}
		return discord.Role{}, false
	}

	return role.Data, true
}

func (c *RoleCache) Roles(guildID snowflake.ID) iter.Seq[discord.Role] {
	ctx, cancel := cacheCtx(c.ctx)
	defer cancel()

	roles, err := c.cache.GetGuildRoles(ctx, guildID)
	if err != nil {
		slog.Error(
			"Failed to get roles from cache",
			slog.String("guild_id", guildID.String()),
			slog.Any("error", err),
		)
	}

	return func(fn func(discord.Role) bool) {
		for _, role := range roles {
			if fn(role.Data) {
				return
			}
		}
	}
}

func (c *RoleCache) RolesLen(groupID snowflake.ID) int {
	ctx, cancel := cacheCtx(c.ctx)
	defer cancel()

	count, err := c.cache.CountGuildRoles(ctx, groupID)
	if err != nil {
		slog.Error(
			"Failed to get roles count from cache",
			slog.String("guild_id", groupID.String()),
			slog.Any("error", err),
		)
	}

	return count
}

func (c *RoleCache) RolesAllLen() int {
	ctx, cancel := cacheCtx(c.ctx)
	defer cancel()

	count, err := c.cache.CountRoles(ctx)
	if err != nil {
		slog.Error(
			"Failed to get roles count from cache",
			slog.Any("error", err),
		)
	}

	return count
}

func (c *RoleCache) AddRole(role discord.Role) {
}

func (c *RoleCache) RemoveRole(guildID snowflake.ID, roleID snowflake.ID) (discord.Role, bool) {
	return discord.Role{}, false
}

func (c *RoleCache) RemoveRolesByGuildID(guildID snowflake.ID) {
}

type anyCache[T any] struct {
	getFunc func(id snowflake.ID) (T, bool)
	allFunc func() iter.Seq[T]
	lenFunc func() int
}

func (c *anyCache[T]) Get(id snowflake.ID) (T, bool) {
	return c.getFunc(id)
}

func (c *anyCache[T]) All() iter.Seq[T] {
	return c.allFunc()
}

func (c *anyCache[T]) Put(_ snowflake.ID, _ T) {}

func (c *anyCache[T]) Remove(_ snowflake.ID) (T, bool) {
	var v T
	return v, false
}

func (c *anyCache[T]) RemoveIf(_ discache.FilterFunc[T]) {}

func (c *anyCache[T]) Len() int {
	return c.lenFunc()
}

func cacheCtx(ctx context.Context) (context.Context, context.CancelFunc) {
	return context.WithTimeout(ctx, 3*time.Second)
}

type groupCache[T any] struct {
	getFunc      func(groupID snowflake.ID, id snowflake.ID) (T, bool)
	allFunc      func() iter.Seq2[snowflake.ID, T]
	lenFunc      func() int
	groupLenFunc func(groupID snowflake.ID) int
	groupAllFunc func(groupID snowflake.ID) iter.Seq[T]
}

func (c *groupCache[T]) Get(groupID snowflake.ID, id snowflake.ID) (T, bool) {
	return c.getFunc(groupID, id)
}

func (c *groupCache[T]) Put(groupID snowflake.ID, id snowflake.ID, entity T) {
}

func (c *groupCache[T]) Remove(groupID snowflake.ID, id snowflake.ID) (T, bool) {
	var v T
	return v, false
}

func (c *groupCache[T]) GroupRemove(groupID snowflake.ID) {
}

func (c *groupCache[T]) RemoveIf(filterFunc discache.GroupedFilterFunc[T]) {
}

func (c *groupCache[T]) GroupRemoveIf(groupID snowflake.ID, filterFunc discache.GroupedFilterFunc[T]) {
}

func (c *groupCache[T]) Len() int {
	return c.lenFunc()
}

func (c *groupCache[T]) GroupLen(groupID snowflake.ID) int {
	return c.groupLenFunc(groupID)
}

func (c *groupCache[T]) All() iter.Seq2[snowflake.ID, T] {
	return c.allFunc()
}

func (c *groupCache[T]) GroupAll(groupID snowflake.ID) iter.Seq[T] {
	return c.groupAllFunc(groupID)
}
