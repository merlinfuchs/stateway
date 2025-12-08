-- name: UpsertAuditConfig :exec
INSERT INTO audit.config (app_id, guild_id, enabled, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5)
ON CONFLICT (app_id, guild_id)
DO UPDATE SET enabled = $3, updated_at = $5;

-- name: GetAuditConfig :one
SELECT * FROM audit.config WHERE app_id = $1 AND guild_id = $2;
