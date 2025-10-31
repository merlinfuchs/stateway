package postgres

import (
	"context"
	"errors"

	"github.com/disgoorg/snowflake/v2"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/merlinfuchs/stateway/stateway-gateway/db/postgres/pgmodel"
	"github.com/merlinfuchs/stateway/stateway-gateway/model"
	"github.com/merlinfuchs/stateway/stateway-gateway/store"
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
	return rowToApp(row), nil
}

func (c *Client) GetApps(ctx context.Context) ([]*model.App, error) {
	rows, err := c.Q.GetApps(ctx)
	if err != nil {
		return nil, err
	}
	apps := make([]*model.App, 0, len(rows))
	for _, row := range rows {
		apps = append(apps, rowToApp(row))
	}
	return apps, nil
}

func (c *Client) GetEnabledApps(ctx context.Context) ([]*model.App, error) {
	rows, err := c.Q.GetEnabledApps(ctx)
	if err != nil {
		return nil, err
	}

	apps := make([]*model.App, 0, len(rows))
	for _, row := range rows {
		apps = append(apps, rowToApp(row))
	}
	return apps, nil
}

func (c *Client) CreateApp(ctx context.Context, params store.CreateAppParams) (*model.App, error) {
	row, err := c.Q.CreateApp(ctx, pgmodel.CreateAppParams{
		ID:               int64(params.ID),
		DisplayName:      params.DisplayName,
		DiscordClientID:  int64(params.DiscordClientID),
		DiscordBotToken:  params.DiscordBotToken,
		DiscordPublicKey: params.DiscordPublicKey,
		DiscordClientSecret: pgtype.Text{
			String: params.DiscordClientSecret.String,
			Valid:  params.DiscordClientSecret.Valid,
		},
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
	return rowToApp(row), nil
}

func (c *Client) UpdateApp(ctx context.Context, params store.UpdateAppParams) (*model.App, error) {
	row, err := c.Q.UpdateApp(ctx, pgmodel.UpdateAppParams{
		ID:               int64(params.ID),
		DisplayName:      params.DisplayName,
		DiscordClientID:  int64(params.DiscordClientID),
		DiscordBotToken:  params.DiscordBotToken,
		DiscordPublicKey: params.DiscordPublicKey,
		DiscordClientSecret: pgtype.Text{
			String: params.DiscordClientSecret.String,
			Valid:  params.DiscordClientSecret.Valid,
		},
		Disabled: params.Disabled,
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
	return rowToApp(row), nil
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
	return rowToApp(row), nil
}

func (c *Client) DeleteApp(ctx context.Context, id snowflake.ID) error {
	err := c.Q.DeleteApp(ctx, int64(id))
	if err != nil {
		return err
	}
	return nil
}

func rowToApp(row pgmodel.GatewayApp) *model.App {
	return &model.App{
		ID:                  snowflake.ID(row.ID),
		DisplayName:         row.DisplayName,
		DiscordClientID:     snowflake.ID(row.DiscordClientID),
		DiscordBotToken:     row.DiscordBotToken,
		DiscordPublicKey:    row.DiscordPublicKey,
		DiscordClientSecret: null.NewString(row.DiscordClientSecret.String, row.DiscordClientSecret.Valid),
		Disabled:            row.Disabled,
		DisabledCode:        model.AppDisabledCode(row.DisabledCode.String),
		DisabledMessage:     null.NewString(row.DisabledMessage.String, row.DisabledMessage.Valid),
		CreatedAt:           row.CreatedAt.Time,
		UpdatedAt:           row.UpdatedAt.Time,
	}
}
