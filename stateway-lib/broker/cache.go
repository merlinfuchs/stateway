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
	CacheMethodGetGuild CacheMethod = "guild.get"
)

type CacheRequest struct {
	EntityID *snowflake.ID      `json:"entity_id"`
	Options  cache.CacheOptions `json:"options"`
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
		EntityID: &id,
		Options:  options,
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
		return s.caches.GetGuild(ctx, *request.EntityID, request.Options.Destructure()...)
	}
	return nil, fmt.Errorf("unknown method: %s", method)
}
