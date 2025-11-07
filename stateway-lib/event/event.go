package event

import (
	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/snowflake/v2"
)

type EventType string

const (
	EventTypeGateway EventType = "gateway"
)

type Event interface {
	EventID() snowflake.ID
	EventType() EventType
}

type GatewayEvent struct {
	ID       snowflake.ID  `json:"id"`
	GroupID  string        `json:"group_id"`
	ClientID snowflake.ID  `json:"client_id"`
	ShardID  int           `json:"shard_id"`
	GuildID  *snowflake.ID `json:"guild_id"`
	Type     string        `json:"type"`
	Data     bot.Event     `json:"data"`
}

func (e *GatewayEvent) EventID() snowflake.ID {
	return e.ID
}

func (e *GatewayEvent) EventType() EventType {
	return EventTypeGateway
}

type EventHandler interface {
	HandleEvent(event Event)
}
