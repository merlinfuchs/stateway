-- name: CreateApp :one
INSERT INTO gateway.apps (
    id,
    group_id,
    display_name, 
    discord_client_id, 
    discord_bot_token, 
    discord_public_key, 
    discord_client_secret,
    shard_count,
    constraints,
    config,
    created_at, 
    updated_at
)
VALUES (
    $1, 
    $2, 
    $3, 
    $4, 
    $5, 
    $6, 
    $7,
    $8,
    $9,
    $10,
    $11,
    $12
)
RETURNING *;

-- name: UpdateApp :one
UPDATE gateway.apps SET 
    group_id = $2,
    display_name = $3, 
    discord_client_id = $4, 
    discord_bot_token = $5, 
    discord_public_key = $6, 
    discord_client_secret = $7,
    shard_count = $8,
    constraints = $9,
    config = $10,
    disabled = $11,  
    disabled_code = $12,
    disabled_message = $13,
    updated_at = $14
WHERE id = $1
RETURNING *;

-- name: UpsertApp :one
INSERT INTO gateway.apps (
    id,
    group_id,
    display_name,
    discord_client_id,
    discord_bot_token,
    discord_public_key,
    discord_client_secret,
    shard_count,
    constraints,
    config,
    created_at,
    updated_at
)
VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
ON CONFLICT (id) DO UPDATE SET
    group_id = EXCLUDED.group_id,
    display_name = EXCLUDED.display_name,
    discord_client_id = EXCLUDED.discord_client_id,
    discord_bot_token = EXCLUDED.discord_bot_token,
    discord_public_key = EXCLUDED.discord_public_key,
    discord_client_secret = EXCLUDED.discord_client_secret,
    shard_count = EXCLUDED.shard_count,
    constraints = EXCLUDED.constraints,
    config = EXCLUDED.config,
    updated_at = EXCLUDED.updated_at,
    disabled = FALSE
RETURNING *;

-- name: DisableApp :exec
UPDATE gateway.apps SET 
    disabled = TRUE,
    disabled_code = $2,
    disabled_message = $3,
    updated_at = $4
WHERE id = $1;

-- name: DeleteApp :exec
DELETE FROM gateway.apps WHERE id = $1;

-- name: GetApp :one
SELECT * FROM gateway.apps WHERE id = $1 LIMIT 1;

-- name: GetApps :many
SELECT * FROM gateway.apps WHERE (group_id = sqlc.narg('group_id') OR sqlc.narg('group_id') IS NULL) LIMIT sqlc.narg('limit') OFFSET sqlc.narg('offset');

-- name: GetEnabledApps :many
SELECT * FROM gateway.apps WHERE disabled = FALSE AND (shard_count > 1 OR id % @gateway_count = @gateway_id);
