CREATE TABLE IF NOT EXISTS cache.guilds (
    app_id BIGINT NOT NULL,
    guild_id BIGINT NOT NULL,
    data JSONB NOT NULL,
    unavailable BOOLEAN NOT NULL DEFAULT FALSE,
    tainted BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,

    PRIMARY KEY (app_id, guild_id)
);
