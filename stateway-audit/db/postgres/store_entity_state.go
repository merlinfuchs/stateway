package postgres

import (
	"context"
	"errors"

	"github.com/disgoorg/snowflake/v2"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/merlinfuchs/stateway/stateway-audit/db/postgres/pgmodel"
	"github.com/merlinfuchs/stateway/stateway-audit/model"
	"github.com/merlinfuchs/stateway/stateway-audit/store"
	"github.com/merlinfuchs/stateway/stateway-lib/audit"
)

var _ store.EntityStateStore = (*Client)(nil)

func (c *Client) GetEntityState(ctx context.Context, appID snowflake.ID, guildID snowflake.ID, entityType audit.EntityType, entityID snowflake.ID) (*model.EntityState, error) {
	row, err := c.Q.GetEntityState(ctx, pgmodel.GetEntityStateParams{
		AppID:      int64(appID),
		GuildID:    int64(guildID),
		EntityType: string(entityType),
		EntityID:   int64(entityID),
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, store.ErrNotFound
		}
		return nil, err
	}

	return rowToEntityState(row)
}

func (c *Client) UpsertEntityState(ctx context.Context, entityState model.EntityState) error {
	err := c.Q.UpsertEntityState(ctx, pgmodel.UpsertEntityStateParams{
		AppID:      int64(entityState.AppID),
		GuildID:    int64(entityState.GuildID),
		EntityType: string(entityState.EntityType),
		EntityID:   int64(entityState.EntityID),
		Data:       entityState.Data,
		Deleted:    entityState.Deleted,
		CreatedAt: pgtype.Timestamp{
			Time:  entityState.CreatedAt,
			Valid: true,
		},
		UpdatedAt: pgtype.Timestamp{
			Time:  entityState.UpdatedAt,
			Valid: true,
		},
	})
	return err
}

func rowToEntityState(row pgmodel.AuditEntityState) (*model.EntityState, error) {
	return &model.EntityState{
		AppID:      snowflake.ID(row.AppID),
		GuildID:    snowflake.ID(row.GuildID),
		EntityType: audit.EntityType(row.EntityType),
		EntityID:   snowflake.ID(row.EntityID),
		Data:       row.Data,
		Deleted:    row.Deleted,
		CreatedAt:  row.CreatedAt.Time,
		UpdatedAt:  row.UpdatedAt.Time,
	}, nil
}
