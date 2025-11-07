package event

import (
	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/snowflake/v2"
	"github.com/merlinfuchs/stateway/stateway-lib/service"
)

type Event interface {
	EventID() snowflake.ID
	ServiceType() service.ServiceType
	EventType() string
}

type GatewayEvent struct {
	ID       snowflake.ID  `json:"id"`
	AppID    snowflake.ID  `json:"app_id"`
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

func (e *GatewayEvent) ServiceType() service.ServiceType {
	return service.ServiceTypeGateway
}

func (e *GatewayEvent) EventType() string {
	return e.Type
}

type EventHandler interface {
	HandleEvent(event Event)
}
