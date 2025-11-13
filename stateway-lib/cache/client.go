package cache

import (
	"context"
	"encoding/json"

	"github.com/disgoorg/snowflake/v2"
	"github.com/merlinfuchs/stateway/stateway-lib/broker"
	"github.com/merlinfuchs/stateway/stateway-lib/service"
)

var _ Cache = &CacheClient{}

type CacheClient struct {
	b       broker.Broker
	options CacheOptions
}

func NewCacheClient(b broker.Broker, opts ...CacheOption) *CacheClient {
	return &CacheClient{b: b, options: ResolveOptions(opts...)}
}

func (c *CacheClient) WithOptions(opts ...CacheOption) *CacheClient {
	return &CacheClient{
		b:       c.b,
		options: ResolveOptions(opts...),
	}
}

func (c *CacheClient) GetGuild(ctx context.Context, id snowflake.ID, opts ...CacheOption) (*Guild, error) {
	options := c.options
	for _, opt := range opts {
		opt(&options)
	}

	return cacheRequest[*Guild](ctx, c.b, CacheMethodGetGuild, GuildGetRequest{
		GuildID: id,
		Options: options,
	})
}

func (c *CacheClient) GetGuilds(ctx context.Context, opts ...CacheOption) ([]*Guild, error) {
	options := c.options
	for _, opt := range opts {
		opt(&options)
	}

	return cacheRequest[[]*Guild](ctx, c.b, CacheMethodListGuilds, GuildListRequest{
		Options: options,
	})
}

func (c *CacheClient) SearchGuilds(ctx context.Context, data json.RawMessage, opts ...CacheOption) ([]*Guild, error) {
	options := c.options
	for _, opt := range opts {
		opt(&options)
	}

	return cacheRequest[[]*Guild](ctx, c.b, CacheMethodSearchGuilds, GuildSearchRequest{
		Data:    data,
		Options: options,
	})
}

func (c *CacheClient) GetChannel(ctx context.Context, guildID snowflake.ID, channelID snowflake.ID, opts ...CacheOption) (*Channel, error) {
	options := c.options
	for _, opt := range opts {
		opt(&options)
	}

	return cacheRequest[*Channel](ctx, c.b, CacheMethodGetChannel, ChannelGetRequest{
		GuildID:   guildID,
		ChannelID: channelID,
		Options:   options,
	})
}

func (c *CacheClient) GetChannels(ctx context.Context, guildID snowflake.ID, opts ...CacheOption) ([]*Channel, error) {
	options := c.options
	for _, opt := range opts {
		opt(&options)
	}

	return cacheRequest[[]*Channel](ctx, c.b, CacheMethodListChannels, ChannelListRequest{
		GuildID: guildID,
		Options: options,
	})
}

func (c *CacheClient) SearchChannels(ctx context.Context, data json.RawMessage, opts ...CacheOption) ([]*Channel, error) {
	options := c.options
	for _, opt := range opts {
		opt(&options)
	}

	return cacheRequest[[]*Channel](ctx, c.b, CacheMethodSearchChannels, ChannelSearchRequest{
		Data:    data,
		Options: options,
	})
}

func (c *CacheClient) GetRole(ctx context.Context, guildID snowflake.ID, roleID snowflake.ID, opts ...CacheOption) (*Role, error) {
	options := c.options
	for _, opt := range opts {
		opt(&options)
	}

	return cacheRequest[*Role](ctx, c.b, CacheMethodGetRole, RoleGetRequest{
		GuildID: guildID,
		RoleID:  roleID,
		Options: options,
	})
}

func (c *CacheClient) GetRoles(ctx context.Context, guildID snowflake.ID, opts ...CacheOption) ([]*Role, error) {
	options := c.options
	for _, opt := range opts {
		opt(&options)
	}

	return cacheRequest[[]*Role](ctx, c.b, CacheMethodListRoles, RoleListRequest{
		GuildID: guildID,
		Options: options,
	})
}

func (c *CacheClient) SearchRoles(ctx context.Context, data json.RawMessage, opts ...CacheOption) ([]*Role, error) {
	options := c.options
	for _, opt := range opts {
		opt(&options)
	}

	return cacheRequest[[]*Role](ctx, c.b, CacheMethodSearchRoles, RoleSearchRequest{
		Data:    data,
		Options: options,
	})
}

func cacheRequest[R any](ctx context.Context, b broker.Broker, method CacheMethod, request CacheRequest) (R, error) {
	var r R

	response, err := b.Request(ctx, service.ServiceTypeCache, string(method), request)
	if err != nil {
		return r, err
	}

	err = json.Unmarshal(response.Data, &r)
	if err != nil {
		return r, err
	}

	return r, nil
}
