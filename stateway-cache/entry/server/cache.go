package server

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/snowflake/v2"
	"github.com/merlinfuchs/stateway/stateway-cache/store"
	"github.com/merlinfuchs/stateway/stateway-lib/cache"
	"github.com/merlinfuchs/stateway/stateway-lib/service"
)

var _ cache.Cache = (*Cache)(nil)

type Cache struct {
	cacheStore store.CacheStore
}

func NewCaches(cacheStore store.CacheStore) *Cache {
	return &Cache{
		cacheStore: cacheStore,
	}
}

func (c *Cache) GetGuild(ctx context.Context, id snowflake.ID, opts ...cache.CacheOption) (*cache.Guild, error) {
	options := cache.ResolveOptions(opts...)

	guild, err := c.cacheStore.GetGuild(ctx, options.AppID, id)
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			return nil, service.ErrNotFound("guild not found")
		}
		return nil, err
	}

	return guild, nil
}

func (c *Cache) GetGuildWithPermissions(
	ctx context.Context,
	guildID snowflake.ID,
	userID snowflake.ID,
	roleIDs []snowflake.ID,
	abortAtPermissions discord.Permissions,
	opts ...cache.CacheOption,
) (*cache.GuildWithPermissions, error) {
	options := cache.ResolveOptions(opts...)

	guild, err := c.cacheStore.GetGuild(ctx, options.AppID, guildID)
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			return nil, service.ErrNotFound("guild not found")
		}
		return nil, err
	}

	guildPermissions, err := c.ComputeGuildPermissions(ctx, guildID, userID, roleIDs, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to compute guild permissions: %w", err)
	}

	if abortAtPermissions != 0 && guildPermissions.Has(abortAtPermissions) {
		return &cache.GuildWithPermissions{
			Guild:                 *guild,
			GuildPermissions:      guildPermissions,
			MaxChannelPermissions: guildPermissions,
		}, nil
	}

	channels, err := c.cacheStore.GetGuildChannels(ctx, options.AppID, guildID, 0, 0)
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			return nil, service.ErrNotFound("channels not found")
		}
		return nil, err
	}

	maxChannelPermissions := discord.PermissionsNone
	minChannelPermissions := discord.PermissionsAll

	for _, channel := range channels {
		guildChannel, ok := channel.Data.(discord.GuildChannel)
		if !ok {
			continue
		}

		permissions := computeChannelPermissions(guildChannel, userID, roleIDs, guildPermissions)
		maxChannelPermissions |= permissions
		minChannelPermissions &= permissions

		if abortAtPermissions != 0 && permissions.Has(abortAtPermissions) {
			break
		}
	}

	return &cache.GuildWithPermissions{
		Guild:                 *guild,
		GuildPermissions:      guildPermissions,
		MaxChannelPermissions: maxChannelPermissions,
		MinChannelPermissions: minChannelPermissions,
	}, nil
}

func (c *Cache) GetGuilds(ctx context.Context, opts ...cache.CacheOption) ([]*cache.Guild, error) {
	options := cache.ResolveOptions(opts...)

	guilds, err := c.cacheStore.GetGuilds(ctx, options.AppID, options.Limit, options.Offset)
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			return nil, service.ErrNotFound("guilds not found")
		}
		return nil, err
	}

	return guilds, nil
}

func (c *Cache) CheckGuildsExist(ctx context.Context, guildIDs []snowflake.ID, opts ...cache.CacheOption) ([]bool, error) {
	options := cache.ResolveOptions(opts...)

	res := make([]bool, len(guildIDs))
	for i, guildID := range guildIDs {
		exists, err := c.cacheStore.CheckGuildExist(ctx, options.AppID, guildID)
		if err != nil {
			return nil, err
		}
		res[i] = exists
	}

	return res, nil
}

