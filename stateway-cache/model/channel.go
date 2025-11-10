package model

import (
	"encoding/json"
	"time"

	"github.com/disgoorg/snowflake/v2"
)

type Channel struct {
	AppID     snowflake.ID
	GuildID   snowflake.ID
	ChannelID snowflake.ID
	Data      json.RawMessage
	CreatedAt time.Time
	UpdatedAt time.Time
}
