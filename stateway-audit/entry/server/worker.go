package server

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/disgoorg/disgo/gateway"
	"github.com/disgoorg/snowflake/v2"
	"github.com/josephburnett/jd/v2"
	"github.com/merlinfuchs/stateway/stateway-audit/batcher"
	"github.com/merlinfuchs/stateway/stateway-audit/model"
	"github.com/merlinfuchs/stateway/stateway-audit/store"
	"github.com/merlinfuchs/stateway/stateway-lib/broker"
	"github.com/merlinfuchs/stateway/stateway-lib/event"
)

type AuditWorkerConfig struct {
	GatewayIDs []int
	NamePrefix string
}

type AuditWorker struct {
	entityStateStore store.EntityStateStore
	batcher          batcher.Batcher

	config AuditWorkerConfig
}

func NewAuditWorker(
	entityStateStore store.EntityStateStore,
	batcher batcher.Batcher,
	config AuditWorkerConfig,
) *AuditWorker {
	return &AuditWorker{
		entityStateStore: entityStateStore,
		batcher:          batcher,
		config:           config,
	}
}

func (l *AuditWorker) BalanceKey() string {
	key := fmt.Sprintf("%s_AUDIT_WORKER", l.config.NamePrefix)
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
			"ready",
			"guild.>",
			"channel.>",
			"thread.>",
		},
	}
}

func (l *AuditWorker) HandleEvent(ctx context.Context, event *event.GatewayEvent) error {
	slog.Debug("Received event:", slog.String("type", event.Type))

	if event.GuildID == nil {
		return nil
	}

	data, err := gateway.UnmarshalEventData(event.Data, gateway.EventType(event.Type))
	if err != nil {
		return fmt.Errorf("failed to unmarshal event data: %w", err)
	}

	switch d := data.(type) {
	case gateway.EventReady:
	case gateway.EventChannelUpdate:
		err = l.handleEntityChange(ctx, event, model.EntityTypeChannel, d.GuildChannel.ID(), d.GuildChannel)
	}

	if err != nil {
		return fmt.Errorf("failed to handle event: %w", err)
	}

	return nil
}

func (l *AuditWorker) handleEntityChange(
	ctx context.Context,
	event *event.GatewayEvent,
	entityType model.EntityType,
	entityID snowflake.ID,
	data any,
) error {
	var oldData json.RawMessage = []byte("{}")

	entityState, err := l.entityStateStore.GetEntityState(ctx, event.AppID, *event.GuildID, entityType, entityID)
	if err == nil {
		oldData = entityState.Data
	} else if !errors.Is(err, store.ErrNotFound) {
		return fmt.Errorf("failed to get entity state: %w", err)
	}

	if data == nil {
		// Entity was deleted
		err = l.batcher.Push(ctx, model.EntityChange{
			AppID:       event.AppID,
			GuildID:     *event.GuildID,
			EntityType:  entityType,
			EntityID:    entityID,
			EventID:     snowflake.New(time.Now().UTC()),
			EventSource: model.EventSourceDiscord,
			Path:        "",
			OldValue:    oldData,
			NewValue:    nil,
			ReceivedAt:  time.Now().UTC(),
			ProcessedAt: time.Now().UTC(),
			IngestedAt:  time.Now().UTC(),
		})
		if err != nil {
			return fmt.Errorf("failed to push entity change: %w", err)
		}

		err = l.entityStateStore.UpsertEntityState(ctx, model.EntityState{
			AppID:      event.AppID,
			GuildID:    *event.GuildID,
			EntityType: entityType,
			EntityID:   entityID,
			Data:       nil,
			Deleted:    true,
			CreatedAt:  time.Now().UTC(),
			UpdatedAt:  time.Now().UTC(),
		})
		if err != nil {
			return fmt.Errorf("failed to upsert entity state: %w", err)
		}

		return nil
	}

	newData, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to marshal entity data: %w", err)
	}

	oldJD, err := jd.ReadJsonString(string(oldData))
	if err != nil {
		return fmt.Errorf("failed to read old value: %w", err)
	}

	newJD, err := jd.ReadJsonString(string(newData))
	if err != nil {
		return fmt.Errorf("failed to read new value: %w", err)
	}

	diff := oldJD.Diff(
		newJD,
		// Ignore ordering of permission overwrites
		jd.PathOption(jd.Path{jd.PathKey("permission_overwrites")}, jd.SET),
	)
	if len(diff) == 0 {
		return nil
	}

	// TODO: Wait for audit log event

	for _, entry := range diff {
		path := entry.Path.JsonNode().Json()

		// Replace operation
		if len(entry.Remove) == 1 && len(entry.Add) == 1 {
			err = l.batcher.Push(ctx, model.EntityChange{
				AppID:       event.AppID,
				GuildID:     *event.GuildID,
				EntityType:  entityType,
				EntityID:    entityID,
				EventID:     snowflake.New(time.Now().UTC()),
				EventSource: model.EventSourceDiscord,
				Path:        path,
				OldValue:    json.RawMessage(entry.Remove[0].Json()),
				NewValue:    json.RawMessage(entry.Add[0].Json()),
				ReceivedAt:  time.Now().UTC(),
				ProcessedAt: time.Now().UTC(),
				IngestedAt:  time.Now().UTC(),
			})
			if err != nil {
				return fmt.Errorf("failed to push entity change: %w", err)
			}
			continue
		}

		// Add operations
		for _, value := range entry.Add {
			err = l.batcher.Push(ctx, model.EntityChange{
				AppID:       event.AppID,
				GuildID:     *event.GuildID,
				EntityType:  entityType,
				EntityID:    entityID,
				EventID:     snowflake.New(time.Now().UTC()),
				EventSource: model.EventSourceDiscord,
				Path:        path,
				OldValue:    nil,
				NewValue:    json.RawMessage(value.Json()),
				ReceivedAt:  time.Now().UTC(),
				ProcessedAt: time.Now().UTC(),
				IngestedAt:  time.Now().UTC(),
			})
			if err != nil {
				return fmt.Errorf("failed to push entity change: %w", err)
			}
		}

		// Remove operations
		for _, value := range entry.Remove {
			err = l.batcher.Push(ctx, model.EntityChange{
				AppID:       event.AppID,
				GuildID:     *event.GuildID,
				EntityType:  entityType,
				EntityID:    entityID,
				EventID:     snowflake.New(time.Now().UTC()),
				EventSource: model.EventSourceDiscord,
				Path:        path,
				OldValue:    json.RawMessage(value.Json()),
				NewValue:    nil,
				ReceivedAt:  time.Now().UTC(),
				ProcessedAt: time.Now().UTC(),
				IngestedAt:  time.Now().UTC(),
			})
			if err != nil {
				return fmt.Errorf("failed to push entity change: %w", err)
			}
		}
	}

	err = l.entityStateStore.UpsertEntityState(ctx, model.EntityState{
		AppID:      event.AppID,
		GuildID:    *event.GuildID,
		EntityType: entityType,
		EntityID:   entityID,
		Data:       newData,
		Deleted:    false,
		CreatedAt:  time.Now().UTC(),
		UpdatedAt:  time.Now().UTC(),
	})
	if err != nil {
		return fmt.Errorf("failed to upsert entity state: %w", err)
	}

	return nil
}
