package server

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/gateway"
	"github.com/disgoorg/snowflake/v2"
	"github.com/go-openapi/jsonpointer"
	"github.com/josephburnett/jd/v2"
	"github.com/merlinfuchs/stateway/stateway-audit/batcher"
	"github.com/merlinfuchs/stateway/stateway-audit/model"
	"github.com/merlinfuchs/stateway/stateway-audit/store"
	"github.com/merlinfuchs/stateway/stateway-lib/audit"
	"github.com/merlinfuchs/stateway/stateway-lib/broker"
	"github.com/merlinfuchs/stateway/stateway-lib/event"
	"github.com/nats-io/nats.go/jetstream"
)

type AuditWorkerConfig struct {
	GatewayIDs []int
	NamePrefix string
}

type AuditWorker struct {
	entityStateStore store.EntityStateStore
	batcher          batcher.Batcher
	auditLogMatcher  *AuditLogMatcher

	config AuditWorkerConfig
}

func NewAuditWorker(
	auditLogMatcher *AuditLogMatcher,
	entityStateStore store.EntityStateStore,
	batcher batcher.Batcher,
	config AuditWorkerConfig,
) *AuditWorker {
	return &AuditWorker{
		entityStateStore: entityStateStore,
		batcher:          batcher,
		auditLogMatcher:  auditLogMatcher,
		config:           config,
	}
}

func (l *AuditWorker) BalanceKey() string {
	key := fmt.Sprintf("%s_AUDIT_CHANGE_WORKER", l.config.NamePrefix)
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
			"guild.create",
			"guild.update",
			"guild.delete",
			"channel.create",
			"channel.delete",
			"channel.update",
			"guild.role.create",
			"guild.role.update",
			"guild.role.delete",
			"invite.create",
			"invite.delete",
		},
	}
}

func (l *AuditWorker) ConsumerConfig() broker.ConsumerConfig {
	return broker.ConsumerConfig{
		AckPolicy:     jetstream.AckExplicitPolicy,
		MaxAckPending: 1000,
		Async:         true,
	}
}

func (l *AuditWorker) HandleEvent(ctx context.Context, event *event.GatewayEvent) (bool, error) {
	slog.Debug("Received event:", slog.String("type", event.Type))

	data, err := gateway.UnmarshalEventData(event.Data, gateway.EventType(event.Type))
	if err != nil {
		return false, fmt.Errorf("failed to unmarshal event data: %w", err)
	}

	switch d := data.(type) {
	case gateway.EventGuildCreate:
		// All events here are synthetic because they are part of the initial guild sync
		err = l.handleEntityChange(
			ctx,
			event,
			event.Type,
			audit.EventSourceDiscord,
			d.Guild.ID,
			audit.EntityTypeGuild,
			d.Guild.ID,
			d.Guild,
			nil,
		)
		if err != nil {
			return false, fmt.Errorf("failed to handle entity change: %w", err)
		}

		for _, channel := range d.Channels {
			err = l.handleEntityChange(
				ctx,
				event,
				event.Type,
				audit.EventSourceDiscord,
				d.Guild.ID,
				audit.EntityTypeChannel,
				channel.ID(),
				channel,
				nil,
			)
			if err != nil {
				return false, fmt.Errorf("failed to handle entity change: %w", err)
			}
		}
		for _, role := range d.Roles {
			err = l.handleEntityChange(
				ctx,
				event,
				event.Type,
				audit.EventSourceDiscord,
				d.Guild.ID,
				audit.EntityTypeRole,
				role.ID,
				role,
				nil,
			)
			if err != nil {
				return false, fmt.Errorf("failed to handle entity change: %w", err)
			}
		}
	case gateway.EventGuildUpdate:
		auditLogInfo := l.auditLogMatcher.WaitForAuditLogAny(
			ctx, d.Guild.ID,
			discord.AuditLogEventGuildUpdate,
		)

		err = l.handleEntityChange(
			ctx,
			event,
			event.Type,
			audit.EventSourceDiscord,
			d.Guild.ID,
			audit.EntityTypeGuild,
			d.Guild.ID,
			d.Guild,
			auditLogInfo,
		)
	case gateway.EventChannelCreate:
		auditLogInfo := l.auditLogMatcher.WaitForAuditLogAny(
			ctx, d.GuildChannel.ID(),
			discord.AuditLogEventChannelCreate,
		)

		err = l.handleEntityChange(
			ctx,
			event,
			event.Type,
			audit.EventSourceDiscord,
			d.GuildID(),
			audit.EntityTypeChannel,
			d.GuildChannel.ID(),
			d.GuildChannel,
			auditLogInfo,
		)
	case gateway.EventChannelUpdate:
		auditLogInfo := l.auditLogMatcher.WaitForAuditLogAny(
			ctx, d.GuildChannel.ID(),
			discord.AuditLogEventChannelUpdate,
			discord.AuditLogEventChannelOverwriteCreate,
			discord.AuditLogEventChannelOverwriteUpdate,
			discord.AuditLogEventChannelOverwriteDelete,
		)

		err = l.handleEntityChange(
			ctx,
			event,
			event.Type,
			audit.EventSourceDiscord,
			d.GuildID(),
			audit.EntityTypeChannel,
			d.GuildChannel.ID(),
			d.GuildChannel,
			auditLogInfo,
		)
	case gateway.EventChannelDelete:
		auditLogInfo := l.auditLogMatcher.WaitForAuditLogAny(
			ctx, d.GuildChannel.ID(),
			discord.AuditLogEventChannelDelete,
		)

		err = l.handleEntityChange(
			ctx,
			event,
			event.Type,
			audit.EventSourceDiscord,
			d.GuildID(),
			audit.EntityTypeChannel,
			d.GuildChannel.ID(),
			nil,
			auditLogInfo,
		)
	case gateway.EventGuildRoleCreate:
		auditLogInfo := l.auditLogMatcher.WaitForAuditLogAny(
			ctx, d.Role.ID,
			discord.AuditLogEventRoleCreate,
		)

		err = l.handleEntityChange(
			ctx,
			event,
			event.Type,
			audit.EventSourceDiscord,
			d.GuildID,
			audit.EntityTypeRole,
			d.Role.ID,
			d.Role,
			auditLogInfo,
		)
	case gateway.EventGuildRoleUpdate:
		auditLogInfo := l.auditLogMatcher.WaitForAuditLogAny(
			ctx, d.Role.ID,
			discord.AuditLogEventRoleUpdate,
		)

		err = l.handleEntityChange(
			ctx,
			event,
			event.Type,
			audit.EventSourceDiscord,
			d.GuildID,
			audit.EntityTypeRole,
			d.Role.ID,
			d.Role,
			auditLogInfo,
		)
	case gateway.EventGuildRoleDelete:
		auditLogInfo := l.auditLogMatcher.WaitForAuditLogAny(
			ctx, d.RoleID,
			discord.AuditLogEventRoleDelete,
		)

		err = l.handleEntityChange(
			ctx,
			event,
			event.Type,
			audit.EventSourceDiscord,
			d.GuildID,
			audit.EntityTypeRole,
			d.RoleID,
			nil,
			auditLogInfo,
		)
	case gateway.EventGuildAuditLogEntryCreate:
		l.auditLogMatcher.HandleAuditLog(d)
		return true, nil
	}

	if err != nil {
		return false, fmt.Errorf("failed to handle event: %w", err)
	}

	return true, nil
}

