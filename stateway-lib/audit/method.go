package audit

import (
	"encoding/json"
	"fmt"

	"github.com/disgoorg/snowflake/v2"
)

type AuditMethod string

const (
	AuditMethodConfigureAuditLogging AuditMethod = "audit.configure"
	AuditMethodListEntityChanges     AuditMethod = "audit.changes.list"
)

func (m AuditMethod) UnmarshalRequest(data json.RawMessage) (AuditRequest, error) {
	switch m {
	case AuditMethodConfigureAuditLogging:
		var req ConfigureAuditLoggingRequest
		err := json.Unmarshal(data, &req)
		return req, err
	case AuditMethodListEntityChanges:
		var req ListEntityChangesRequest
		err := json.Unmarshal(data, &req)
		return req, err
	default:
		return nil, fmt.Errorf("unknown audit method: %v", m)
	}
}

type AuditRequest interface {
	auditRequest()
}

type ConfigureAuditLoggingRequest struct {
	GuildID snowflake.ID `json:"guild_id"`
	Enabled bool         `json:"enabled"`
	Options AuditOptions `json:"options"`
}

func (r ConfigureAuditLoggingRequest) auditRequest() {}

type ListEntityChangesRequest struct {
	GuildID    snowflake.ID  `json:"guild_id"`
	EntityType *EntityType   `json:"entity_type"`
	EntityID   *snowflake.ID `json:"entity_id"`
	Path       *string       `json:"path"`
	Options    AuditOptions  `json:"options"`
}

func (r ListEntityChangesRequest) auditRequest() {}
