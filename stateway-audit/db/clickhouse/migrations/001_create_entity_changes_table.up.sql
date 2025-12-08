CREATE TABLE audit_entity_changes (
    -- Entity identification
    app_id      UInt64,
    guild_id    UInt64,
    entity_type LowCardinality(String),   -- 'guild', 'channel', 'role', ...
    entity_id   String,                   -- Discord ID as string

    -- Event identification
    event_id     String,                 -- unique ID for this logical event
    event_type   LowCardinality(String), -- 'GUILD_CREATE', 'GUILD_UPDATE', 'CHANNEL_CREATE', ...
    event_source LowCardinality(String), -- 'dispatch', 'guild_sync'

    -- Information from the audit log
    audit_log_id      Nullable(UInt64),
    audit_log_action  Nullable(UInt16),
    audit_log_user_id Nullable(UInt64),
    audit_log_reason  Nullable(String),

    -- Values as JSON strings (can be scalar, object, or array)
    path       String,           -- JSON path that changed
    operation  LowCardinality(String), -- 'add', 'remove', 'replace'
    old_value Nullable(String), -- JSON-encoded "before" value (Null when entity was created)
    new_value Nullable(String), -- JSON-encoded "after" value (Null when entity was deleted)

    received_at DateTime, -- When the change was received from the gateway or internally produced
    processed_at DateTime, -- When the change was processed and enqueued for ingestion
    ingested_at DateTime DEFAULT now() -- When the change was ingested into the database
)
ENGINE = MergeTree
ORDER BY (
    app_id,
    guild_id,
    entity_type,
    entity_id,
    received_at,
    event_id,
    path
);
