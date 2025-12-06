package clickhouse

import (
	"context"
	"fmt"

	"github.com/merlinfuchs/stateway/stateway-audit/model"
)

func (c *Client) InsertEntityChanges(ctx context.Context, entityChanges ...model.EntityChange) error {
	if len(entityChanges) == 0 {
		return nil
	}

	batch, err := c.Conn.PrepareBatch(ctx, "INSERT INTO entity_changes")
	if err != nil {
		return fmt.Errorf("failed to prepare batch: %w", err)
	}

	for _, change := range entityChanges {
		// Convert nullable fields to pointers for ClickHouse Nullable types
		// Empty strings are converted to nil (null in ClickHouse)
		err := batch.Append(
			change.GuildID,
			change.EntityType,
			change.EntityID,
			change.EventID,
			change.Source, // event_source in table
			nullableUint64(change.AuditLogID),
			nullableUint64(change.AuditLogUserID),
			nullableString(change.AuditLogReason),
			change.Key,
			nullableString(change.OldValue), // Nullable(String) - null when entity was created
			nullableString(change.NewValue), // Nullable(String) - null when entity was deleted
			change.ReceivedAt,
			change.ProcessedAt,
			change.IngestedAt,
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
