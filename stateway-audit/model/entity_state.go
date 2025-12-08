package model

import (
	"encoding/json"
	"time"

	"github.com/disgoorg/snowflake/v2"
	"github.com/merlinfuchs/stateway/stateway-lib/audit"
)

type EntityState struct {
	AppID      snowflake.ID
	GuildID    snowflake.ID
	EntityType audit.EntityType
	EntityID   snowflake.ID
	Data       json.RawMessage
	Deleted    bool
	CreatedAt  time.Time
	UpdatedAt  time.Time
}
