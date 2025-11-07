CREATE TABLE IF NOT EXISTS gateway.groups (
    id TEXT PRIMARY KEY,
    display_name TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL
);

INSERT INTO gateway.groups (id, display_name, created_at, updated_at) VALUES ('default', 'Default', NOW(), NOW());