func (l *AuditWorker) handleEntityChange(
	ctx context.Context,
	event *event.GatewayEvent,
	eventType string,
	eventSource audit.EventSource,
	guildID snowflake.ID,
	entityType audit.EntityType,
	entityID snowflake.ID,
	data any,
	auditLogInfo *AuditLogInfo,
) error {
	slog.Debug("Handling entity change", slog.String("entity_type", string(entityType)), slog.String("entity_id", entityID.String()))

	eventID := snowflake.New(time.Now().UTC())

	entityState, err := l.entityStateStore.GetEntityState(ctx, event.AppID, guildID, entityType, entityID)
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			// We never knew about this entity, so we can't create a change for it
			if data == nil {
				return nil
			}

			// Entity was created or previously unknown
			newData, err := json.Marshal(data)
			if err != nil {
				return fmt.Errorf("failed to marshal entity data: %w", err)
			}

			fmt.Printf("entity state not found, creating new entity change\n")
			entityChange := model.EntityChange{
				AppID:       event.AppID,
				GuildID:     guildID,
				EntityType:  entityType,
				EntityID:    entityID,
				EventID:     eventID,
				EventType:   eventType,
				EventSource: eventSource,
				Path:        "",
				Operation:   audit.JSONOperationAdd,
				OldValue:    nil,
				NewValue:    newData,
				ReceivedAt:  time.Now().UTC(),
				ProcessedAt: time.Now().UTC(),
			}
			if auditLogInfo != nil {
				entityChange.AuditLogID = &auditLogInfo.ID
				entityChange.AuditLogAction = &auditLogInfo.Action
				entityChange.AuditLogUserID = &auditLogInfo.UserID
				entityChange.AuditLogReason = auditLogInfo.Reason
			}

			err = l.batcher.Push(ctx, entityChange)
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
		entityChange := model.EntityChange{
			AppID:       event.AppID,
			GuildID:     guildID,
			EntityType:  entityType,
			EntityID:    entityID,
			EventID:     eventID,
			EventType:   eventType,
			EventSource: eventSource,
			Path:        "",
			Operation:   audit.JSONOperationRemove,
			OldValue:    entityState.Data,
			NewValue:    nil,
			ReceivedAt:  time.Now().UTC(),
			ProcessedAt: time.Now().UTC(),
		}
		if auditLogInfo != nil {
			entityChange.AuditLogID = &auditLogInfo.ID
			entityChange.AuditLogUserID = &auditLogInfo.UserID
			entityChange.AuditLogAction = &auditLogInfo.Action
			entityChange.AuditLogReason = auditLogInfo.Reason
		}

		err = l.batcher.Push(ctx, entityChange)
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

	diff := oldJD.Diff(newJD, diffOptions(entityType)...)
	if len(diff) == 0 {
		return nil
	}

	for _, entry := range diff {
		path, err := formatDiffPath(entry.Path)
		if err != nil {
			return fmt.Errorf("failed to format change path: %w", err)
		}

		// Replace operation
		if len(entry.Remove) == 1 && len(entry.Add) == 1 {
			entityChange := model.EntityChange{
				AppID:       event.AppID,
				GuildID:     guildID,
				EntityType:  entityType,
				EntityID:    entityID,
				EventID:     eventID,
				EventType:   eventType,
				EventSource: eventSource,
				Path:        path,
				Operation:   audit.JSONOperationReplace,
				OldValue:    json.RawMessage(entry.Remove[0].Json()),
				NewValue:    json.RawMessage(entry.Add[0].Json()),
				ReceivedAt:  time.Now().UTC(),
				ProcessedAt: time.Now().UTC(),
			}
			if auditLogInfo != nil {
				entityChange.AuditLogID = &auditLogInfo.ID
				entityChange.AuditLogUserID = &auditLogInfo.UserID
				entityChange.AuditLogAction = &auditLogInfo.Action
				entityChange.AuditLogReason = auditLogInfo.Reason
			}
			err = l.batcher.Push(ctx, entityChange)
			if err != nil {
				return fmt.Errorf("failed to push entity change: %w", err)
			}
			continue
		}

		// Add operations
		for _, value := range entry.Add {
			entityChange := model.EntityChange{
				AppID:       event.AppID,
				GuildID:     guildID,
				EntityType:  entityType,
				EntityID:    entityID,
				EventID:     eventID,
				EventType:   eventType,
				EventSource: eventSource,
				Path:        path,
				Operation:   audit.JSONOperationAdd,
				OldValue:    nil,
				NewValue:    json.RawMessage(value.Json()),
				ReceivedAt:  time.Now().UTC(),
				ProcessedAt: time.Now().UTC(),
			}
			if auditLogInfo != nil {
				entityChange.AuditLogID = &auditLogInfo.ID
				entityChange.AuditLogUserID = &auditLogInfo.UserID
				entityChange.AuditLogAction = &auditLogInfo.Action
				entityChange.AuditLogReason = auditLogInfo.Reason
			}
			err = l.batcher.Push(ctx, entityChange)
			if err != nil {
				return fmt.Errorf("failed to push entity change: %w", err)
			}
		}

		// Remove operations
		for _, value := range entry.Remove {
			entityChange := model.EntityChange{
				AppID:       event.AppID,
				GuildID:     guildID,
				EntityType:  entityType,
				EntityID:    entityID,
				EventID:     eventID,
				EventType:   eventType,
				EventSource: eventSource,
				Path:        path,
				Operation:   audit.JSONOperationRemove,
				OldValue:    json.RawMessage(value.Json()),
				NewValue:    nil,
				ReceivedAt:  time.Now().UTC(),
				ProcessedAt: time.Now().UTC(),
			}
			if auditLogInfo != nil {
				entityChange.AuditLogID = &auditLogInfo.ID
				entityChange.AuditLogUserID = &auditLogInfo.UserID
				entityChange.AuditLogAction = &auditLogInfo.Action
				entityChange.AuditLogReason = auditLogInfo.Reason
			}
			err = l.batcher.Push(ctx, entityChange)
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

func formatDiffPath(path jd.Path) (string, error) {
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
			for key, value := range element {
				formattedPath += jsonpointer.Escape(fmt.Sprintf("%v:%v", key, value))
			}
		default:
			return "", fmt.Errorf("unknown path element type: %T", element)
		}
	}

	return formattedPath, nil
}

func diffOptions(entityType audit.EntityType) []jd.Option {
	switch entityType {
	case audit.EntityTypeChannel:
		return []jd.Option{
			jd.PathOption(jd.Path{jd.PathKey("permission_overwrites")}, jd.SET, jd.SetKeys("id")),
			jd.PathOption(jd.Path{jd.PathKey("guild_id")}, jd.DIFF_OFF),
			jd.PathOption(jd.Path{jd.PathKey("last_message_id")}, jd.DIFF_OFF),
			jd.PathOption(jd.Path{jd.PathKey("last_pin_timestamp")}, jd.DIFF_OFF),
		}
	case audit.EntityTypeRole:
		return []jd.Option{
			jd.PathOption(jd.Path{jd.PathKey("position")}, jd.DIFF_OFF),
			jd.PathOption(jd.Path{jd.PathKey("guild_id")}, jd.DIFF_OFF),
			jd.PathOption(jd.Path{jd.PathKey("tags")}, jd.DIFF_OFF),
		}
	case audit.EntityTypeGuild:
		return []jd.Option{
			jd.PathOption(jd.Path{jd.PathKey("features")}, jd.DIFF_OFF),
			jd.PathOption(jd.Path{jd.PathKey("joined_at")}, jd.DIFF_OFF),
			jd.PathOption(jd.Path{jd.PathKey("member_count")}, jd.DIFF_OFF),
			jd.PathOption(jd.Path{jd.PathKey("premium_subscription_count")}, jd.DIFF_OFF),
		}
	default:
		return []jd.Option{}
	}
}
