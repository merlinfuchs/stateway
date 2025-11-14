CREATE TABLE IF NOT EXISTS cache.emojis (
    app_id BIGINT NOT NULL,
    guild_id BIGINT NOT NULL,
    emoji_id BIGINT NOT NULL,
    data JSONB NOT NULL,
    tainted BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,

    PRIMARY KEY (app_id, guild_id, emoji_id)
);

CREATE INDEX IF NOT EXISTS idx_cache_emojis_app_id_emoji_id ON cache.emojis (app_id, emoji_id);
