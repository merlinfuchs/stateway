package postgres

import (
	"context"
	"errors"
	"time"

	"github.com/disgoorg/snowflake/v2"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/merlinfuchs/stateway/stateway-gateway/db/postgres/pgmodel"
	"github.com/merlinfuchs/stateway/stateway-gateway/model"
	"github.com/merlinfuchs/stateway/stateway-gateway/store"
	"gopkg.in/guregu/null.v4"
)

func (c *Client) UpsertShardSession(ctx context.Context, params store.UpsertShardSessionParams) error {
	return c.Q.UpsertShardSession(ctx, pgmodel.UpsertShardSessionParams{
		ID:           params.ID,
		AppID:        int64(params.AppID),
		ShardID:      int32(params.ShardID),
		ShardCount:   int32(params.ShardCount),
		LastSequence: int32(params.LastSequence),
		ResumeUrl:    params.ResumeURL,
		CreatedAt:    pgtype.Timestamp{Time: params.CreatedAt, Valid: true},
		UpdatedAt:    pgtype.Timestamp{Time: params.UpdatedAt, Valid: true},
	})
}

func (c *Client) GetLastShardSession(ctx context.Context, appID snowflake.ID, shardID int, shardCount int) (*model.ShardSession, error) {
	row, err := c.Q.GetLastShardSession(ctx, pgmodel.GetLastShardSessionParams{
		AppID:      int64(appID),
		ShardID:    int32(shardID),
		ShardCount: int32(shardCount),
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, store.ErrNotFound
		}
		return nil, err
	}
	return rowToShardSession(row), nil
}

func (c *Client) InvalidateShardSession(ctx context.Context, appID snowflake.ID, shardID int, shardCount int) error {
	return c.Q.InvalidateShardSession(ctx, pgmodel.InvalidateShardSessionParams{
		AppID:         int64(appID),
		ShardID:       int32(shardID),
		InvalidatedAt: pgtype.Timestamp{Time: time.Now().UTC(), Valid: true},
	})
}

func rowToShardSession(row pgmodel.GatewayShardSession) *model.ShardSession {
	return &model.ShardSession{
		ID:            row.ID,
		AppID:         snowflake.ID(row.AppID),
		ShardID:       int(row.ShardID),
		ShardCount:    int(row.ShardCount),
		LastSequence:  int(row.LastSequence),
		ResumeURL:     row.ResumeUrl,
		CreatedAt:     row.CreatedAt.Time,
		UpdatedAt:     row.UpdatedAt.Time,
		InvalidatedAt: null.NewTime(row.InvalidatedAt.Time, row.InvalidatedAt.Valid),
	}
}

func (c *Client) PurgeSessions(ctx context.Context) error {
	return c.Q.PurgeSessions(ctx)
}
