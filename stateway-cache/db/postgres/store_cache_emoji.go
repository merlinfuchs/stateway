package postgres

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/snowflake/v2"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/merlinfuchs/stateway/stateway-cache/db/postgres/pgmodel"
	"github.com/merlinfuchs/stateway/stateway-cache/model"
	"github.com/merlinfuchs/stateway/stateway-cache/store"
)

func (c *Client) GetEmoji(ctx context.Context, appID snowflake.ID, guildID snowflake.ID, emojiID snowflake.ID) (*model.Emoji, error) {
	row, err := c.Q.GetEmoji(ctx, pgmodel.GetEmojiParams{
		AppID:   int64(appID),
		GuildID: int64(guildID),
		EmojiID: int64(emojiID),
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, store.ErrNotFound
		}
		return nil, err
	}
	return rowToEmoji(row)
}

func (c *Client) GetEmojis(ctx context.Context, appID snowflake.ID, guildID snowflake.ID, limit int, offset int) ([]*model.Emoji, error) {
	rows, err := c.Q.GetEmojis(ctx, pgmodel.GetEmojisParams{
		AppID:   int64(appID),
		GuildID: int64(guildID),
		Limit:   int32(limit),
		Offset:  int32(offset),
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, store.ErrNotFound
		}
		return nil, err
	}

	emojis := make([]*model.Emoji, len(rows))
	for i, row := range rows {
		emoji, err := rowToEmoji(row)
		if err != nil {
			return nil, err
		}
		emojis[i] = emoji
	}
	return emojis, nil
}

func (c *Client) SearchEmojis(ctx context.Context, params store.SearchEmojisParams) ([]*model.Emoji, error) {
	rows, err := c.Q.SearchEmojis(ctx, pgmodel.SearchEmojisParams{
		AppID:   int64(params.AppID),
		GuildID: int64(params.GuildID),
		Data:    params.Data,
		Limit:   int32(params.Limit),
		Offset:  int32(params.Offset),
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, store.ErrNotFound
		}
		return nil, err
	}

	emojis := make([]*model.Emoji, len(rows))
	for i, row := range rows {
		emoji, err := rowToEmoji(row)
		if err != nil {
			return nil, err
		}
		emojis[i] = emoji
	}
	return emojis, nil
}

func (c *Client) UpsertEmojis(ctx context.Context, emojis ...store.UpsertEmojiParams) error {
	if len(emojis) == 0 {
		return nil
	}

	params := make([]pgmodel.UpsertEmojisParams, len(emojis))
	for i, emoji := range emojis {
		data, err := json.Marshal(emoji.Data)
		if err != nil {
			return fmt.Errorf("failed to marshal emoji data: %w", err)
		}

		params[i] = pgmodel.UpsertEmojisParams{
			AppID:   int64(emoji.AppID),
			GuildID: int64(emoji.GuildID),
			EmojiID: int64(emoji.EmojiID),
			Data:    data,
			CreatedAt: pgtype.Timestamp{
				Time:  emoji.CreatedAt,
				Valid: true,
			},
			UpdatedAt: pgtype.Timestamp{
				Time:  emoji.UpdatedAt,
				Valid: true,
			},
		}
	}
	res := c.Q.UpsertEmojis(ctx, params)
	return res.Close()
}

func (c *Client) DeleteEmoji(ctx context.Context, appID snowflake.ID, guildID snowflake.ID, emojiID snowflake.ID) error {
	return c.Q.DeleteEmoji(ctx, pgmodel.DeleteEmojiParams{
		AppID:   int64(appID),
		GuildID: int64(guildID),
		EmojiID: int64(emojiID),
	})
}

func rowToEmoji(row pgmodel.CacheEmoji) (*model.Emoji, error) {
	var data discord.Emoji
	err := json.Unmarshal(row.Data, &data)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal emoji data: %w", err)
	}

	return &model.Emoji{
		AppID:     snowflake.ID(row.AppID),
		GuildID:   snowflake.ID(row.GuildID),
		EmojiID:   snowflake.ID(row.EmojiID),
		Data:      data,
		Tainted:   row.Tainted,
		CreatedAt: row.CreatedAt.Time,
		UpdatedAt: row.UpdatedAt.Time,
	}, nil
}
