package model

import (
	"encoding/json"
	"time"

	"github.com/disgoorg/snowflake/v2"
)

type Guild struct {
	AppID       snowflake.ID
	GuildID     snowflake.ID
	Data        json.RawMessage
	Unavailable bool
	Tainted     bool
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
