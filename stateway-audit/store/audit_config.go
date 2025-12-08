package store

import (
	"context"

	"github.com/disgoorg/snowflake/v2"
	"github.com/merlinfuchs/stateway/stateway-audit/model"
)

type AuditConfigStore interface {
	GetAuditConfig(ctx context.Context, appID snowflake.ID, guildID snowflake.ID) (*model.AuditConfig, error)
	UpsertAuditConfig(ctx context.Context, auditConfig model.AuditConfig) error
}
