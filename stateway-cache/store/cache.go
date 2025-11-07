package store

import (
	"context"

	"github.com/disgoorg/snowflake/v2"
)

type MarkShardEntitiesTaintedParams struct {
	GroupID    string
	ClientID   snowflake.ID
	ShardCount int
	ShardID    int
}

type CacheStore interface {
	CacheGuildStore
	CacheRoleStore
	CacheChannelStore

	MarkShardEntitiesTainted(ctx context.Context, params MarkShardEntitiesTaintedParams) error
}
