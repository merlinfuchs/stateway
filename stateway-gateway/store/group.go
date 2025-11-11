package store

import (
	"context"
	"time"

	"github.com/merlinfuchs/stateway/stateway-gateway/model"
	"github.com/merlinfuchs/stateway/stateway-lib/gateway"
)

type CreateGroupParams struct {
	ID                 string
	DisplayName        string
	DefaultConfig      gateway.AppConfig
	DefaultConstraints gateway.AppConstraints
	CreatedAt          time.Time
	UpdatedAt          time.Time
}

type UpsertGroupParams = CreateGroupParams

type UpdateGroupParams struct {
	ID                 string
	DisplayName        string
	DefaultConfig      gateway.AppConfig
	DefaultConstraints gateway.AppConstraints
	UpdatedAt          time.Time
}

type GroupStore interface {
	GetGroup(ctx context.Context, id string) (*model.Group, error)
	GetGroups(ctx context.Context) ([]*model.Group, error)
	CreateGroup(ctx context.Context, params CreateGroupParams) (*model.Group, error)
	UpdateGroup(ctx context.Context, params UpdateGroupParams) (*model.Group, error)
	UpsertGroup(ctx context.Context, params UpsertGroupParams) (*model.Group, error)
	DeleteGroup(ctx context.Context, id string) error
}
