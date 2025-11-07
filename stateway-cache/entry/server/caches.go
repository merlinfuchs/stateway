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

	mod, err := c.cacheStore.GetGuild(ctx, store.GuildIdentifier{
		AppID:   options.AppID,
		GuildID: id,
	})
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
