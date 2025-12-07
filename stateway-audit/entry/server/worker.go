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
	"github.com/go-openapi/jsonpointer"
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

	data, err := gateway.UnmarshalEventData(event.Data, gateway.EventType(event.Type))
	if err != nil {
		return fmt.Errorf("failed to unmarshal event data: %w", err)
	}

	fmt.Printf("data: TYPE %T\n", data)

	switch d := data.(type) {
	case gateway.EventChannelUpdate:
		err = l.handleEntityChange(ctx, event, d.GuildID(), model.EntityTypeChannel, d.GuildChannel.ID(), d.GuildChannel)
	case gateway.EventGuildAuditLogEntryCreate:
		// TODO: Match audit log entry to entity change
	}

	if err != nil {
		return fmt.Errorf("failed to handle event: %w", err)
	}

	return nil
}

func (l *AuditWorker) handleEntityChange(
	ctx context.Context,
	event *event.GatewayEvent,
	guildID snowflake.ID,
	entityType model.EntityType,
	entityID snowflake.ID,
	data any,
) error {
	slog.Debug("Handling entity change", slog.String("entity_type", string(entityType)), slog.String("entity_id", entityID.String()))

	eventID := snowflake.New(time.Now().UTC())

	entityState, err := l.entityStateStore.GetEntityState(ctx, event.AppID, guildID, entityType, entityID)
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			// Entity was created or previously unknown
			newData, err := json.Marshal(data)
			if err != nil {
				return fmt.Errorf("failed to marshal entity data: %w", err)
			}

			err = l.batcher.Push(ctx, model.EntityChange{
				AppID:       event.AppID,
				GuildID:     guildID,
				EntityType:  entityType,
				EntityID:    entityID,
				EventID:     eventID,
				EventSource: model.EventSourceDiscord,
				Path:        "",
				Operation:   model.OperationAdd,
				OldValue:    nil,
				NewValue:    newData,
				ReceivedAt:  time.Now().UTC(),
				ProcessedAt: time.Now().UTC(),
				IngestedAt:  time.Now().UTC(),
			})
			if err != nil {
				return fmt.Errorf("failed to push entity change: %w", err)
			}

			err = l.entityStateStore.UpsertEntityState(ctx, model.EntityState{
				AppID:      event.AppID,
				GuildID:    guildID,
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
		return fmt.Errorf("failed to get entity state: %w", err)
	}

	if data == nil {
		// Entity was deleted
		err = l.batcher.Push(ctx, model.EntityChange{
			AppID:       event.AppID,
			GuildID:     guildID,
			EntityType:  entityType,
			EntityID:    entityID,
			EventID:     eventID,
			EventSource: model.EventSourceDiscord,
			Path:        "",
			Operation:   model.OperationRemove,
			OldValue:    entityState.Data,
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
			GuildID:    guildID,
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

	oldJD, err := jd.ReadJsonString(string(entityState.Data))
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
		jd.PathOption(jd.Path{jd.PathKey("permission_overwrites")}, jd.SET, jd.SetKeys("id")),
	)
	if len(diff) == 0 {
		return nil
	}

	// TODO: Wait for audit log event

	for _, entry := range diff {
		path, err := formatChangePath(entry.Path)
		if err != nil {
			return fmt.Errorf("failed to format change path: %w", err)
		}

		// Replace operation
		if len(entry.Remove) == 1 && len(entry.Add) == 1 {
			err = l.batcher.Push(ctx, model.EntityChange{
				AppID:       event.AppID,
				GuildID:     guildID,
				EntityType:  entityType,
				EntityID:    entityID,
				EventID:     eventID,
				EventSource: model.EventSourceDiscord,
				Path:        path,
				Operation:   model.OperationReplace,
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
				GuildID:     guildID,
				EntityType:  entityType,
				EntityID:    entityID,
				EventID:     eventID,
				EventSource: model.EventSourceDiscord,
				Path:        path,
				Operation:   model.OperationAdd,
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
				GuildID:     guildID,
				EntityType:  entityType,
				EntityID:    entityID,
				EventID:     eventID,
				EventSource: model.EventSourceDiscord,
				Path:        path,
				Operation:   model.OperationRemove,
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
		GuildID:    guildID,
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

func formatChangePath(path jd.Path) (string, error) {
	formattedPath := ""
	for _, element := range path {
		if formattedPath != "" {
			formattedPath += "/"
		}

		switch element := element.(type) {
		case jd.PathIndex:
			formattedPath += fmt.Sprintf("%d", element)
		case jd.PathKey:
			formattedPath += jsonpointer.Escape(string(element))
		case jd.PathSet:
			formattedPath += "#"
		case jd.PathSetKeys:
			for _, value := range element {
				formattedPath += jsonpointer.Escape(fmt.Sprintf("#%v", value))
			}
		default:
			return "", fmt.Errorf("unknown path element type: %T", element)
		}
	}

	return formattedPath, nil
}
