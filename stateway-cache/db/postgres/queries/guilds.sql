-- name: GetGuild :one
SELECT * FROM cache.guilds WHERE group_id = $1 AND client_id = $2 AND guild_id = $3 LIMIT 1;

-- name: UpsertGuilds :batchexec
INSERT INTO cache.guilds (
    group_id, 
    client_id, 
    guild_id, 
    data, 
    created_at, 
    updated_at
) VALUES ($1, $2, $3, $4, $5, $6) 
ON CONFLICT (group_id, client_id, guild_id) DO UPDATE SET 
    data = EXCLUDED.data, 
    tainted = FALSE,
    updated_at = EXCLUDED.updated_at;

-- name: DeleteGuild :exec
DELETE FROM cache.guilds WHERE group_id = $1 AND client_id = $2 AND guild_id = $3;

-- name: MarkGuildUnavailable :exec
UPDATE cache.guilds SET unavailable = TRUE WHERE group_id = $1 AND client_id = $2 AND guild_id = $3;

-- name: MarkShardGuildsTainted :exec
UPDATE cache.guilds SET tainted = TRUE WHERE group_id = $1 AND client_id = $2 AND guild_id % @shard_count = @shard_id;
