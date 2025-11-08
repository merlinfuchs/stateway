package app

import (
	"context"
	"io"
	"log/slog"
	"time"

	"github.com/disgoorg/disgo"
	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/disgo/gateway"
	"github.com/disgoorg/disgo/sharding"
	"github.com/disgoorg/snowflake/v2"
	"github.com/merlinfuchs/stateway/stateway-gateway/model"
	"github.com/merlinfuchs/stateway/stateway-gateway/store"
	"github.com/merlinfuchs/stateway/stateway-lib/event"
	"gopkg.in/guregu/null.v4"
)

type AppConfig struct {
	GatewayCount int
	GatewayID    int
}

type App struct {
	cfg   AppConfig
	model *model.App

	appStore     store.AppStore
	eventHandler event.EventHandler

	client *bot.Client
}

func NewApp(
	cfg AppConfig,
	model *model.App,
	appStore store.AppStore,
	eventHandler event.EventHandler,
) *App {
	return &App{
		cfg:          cfg,
		model:        model,
		appStore:     appStore,
		eventHandler: eventHandler,
	}
}

func (a *App) Run(ctx context.Context) {
	shardCount, shardIDs := shardsFromApp(a.model, a.cfg.GatewayCount, a.cfg.GatewayID)
	intents := intentsFromApp(a.model)
	presenceOpts := presenceOptsFromApp(a.model)

	client, err := disgo.New(a.model.DiscordBotToken,
		bot.WithLogger(slog.Default()),
		bot.WithShardManagerConfigOpts(
			sharding.WithAutoScaling(false),
			sharding.WithShardCount(shardCount),
			sharding.WithShardIDs(shardIDs...),
			sharding.WithGatewayConfigOpts(
				gateway.WithIntents(intents),
				gateway.WithCompression(gateway.CompressionZstdStream),
				gateway.WithEnableRawEvents(true),
				gateway.WithPresenceOpts(presenceOpts...),
			),
		),
		bot.WithEventListenerFunc(func(event *events.Ready) {
			slog.Info(
				"Discord shard READY",
				slog.String("group_id", a.model.GroupID),
				slog.String("app_id", a.model.ID.String()),
				slog.Int("shard_id", event.ShardID()),
				slog.String("display_name", a.model.DisplayName),
			)
		}),
		bot.WithEventListenerFunc(func(event *events.Resumed) {
			slog.Info(
				"Discord shard RESUMED",
				slog.String("group_id", a.model.GroupID),
				slog.String("app_id", a.model.ID.String()),
				slog.Int("shard_id", event.ShardID()),
				slog.String("display_name", a.model.DisplayName),
			)
		}),
		bot.WithEventListenerFunc(func(e *events.Raw) {
			data, err := io.ReadAll(e.Payload)
			if err != nil {
				slog.Error("Failed to read event payload", slog.Any("error", err))
				return
			}

			a.eventHandler.HandleEvent(&event.GatewayEvent{
				ID:        snowflake.New(time.Now().UTC()),
				GatewayID: a.cfg.GatewayID,
				AppID:     a.model.ID,
				GroupID:   a.model.GroupID,
				ShardID:   e.ShardID(),
				Type:      string(e.EventType),
				Data:      data,
			})
		}),
	)
	if err != nil {
		// TODO: Detect invalid token
		a.disable(ctx, model.AppDisabledCodeUnknown, err.Error())
		slog.Error("Failed to create Discord client", slog.Any("error", err))
		return
	}

	a.client = client

	err = client.OpenShardManager(ctx)
	if err != nil {
		// TODO: Detect invalid token
		a.disable(ctx, model.AppDisabledCodeUnknown, err.Error())
		slog.Error("Failed to open Discord gateway", slog.Any("error", err))
		return
	}
}

func (a *App) Close(ctx context.Context) {
	a.client.Close(ctx)
	a.client = nil
}

func (a *App) Update(ctx context.Context, model *model.App) {
	a.Close(ctx)
	// TODO: Re-use client if possible
	a.model = model
	go a.Run(ctx)
}

func (a *App) disable(ctx context.Context, code model.AppDisabledCode, message string) {
	_, err := a.appStore.DisableApp(ctx, store.DisableAppParams{
		ID:              a.model.ID,
		DisabledCode:    code,
		DisabledMessage: null.NewString(message, message != ""),
		UpdatedAt:       time.Now().UTC(),
	})
	if err != nil {
		slog.Error(
			"Failed to disable app",
			slog.String("app_id", a.model.ID.String()),
			slog.String("group_id", a.model.GroupID),
			slog.String("code", string(code)),
			slog.String("message", message),
			slog.Any("error", err),
		)
		return
	}
}
