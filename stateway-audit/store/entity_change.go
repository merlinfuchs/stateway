package store

import (
	"context"
	"time"

	"github.com/disgoorg/snowflake/v2"
	"github.com/merlinfuchs/stateway/stateway-audit/model"
	"github.com/merlinfuchs/stateway/stateway-lib/audit"
)

type GetEntityChangesParams struct {
	AppID      snowflake.ID
	GuildID    snowflake.ID
	EntityType *audit.EntityType
	EntityID   *snowflake.ID
	Path       *string
	Before     *time.Time
	After      *time.Time
}

type EntityChangeStore interface {
	InsertEntityChanges(ctx context.Context, entityChanges ...model.EntityChange) error
	GetEntityChanges(ctx context.Context, params GetEntityChangesParams) ([]*model.EntityChange, error)
}
