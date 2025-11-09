package admin

import (
	"context"

	"github.com/merlinfuchs/stateway/stateway-gateway/store"
)

func PurgeSessions(ctx context.Context, shardSessionStore store.ShardSessionStore) error {
	return shardSessionStore.PurgeSessions(ctx)
}
