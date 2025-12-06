package store

import (
	"context"

	"github.com/merlinfuchs/stateway/stateway-audit/model"
)

type EntityChangeStore interface {
	InsertEntityChanges(ctx context.Context, entityChanges ...model.EntityChange) error
}
