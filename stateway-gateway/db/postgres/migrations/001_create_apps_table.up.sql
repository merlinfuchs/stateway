CREATE TABLE IF NOT EXISTS gateway.apps (
    id BIGINT PRIMARY KEY,
    display_name TEXT NOT NULL,
    discord_client_id BIGINT NOT NULL,
    discord_bot_token TEXT NOT NULL,
    discord_public_key TEXT NOT NULL,
    discord_client_secret TEXT,
    disabled BOOLEAN NOT NULL DEFAULT FALSE,
    disabled_code TEXT,
    disabled_message TEXT,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL
);
