package batcher

import (
	"context"

	"github.com/merlinfuchs/stateway/stateway-audit/model"
)

type Batcher interface {
	Push(ctx context.Context, entityChange model.EntityChange) error
	Start(ctx context.Context) error
}