func (c *Cache) SearchGuilds(ctx context.Context, data json.RawMessage, opts ...cache.CacheOption) ([]*cache.Guild, error) {
	options := cache.ResolveOptions(opts...)

	guilds, err := c.cacheStore.SearchGuilds(ctx, store.SearchGuildsParams{
		AppID:  options.AppID,
		Limit:  options.Limit,
		Offset: options.Offset,
		Data:   data,
	})
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			return nil, service.ErrNotFound("guilds not found")
		}
		return nil, err
	}

	return guilds, nil
}

func (c *Cache) ComputeGuildPermissions(
	ctx context.Context,
	guildID snowflake.ID,
	userID snowflake.ID,
	roleIDs []snowflake.ID,
	opts ...cache.CacheOption,
) (discord.Permissions, error) {
	options := cache.ResolveOptions(opts...)

	ownerID, err := c.cacheStore.GetGuildOwnerID(ctx, options.AppID, guildID)
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			return 0, service.ErrNotFound("guild not found")
		}
		return 0, fmt.Errorf("failed to get guild owner ID: %w", err)
	}

	if ownerID == userID {
		return discord.PermissionsAll, nil
	}

	defaultRole, err := c.cacheStore.GetGuildRole(ctx, options.AppID, guildID, guildID)
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			return 0, service.ErrNotFound("default role not found")
		}
		return 0, fmt.Errorf("failed to get default role: %w", err)
	}

	permissions := defaultRole.Data.Permissions

	roles, err := c.cacheStore.GetGuildRolesByIDs(ctx, options.AppID, guildID, roleIDs)
	if err != nil {
		return 0, fmt.Errorf("failed to get guild roles by IDs: %w", err)
	}

	for _, role := range roles {
		permissions |= role.Data.Permissions
		if permissions.Has(discord.PermissionAdministrator) {
			return discord.PermissionsAll, nil
		}
	}

	return permissions, nil
}

func (c *Cache) GetChannel(ctx context.Context, channelID snowflake.ID, opts ...cache.CacheOption) (*cache.Channel, error) {
	options := cache.ResolveOptions(opts...)

	channel, err := c.cacheStore.GetChannel(ctx, options.AppID, channelID)
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			return nil, service.ErrNotFound("channel not found")
		}
		return nil, err
	}

	return channel, nil
}

func (c *Cache) GetChannels(ctx context.Context, opts ...cache.CacheOption) ([]*cache.Channel, error) {
	options := cache.ResolveOptions(opts...)

	channels, err := c.cacheStore.GetChannels(ctx, options.AppID, options.Limit, options.Offset)
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			return nil, service.ErrNotFound("channels not found")
		}
		return nil, err
	}

	return channels, nil
}

func (c *Cache) CountChannels(ctx context.Context, opts ...cache.CacheOption) (int, error) {
	options := cache.ResolveOptions(opts...)

	count, err := c.cacheStore.CountChannels(ctx, options.AppID)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (c *Cache) GetGuildChannel(ctx context.Context, guildID snowflake.ID, channelID snowflake.ID, opts ...cache.CacheOption) (*cache.Channel, error) {
	options := cache.ResolveOptions(opts...)

	channel, err := c.cacheStore.GetGuildChannel(ctx, options.AppID, guildID, channelID)
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			return nil, service.ErrNotFound("channel not found")
		}
		return nil, err
	}

	return channel, nil
}

func (c *Cache) GetGuildChannels(ctx context.Context, guildID snowflake.ID, opts ...cache.CacheOption) ([]*cache.Channel, error) {
	options := cache.ResolveOptions(opts...)

	channels, err := c.cacheStore.GetGuildChannels(ctx, options.AppID, guildID, options.Limit, options.Offset)
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			return nil, service.ErrNotFound("channels not found")
		}
		return nil, err
	}

	return channels, nil
}

