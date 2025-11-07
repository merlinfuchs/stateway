package store

import (
	"context"
	"time"

	"github.com/disgoorg/snowflake/v2"
	"github.com/merlinfuchs/stateway/stateway-gateway/model"
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
	CreatedAt           time.Time
	UpdatedAt           time.Time
}

type UpdateAppParams struct {
	GroupID             string
	DisplayName         string
	DiscordClientID     snowflake.ID
	DiscordBotToken     string
	DiscordPublicKey    string
	DiscordClientSecret null.String
	Disabled            bool
	DisabledCode        null.String
	DisabledMessage     null.String
	UpdatedAt           time.Time
}

type DisableAppParams struct {
	ID              snowflake.ID
	DisabledCode    model.AppDisabledCode
	DisabledMessage null.String
	UpdatedAt       time.Time
}

type AppStore interface {
	GetApp(ctx context.Context, id snowflake.ID) (*model.App, error)
	GetApps(ctx context.Context) ([]*model.App, error)
	GetEnabledApps(ctx context.Context) ([]*model.App, error)
	CreateApp(ctx context.Context, params CreateAppParams) (*model.App, error)
	UpdateApp(ctx context.Context, params UpdateAppParams) (*model.App, error)
	DisableApp(ctx context.Context, params DisableAppParams) (*model.App, error)
	DeleteApp(ctx context.Context, id snowflake.ID) error
}
