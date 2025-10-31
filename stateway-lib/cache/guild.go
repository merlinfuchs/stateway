package cache

import (
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/snowflake/v2"
)

type GuildCache interface {
	GetGuild(id snowflake.ID) (discord.Guild, error)
}
