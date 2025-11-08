package event

import (
	"encoding/json"
	"fmt"

	"github.com/disgoorg/snowflake/v2"
	"github.com/merlinfuchs/stateway/stateway-lib/service"
)

type Event interface {
	EventID() snowflake.ID
	ServiceType() service.ServiceType
	EventType() string
}

type GatewayEvent struct {
	ID        snowflake.ID    `json:"id"`
	GatewayID int             `json:"gateway_index"`
	GroupID   string          `json:"group_id"`
	AppID     snowflake.ID    `json:"app_id"`
	ShardID   int             `json:"shard_id"`
	GuildID   *snowflake.ID   `json:"guild_id"`
	Type      string          `json:"type"`
	Data      json.RawMessage `json:"data"`
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

type unmarshalEvent struct {
	EventID     snowflake.ID        `json:"event_id"`
	ServiceType service.ServiceType `json:"service_type"`
	EventType   string              `json:"event_type"`
	Data        json.RawMessage     `json:"data"`
}

func UnmarshalEvent(data []byte) (Event, error) {
	var event unmarshalEvent
	err := json.Unmarshal(data, &event)
	if err != nil {
		return nil, err
	}

	switch event.ServiceType {
	case service.ServiceTypeGateway:
		var gatewayEvent GatewayEvent
		err := json.Unmarshal(event.Data, &gatewayEvent)
		if err != nil {
			return nil, err
		}
		return &gatewayEvent, nil
	}

	return nil, fmt.Errorf("unknown service type: %s", event.ServiceType)
}

type marshalEvent struct {
	EventID     snowflake.ID        `json:"event_id"`
	ServiceType service.ServiceType `json:"service_type"`
	EventType   string              `json:"event_type"`
	Data        any                 `json:"data"`
}

func MarshalEvent(event Event) ([]byte, error) {
	return json.Marshal(marshalEvent{
		EventID:     event.EventID(),
		ServiceType: event.ServiceType(),
		EventType:   event.EventType(),
		Data:        event,
	})
}
