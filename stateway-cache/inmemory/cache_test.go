package inmemory

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/snowflake/v2"
	"github.com/merlinfuchs/stateway/stateway-cache/store"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var testGuilds = []store.UpsertGuildParams{
	{
		AppID:   1,
		GuildID: 1,
		Data:    discord.Guild{},
	},
	{
		AppID:   1,
		GuildID: 2,
		Data:    discord.Guild{},
	},
	{
		AppID:   1,
		GuildID: 3,
		Data:    discord.Guild{},
	},
	{
		AppID:   1,
		GuildID: 4,
		Data:    discord.Guild{},
	},
	{
		AppID:   2,
		GuildID: 1,
		Data:    discord.Guild{},
	},
	{
		AppID:   2,
		GuildID: 2,
		Data:    discord.Guild{},
	},
	{
		AppID:   2,
		GuildID: 3,
		Data:    discord.Guild{},
	},
	{
		AppID:   2,
		GuildID: 5,
		Data:    discord.Guild{},
	},
	{
		AppID:   2,
		GuildID: 6,
		Data:    discord.Guild{},
	},
	{
		AppID:   2,
		GuildID: 7,
		Data:    discord.Guild{},
	},
}

func TestInMemoryGuildCache(t *testing.T) {
	cache, err := NewInMemoryCacheStore()
	if err != nil {
		t.Fatalf("failed to create in-memory cache store: %v", err)
	}

	err = cache.UpsertGuilds(context.Background(), testGuilds...)
	assert.NoError(t, err)

	guild, err := cache.GetGuild(context.Background(), 1, 1)
	assert.NoError(t, err)
	assert.Equal(t, guild.GuildID, snowflake.ID(1))
	assert.Equal(t, guild.AppID, snowflake.ID(1))
	assert.Equal(t, guild.Data, discord.Guild{})

	guilds, err := cache.GetGuilds(context.Background(), 1, 10, 0)
	assert.NoError(t, err)
	assert.Equal(t, len(guilds), 4)
	assert.Equal(t, guilds[0].GuildID, snowflake.ID(1))
	assert.Equal(t, guilds[0].AppID, snowflake.ID(1))
	assert.Equal(t, guilds[0].Data, discord.Guild{})

	guilds, err = cache.GetGuilds(context.Background(), 2, 10, 3)
	assert.NoError(t, err)
	assert.Equal(t, len(guilds), 3)
	assert.Equal(t, guilds[0].GuildID, snowflake.ID(5))
	assert.Equal(t, guilds[0].AppID, snowflake.ID(2))
	assert.Equal(t, guilds[0].Data, discord.Guild{})
}

func TestInMemoryCacheInsertPerformance(t *testing.T) {
	cache, err := NewInMemoryCacheStore()
	if err != nil {
		t.Fatalf("failed to create in-memory cache store: %v", err)
	}

	insertCount := 1000
	batchSize := 100
	totalCount := insertCount * len(testGuilds)

	roundCount := 100
	for r := 0; r < roundCount; r += 1 {
		start := time.Now()
		for i := 0; i < insertCount; i++ {
			data := make([]store.UpsertGuildParams, 0, len(testGuilds))
			for b := 0; b < batchSize; b++ {
				data = append(data, store.UpsertGuildParams{
					AppID:   snowflake.ID(r),
					GuildID: snowflake.ID(i),
					Data:    discord.Guild{},
				})
			}
			err = cache.UpsertGuilds(context.Background(), data...)
			require.NoError(t, err)
		}
		duration := time.Since(start)
		fmt.Printf("Upserted %d guilds in %v\n", totalCount, duration)
	}
}
