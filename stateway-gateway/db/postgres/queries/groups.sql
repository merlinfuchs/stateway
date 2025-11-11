-- name: CreateGroup :one
INSERT INTO gateway.groups (
    id, 
    display_name, 
    default_config, 
    default_constraints, 
    created_at, 
    updated_at
) VALUES ($1, $2, $3, $4, $5, $6) 
RETURNING *;

-- name: UpdateGroup :one
UPDATE gateway.groups SET 
    display_name = $2, 
    default_config = $3, 
    default_constraints = $4, 
    updated_at = $5 
WHERE id = $1
RETURNING *;

-- name: UpsertGroup :one
INSERT INTO gateway.groups (
    id, 
    display_name, 
    default_config, 
    default_constraints, 
    created_at, 
    updated_at
) VALUES ($1, $2, $3, $4, $5, $6) ON CONFLICT (id) DO UPDATE SET 
    display_name = EXCLUDED.display_name, 
    default_config = EXCLUDED.default_config, 
    default_constraints = EXCLUDED.default_constraints, 
    updated_at = EXCLUDED.updated_at
RETURNING *;

-- name: DeleteGroup :exec
DELETE FROM gateway.groups WHERE id = $1;

-- name: GetGroup :one
SELECT * FROM gateway.groups WHERE id = $1 LIMIT 1;

-- name: GetGroups :many
SELECT * FROM gateway.groups;
