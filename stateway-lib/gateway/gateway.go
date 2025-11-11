package gateway

import (
	"context"

	"github.com/disgoorg/snowflake/v2"
	"gopkg.in/guregu/null.v4"
)

type Gateway interface {
	GetApp(ctx context.Context, appID snowflake.ID) (*App, error)
	GetApps(ctx context.Context, groupID null.String, limit int, offset int) ([]*App, error)
	UpsertApp(ctx context.Context, app UpsertAppRequest) error
	DisableApp(ctx context.Context, appID snowflake.ID) error
	DeleteApp(ctx context.Context, appID snowflake.ID) error
	GetGroup(ctx context.Context, groupID string) (*Group, error)
	GetGroups(ctx context.Context) ([]*Group, error)
	UpsertGroup(ctx context.Context, group UpsertGroupRequest) error
	DeleteGroup(ctx context.Context, groupID string) error
}
