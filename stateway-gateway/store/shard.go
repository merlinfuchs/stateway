package store

import (
	"context"
	"time"

	"github.com/disgoorg/snowflake/v2"
	"github.com/merlinfuchs/stateway/stateway-gateway/model"
)

type UpsertShardSessionParams struct {
	ID           string
	AppID        snowflake.ID
	ShardID      int
	ShardCount   int
	LastSequence int
	ResumeURL    string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type ShardSessionStore interface {
	UpsertShardSession(ctx context.Context, params UpsertShardSessionParams) error
	GetLastShardSession(ctx context.Context, appID snowflake.ID, shardID int, shardCount int) (*model.ShardSession, error)
	InvalidateShardSession(ctx context.Context, appID snowflake.ID, shardID int, shardCount int) error
	PurgeSessions(ctx context.Context) error
}
