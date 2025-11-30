CREATE UNLOGGED TABLE IF NOT EXISTS cache.roles (
    app_id BIGINT NOT NULL,
    guild_id BIGINT NOT NULL,
    role_id BIGINT NOT NULL,
    data JSONB NOT NULL,
    tainted BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,

    PRIMARY KEY (app_id, guild_id, role_id) WITH (fillfactor = 80)
) WITH (fillfactor = 90);

CREATE INDEX IF NOT EXISTS idx_cache_roles_app_id_role_id ON cache.roles (app_id, role_id) WITH (fillfactor = 80);
