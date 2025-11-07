package postgres

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/merlinfuchs/stateway/stateway-cache/db/postgres/pgmodel"
	"github.com/merlinfuchs/stateway/stateway-cache/store"
)

func (c *Client) UpsertGuild(ctx context.Context, params store.UpsertGuildParams) error {
	err := c.Q.UpsertGuild(ctx, pgmodel.UpsertGuildParams{
		ID:    int64(params.ID),
		AppID: int64(params.AppID),
		Data:  params.Data,
		CreatedAt: pgtype.Timestamp{
			Time:  params.CreatedAt,
			Valid: true,
		},
		UpdatedAt: pgtype.Timestamp{
			Time:  params.UpdatedAt,
			Valid: true,
		},
	})
	return err
}
