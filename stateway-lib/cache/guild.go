package cache

import (
	"context"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/snowflake/v2"
)

type GuildCache interface {
	GetGuild(ctx context.Context, id snowflake.ID, opts ...CacheOption) (*discord.Guild, error)
}
