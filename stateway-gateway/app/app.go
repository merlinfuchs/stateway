package app

import (
	"context"
	"io"
	"log/slog"
	"time"

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

	shardManager sharding.ShardManager
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

	shardManager := sharding.New(
		a.model.DiscordBotToken,
		func(g gateway.Gateway, _ gateway.EventType, _ int, ev gateway.EventData) {
			switch e := ev.(type) {
			case gateway.EventRaw:
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
					ShardID:   g.ShardID(),
					Type:      string(e.EventType),
					Data:      data,
				})
			case gateway.EventReady:
				slog.Info(
					"Discord shard READY",
					slog.String("group_id", a.model.GroupID),
					slog.String("app_id", a.model.ID.String()),
					slog.Int("shard_id", g.ShardID()),
					slog.String("display_name", a.model.DisplayName),
				)
			case gateway.EventResumed:
				slog.Info(
					"Discord shard RESUMED",
					slog.String("group_id", a.model.GroupID),
					slog.String("app_id", a.model.ID.String()),
					slog.Int("shard_id", g.ShardID()),
					slog.String("display_name", a.model.DisplayName),
				)
			case gateway.EventRateLimited:
				slog.Info(
					"Discord shard RATE_LIMITED",
					slog.String("group_id", a.model.GroupID),
					slog.String("app_id", a.model.ID.String()),
					slog.Int("shard_id", g.ShardID()),
					slog.String("display_name", a.model.DisplayName),
				)
			}
		},
		sharding.WithAutoScaling(false),
		sharding.WithShardCount(shardCount),
		sharding.WithShardIDs(shardIDs...),
		sharding.WithLogger(slog.Default()),
		sharding.WithGatewayConfigOpts(
			gateway.WithIntents(intents),
			gateway.WithEnableRawEvents(true),
			gateway.WithPresenceOpts(presenceOpts...),
			gateway.WithLogger(slog.Default()),
		),
		// TODO: Add close handler that disables the app on some errors
	)

	a.shardManager = shardManager
	shardManager.Open(ctx)
}

func (a *App) Close(ctx context.Context) {
	a.shardManager.Close(ctx)
	a.shardManager = nil
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
