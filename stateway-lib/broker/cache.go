package broker

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/snowflake/v2"
	"github.com/merlinfuchs/stateway/stateway-lib/cache"
	"github.com/merlinfuchs/stateway/stateway-lib/service"
)

type CacheMethod string

const (
	CacheMethodGetGuild       CacheMethod = "guild.get"
	CacheMethodListGuilds     CacheMethod = "guild.list"
	CacheMethodSearchGuilds   CacheMethod = "guild.search"
	CacheMethodGetChannel     CacheMethod = "channel.get"
	CacheMethodListChannels   CacheMethod = "channel.list"
	CacheMethodSearchChannels CacheMethod = "channel.search"
	CacheMethodGetRole        CacheMethod = "role.get"
	CacheMethodListRoles      CacheMethod = "role.list"
	CacheMethodSearchRoles    CacheMethod = "role.search"
)

type CacheRequest struct {
	GuildID   snowflake.ID       `json:"guild_id,omitempty"`
	ChannelID snowflake.ID       `json:"channel_id,omitempty"`
	RoleID    snowflake.ID       `json:"role_id,omitempty"`
	Types     []int              `json:"types,omitempty"`
	Data      json.RawMessage    `json:"data,omitempty"`
	Options   cache.CacheOptions `json:"options,omitempty"`
}

type CacheClient struct {
	b       Broker
	options cache.CacheOptions
}

func (c *CacheClient) WithOptions(opts ...cache.CacheOption) *CacheClient {
	return &CacheClient{
		b:       c.b,
		options: cache.ResolveOptions(opts...),
	}
}

func (c *CacheClient) GetGuild(ctx context.Context, id snowflake.ID, opts ...cache.CacheOption) (*discord.Guild, error) {
	options := c.options
	for _, opt := range opts {
		opt(&options)
	}

	return cacheRequest[discord.Guild](ctx, c.b, CacheMethodGetGuild, CacheRequest{
		GuildID: id,
		Options: options,
	})
}

func cacheRequest[R any](ctx context.Context, b Broker, method CacheMethod, request CacheRequest) (*R, error) {
	response, err := b.Request(ctx, service.ServiceTypeCache, string(method), request)
	if err != nil {
		return new(R), err
	}
	var r R
	err = json.Unmarshal(response.Data, &r)
	if err != nil {
		return new(R), err
	}
	return &r, nil
}

type CacheService struct {
	caches cache.Caches
}

func NewCacheService(caches cache.Caches) *CacheService {
	return &CacheService{caches: caches}
}

func (s *CacheService) ServiceType() service.ServiceType {
	return service.ServiceTypeCache
}

func (s *CacheService) HandleRequest(ctx context.Context, method string, request CacheRequest) (any, error) {
	switch CacheMethod(method) {
	case CacheMethodGetGuild:
		return s.caches.GetGuild(ctx, request.GuildID, request.Options.Destructure()...)
	case CacheMethodListGuilds:
		return s.caches.GetGuilds(ctx, request.Options.Destructure()...)
	case CacheMethodSearchGuilds:
		return s.caches.SearchGuilds(ctx, request.Data, request.Options.Destructure()...)
	case CacheMethodGetChannel:
		return s.caches.GetChannel(ctx, request.GuildID, request.ChannelID, request.Options.Destructure()...)
	case CacheMethodListChannels:
		return s.caches.GetChannels(ctx, request.GuildID, request.Options.Destructure()...)
	case CacheMethodSearchChannels:
		return s.caches.SearchChannels(ctx, request.Data, request.Options.Destructure()...)
	case CacheMethodGetRole:
		return s.caches.GetRole(ctx, request.GuildID, request.RoleID, request.Options.Destructure()...)
	case CacheMethodListRoles:
		return s.caches.GetRoles(ctx, request.GuildID, request.Options.Destructure()...)
	case CacheMethodSearchRoles:
		return s.caches.SearchRoles(ctx, request.Data, request.Options.Destructure()...)
	}
	return nil, fmt.Errorf("unknown method: %s", method)
}
