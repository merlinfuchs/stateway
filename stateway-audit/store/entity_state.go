package store

import (
	"context"

	"github.com/disgoorg/snowflake/v2"
	"github.com/merlinfuchs/stateway/stateway-audit/model"
	"github.com/merlinfuchs/stateway/stateway-lib/audit"
)

type EntityStateStore interface {
	GetEntityState(ctx context.Context, appID snowflake.ID, guildID snowflake.ID, entityType audit.EntityType, entityID snowflake.ID) (*model.EntityState, error)
	UpsertEntityState(ctx context.Context, entityState model.EntityState) error
}
