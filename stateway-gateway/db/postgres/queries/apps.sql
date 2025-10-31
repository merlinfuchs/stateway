-- name: CreateApp :one
INSERT INTO gateway.apps (
    id, 
    display_name, 
    discord_client_id, 
    discord_bot_token, 
    discord_public_key, 
    discord_client_secret,
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
    $8
)
RETURNING *;

-- name: UpdateApp :one
UPDATE gateway.apps SET 
    display_name = $2, 
    discord_client_id = $3, 
    discord_bot_token = $4, 
    discord_public_key = $5, 
    discord_client_secret = $6,
    disabled = $7,
    disabled_code = $8,
    disabled_message = $9,
    updated_at = $10
WHERE id = $1
RETURNING *;

-- name: DisableApp :one
UPDATE gateway.apps SET 
    disabled = TRUE,
    disabled_code = $2,
    disabled_message = $3,
    updated_at = $4
WHERE id = $1
RETURNING *;

-- name: DeleteApp :exec
DELETE FROM gateway.apps WHERE id = $1;

-- name: GetApp :one
SELECT * FROM gateway.apps WHERE id = $1 LIMIT 1;

-- name: GetApps :many
SELECT * FROM gateway.apps;

-- name: GetEnabledApps :many
SELECT * FROM gateway.apps WHERE disabled = FALSE;