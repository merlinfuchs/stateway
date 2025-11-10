package store

import (
	"context"
	"encoding/json"
	"time"

	"github.com/disgoorg/snowflake/v2"
	"github.com/merlinfuchs/stateway/stateway-cache/model"
)

type UpsertStickerParams struct {
	AppID     snowflake.ID
	GuildID   snowflake.ID
	StickerID snowflake.ID
	Data      json.RawMessage
	CreatedAt time.Time
	UpdatedAt time.Time
}

type SearchStickersParams struct {
	AppID   snowflake.ID
	GuildID snowflake.ID
	Limit   int
	Offset  int
	Data    json.RawMessage
}

type CacheStickerStore interface {
	GetSticker(ctx context.Context, appID snowflake.ID, guildID snowflake.ID, stickerID snowflake.ID) (*model.Sticker, error)
	GetStickers(ctx context.Context, appID snowflake.ID, guildID snowflake.ID, limit int, offset int) ([]*model.Sticker, error)
	SearchStickers(ctx context.Context, params SearchStickersParams) ([]*model.Sticker, error)
	UpsertStickers(ctx context.Context, stickers ...UpsertStickerParams) error
	DeleteSticker(ctx context.Context, appID snowflake.ID, guildID snowflake.ID, stickerID snowflake.ID) error
}
