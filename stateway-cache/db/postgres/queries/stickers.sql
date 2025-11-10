-- name: GetSticker :one
SELECT * FROM cache.stickers WHERE app_id = $1 AND guild_id = $2 AND sticker_id = $3 LIMIT 1;

-- name: GetStickers :many
SELECT * FROM cache.stickers WHERE app_id = $1 AND guild_id = $2 ORDER BY sticker_id LIMIT $3 OFFSET $4;

-- name: SearchStickers :many
SELECT * FROM cache.stickers WHERE app_id = $1 AND guild_id = $2 AND data @> $3 ORDER BY sticker_id LIMIT $4 OFFSET $5;

-- name: UpsertStickers :batchexec
INSERT INTO cache.stickers (
    app_id, 
    guild_id, 
    sticker_id, 
    data, 
    created_at, 
    updated_at
) VALUES ($1, $2, $3, $4, $5, $6) 
ON CONFLICT (app_id, guild_id, sticker_id) DO UPDATE SET 
    data = EXCLUDED.data, 
    tainted = FALSE,
    updated_at = EXCLUDED.updated_at;

-- name: DeleteSticker :exec
DELETE FROM cache.stickers WHERE app_id = $1 AND guild_id = $2 AND sticker_id = $3;

-- name: MarkShardStickersTainted :exec
UPDATE cache.stickers SET tainted = TRUE WHERE app_id = $1 AND guild_id % @shard_count = @shard_id;
