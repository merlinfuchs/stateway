package app

import (
	"context"
	"fmt"
	"log/slog"
	"sync"
	"time"

	"github.com/disgoorg/snowflake/v2"
	"github.com/merlinfuchs/stateway/stateway-gateway/store"
)

type IdentifyRateLimiter struct {
	store          store.IdentifyRateLimitStore
	appID          snowflake.ID
	maxConcurrency int

	mu          sync.Mutex
	unlockFuncs map[int]func(ctx context.Context) error
}

func NewIdentifyRateLimiter(
	store store.IdentifyRateLimitStore,
	appID snowflake.ID,
	maxConcurrency int,
) *IdentifyRateLimiter {
	if maxConcurrency < 1 {
		maxConcurrency = 1
	}
	return &IdentifyRateLimiter{
		store:          store,
		appID:          appID,
		maxConcurrency: maxConcurrency,
		unlockFuncs:    make(map[int]func(ctx context.Context) error),
	}
}

func (l *IdentifyRateLimiter) Close(ctx context.Context) {
	l.mu.Lock()
	defer l.mu.Unlock()

	for _, unlockFunc := range l.unlockFuncs {
		err := unlockFunc(ctx)
		if err != nil {
			slog.Error(
				"Failed to unlock bucket",
				slog.String("app_id", l.appID.String()),
				slog.Any("error", err),
			)
		}
	}

	l.unlockFuncs = make(map[int]func(ctx context.Context) error)
}

func (l *IdentifyRateLimiter) Wait(ctx context.Context, shardID int) error {
	bucket := l.shardIdentifyBucket(shardID)

	for {
		res, err := l.store.TryLockBucket(ctx, l.appID, bucket)
		if err != nil {
			return fmt.Errorf("failed to try lock bucket: %w", err)
		}

		if res.Locked {
			l.mu.Lock()
			l.unlockFuncs[bucket] = res.Unlock
			l.mu.Unlock()
			return nil
		}

		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(10 * time.Millisecond):
		}
	}
}

func (l *IdentifyRateLimiter) Unlock(shardID int) {
	unlockCtx, unlockCtxCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer unlockCtxCancel()

	bucket := l.shardIdentifyBucket(shardID)

	l.mu.Lock()
	unlockFunc, ok := l.unlockFuncs[bucket]
	if ok {
		delete(l.unlockFuncs, bucket)
	}
	l.mu.Unlock()

	if ok {
		err := unlockFunc(unlockCtx)
		if err != nil {
			slog.Error(
				"Failed to unlock bucket",
				slog.String("app_id", l.appID.String()),
				slog.Any("error", err),
			)
		}
		return
	}
}

func (l *IdentifyRateLimiter) shardIdentifyBucket(shardID int) int {
	return shardID % l.maxConcurrency
}
