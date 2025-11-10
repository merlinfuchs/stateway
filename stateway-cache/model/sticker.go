package model

import (
	"encoding/json"
	"time"

	"github.com/disgoorg/snowflake/v2"
)

type Sticker struct {
	AppID     snowflake.ID
	GuildID   snowflake.ID
	StickerID snowflake.ID
	Data      json.RawMessage
	CreatedAt time.Time
	UpdatedAt time.Time
}
