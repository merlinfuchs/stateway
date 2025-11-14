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

func (c *Client) GetGuildSticker(ctx context.Context, appID snowflake.ID, guildID snowflake.ID, stickerID snowflake.ID) (*model.Sticker, error) {
	row, err := c.Q.GetGuildSticker(ctx, pgmodel.GetGuildStickerParams{
		AppID:     int64(appID),
		GuildID:   int64(guildID),
		StickerID: int64(stickerID),
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, store.ErrNotFound
		}
		return nil, err
	}
	return rowToSticker(row)
}

func (c *Client) GetSticker(ctx context.Context, appID snowflake.ID, guildID snowflake.ID, stickerID snowflake.ID) (*model.Sticker, error) {
	row, err := c.Q.GetSticker(ctx, pgmodel.GetStickerParams{
		AppID:     int64(appID),
		StickerID: int64(stickerID),
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, store.ErrNotFound
		}
		return nil, err
	}
	return rowToSticker(row)
}

func (c *Client) GetGuildStickers(ctx context.Context, appID snowflake.ID, guildID snowflake.ID, limit int, offset int) ([]*model.Sticker, error) {
	rows, err := c.Q.GetGuildStickers(ctx, pgmodel.GetGuildStickersParams{
		AppID:   int64(appID),
		GuildID: int64(guildID),
		Limit: pgtype.Int4{
			Int32: int32(limit),
			Valid: limit != 0,
		},
		Offset: pgtype.Int4{
			Int32: int32(offset),
			Valid: offset != 0,
		},
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, store.ErrNotFound
		}
		return nil, err
	}

	stickers := make([]*model.Sticker, len(rows))
	for i, row := range rows {
		sticker, err := rowToSticker(row)
		if err != nil {
			return nil, err
		}
		stickers[i] = sticker
	}
	return stickers, nil
}

func (c *Client) GetStickers(ctx context.Context, appID snowflake.ID, limit int, offset int) ([]*model.Sticker, error) {
	rows, err := c.Q.GetStickers(ctx, pgmodel.GetStickersParams{
		AppID: int64(appID),
		Limit: pgtype.Int4{
			Int32: int32(limit),
			Valid: limit != 0,
		},
		Offset: pgtype.Int4{
			Int32: int32(offset),
			Valid: offset != 0,
		},
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, store.ErrNotFound
		}
		return nil, err
	}

	stickers := make([]*model.Sticker, len(rows))
	for i, row := range rows {
		sticker, err := rowToSticker(row)
		if err != nil {
			return nil, err
		}
		stickers[i] = sticker
	}
	return stickers, nil
}

func (c *Client) SearchGuildStickers(ctx context.Context, params store.SearchGuildStickersParams) ([]*model.Sticker, error) {
	rows, err := c.Q.SearchGuildStickers(ctx, pgmodel.SearchGuildStickersParams{
		AppID:   int64(params.AppID),
		GuildID: int64(params.GuildID),
		Data:    params.Data,
		Limit: pgtype.Int4{
			Int32: int32(params.Limit),
			Valid: params.Limit != 0,
		},
		Offset: pgtype.Int4{
			Int32: int32(params.Offset),
			Valid: params.Offset != 0,
		},
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, store.ErrNotFound
		}
		return nil, err
	}

	stickers := make([]*model.Sticker, len(rows))
	for i, row := range rows {
		sticker, err := rowToSticker(row)
		if err != nil {
			return nil, err
		}
		stickers[i] = sticker
	}
	return stickers, nil
}

func (c *Client) SearchStickers(ctx context.Context, params store.SearchStickersParams) ([]*model.Sticker, error) {
	rows, err := c.Q.SearchStickers(ctx, pgmodel.SearchStickersParams{
		AppID: int64(params.AppID),
		Data:  params.Data,
		Limit: pgtype.Int4{
			Int32: int32(params.Limit),
			Valid: params.Limit != 0,
		},
		Offset: pgtype.Int4{
			Int32: int32(params.Offset),
			Valid: params.Offset != 0,
		},
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, store.ErrNotFound
		}
		return nil, err
	}

	stickers := make([]*model.Sticker, len(rows))
	for i, row := range rows {
		sticker, err := rowToSticker(row)
		if err != nil {
			return nil, err
		}
		stickers[i] = sticker
	}
	return stickers, nil
}

func (c *Client) CountGuildStickers(ctx context.Context, appID snowflake.ID, guildID snowflake.ID) (int, error) {
	res, err := c.Q.CountGuildStickers(ctx, pgmodel.CountGuildStickersParams{
		AppID:   int64(appID),
		GuildID: int64(guildID),
	})
	if err != nil {
		return 0, err
	}
	return int(res), nil
}

func (c *Client) CountStickers(ctx context.Context, appID snowflake.ID) (int, error) {
	res, err := c.Q.CountStickers(ctx, int64(appID))
	if err != nil {
		return 0, err
	}
	return int(res), nil
}

func (c *Client) UpsertStickers(ctx context.Context, stickers ...store.UpsertStickerParams) error {
	if len(stickers) == 0 {
		return nil
	}

	params := make([]pgmodel.UpsertStickersParams, len(stickers))
	for i, sticker := range stickers {
		data, err := json.Marshal(sticker.Data)
		if err != nil {
			return fmt.Errorf("failed to marshal sticker data: %w", err)
		}

		params[i] = pgmodel.UpsertStickersParams{
			AppID:     int64(sticker.AppID),
			GuildID:   int64(sticker.GuildID),
			StickerID: int64(sticker.StickerID),
			Data:      data,
			CreatedAt: pgtype.Timestamp{
				Time:  sticker.CreatedAt,
				Valid: true,
			},
			UpdatedAt: pgtype.Timestamp{
				Time:  sticker.UpdatedAt,
				Valid: true,
			},
		}
	}
	res := c.Q.UpsertStickers(ctx, params)
	return res.Close()
}

func (c *Client) DeleteSticker(ctx context.Context, appID snowflake.ID, guildID snowflake.ID, stickerID snowflake.ID) error {
	return c.Q.DeleteSticker(ctx, pgmodel.DeleteStickerParams{
		AppID:     int64(appID),
		GuildID:   int64(guildID),
		StickerID: int64(stickerID),
	})
}

func rowToSticker(row pgmodel.CacheSticker) (*model.Sticker, error) {
	var data discord.Sticker
	err := json.Unmarshal(row.Data, &data)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal sticker data: %w", err)
	}

	return &model.Sticker{
		AppID:     snowflake.ID(row.AppID),
		GuildID:   snowflake.ID(row.GuildID),
		StickerID: snowflake.ID(row.StickerID),
		Data:      data,
		Tainted:   row.Tainted,
		CreatedAt: row.CreatedAt.Time,
		UpdatedAt: row.UpdatedAt.Time,
	}, nil
}
