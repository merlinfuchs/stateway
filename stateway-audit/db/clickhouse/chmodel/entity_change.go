package chmodel

import "time"

type EntityChange struct {
	GuildID        uint64    `ch:"guild_id"`
	EntityType     string    `ch:"entity_type"`
	EntityID       string    `ch:"entity_id"`
	EventID        string    `ch:"event_id"`
	EventKind      string    `ch:"event_kind"`
	Source         string    `ch:"source"`
	AuditLogID     uint64    `ch:"audit_log_id"`
	AuditLogUserID uint64    `ch:"audit_log_user_id"`
	AuditLogReason string    `ch:"audit_log_reason"`
	Key            string    `ch:"key"`
	OldValue       string    `ch:"old_value"`
	NewValue       string    `ch:"new_value"`
	ReceivedAt     time.Time `ch:"received_at"`
	ProcessedAt    time.Time `ch:"processed_at"`
	IngestedAt     time.Time `ch:"ingested_at"`
}
