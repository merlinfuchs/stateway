package clickhouse

import (
	"context"
	"fmt"
	"time"

	"github.com/disgoorg/snowflake/v2"
	"github.com/merlinfuchs/stateway/stateway-audit/model"
	"github.com/merlinfuchs/stateway/stateway-audit/store"
)

func (c *Client) InsertEntityChanges(ctx context.Context, entityChanges ...model.EntityChange) error {
	if len(entityChanges) == 0 {
		return nil
	}

	batch, err := c.Conn.PrepareBatch(ctx, "INSERT INTO audit_entity_changes")
	if err != nil {
		return fmt.Errorf("failed to prepare batch: %w", err)
	}

	for _, change := range entityChanges {
		var oldValue *string
		if len(change.OldValue) > 0 && string(change.OldValue) != "null" {
			oldValue = nullableString(string(change.OldValue))
		}
		var newValue *string
		if len(change.NewValue) > 0 && string(change.NewValue) != "null" {
			newValue = nullableString(string(change.NewValue))
		}

		// Convert nullable fields to pointers for ClickHouse Nullable types
		// Empty strings are converted to nil (null in ClickHouse)
		err := batch.Append(
			change.AppID,
			change.GuildID,
			string(change.EntityType),
			change.EntityID,
			change.EventID,
			change.EventType,
			string(change.EventSource),
			nullableSnowflake(change.AuditLogID),
			nullableAuditLogAction(change.AuditLogAction),
			nullableSnowflake(change.AuditLogUserID),
			change.AuditLogReason,
			change.Path,
			string(change.Operation),
			oldValue,
			newValue,
			change.ReceivedAt,
			change.ProcessedAt,
			time.Now().UTC(),
		)
		if err != nil {
			return fmt.Errorf("failed to append entity change to batch: %w", err)
		}
	}

	if err := batch.Send(); err != nil {
		return fmt.Errorf("failed to send batch: %w", err)
	}

	return nil
}

func (c *Client) GetEntityChanges(ctx context.Context, params store.GetEntityChangesParams) ([]*model.EntityChange, error) {
	return nil, nil
}

func nullableAuditLogAction(v *int) *uint16 {
	if v == nil {
		return nil
	}
	action := uint16(*v)
	return &action
}

func nullableSnowflake(v *snowflake.ID) *uint64 {
	if v == nil {
		return nil
	}
	return nullableUint64(uint64(*v))
}

func nullableUint64(v uint64) *uint64 {
	if v == 0 {
		return nil
	}
	return &v
}

func nullableString(v string) *string {
	if v == "" {
		return nil
	}
	return &v
}