func (c *Cache) GetGuildChannelsWithPermissions(
	ctx context.Context,
	guildID snowflake.ID,
	userID snowflake.ID,
	roleIDs []snowflake.ID,
	opts ...cache.CacheOption,
) ([]*cache.ChannelWithPermissions, error) {
	options := cache.ResolveOptions(opts...)

	channels, err := c.cacheStore.GetGuildChannels(ctx, options.AppID, guildID, options.Limit, options.Offset)
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			return nil, service.ErrNotFound("channels not found")
		}
		return nil, err
	}

	guildPermissions, err := c.ComputeGuildPermissions(ctx, guildID, userID, roleIDs, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to compute guild permissions: %w", err)
	}

	res := make([]*cache.ChannelWithPermissions, 0, len(channels))
	for _, channel := range channels {
		guildChannel, ok := channel.Data.(discord.GuildChannel)
		if !ok {
			continue
		}

		permissions := computeChannelPermissions(guildChannel, userID, roleIDs, guildPermissions)

		res = append(res, &cache.ChannelWithPermissions{
			Channel:     *channel,
			Permissions: permissions,
		})
	}

	return res, nil
}

func (c *Cache) CountGuildChannels(ctx context.Context, guildID snowflake.ID, opts ...cache.CacheOption) (int, error) {
	options := cache.ResolveOptions(opts...)

	count, err := c.cacheStore.CountGuildChannels(ctx, options.AppID, guildID)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (c *Cache) SearchChannels(ctx context.Context, data json.RawMessage, opts ...cache.CacheOption) ([]*cache.Channel, error) {
	options := cache.ResolveOptions(opts...)

	channels, err := c.cacheStore.SearchChannels(ctx, store.SearchChannelsParams{
		AppID:  options.AppID,
		Limit:  options.Limit,
		Offset: options.Offset,
		Data:   data,
	})
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			return nil, service.ErrNotFound("channels not found")
		}
		return nil, err
	}

	return channels, nil
}

func (c *Cache) SearchGuildChannels(ctx context.Context, guildID snowflake.ID, data json.RawMessage, opts ...cache.CacheOption) ([]*cache.Channel, error) {
	options := cache.ResolveOptions(opts...)

	channels, err := c.cacheStore.SearchGuildChannels(ctx, store.SearchGuildChannelsParams{
		AppID:   options.AppID,
		GuildID: guildID,
		Limit:   options.Limit,
		Offset:  options.Offset,
		Data:    data,
	})
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			return nil, service.ErrNotFound("channels not found")
		}
		return nil, err
	}

	return channels, nil
}

func (c *Cache) ComputeChannelPermissions(
	ctx context.Context,
	channelID snowflake.ID,
	userID snowflake.ID,
	roleIDs []snowflake.ID,
	opts ...cache.CacheOption,
) (discord.Permissions, error) {
	options := cache.ResolveOptions(opts...)

	channel, err := c.cacheStore.GetChannel(ctx, options.AppID, channelID)
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			return 0, service.ErrNotFound("channel not found")
		}
		return 0, fmt.Errorf("failed to get channel: %w", err)
	}

	guildPermissions, err := c.ComputeGuildPermissions(ctx, channel.GuildID, userID, roleIDs, opts...)
	if err != nil {
		return 0, fmt.Errorf("failed to compute guild permissions: %w", err)
	}

	if guildPermissions.Has(discord.PermissionAdministrator) {
		return discord.PermissionsAll, nil
	}

	guildChannel, ok := channel.Data.(discord.GuildChannel)
	if !ok {
		return 0, fmt.Errorf("channel is not a guild channel: %w", err)
	}

	return computeChannelPermissions(guildChannel, userID, roleIDs, guildPermissions), nil
}

