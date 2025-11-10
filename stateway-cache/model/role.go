package model

import (
	"encoding/json"
	"time"

	"github.com/disgoorg/snowflake/v2"
)

type Role struct {
	AppID     snowflake.ID
	GuildID   snowflake.ID
	RoleID    snowflake.ID
	Data      json.RawMessage
	CreatedAt time.Time
	UpdatedAt time.Time
}
