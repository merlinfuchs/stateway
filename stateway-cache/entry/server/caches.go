package server

import (
	"context"
	"encoding/json"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/snowflake/v2"
	"github.com/merlinfuchs/stateway/stateway-cache/store"
	"github.com/merlinfuchs/stateway/stateway-lib/cache"
)

type Caches struct {
	cacheStore store.CacheStore
}

func NewCaches(cacheStore store.CacheStore) *Caches {
	return &Caches{
		cacheStore: cacheStore,
	}
}

func (c *Caches) GetGuild(ctx context.Context, id snowflake.ID, opts ...cache.CacheOption) (*discord.Guild, error) {
	options := cache.ResolveOptions(opts...)

	mod, err := c.cacheStore.GetGuild(ctx, options.AppID, id)
	if err != nil {
		return nil, err
	}

	var guild discord.Guild
	err = json.Unmarshal(mod.Data, &guild)
	if err != nil {
		return nil, err
	}
	return &guild, nil
}

func (c *Caches) GetGuilds(ctx context.Context, opts ...cache.CacheOption) ([]*discord.Guild, error) {
	options := cache.ResolveOptions(opts...)

	mods, err := c.cacheStore.GetGuilds(ctx, options.AppID, options.Limit, options.Offset)
	if err != nil {
		return nil, err
	}

	guilds := make([]*discord.Guild, len(mods))
	for i, mod := range mods {
		var guild discord.Guild
		err = json.Unmarshal(mod.Data, &guild)
		if err != nil {
			return nil, err
		}
		guilds[i] = &guild
	}
	return guilds, nil
}

func (c *Caches) SearchGuilds(ctx context.Context, data json.RawMessage, opts ...cache.CacheOption) ([]*discord.Guild, error) {
	options := cache.ResolveOptions(opts...)

	mods, err := c.cacheStore.SearchGuilds(ctx, store.SearchGuildsParams{
		AppID:  options.AppID,
		Limit:  options.Limit,
		Offset: options.Offset,
		Data:   data,
	})
	if err != nil {
		return nil, err
	}

	guilds := make([]*discord.Guild, len(mods))
	for i, mod := range mods {
		var guild discord.Guild
		err = json.Unmarshal(mod.Data, &guild)
		if err != nil {
			return nil, err
		}
		guilds[i] = &guild
	}
	return guilds, nil
}

func (c *Caches) GetChannel(ctx context.Context, guildID snowflake.ID, channelID snowflake.ID, opts ...cache.CacheOption) (*discord.Channel, error) {
	options := cache.ResolveOptions(opts...)

	mod, err := c.cacheStore.GetChannel(ctx, options.AppID, guildID, channelID)
	if err != nil {
		return nil, err
	}

	var channel discord.Channel
	err = json.Unmarshal(mod.Data, &channel)
	if err != nil {
		return nil, err
	}
	return &channel, nil
}

func (c *Caches) GetChannels(ctx context.Context, guildID snowflake.ID, opts ...cache.CacheOption) ([]*discord.Channel, error) {
	options := cache.ResolveOptions(opts...)

	mods, err := c.cacheStore.GetChannels(ctx, options.AppID, guildID, options.Limit, options.Offset)
	if err != nil {
		return nil, err
	}

	channels := make([]*discord.Channel, len(mods))
	for i, mod := range mods {
		var channel discord.Channel
		err = json.Unmarshal(mod.Data, &channel)
		if err != nil {
			return nil, err
		}
		channels[i] = &channel
	}
	return channels, nil
}

func (c *Caches) GetChannelsByType(ctx context.Context, guildID snowflake.ID, types []int, opts ...cache.CacheOption) ([]*discord.Channel, error) {
	options := cache.ResolveOptions(opts...)

	mods, err := c.cacheStore.GetChannelsByType(ctx, options.AppID, guildID, types, options.Limit, options.Offset)
	if err != nil {
		return nil, err
	}

	channels := make([]*discord.Channel, len(mods))
	for i, mod := range mods {
		var channel discord.Channel
		err = json.Unmarshal(mod.Data, &channel)
		if err != nil {
			return nil, err
		}
		channels[i] = &channel
	}

	return channels, nil
}

func (c *Caches) SearchChannels(ctx context.Context, data json.RawMessage, opts ...cache.CacheOption) ([]*discord.Channel, error) {
	options := cache.ResolveOptions(opts...)

	mods, err := c.cacheStore.SearchChannels(ctx, store.SearchChannelsParams{
		AppID:  options.AppID,
		Limit:  options.Limit,
		Offset: options.Offset,
		Data:   data,
	})
	if err != nil {
		return nil, err
	}

	channels := make([]*discord.Channel, len(mods))
	for i, mod := range mods {
		var channel discord.Channel
		err = json.Unmarshal(mod.Data, &channel)
		if err != nil {
			return nil, err
		}
		channels[i] = &channel
	}
	return channels, nil
}

func (c *Caches) GetRole(ctx context.Context, guildID snowflake.ID, roleID snowflake.ID, opts ...cache.CacheOption) (*discord.Role, error) {
	options := cache.ResolveOptions(opts...)

	mod, err := c.cacheStore.GetRole(ctx, options.AppID, guildID, roleID)
	if err != nil {
		return nil, err
	}

	var role discord.Role
	err = json.Unmarshal(mod.Data, &role)
	if err != nil {
		return nil, err
	}
	return &role, nil
}

func (c *Caches) GetRoles(ctx context.Context, guildID snowflake.ID, opts ...cache.CacheOption) ([]*discord.Role, error) {
	options := cache.ResolveOptions(opts...)

	mods, err := c.cacheStore.GetRoles(ctx, options.AppID, guildID, options.Limit, options.Offset)
	if err != nil {
		return nil, err
	}

	roles := make([]*discord.Role, len(mods))
	for i, mod := range mods {
		var role discord.Role
		err = json.Unmarshal(mod.Data, &role)
		if err != nil {
			return nil, err
		}
		roles[i] = &role
	}
	return roles, nil
}

func (c *Caches) SearchRoles(ctx context.Context, data json.RawMessage, opts ...cache.CacheOption) ([]*discord.Role, error) {
	options := cache.ResolveOptions(opts...)

	mods, err := c.cacheStore.SearchRoles(ctx, store.SearchRolesParams{
		AppID:  options.AppID,
		Limit:  options.Limit,
		Offset: options.Offset,
		Data:   data,
	})
	if err != nil {
		return nil, err
	}

	roles := make([]*discord.Role, len(mods))
	for i, mod := range mods {
		var role discord.Role
		err = json.Unmarshal(mod.Data, &role)
		if err != nil {
			return nil, err
		}
		roles[i] = &role
	}
	return roles, nil
}
