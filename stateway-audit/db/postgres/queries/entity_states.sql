-- name: UpsertEntityState :exec
INSERT INTO audit.entity_states (app_id, guild_id, entity_type, entity_id, data, deleted, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
ON CONFLICT (app_id, guild_id, entity_type, entity_id) DO UPDATE SET
    data = $5,
    deleted = $6,
    updated_at = $8;

-- name: GetEntityState :one
SELECT * FROM audit.entity_states WHERE app_id = $1 AND guild_id = $2 AND entity_type = $3 AND entity_id = $4;
