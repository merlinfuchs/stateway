package store

import (
	"context"

	"github.com/disgoorg/snowflake/v2"
)

type MarkShardEntitiesTaintedParams struct {
	AppID      snowflake.ID
	ShardCount int
	ShardID    int
}

type MassUpsertEntitiesParams struct {
	AppID    snowflake.ID
	Guilds   []UpsertGuildParams
	Roles    []UpsertRoleParams
	Channels []UpsertChannelParams
	Emojis   []UpsertEmojiParams
	Stickers []UpsertStickerParams
}

type CacheStore interface {
	CacheGuildStore
	CacheRoleStore
	CacheChannelStore
	CacheEmojiStore
	CacheStickerStore

	MarkShardEntitiesTainted(ctx context.Context, params MarkShardEntitiesTaintedParams) error
	MassUpsertEntities(ctx context.Context, params MassUpsertEntitiesParams) error
}
