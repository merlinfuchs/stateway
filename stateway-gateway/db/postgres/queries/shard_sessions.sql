-- name: UpsertShardSession :exec
INSERT INTO gateway.shard_sessions (
    id,
    app_id,
    shard_id,
    last_sequence,
    resume_url,
    created_at,
    updated_at
)
VALUES ($1, $2, $3, $4, $5, $6, $7)
ON CONFLICT (id) DO UPDATE SET
    last_sequence = EXCLUDED.last_sequence,
    resume_url = EXCLUDED.resume_url,
    updated_at = EXCLUDED.updated_at;

-- name: GetLastShardSession :one
SELECT * FROM gateway.shard_sessions WHERE app_id = $1 AND shard_id = $2 ORDER BY updated_at DESC LIMIT 1;

-- name: InvalidateShardSession :exec
UPDATE gateway.shard_sessions SET invalidated_at = NOW() WHERE app_id = $1 AND shard_id = $2;

-- name: PurgeSessions :exec
DELETE FROM gateway.shard_sessions;
