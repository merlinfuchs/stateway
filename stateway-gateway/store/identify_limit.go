package store

import (
	"context"

	"github.com/disgoorg/snowflake/v2"
)

type IdentifyRateLimitStore interface {
	TryLockBucket(ctx context.Context, appID snowflake.ID, bucketKey int) (IdentifyRateLimitLock, error)
}

type IdentifyRateLimitLock struct {
	Locked bool
	Unlock func(ctx context.Context) error
}
