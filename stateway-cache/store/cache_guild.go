package store

import (
	"context"
	"encoding/json"
	"time"

	"github.com/disgoorg/snowflake/v2"
)

type UpsertGuildParams struct {
	ID        snowflake.ID
	AppID     snowflake.ID
	Data      json.RawMessage
	CreatedAt time.Time
	UpdatedAt time.Time
}

type CacheGuildStore interface {
	UpsertGuild(ctx context.Context, params UpsertGuildParams) error
}
