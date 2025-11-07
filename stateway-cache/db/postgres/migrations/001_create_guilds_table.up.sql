CREATE TABLE IF NOT EXISTS cache.guilds (
    id BIGINT NOT NULL,
    app_id BIGINT NOT NULL,
    data JSONB NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,

    PRIMARY KEY (id, app_id)
);

CREATE INDEX IF NOT EXISTS idx_guilds_app_id ON cache.guilds (app_id);
