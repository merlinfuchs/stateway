package audit

import (
	"context"
	"time"

	"github.com/disgoorg/snowflake/v2"
)

type ConfigureAuditLoggingParams struct {
	GuildID snowflake.ID
	Enabled bool
}

type GetEntityChangesParams struct {
	GuildID    snowflake.ID
	Before     *time.Time
	After      *time.Time
	EntityType *EntityType
	EntityID   *snowflake.ID
	Path       *string
}

type Auditor interface {
	ConfigureAuditLogging(ctx context.Context, params ConfigureAuditLoggingParams, opts ...AuditOption) error
	GetEntityChanges(ctx context.Context, params GetEntityChangesParams, opts ...AuditOption) ([]*EntityChange, error)
}
