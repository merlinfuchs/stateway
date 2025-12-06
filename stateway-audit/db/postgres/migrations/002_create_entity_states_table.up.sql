CREATE UNLOGGED TABLE IF NOT EXISTS audit.entity_states (
    app_id BIGINT NOT NULL,
    guild_id BIGINT NOT NULL,
    entity_type TEXT NOT NULL,
    entity_id BIGINT NOT NULL,
    data JSONB NOT NULL,
    deleted BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,

    PRIMARY KEY (app_id, guild_id, entity_type, entity_id)
);
