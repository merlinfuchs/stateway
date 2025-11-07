CREATE TABLE IF NOT EXISTS cache.channels (
    group_id TEXT NOT NULL,
    client_id BIGINT NOT NULL,
    guild_id BIGINT NOT NULL,
    channel_id BIGINT NOT NULL,
    data JSONB NOT NULL,
    tainted BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,

    PRIMARY KEY (group_id, client_id, guild_id, channel_id)
);
