package broker

import (
	"encoding/json"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/snowflake/v2"
)

type CacheMethod string

const (
	CacheMethodGetGuild CacheMethod = "get_guild"
)

type CacheRequest struct {
	EntityID *snowflake.ID
}

type BrokerCaches struct {
	b Broker
}

func (c *BrokerCaches) GetGuild(id snowflake.ID) (*discord.Guild, error) {
	return cacheRequest[discord.Guild](c.b, CacheMethodGetGuild, CacheRequest{EntityID: &id})
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
