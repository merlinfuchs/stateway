CREATE TABLE IF NOT EXISTS gateway.shard_sessions (
    id TEXT PRIMARY KEY,
    app_id BIGINT NOT NULL REFERENCES gateway.apps(id) ON DELETE CASCADE,
    shard_id INTEGER NOT NULL,
    shard_count INTEGER NOT NULL,
    last_sequence INTEGER NOT NULL,
    resume_url TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    invalidated_at TIMESTAMP
);
