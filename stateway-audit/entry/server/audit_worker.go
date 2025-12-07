package server

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/disgoorg/disgo/gateway"
	"github.com/merlinfuchs/stateway/stateway-lib/broker"
	"github.com/merlinfuchs/stateway/stateway-lib/event"
)

type AuditWorkerConfig struct {
	GatewayIDs []int
	NamePrefix string
}

type AuditWorker struct {
	auditLogMatcher *AuditLogMatcher

	config AuditWorkerConfig
}

func NewAuditWorker(
	auditLogMatcher *AuditLogMatcher,
	config AuditWorkerConfig,
) *AuditWorker {
	return &AuditWorker{
		auditLogMatcher: auditLogMatcher,
		config:          config,
	}
}

func (l *AuditWorker) BalanceKey() string {
	key := fmt.Sprintf("%s_AUDIT_LOG_WORKER", l.config.NamePrefix)
	for _, gatewayID := range l.config.GatewayIDs {
		key += fmt.Sprintf("_%d", gatewayID)
	}
	if len(l.config.GatewayIDs) == 0 {
		key += "_ALL"
	}
	return key
}

func (l *AuditWorker) EventFilter() broker.EventFilter {
	return broker.EventFilter{
		GatewayIDs: l.config.GatewayIDs,
		EventTypes: []string{
			"guild.audit.log.entry.create",
		},
	}
}

func (l *AuditWorker) HandleEvent(ctx context.Context, event *event.GatewayEvent) error {
	slog.Debug("Received event:", slog.String("type", event.Type))

	data, err := gateway.UnmarshalEventData(event.Data, gateway.EventType(event.Type))
	if err != nil {
		return fmt.Errorf("failed to unmarshal event data: %w", err)
	}

	if d, ok := data.(gateway.EventGuildAuditLogEntryCreate); ok {
		l.auditLogMatcher.HandleAuditLog(d)
	}

	return nil
}
