package model

import (
	"encoding/json"
	"time"

	"github.com/disgoorg/snowflake/v2"
)

type EventSource string

const (
	EventSourceDiscord   EventSource = "discord"
	EventSourceSynthetic EventSource = "synthetic"
)

type EntityChange struct {
	AppID          snowflake.ID    `json:"app_id"`
	GuildID        snowflake.ID    `json:"guild_id"`
	EntityType     EntityType      `json:"entity_type"`
	EntityID       snowflake.ID    `json:"entity_id"`
	EventID        snowflake.ID    `json:"event_id"`
	EventSource    EventSource     `json:"event_source"`
	AuditLogID     snowflake.ID    `json:"audit_log_id"`
	AuditLogUserID snowflake.ID    `json:"audit_log_user_id"`
	AuditLogReason string          `json:"audit_log_reason"`
	Path           string          `json:"path"`
	OldValue       json.RawMessage `json:"old_value"`
	NewValue       json.RawMessage `json:"new_value"`
	ReceivedAt     time.Time       `json:"received_at"`
	ProcessedAt    time.Time       `json:"processed_at"`
	IngestedAt     time.Time       `json:"ingested_at"`
}
