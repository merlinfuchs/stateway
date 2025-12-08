package server

import (
	"context"
	"fmt"
	"time"

	"github.com/merlinfuchs/stateway/stateway-audit/model"
	"github.com/merlinfuchs/stateway/stateway-audit/store"
	"github.com/merlinfuchs/stateway/stateway-lib/audit"
)

type Auditor struct {
	entityChangeStore store.EntityChangeStore
	auditConfigStore  store.AuditConfigStore
}

func NewAuditor(entityChangeStore store.EntityChangeStore, auditConfigStore store.AuditConfigStore) *Auditor {
	return &Auditor{
		entityChangeStore: entityChangeStore,
		auditConfigStore:  auditConfigStore,
	}
}

func (l *Auditor) ConfigureAuditLogging(ctx context.Context, params audit.ConfigureAuditLoggingParams, opts ...audit.AuditOption) error {
	options := audit.ResolveOptions(opts...)

	err := l.auditConfigStore.UpsertAuditConfig(ctx, model.AuditConfig{
		AppID:     options.AppID,
		GuildID:   params.GuildID,
		Enabled:   params.Enabled,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	})
	if err != nil {
		return fmt.Errorf("failed to get audit config: %w", err)
	}

	return nil
}

func (l *Auditor) GetEntityChanges(ctx context.Context, params audit.GetEntityChangesParams, opts ...audit.AuditOption) ([]*audit.EntityChange, error) {
	options := audit.ResolveOptions(opts...)

	entityChanges, err := l.entityChangeStore.GetEntityChanges(ctx, store.GetEntityChangesParams{
		AppID:      options.AppID,
		GuildID:    params.GuildID,
		EntityType: params.EntityType,
		EntityID:   params.EntityID,
		Path:       params.Path,
		Before:     params.Before,
		After:      params.After,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get entity changes: %w", err)
	}

	return entityChanges, nil
}
