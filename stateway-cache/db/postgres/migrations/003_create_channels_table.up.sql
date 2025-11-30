CREATE UNLOGGED TABLE IF NOT EXISTS cache.channels (
    app_id BIGINT NOT NULL,
    guild_id BIGINT NOT NULL,
    channel_id BIGINT NOT NULL,
    data JSONB NOT NULL,
    tainted BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,

    PRIMARY KEY (app_id, guild_id, channel_id) WITH (fillfactor = 80)
) WITH (fillfactor = 90);

/* CREATE INDEX IF NOT EXISTS idx_cache_channels_data_type ON cache.channels ((data->>'type')) WITH (fillfactor = 80); */
CREATE INDEX IF NOT EXISTS idx_cache_channels_app_id_channel_id ON cache.channels (app_id, channel_id) WITH (fillfactor = 80);
