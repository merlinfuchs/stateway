package store

import (
	"context"
	"time"

	"github.com/disgoorg/snowflake/v2"
	"github.com/merlinfuchs/stateway/stateway-gateway/model"
	"github.com/merlinfuchs/stateway/stateway-lib/gateway"
	"gopkg.in/guregu/null.v4"
)

type CreateAppParams struct {
	ID                  snowflake.ID
	GroupID             string
	DisplayName         string
	DiscordClientID     snowflake.ID
	DiscordBotToken     string
	DiscordPublicKey    string
	DiscordClientSecret null.String
	ShardCount          int
	Constraints         gateway.AppConstraints
	Config              gateway.AppConfig
	CreatedAt           time.Time
	UpdatedAt           time.Time
}

type UpdateAppParams struct {
	ID                  snowflake.ID
	GroupID             string
	DisplayName         string
	DiscordClientID     snowflake.ID
	DiscordBotToken     string
	DiscordPublicKey    string
	DiscordClientSecret null.String
	ShardCount          int
	Constraints         gateway.AppConstraints
	Config              gateway.AppConfig
	Disabled            bool
	DisabledCode        null.String
	DisabledMessage     null.String
	UpdatedAt           time.Time
}

type UpsertAppParams = CreateAppParams

type DisableAppParams struct {
	ID              snowflake.ID
	DisabledCode    gateway.AppDisabledCode
	DisabledMessage null.String
	UpdatedAt       time.Time
}

type GetEnabledAppsParams struct {
	GatewayCount int
	GatewayID    int
}

type GetAppsParams struct {
	GroupID null.String
	Limit   null.Int
	Offset  null.Int
}

type AppStore interface {
	GetApp(ctx context.Context, id snowflake.ID) (*model.App, error)
	GetApps(ctx context.Context, params GetAppsParams) ([]*model.App, error)
	GetEnabledApps(ctx context.Context, params GetEnabledAppsParams) ([]*model.App, error)
	CreateApp(ctx context.Context, params CreateAppParams) (*model.App, error)
	UpdateApp(ctx context.Context, params UpdateAppParams) (*model.App, error)
	UpsertApp(ctx context.Context, params UpsertAppParams) (*model.App, error)
	DisableApp(ctx context.Context, params DisableAppParams) error
	DeleteApp(ctx context.Context, id snowflake.ID) error
}
