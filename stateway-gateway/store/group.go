package store

import (
	"context"

	"github.com/merlinfuchs/stateway/stateway-gateway/model"
)

type GroupStore interface {
	GetGroup(ctx context.Context, id string) (*model.Group, error)
	GetGroups(ctx context.Context) ([]*model.Group, error)
	CreateGroup(ctx context.Context, group *model.Group) error
	UpdateGroup(ctx context.Context, group *model.Group) error
	UpsertGroup(ctx context.Context, group *model.Group) error
	DeleteGroup(ctx context.Context, id string) error
}
