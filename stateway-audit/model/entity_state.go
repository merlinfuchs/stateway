package model

import (
	"encoding/json"
	"time"

	"github.com/disgoorg/snowflake/v2"
)

type EntityType string

const (
	EntityTypeGuild   EntityType = "guild"
	EntityTypeChannel EntityType = "channel"
	EntityTypeRole    EntityType = "role"
)

type EntityState struct {
	AppID      snowflake.ID
	GuildID    snowflake.ID
	EntityType EntityType
	EntityID   snowflake.ID
	Data       json.RawMessage
	Deleted    bool
	CreatedAt  time.Time
	UpdatedAt  time.Time
}
