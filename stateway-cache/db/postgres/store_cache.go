package postgres

import (
	"context"
	"fmt"

	"github.com/merlinfuchs/stateway/stateway-cache/db/postgres/pgmodel"
	"github.com/merlinfuchs/stateway/stateway-cache/store"
)

func (c *Client) MarkShardEntitiesTainted(ctx context.Context, params store.MarkShardEntitiesTaintedParams) error {
	tx, err := c.DB.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	q := c.Q.WithTx(tx)
	err = q.MarkShardGuildsTainted(ctx, pgmodel.MarkShardGuildsTaintedParams{
		GroupID:    params.GroupID,
		ClientID:   int64(params.ClientID),
		ShardCount: int64(params.ShardCount),
		ShardID:    int64(params.ShardID),
	})
	if err != nil {
		return fmt.Errorf("failed to mark shard guilds tainted: %w", err)
	}

	err = q.MarkShardRolesTainted(ctx, pgmodel.MarkShardRolesTaintedParams{
		GroupID:    params.GroupID,
		ClientID:   int64(params.ClientID),
		ShardCount: int64(params.ShardCount),
		ShardID:    int64(params.ShardID),
	})
	if err != nil {
		return fmt.Errorf("failed to mark shard roles tainted: %w", err)
	}

	err = q.MarkShardChannelsTainted(ctx, pgmodel.MarkShardChannelsTaintedParams{
		GroupID:    params.GroupID,
		ClientID:   int64(params.ClientID),
		ShardCount: int64(params.ShardCount),
		ShardID:    int64(params.ShardID),
	})
	if err != nil {
		return fmt.Errorf("failed to mark shard channels tainted: %w", err)
	}

	err = tx.Commit(ctx)
	if err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}
