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
	InstanceCount int
	InstanceIndex int
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
	shardCount := a.model.ShardCount
	if shardCount == 0 {
		shardCount = 1
	}

	shardIDs := make([]int, 0, shardCount)
	for shardID := 0; shardID < shardCount; shardID++ {
		// if shardCount == 1, this instance is the only one that should run the app
		// otherwise, we are splitting the app shards across the instances
		// shardID % instanceCount gives us the index of the instance that should run the shard
		if shardCount == 1 || shardID%a.cfg.InstanceCount == a.cfg.InstanceIndex {
			shardIDs = append(shardIDs, shardID)
		}
	}

	client, err := disgo.New(a.model.DiscordBotToken,
		bot.WithShardManagerConfigOpts(
			sharding.WithAutoScaling(false),
			sharding.WithShardCount(shardCount),
			sharding.WithShardIDs(shardIDs...),
			sharding.WithGatewayConfigOpts(
				gateway.WithIntents(gateway.IntentGuilds, gateway.IntentGuildMessages, gateway.IntentDirectMessages),
				gateway.WithCompression(gateway.CompressionZstdStream),
				gateway.WithEnableRawEvents(true),
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
				ID:       snowflake.New(time.Now().UTC()),
				AppID:    a.model.ID,
				GroupID:  a.model.GroupID,
				ClientID: a.model.DiscordClientID,
				ShardID:  e.ShardID(),
				Type:     string(e.EventType),
				Data:     data,
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
