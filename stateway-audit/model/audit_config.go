package model

import (
	"time"

	"github.com/disgoorg/snowflake/v2"
)

type AuditConfig struct {
	AppID     snowflake.ID
	GuildID   snowflake.ID
	Enabled   bool
	CreatedAt time.Time
	UpdatedAt time.Time
}
