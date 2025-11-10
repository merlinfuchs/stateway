package postgres

import (
	"context"

	"github.com/disgoorg/snowflake/v2"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/merlinfuchs/stateway/stateway-cache/db/postgres/pgmodel"
	"github.com/merlinfuchs/stateway/stateway-cache/model"
	"github.com/merlinfuchs/stateway/stateway-cache/store"
)

func (c *Client) GetSticker(ctx context.Context, appID snowflake.ID, guildID snowflake.ID, stickerID snowflake.ID) (*model.Sticker, error) {
	row, err := c.Q.GetSticker(ctx, pgmodel.GetStickerParams{
		AppID:     int64(appID),
		GuildID:   int64(guildID),
		StickerID: int64(stickerID),
	})
	if err != nil {
		return nil, err
	}
	return rowToSticker(row), nil
}

func (c *Client) GetStickers(ctx context.Context, appID snowflake.ID, guildID snowflake.ID, limit int, offset int) ([]*model.Sticker, error) {
	rows, err := c.Q.GetStickers(ctx, pgmodel.GetStickersParams{
		AppID:   int64(appID),
		GuildID: int64(guildID),
		Limit:   int32(limit),
		Offset:  int32(offset),
	})
	if err != nil {
		return nil, err
	}

	stickers := make([]*model.Sticker, len(rows))
	for i, row := range rows {
		stickers[i] = rowToSticker(row)
	}
	return stickers, nil
}

func (c *Client) SearchStickers(ctx context.Context, params store.SearchStickersParams) ([]*model.Sticker, error) {
	rows, err := c.Q.SearchStickers(ctx, pgmodel.SearchStickersParams{
		AppID:   int64(params.AppID),
		GuildID: int64(params.GuildID),
		Data:    params.Data,
		Limit:   int32(params.Limit),
		Offset:  int32(params.Offset),
	})
	if err != nil {
		return nil, err
	}

	stickers := make([]*model.Sticker, len(rows))
	for i, row := range rows {
		stickers[i] = rowToSticker(row)
	}
	return stickers, nil
}

func (c *Client) UpsertStickers(ctx context.Context, stickers ...store.UpsertStickerParams) error {
	if len(stickers) == 0 {
		return nil
	}

	params := make([]pgmodel.UpsertStickersParams, len(stickers))
	for i, sticker := range stickers {
		params[i] = pgmodel.UpsertStickersParams{
			AppID:     int64(sticker.AppID),
			GuildID:   int64(sticker.GuildID),
			StickerID: int64(sticker.StickerID),
			Data:      sticker.Data,
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

func rowToSticker(row pgmodel.CacheSticker) *model.Sticker {
	return &model.Sticker{
		AppID:     snowflake.ID(row.AppID),
		GuildID:   snowflake.ID(row.GuildID),
		StickerID: snowflake.ID(row.StickerID),
		Data:      row.Data,
		CreatedAt: row.CreatedAt.Time,
		UpdatedAt: row.UpdatedAt.Time,
	}
}
