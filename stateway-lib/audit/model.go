package audit

import (
	"encoding/json"
	"time"

	"github.com/disgoorg/snowflake/v2"
)

type EntityType string

const (
	EntityTypeGuild   EntityType = "guild"
	EntityTypeChannel EntityType = "channel"
	EntityTypeRole    EntityType = "role"
)

type EventSource string

const (
	EventSourceDiscord EventSource = "discord"
)

type JSONOperation string

const (
	JSONOperationAdd     JSONOperation = "ADD"
	JSONOperationRemove  JSONOperation = "REMOVE"
	JSONOperationReplace JSONOperation = "REPLACE"
)

type EntityChange struct {
	AppID          snowflake.ID    `json:"app_id"`
	GuildID        snowflake.ID    `json:"guild_id"`
	EntityType     EntityType      `json:"entity_type"`
	EntityID       snowflake.ID    `json:"entity_id"`
	EventID        snowflake.ID    `json:"event_id"`
	EventType      string          `json:"event_type"`
	EventSource    EventSource     `json:"event_source"`
	AuditLogID     *snowflake.ID   `json:"audit_log_id"`
	AuditLogAction *int            `json:"audit_log_action"`
	AuditLogUserID *snowflake.ID   `json:"audit_log_user_id"`
	AuditLogReason *string         `json:"audit_log_reason"`
	Path           string          `json:"path"`
	Operation      JSONOperation   `json:"operation"`
	OldValue       json.RawMessage `json:"old_value"`
	NewValue       json.RawMessage `json:"new_value"`
	ReceivedAt     time.Time       `json:"received_at"`
	ProcessedAt    time.Time       `json:"processed_at"`
}

type AuditOptions struct {
	AppID  snowflake.ID `json:"app_id"`
	Limit  int          `json:"limit,omitempty"`
	Offset int          `json:"offset,omitempty"`
}

func ResolveOptions(opts ...AuditOption) AuditOptions {
	options := AuditOptions{}
	for _, opt := range opts {
		opt(&options)
	}
	return options
}

func (o AuditOptions) Destructure() []AuditOption {
	res := []AuditOption{}
	if o.AppID != 0 {
		res = append(res, WithAppID(o.AppID))
	}
	return res
}

type AuditOption func(*AuditOptions)

func WithAppID(appID snowflake.ID) AuditOption {
	return func(o *AuditOptions) {
		o.AppID = appID
	}
}

func WithLimit(limit int) AuditOption {
	return func(o *AuditOptions) {
		o.Limit = limit
	}
}

func WithOffset(offset int) AuditOption {
	return func(o *AuditOptions) {
		o.Offset = offset
	}
}
