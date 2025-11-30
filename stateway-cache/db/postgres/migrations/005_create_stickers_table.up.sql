CREATE UNLOGGED TABLE IF NOT EXISTS cache.stickers (
    app_id BIGINT NOT NULL,
    guild_id BIGINT NOT NULL,
    sticker_id BIGINT NOT NULL,
    data JSONB NOT NULL,
    tainted BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,

    PRIMARY KEY (app_id, guild_id, sticker_id) WITH (fillfactor = 80)
) WITH (fillfactor = 90);

CREATE INDEX IF NOT EXISTS idx_cache_stickers_app_id_sticker_id ON cache.stickers (app_id, sticker_id) WITH (fillfactor = 80);
