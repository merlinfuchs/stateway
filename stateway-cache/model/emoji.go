package model

import (
	"encoding/json"
	"time"

	"github.com/disgoorg/snowflake/v2"
)

type Emoji struct {
	AppID     snowflake.ID
	GuildID   snowflake.ID
	EmojiID   snowflake.ID
	Data      json.RawMessage
	CreatedAt time.Time
	UpdatedAt time.Time
}
