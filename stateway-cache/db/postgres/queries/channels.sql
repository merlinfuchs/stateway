-- name: UpsertChannels :batchexec
INSERT INTO cache.channels (
    group_id, 
    client_id, 
    guild_id, 
    channel_id, 
    data, 
    created_at, 
    updated_at
) VALUES ($1, $2, $3, $4, $5, $6, $7) 
ON CONFLICT (group_id, client_id, guild_id, channel_id) DO UPDATE SET 
    data = EXCLUDED.data, 
    tainted = FALSE,
    updated_at = EXCLUDED.updated_at;

-- name: DeleteChannel :exec
DELETE FROM cache.channels WHERE group_id = $1 AND client_id = $2 AND guild_id = $3 AND channel_id = $4;

-- name: MarkShardChannelsTainted :exec
UPDATE cache.channels SET tainted = TRUE WHERE group_id = $1 AND client_id = $2 AND guild_id % @shard_count = @shard_id;
