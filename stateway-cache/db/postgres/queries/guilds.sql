-- name: UpsertGuild :exec
INSERT INTO cache.guilds (
    id, 
    app_id, 
    data, 
    created_at, 
    updated_at
) VALUES ($1, $2, $3, $4, $5) 
ON CONFLICT (id, app_id) DO UPDATE SET 
    data = EXCLUDED.data, 
    updated_at = EXCLUDED.updated_at;