func (c *Cache) MassComputeChannelPermissions(
	ctx context.Context,
	guildID snowflake.ID,
	channelIDs []snowflake.ID,
	userID snowflake.ID,
	roleIDs []snowflake.ID,
	opts ...cache.CacheOption,
) ([]discord.Permissions, error) {
	options := cache.ResolveOptions(opts...)

	guildPermissions, err := c.ComputeGuildPermissions(ctx, guildID, userID, roleIDs, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to compute guild permissions: %w", err)
	}

	res := make([]discord.Permissions, len(channelIDs))
	for i, channelID := range channelIDs {
		if guildPermissions.Has(discord.PermissionAdministrator) {
			res[i] = discord.PermissionsAll
			continue
		}

		channel, err := c.cacheStore.GetChannel(ctx, options.AppID, channelID)
		if err != nil {
			continue
		}

		guildChannel, ok := channel.Data.(discord.GuildChannel)
		if !ok {
			continue
		}

		res[i] = computeChannelPermissions(guildChannel, userID, roleIDs, guildPermissions)
	}

	return res, nil
}

func (c *Cache) GetRole(ctx context.Context, roleID snowflake.ID, opts ...cache.CacheOption) (*cache.Role, error) {
	options := cache.ResolveOptions(opts...)

	role, err := c.cacheStore.GetRole(ctx, options.AppID, roleID)
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			return nil, service.ErrNotFound("role not found")
		}
		return nil, err
	}

	return role, nil
}

func (c *Cache) GetRoles(ctx context.Context, opts ...cache.CacheOption) ([]*cache.Role, error) {
	options := cache.ResolveOptions(opts...)

	roles, err := c.cacheStore.GetRoles(ctx, options.AppID, options.Limit, options.Offset)
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			return nil, service.ErrNotFound("roles not found")
		}
		return nil, err
	}

	return roles, nil
}

func (c *Cache) CountRoles(ctx context.Context, opts ...cache.CacheOption) (int, error) {
	options := cache.ResolveOptions(opts...)

	count, err := c.cacheStore.CountRoles(ctx, options.AppID)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (c *Cache) GetGuildRole(ctx context.Context, guildID snowflake.ID, roleID snowflake.ID, opts ...cache.CacheOption) (*cache.Role, error) {
	options := cache.ResolveOptions(opts...)

	role, err := c.cacheStore.GetGuildRole(ctx, options.AppID, guildID, roleID)
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			return nil, service.ErrNotFound("role not found")
		}
		return nil, err
	}

	return role, nil
}

func (c *Cache) GetGuildRoles(ctx context.Context, guildID snowflake.ID, opts ...cache.CacheOption) ([]*cache.Role, error) {
	options := cache.ResolveOptions(opts...)

	roles, err := c.cacheStore.GetGuildRoles(ctx, options.AppID, guildID, options.Limit, options.Offset)
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			return nil, service.ErrNotFound("roles not found")
		}
		return nil, err
	}

	return roles, nil
}

func (c *Cache) CountGuildRoles(ctx context.Context, guildID snowflake.ID, opts ...cache.CacheOption) (int, error) {
	options := cache.ResolveOptions(opts...)

	count, err := c.cacheStore.CountGuildRoles(ctx, options.AppID, guildID)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (c *Cache) SearchRoles(ctx context.Context, data json.RawMessage, opts ...cache.CacheOption) ([]*cache.Role, error) {
	options := cache.ResolveOptions(opts...)

	roles, err := c.cacheStore.SearchRoles(ctx, store.SearchRolesParams{
		AppID:  options.AppID,
		Limit:  options.Limit,
		Offset: options.Offset,
		Data:   data,
	})
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			return nil, service.ErrNotFound("roles not found")
		}
		return nil, err
	}

	return roles, nil
}

func (c *Cache) SearchGuildRoles(ctx context.Context, guildID snowflake.ID, data json.RawMessage, opts ...cache.CacheOption) ([]*cache.Role, error) {
	options := cache.ResolveOptions(opts...)

	roles, err := c.cacheStore.SearchGuildRoles(ctx, store.SearchGuildRolesParams{
		AppID:   options.AppID,
		GuildID: guildID,
		Limit:   options.Limit,
		Offset:  options.Offset,
		Data:    data,
	})
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			return nil, service.ErrNotFound("roles not found")
		}
		return nil, err
	}

	return roles, nil
}
