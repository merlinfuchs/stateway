package postgres

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/merlinfuchs/stateway/stateway-gateway/db/postgres/pgmodel"
	"github.com/merlinfuchs/stateway/stateway-gateway/model"
	"github.com/merlinfuchs/stateway/stateway-gateway/store"
	"github.com/merlinfuchs/stateway/stateway-lib/gateway"
)

func (c *Client) GetGroup(ctx context.Context, id string) (*model.Group, error) {
	row, err := c.Q.GetGroup(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, store.ErrNotFound
		}
		return nil, err
	}
	return rowToGroup(row)
}

func (c *Client) GetGroups(ctx context.Context) ([]*model.Group, error) {
	rows, err := c.Q.GetGroups(ctx)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, store.ErrNotFound
		}
		return nil, err
	}

	var groups []*model.Group
	for _, row := range rows {
		group, err := rowToGroup(row)
		if err != nil {
			return nil, err
		}
		groups = append(groups, group)
	}
	return groups, nil
}

func (c *Client) CreateGroup(ctx context.Context, params store.CreateGroupParams) (*model.Group, error) {
	rawConfig, err := json.Marshal(params.DefaultConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal default config: %w", err)
	}
	rawConstraints, err := json.Marshal(params.DefaultConstraints)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal default constraints: %w", err)
	}

	row, err := c.Q.CreateGroup(ctx, pgmodel.CreateGroupParams{
		ID:                 params.ID,
		DisplayName:        params.DisplayName,
		DefaultConfig:      rawConfig,
		DefaultConstraints: rawConstraints,
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
	return rowToGroup(row)
}

func (c *Client) UpdateGroup(ctx context.Context, params store.UpdateGroupParams) (*model.Group, error) {
	rawConfig, err := json.Marshal(params.DefaultConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal default config: %w", err)
	}
	rawConstraints, err := json.Marshal(params.DefaultConstraints)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal default constraints: %w", err)
	}

	row, err := c.Q.UpdateGroup(ctx, pgmodel.UpdateGroupParams{
		ID:                 params.ID,
		DisplayName:        params.DisplayName,
		DefaultConfig:      rawConfig,
		DefaultConstraints: rawConstraints,
		UpdatedAt: pgtype.Timestamp{
			Time:  params.UpdatedAt,
			Valid: true,
		},
	})
	if err != nil {
		return nil, err
	}
	return rowToGroup(row)
}

func (c *Client) UpsertGroup(ctx context.Context, params store.UpsertGroupParams) (*model.Group, error) {
	rawConfig, err := json.Marshal(params.DefaultConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal default config: %w", err)
	}
	rawConstraints, err := json.Marshal(params.DefaultConstraints)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal default constraints: %w", err)
	}

	row, err := c.Q.UpsertGroup(ctx, pgmodel.UpsertGroupParams{
		ID:                 params.ID,
		DisplayName:        params.DisplayName,
		DefaultConfig:      rawConfig,
		DefaultConstraints: rawConstraints,
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

	return rowToGroup(row)
}

func (c *Client) DeleteGroup(ctx context.Context, id string) error {
	err := c.Q.DeleteGroup(ctx, id)
	if err != nil {
		return err
	}
	return nil
}

func rowToGroup(row pgmodel.GatewayGroup) (*model.Group, error) {
	var defaultConfig gateway.AppConfig
	var defaultConstraints gateway.AppConstraints
	if row.DefaultConfig != nil {
		err := json.Unmarshal(row.DefaultConfig, &defaultConfig)
		if err != nil {
			return nil, err
		}
	}
	if row.DefaultConstraints != nil {
		err := json.Unmarshal(row.DefaultConstraints, &defaultConstraints)
		if err != nil {
			return nil, err
		}
	}

	return &model.Group{
		ID:                 row.ID,
		DisplayName:        row.DisplayName,
		DefaultConfig:      defaultConfig,
		DefaultConstraints: defaultConstraints,
		CreatedAt:          row.CreatedAt.Time,
		UpdatedAt:          row.UpdatedAt.Time,
	}, nil
}
