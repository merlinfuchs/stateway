CREATE TABLE IF NOT EXISTS audit.config (
    app_id BIGINT NOT NULL,
    guild_id BIGINT NOT NULL,
    audit_enabled BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,

    PRIMARY KEY (app_id, guild_id)
);
