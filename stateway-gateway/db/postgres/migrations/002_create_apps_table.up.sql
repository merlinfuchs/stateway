CREATE TABLE IF NOT EXISTS gateway.apps (
    id BIGINT PRIMARY KEY,
    group_id TEXT NOT NULL DEFAULT 'default' REFERENCES gateway.groups(id) ON DELETE RESTRICT,
    display_name TEXT NOT NULL,
    discord_client_id BIGINT NOT NULL,
    discord_bot_token TEXT NOT NULL,
    discord_public_key TEXT NOT NULL,
    discord_client_secret TEXT,
    shard_count INTEGER NOT NULL DEFAULT 1,
    constraints JSONB,
    config JSONB,
    disabled BOOLEAN NOT NULL DEFAULT FALSE,
    disabled_code TEXT,
    disabled_message TEXT,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL
);
