package store

import (
	"context"
	"encoding/json"
	"time"

	"github.com/disgoorg/snowflake/v2"
	"github.com/merlinfuchs/stateway/stateway-cache/model"
)

type UpsertEmojiParams struct {
	AppID     snowflake.ID
	GuildID   snowflake.ID
	EmojiID   snowflake.ID
	Data      json.RawMessage
	CreatedAt time.Time
	UpdatedAt time.Time
}

type SearchEmojisParams struct {
	AppID   snowflake.ID
	GuildID snowflake.ID
	Limit   int
	Offset  int
	Data    json.RawMessage
}

type CacheEmojiStore interface {
	GetEmoji(ctx context.Context, appID snowflake.ID, guildID snowflake.ID, emojiID snowflake.ID) (*model.Emoji, error)
	GetEmojis(ctx context.Context, appID snowflake.ID, guildID snowflake.ID, limit int, offset int) ([]*model.Emoji, error)
	SearchEmojis(ctx context.Context, params SearchEmojisParams) ([]*model.Emoji, error)
	UpsertEmojis(ctx context.Context, emojis ...UpsertEmojiParams) error
	DeleteEmoji(ctx context.Context, appID snowflake.ID, guildID snowflake.ID, emojiID snowflake.ID) error
}
