package postgres

import (
	"context"
	"fmt"

	"github.com/disgoorg/snowflake/v2"
	"github.com/merlinfuchs/stateway/stateway-gateway/db/postgres/pgmodel"
	"github.com/merlinfuchs/stateway/stateway-gateway/store"
)

func (c *Client) TryLockBucket(ctx context.Context, appID snowflake.ID, bucketKey int) (store.IdentifyRateLimitLock, error) {
	var res store.IdentifyRateLimitLock

	tx, err := c.DB.Begin(ctx)
	if err != nil {
		return res, fmt.Errorf("failed to begin transaction: %w", err)
	}

	txQ := c.Q.WithTx(tx)
	locked, err := txQ.TryLockBucket(ctx, pgmodel.TryLockBucketParams{
		PgTryAdvisoryXactLock:   int32(appID << 32),
		PgTryAdvisoryXactLock_2: int32(bucketKey),
	})
	if err != nil {
		return res, fmt.Errorf("failed to try lock bucket: %w", err)
	}

	if locked {
		res.Locked = true
		res.Unlock = func(ctx context.Context) error {
			err := tx.Rollback(ctx)
			if err != nil {
				return fmt.Errorf("failed to rollback transaction: %w", err)
			}
			return nil
		}
		return res, nil
	}

	err = tx.Rollback(ctx)
	if err != nil {
		return res, fmt.Errorf("failed to rollback transaction: %w", err)
	}

	return res, nil
}
