-- name: UpsertRoles :batchexec
INSERT INTO cache.roles (
    group_id, 
    client_id, 
    guild_id, 
    role_id, 
    data, 
    created_at, 
    updated_at
) VALUES ($1, $2, $3, $4, $5, $6, $7) 
ON CONFLICT (group_id, client_id, guild_id, role_id) DO UPDATE SET 
    data = EXCLUDED.data, 
    tainted = FALSE,
    updated_at = EXCLUDED.updated_at;

-- name: DeleteRole :exec
DELETE FROM cache.roles WHERE group_id = $1 AND client_id = $2 AND guild_id = $3 AND role_id = $4;

-- name: MarkShardRolesTainted :exec
UPDATE cache.roles SET tainted = TRUE WHERE group_id = $1 AND client_id = $2 AND guild_id % @shard_count = @shard_id;
