package event

import (
	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/snowflake/v2"
)

type EventType string

const (
	EventTypeDiscordDispatch EventType = "discord_dispatch"
)

type Event interface {
	EventType() EventType
}

type DiscordDispatchEvent struct {
	AppID   snowflake.ID  `json:"app_id"`
	ShardID int           `json:"shard_id"`
	GuildID *snowflake.ID `json:"guild_id"`
	Data    bot.Event     `json:"data"`
}

func (e *DiscordDispatchEvent) EventType() EventType {
	return EventTypeDiscordDispatch
}

type EventHandler interface {
	HandleEvent(event Event)
}
