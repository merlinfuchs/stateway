-- name: GetGuild :one
SELECT * FROM cache.guilds WHERE app_id = $1 AND guild_id = $2 LIMIT 1;

-- name: GetGuilds :many
SELECT * FROM cache.guilds WHERE app_id = $1 ORDER BY guild_id LIMIT sqlc.narg('limit') OFFSET sqlc.narg('offset');

-- name: SearchGuilds :many
SELECT * FROM cache.guilds WHERE app_id = $1 AND data @> $2 ORDER BY guild_id LIMIT sqlc.narg('limit') OFFSET sqlc.narg('offset');

-- name: UpsertGuilds :batchexec
INSERT INTO cache.guilds (
    app_id, 
    guild_id, 
    data, 
    created_at, 
    updated_at
) VALUES ($1, $2, $3, $4, $5) 
ON CONFLICT (app_id, guild_id) DO UPDATE SET 
    data = EXCLUDED.data, 
    tainted = FALSE,
    updated_at = EXCLUDED.updated_at;

-- name: DeleteGuild :exec
DELETE FROM cache.guilds WHERE app_id = $1 AND guild_id = $2;

-- name: MarkGuildUnavailable :exec
UPDATE cache.guilds SET unavailable = TRUE WHERE app_id = $1 AND guild_id = $2;

-- name: MarkShardGuildsTainted :exec
UPDATE cache.guilds SET tainted = TRUE WHERE app_id = $1 AND guild_id % @shard_count = @shard_id;
