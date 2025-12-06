package model

import "time"

type EntityChange struct {
	GuildID        uint64    `json:"guild_id"`
	EntityType     string    `json:"entity_type"`
	EntityID       string    `json:"entity_id"`
	EventID        string    `json:"event_id"`
	EventKind      string    `json:"event_kind"`
	Source         string    `json:"source"`
	AuditLogID     uint64    `json:"audit_log_id"`
	AuditLogUserID uint64    `json:"audit_log_user_id"`
	AuditLogReason string    `json:"audit_log_reason"`
	Key            string    `json:"key"`
	OldValue       string    `json:"old_value"`
	NewValue       string    `json:"new_value"`
	ReceivedAt     time.Time `json:"received_at"`
	ProcessedAt    time.Time `json:"processed_at"`
	IngestedAt     time.Time `json:"ingested_at"`
}
