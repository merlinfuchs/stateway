package postgres

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/disgoorg/snowflake/v2"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/merlinfuchs/stateway/stateway-gateway/db/postgres/pgmodel"
	"github.com/merlinfuchs/stateway/stateway-gateway/model"
	"github.com/merlinfuchs/stateway/stateway-gateway/store"
	"github.com/merlinfuchs/stateway/stateway-lib/gateway"
	"gopkg.in/guregu/null.v4"
)

func (c *Client) GetApp(ctx context.Context, id snowflake.ID) (*model.App, error) {
	row, err := c.Q.GetApp(ctx, int64(id))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, store.ErrNotFound
		}
		return nil, err
	}
	return rowToApp(row)
}

func (c *Client) GetApps(ctx context.Context) ([]*model.App, error) {
	rows, err := c.Q.GetApps(ctx)
	if err != nil {
		return nil, err
	}
	apps := make([]*model.App, 0, len(rows))
	for _, row := range rows {
		app, err := rowToApp(row)
		if err != nil {
			return nil, err
		}
		apps = append(apps, app)
	}
	return apps, nil
}

func (c *Client) GetEnabledApps(ctx context.Context, params store.GetEnabledAppsParams) ([]*model.App, error) {
	if params.GatewayCount == 0 {
		params.GatewayCount = 1
	}

	rows, err := c.Q.GetEnabledApps(ctx, pgmodel.GetEnabledAppsParams{
		GatewayCount: int64(params.GatewayCount),
		GatewayID:    int64(params.GatewayID),
	})
	if err != nil {
		return nil, err
	}

	apps := make([]*model.App, 0, len(rows))
	for _, row := range rows {
		app, err := rowToApp(row)
		if err != nil {
			return nil, err
		}
		apps = append(apps, app)
	}
	return apps, nil
}

func (c *Client) CreateApp(ctx context.Context, params store.CreateAppParams) (*model.App, error) {
	rawConstraints, err := json.Marshal(params.Constraints)
	if err != nil {
		return nil, err
	}
	rawConfig, err := json.Marshal(params.Config)
	if err != nil {
		return nil, err
	}

	row, err := c.Q.CreateApp(ctx, pgmodel.CreateAppParams{
		ID:               int64(params.ID),
		GroupID:          params.GroupID,
		DisplayName:      params.DisplayName,
		DiscordClientID:  int64(params.DiscordClientID),
		DiscordBotToken:  params.DiscordBotToken,
		DiscordPublicKey: params.DiscordPublicKey,
		DiscordClientSecret: pgtype.Text{
			String: params.DiscordClientSecret.String,
			Valid:  params.DiscordClientSecret.Valid,
		},
		ShardCount:  int32(params.ShardCount),
		Constraints: rawConstraints,
		Config:      rawConfig,
		CreatedAt: pgtype.Timestamp{
			Time:  params.CreatedAt,
			Valid: true,
		},
		UpdatedAt: pgtype.Timestamp{
			Time:  params.UpdatedAt,
			Valid: true,
		},
	})
	if err != nil {
		return nil, err
	}
	return rowToApp(row)
}

func (c *Client) UpdateApp(ctx context.Context, params store.UpdateAppParams) (*model.App, error) {
	rawConstraints, err := json.Marshal(params.Constraints)
	if err != nil {
		return nil, err
	}
	rawConfig, err := json.Marshal(params.Config)
	if err != nil {
		return nil, err
	}

	row, err := c.Q.UpdateApp(ctx, pgmodel.UpdateAppParams{
		ID:               int64(params.ID),
		GroupID:          params.GroupID,
		DisplayName:      params.DisplayName,
		DiscordClientID:  int64(params.DiscordClientID),
		DiscordBotToken:  params.DiscordBotToken,
		DiscordPublicKey: params.DiscordPublicKey,
		DiscordClientSecret: pgtype.Text{
			String: params.DiscordClientSecret.String,
			Valid:  params.DiscordClientSecret.Valid,
		},
		ShardCount:  int32(params.ShardCount),
		Constraints: rawConstraints,
		Config:      rawConfig,
		Disabled:    params.Disabled,
		DisabledCode: pgtype.Text{
			String: string(params.DisabledCode.String),
			Valid:  params.DisabledCode.Valid,
		},
		DisabledMessage: pgtype.Text{
			String: params.DisabledMessage.String,
			Valid:  params.DisabledMessage.Valid,
		},
		UpdatedAt: pgtype.Timestamp{
			Time:  params.UpdatedAt,
			Valid: true,
		},
	})
	if err != nil {
		return nil, err
	}
	return rowToApp(row)
}

