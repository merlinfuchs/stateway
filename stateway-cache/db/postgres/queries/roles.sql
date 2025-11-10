-- name: GetRole :one
SELECT * FROM cache.roles WHERE app_id = $1 AND guild_id = $2 AND role_id = $3 LIMIT 1;

-- name: GetRoles :many
SELECT * FROM cache.roles WHERE app_id = $1 AND guild_id = $2 ORDER BY role_id LIMIT $3 OFFSET $4;

-- name: SearchRoles :many
SELECT * FROM cache.roles WHERE app_id = $1 AND guild_id = $2 AND data @> $3 ORDER BY role_id LIMIT $4 OFFSET $5;

-- name: UpsertRoles :batchexec
INSERT INTO cache.roles (
    app_id, 
    guild_id, 
    role_id, 
    data, 
    created_at, 
    updated_at
) VALUES ($1, $2, $3, $4, $5, $6) 
ON CONFLICT (app_id, guild_id, role_id) DO UPDATE SET 
    data = EXCLUDED.data, 
    tainted = FALSE,
    updated_at = EXCLUDED.updated_at;

-- name: DeleteRole :exec
DELETE FROM cache.roles WHERE app_id = $1 AND guild_id = $2 AND role_id = $3;

-- name: MarkShardRolesTainted :exec
UPDATE cache.roles SET tainted = TRUE WHERE app_id = $1 AND guild_id % @shard_count = @shard_id;
