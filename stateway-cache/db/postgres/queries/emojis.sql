-- name: GetGuildEmoji :one
SELECT * FROM cache.emojis WHERE app_id = $1 AND guild_id = $2 AND emoji_id = $3 LIMIT 1;

-- name: GetEmoji :one
SELECT * FROM cache.emojis WHERE app_id = $1 AND emoji_id = $2 LIMIT 1;

-- name: GetGuildEmojis :many
SELECT * FROM cache.emojis WHERE app_id = $1 AND guild_id = $2 ORDER BY emoji_id LIMIT sqlc.narg('limit') OFFSET sqlc.narg('offset');

-- name: GetEmojis :many
SELECT * FROM cache.emojis WHERE app_id = $1 ORDER BY emoji_id LIMIT sqlc.narg('limit') OFFSET sqlc.narg('offset');

-- name: SearchGuildEmojis :many
SELECT * FROM cache.emojis WHERE app_id = $1 AND guild_id = $2 AND data @> $3 ORDER BY emoji_id LIMIT sqlc.narg('limit') OFFSET sqlc.narg('offset');

-- name: SearchEmojis :many
SELECT * FROM cache.emojis WHERE app_id = $1 AND data @> $2 ORDER BY emoji_id LIMIT sqlc.narg('limit') OFFSET sqlc.narg('offset');

-- name: CountGuildEmojis :one
SELECT COUNT(*) FROM cache.emojis WHERE app_id = $1 AND guild_id = $2;

-- name: CountEmojis :one
SELECT COUNT(*) FROM cache.emojis WHERE app_id = $1;

-- name: UpsertEmojis :batchexec
INSERT INTO cache.emojis (
    app_id, 
    guild_id, 
    emoji_id, 
    data, 
    created_at, 
    updated_at
) VALUES ($1, $2, $3, $4, $5, $6) 
ON CONFLICT (app_id, guild_id, emoji_id) DO UPDATE SET 
    data = EXCLUDED.data, 
    tainted = FALSE,
    updated_at = EXCLUDED.updated_at;

-- name: DeleteEmoji :exec
DELETE FROM cache.emojis WHERE app_id = $1 AND guild_id = $2 AND emoji_id = $3;

-- name: MarkShardEmojisTainted :exec
UPDATE cache.emojis SET tainted = TRUE WHERE app_id = $1 AND guild_id % @shard_count = @shard_id;