func (c *Client) UpsertApp(ctx context.Context, params store.UpsertAppParams) error {
	rawConstraints, err := json.Marshal(params.Constraints)
	if err != nil {
		return err
	}
	rawConfig, err := json.Marshal(params.Config)
	if err != nil {
		return err
	}

	err = c.Q.UpsertApp(ctx, pgmodel.UpsertAppParams{
		ID:               int64(params.ID),
		GroupID:          params.GroupID,
		DisplayName:      params.DisplayName,
		DiscordClientID:  int64(params.DiscordClientID),
		DiscordBotToken:  params.DiscordBotToken,
		DiscordPublicKey: params.DiscordPublicKey,
		DiscordClientSecret: pgtype.Text{
			String: params.DiscordClientSecret.String,
			Valid:  params.DiscordClientSecret.Valid,
		},
		ShardCount:  int32(params.ShardCount),
		Constraints: rawConstraints,
		Config:      rawConfig,
		CreatedAt: pgtype.Timestamp{
			Time:  params.CreatedAt,
			Valid: true,
		},
		UpdatedAt: pgtype.Timestamp{
			Time:  params.UpdatedAt,
			Valid: true,
		},
	})
	if err != nil {
		return err
	}
	return nil
}

func (c *Client) DisableApp(ctx context.Context, params store.DisableAppParams) (*model.App, error) {
	row, err := c.Q.DisableApp(ctx, pgmodel.DisableAppParams{
		ID: int64(params.ID),
		DisabledCode: pgtype.Text{
			String: string(params.DisabledCode),
			Valid:  params.DisabledCode != "",
		},
		DisabledMessage: pgtype.Text{
			String: params.DisabledMessage.String,
			Valid:  params.DisabledMessage.Valid,
		},
		UpdatedAt: pgtype.Timestamp{
			Time:  params.UpdatedAt,
			Valid: true,
		},
	})
	if err != nil {
		return nil, err
	}
	return rowToApp(row)
}

func (c *Client) DeleteApp(ctx context.Context, id snowflake.ID) error {
	err := c.Q.DeleteApp(ctx, int64(id))
	if err != nil {
		return err
	}
	return nil
}

func rowToApp(row pgmodel.GatewayApp) (*model.App, error) {
	var constraints gateway.AppConstraints
	var config gateway.AppConfig
	if row.Constraints != nil {
		err := json.Unmarshal(row.Constraints, &constraints)
		if err != nil {
			return nil, err
		}
	}
	if row.Config != nil {
		err := json.Unmarshal(row.Config, &config)
		if err != nil {
			return nil, err
		}
	}

	return &model.App{
		ID:                  snowflake.ID(row.ID),
		GroupID:             row.GroupID,
		DisplayName:         row.DisplayName,
		DiscordClientID:     snowflake.ID(row.DiscordClientID),
		DiscordBotToken:     row.DiscordBotToken,
		DiscordPublicKey:    row.DiscordPublicKey,
		DiscordClientSecret: null.NewString(row.DiscordClientSecret.String, row.DiscordClientSecret.Valid),
		ShardCount:          int(row.ShardCount),
		Constraints:         constraints,
		Config:              config,
		Disabled:            row.Disabled,
		DisabledCode:        gateway.AppDisabledCode(row.DisabledCode.String),
		DisabledMessage:     null.NewString(row.DisabledMessage.String, row.DisabledMessage.Valid),
		CreatedAt:           row.CreatedAt.Time,
		UpdatedAt:           row.UpdatedAt.Time,
	}, nil
}
