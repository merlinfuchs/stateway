package store

import (
	"context"
	"encoding/json"
	"time"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/snowflake/v2"
	"github.com/merlinfuchs/stateway/stateway-cache/model"
)

type UpsertEmojiParams struct {
	AppID     snowflake.ID
	GuildID   snowflake.ID
	EmojiID   snowflake.ID
	Data      discord.Emoji
	CreatedAt time.Time
	UpdatedAt time.Time
}

type SearchEmojisParams struct {
	AppID  snowflake.ID
	Limit  int
	Offset int
	Data   json.RawMessage
}

type SearchGuildEmojisParams struct {
	AppID   snowflake.ID
	GuildID snowflake.ID
	Limit   int
	Offset  int
	Data    json.RawMessage
}

type CacheEmojiStore interface {
	GetGuildEmoji(ctx context.Context, appID snowflake.ID, guildID snowflake.ID, emojiID snowflake.ID) (*model.Emoji, error)
	GetEmoji(ctx context.Context, appID snowflake.ID, emojiID snowflake.ID) (*model.Emoji, error)
	GetGuildEmojis(ctx context.Context, appID snowflake.ID, guildID snowflake.ID, limit int, offset int) ([]*model.Emoji, error)
	GetEmojis(ctx context.Context, appID snowflake.ID, limit int, offset int) ([]*model.Emoji, error)
	SearchGuildEmojis(ctx context.Context, params SearchGuildEmojisParams) ([]*model.Emoji, error)
	SearchEmojis(ctx context.Context, params SearchEmojisParams) ([]*model.Emoji, error)
	CountGuildEmojis(ctx context.Context, appID snowflake.ID, guildID snowflake.ID) (int, error)
	CountEmojis(ctx context.Context, appID snowflake.ID) (int, error)
	UpsertEmojis(ctx context.Context, emojis ...UpsertEmojiParams) error
	DeleteEmoji(ctx context.Context, appID snowflake.ID, guildID snowflake.ID, emojiID snowflake.ID) error
}
