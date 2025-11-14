package server

import (
	"context"
	"encoding/json"
	"errors"

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
