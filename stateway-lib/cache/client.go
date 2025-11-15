package cache

import (
	"context"
	"encoding/json"

	"github.com/disgoorg/disgo/discord"
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

func (c *CacheClient) CountGuilds(ctx context.Context, opts ...CacheOption) (int, error) {
	options := c.options
	for _, opt := range opts {
		opt(&options)
	}

	return cacheRequest[int](ctx, c.b, CacheMethodCountGuilds, GuildCountRequest{
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

func (c *CacheClient) ComputeGuildPermissions(
	ctx context.Context,
	guildID snowflake.ID,
	userID snowflake.ID,
	roleIDs []snowflake.ID,
	opts ...CacheOption,
) (discord.Permissions, error) {
	options := c.options
	for _, opt := range opts {
		opt(&options)
	}

	return cacheRequest[discord.Permissions](ctx, c.b, CacheMethodComputePermissions, PermissionsComputeRequest{
		GuildID: &guildID,
		UserID:  userID,
		RoleIDs: roleIDs,
		Options: options,
	})
}

func (c *CacheClient) GetChannel(ctx context.Context, channelID snowflake.ID, opts ...CacheOption) (*Channel, error) {
	options := c.options
	for _, opt := range opts {
		opt(&options)
	}

	return cacheRequest[*Channel](ctx, c.b, CacheMethodGetChannel, ChannelGetRequest{
		ChannelID: channelID,
		Options:   options,
	})
}

func (c *CacheClient) GetChannels(ctx context.Context, opts ...CacheOption) ([]*Channel, error) {
	options := c.options
	for _, opt := range opts {
		opt(&options)
	}

	return cacheRequest[[]*Channel](ctx, c.b, CacheMethodListChannels, ChannelListRequest{
		Options: options,
	})
}

func (c *CacheClient) CountChannels(ctx context.Context, opts ...CacheOption) (int, error) {
	options := c.options
	for _, opt := range opts {
		opt(&options)
	}

	return cacheRequest[int](ctx, c.b, CacheMethodCountChannels, ChannelCountRequest{
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

func (c *CacheClient) GetGuildChannel(ctx context.Context, guildID snowflake.ID, channelID snowflake.ID, opts ...CacheOption) (*Channel, error) {
	options := c.options
	for _, opt := range opts {
		opt(&options)
	}

	return cacheRequest[*Channel](ctx, c.b, CacheMethodGetChannel, ChannelGetRequest{
		GuildID:   &guildID,
		ChannelID: channelID,
		Options:   options,
	})
}

func (c *CacheClient) GetGuildChannels(ctx context.Context, guildID snowflake.ID, opts ...CacheOption) ([]*Channel, error) {
	options := c.options
	for _, opt := range opts {
		opt(&options)
	}

	return cacheRequest[[]*Channel](ctx, c.b, CacheMethodListChannels, ChannelListRequest{
		GuildID: &guildID,
		Options: options,
	})
}

func (c *CacheClient) SearchGuildChannels(ctx context.Context, guildID snowflake.ID, data json.RawMessage, opts ...CacheOption) ([]*Channel, error) {
	options := c.options
	for _, opt := range opts {
		opt(&options)
	}

	return cacheRequest[[]*Channel](ctx, c.b, CacheMethodSearchChannels, ChannelSearchRequest{
		GuildID: &guildID,
		Data:    data,
		Options: options,
	})
}

func (c *CacheClient) CountGuildChannels(ctx context.Context, guildID snowflake.ID, opts ...CacheOption) (int, error) {
	options := c.options
	for _, opt := range opts {
		opt(&options)
	}

	return cacheRequest[int](ctx, c.b, CacheMethodCountChannels, ChannelCountRequest{
		GuildID: &guildID,
		Options: options,
	})
}

func (c *CacheClient) ComputeChannelPermissions(
	ctx context.Context,
	channelID snowflake.ID,
	userID snowflake.ID,
	roleIDs []snowflake.ID,
	opts ...CacheOption,
) (discord.Permissions, error) {
	options := c.options
	for _, opt := range opts {
		opt(&options)
	}

	return cacheRequest[discord.Permissions](ctx, c.b, CacheMethodComputePermissions, PermissionsComputeRequest{
		ChannelID: &channelID,
		UserID:    userID,
		RoleIDs:   roleIDs,
		Options:   options,
	})
}

func (c *CacheClient) GetRole(ctx context.Context, roleID snowflake.ID, opts ...CacheOption) (*Role, error) {
	options := c.options
	for _, opt := range opts {
		opt(&options)
	}

	return cacheRequest[*Role](ctx, c.b, CacheMethodGetRole, RoleGetRequest{
		RoleID:  roleID,
		Options: options,
	})
}

func (c *CacheClient) GetRoles(ctx context.Context, opts ...CacheOption) ([]*Role, error) {
	options := c.options
	for _, opt := range opts {
		opt(&options)
	}

	return cacheRequest[[]*Role](ctx, c.b, CacheMethodListRoles, RoleListRequest{
		Options: options,
	})
}

func (c *CacheClient) CountRoles(ctx context.Context, opts ...CacheOption) (int, error) {
	options := c.options
	for _, opt := range opts {
		opt(&options)
	}

	return cacheRequest[int](ctx, c.b, CacheMethodCountRoles, RoleCountRequest{
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

func (c *CacheClient) GetGuildRole(ctx context.Context, guildID snowflake.ID, roleID snowflake.ID, opts ...CacheOption) (*Role, error) {
	options := c.options
	for _, opt := range opts {
		opt(&options)
	}

	return cacheRequest[*Role](ctx, c.b, CacheMethodGetRole, RoleGetRequest{
		GuildID: &guildID,
		RoleID:  roleID,
		Options: options,
	})
}

func (c *CacheClient) GetGuildRoles(ctx context.Context, guildID snowflake.ID, opts ...CacheOption) ([]*Role, error) {
	options := c.options
	for _, opt := range opts {
		opt(&options)
	}

	return cacheRequest[[]*Role](ctx, c.b, CacheMethodListRoles, RoleListRequest{
		GuildID: &guildID,
		Options: options,
	})
}

func (c *CacheClient) CountGuildRoles(ctx context.Context, guildID snowflake.ID, opts ...CacheOption) (int, error) {
	options := c.options
	for _, opt := range opts {
		opt(&options)
	}

	return cacheRequest[int](ctx, c.b, CacheMethodCountRoles, RoleCountRequest{
		GuildID: &guildID,
		Options: options,
	})
}

func (c *CacheClient) SearchGuildRoles(ctx context.Context, guildID snowflake.ID, data json.RawMessage, opts ...CacheOption) ([]*Role, error) {
	options := c.options
	for _, opt := range opts {
		opt(&options)
	}

	return cacheRequest[[]*Role](ctx, c.b, CacheMethodSearchRoles, RoleSearchRequest{
		GuildID: &guildID,
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

	if !response.Success {
		if response.Error != nil {
			return r, response.Error
		}
		return r, service.ErrUnknown("unknown error")
	}

	err = json.Unmarshal(response.Data, &r)
	if err != nil {
		return r, err
	}

	return r, nil
}
