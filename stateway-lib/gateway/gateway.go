package gateway

import (
	"context"

	"github.com/disgoorg/snowflake/v2"
)

type Gateway interface {
	GetApp(ctx context.Context, appID snowflake.ID) (*App, error)
	GetApps(ctx context.Context, params ListAppsRequest) ([]*App, error)
	UpsertApp(ctx context.Context, app UpsertAppRequest) (*App, error)
	DisableApp(ctx context.Context, appID snowflake.ID) error
	DeleteApp(ctx context.Context, appID snowflake.ID) error
	GetGroup(ctx context.Context, groupID string) (*Group, error)
	GetGroups(ctx context.Context) ([]*Group, error)
	UpsertGroup(ctx context.Context, group UpsertGroupRequest) (*Group, error)
	DeleteGroup(ctx context.Context, groupID string) error
}
