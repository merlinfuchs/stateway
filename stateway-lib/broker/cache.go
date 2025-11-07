package broker

import (
	"encoding/json"
	"fmt"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/snowflake/v2"
	"github.com/merlinfuchs/stateway/stateway-lib/cache"
)

type CacheMethod string

const (
	CacheMethodGetGuild CacheMethod = "get_guild"
)

type CacheRequest struct {
	EntityID *snowflake.ID      `json:"entity_id"`
	Options  cache.CacheOptions `json:"options"`
}

type CacheClient struct {
	b Broker
}

func (c *CacheClient) GetGuild(id snowflake.ID, opts ...cache.CacheOption) (*discord.Guild, error) {
	options := resolveCacheOptions(opts...)

	return cacheRequest[discord.Guild](c.b, CacheMethodGetGuild, CacheRequest{
		EntityID: &id,
		Options:  options,
	})
}

func resolveCacheOptions(opts ...cache.CacheOption) cache.CacheOptions {
	options := cache.CacheOptions{}
	for _, opt := range opts {
		opt(&options)
	}
	return options
}

func cacheRequest[R any](b Broker, method CacheMethod, request any) (*R, error) {
	response, err := b.Request(ServiceTypeCache, string(method), request)
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

func (s *CacheService) ServiceType() ServiceType {
	return ServiceTypeCache
}

func (s *CacheService) HandleRequest(method CacheMethod, request CacheRequest) (any, error) {
	switch method {
	case CacheMethodGetGuild:
		return s.caches.GetGuild(*request.EntityID)
	}
	return nil, fmt.Errorf("unknown method: %s", method)
}
