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
)

type AppConfig struct {
	GatewayCount int
	GatewayID    int
}

type App struct {
	cfg   AppConfig
	model *model.App

	appStore          store.AppStore
	shardSessionStore store.ShardSessionStore
	eventHandler      event.EventHandler

	shardManager sharding.ShardManager
}

func NewApp(
	cfg AppConfig,
	model *model.App,
	appStore store.AppStore,
	shardSessionStore store.ShardSessionStore,
	eventHandler event.EventHandler,
) *App {
	return &App{
		cfg:               cfg,
		model:             model,
		appStore:          appStore,
		shardSessionStore: shardSessionStore,
		eventHandler:      eventHandler,
	}
}

func (a *App) Run(ctx context.Context) {
	intents := intentsFromApp(a.model)
	presenceOpts := presenceOptsFromApp(a.model)

	shardCount, shards, err := a.shardsFromApp(ctx, a.cfg.GatewayCount, a.cfg.GatewayID)
	if err != nil {
		slog.Error("Failed to get shards", slog.Any("error", err))
		return
	}

	shardManager := sharding.New(
		a.model.DiscordBotToken,
		func(g gateway.Gateway, eventType gateway.EventType, sequenceNumber int, ev gateway.EventData) {
			a.handleEvent(ctx, g, eventType, sequenceNumber, ev)
		},
		sharding.WithAutoScaling(false),
		sharding.WithShardCount(shardCount),
		sharding.WithShardIDsWithStates(shards),
		sharding.WithLogger(slog.Default()),
		sharding.WithGatewayConfigOpts(
			gateway.WithIntents(intents),
			gateway.WithEnableRawEvents(true),
			gateway.WithPresenceOpts(presenceOpts...),
			gateway.WithLogger(slog.Default()),
			gateway.WithAutoReconnect(true),
		),
		sharding.WithGatewayCreateFunc(func(
			token string,
			eventHandlerFunc gateway.EventHandlerFunc,
			closeHandlerFunc gateway.CloseHandlerFunc,
			opts ...gateway.ConfigOpt,
		) gateway.Gateway {
			return a.createGateway(ctx, token, eventHandlerFunc, closeHandlerFunc, opts...)
		}),
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

func (a *App) createGateway(
	ctx context.Context,
	token string,
	eventHandlerFunc gateway.EventHandlerFunc,
	originalCloseHandlerFunc gateway.CloseHandlerFunc,
	opts ...gateway.ConfigOpt,
) gateway.Gateway {
	closeHandlerFunc := func(g gateway.Gateway, err error, reconnect bool) {
		a.handleClose(ctx, g, originalCloseHandlerFunc, err, reconnect)
	}

	return gateway.New(token, eventHandlerFunc, closeHandlerFunc, opts...)
}

func (a *App) handleClose(ctx context.Context, g gateway.Gateway, originalCloseHandlerFunc gateway.CloseHandlerFunc, err error, reconnect bool) {
	slog.Info(
		"Discord shard CLOSED",
		slog.String("group_id", a.model.GroupID),
		slog.String("app_id", a.model.ID.String()),
		slog.Int("shard_id", g.ShardID()),
		slog.String("display_name", a.model.DisplayName),
		slog.Bool("reconnect", reconnect),
		slog.Any("error", err),
	)

	a.disableIfFatal(ctx, err)
	a.invalidateSession(ctx, g)

	originalCloseHandlerFunc(g, err, reconnect)
}

func (a *App) handleEvent(ctx context.Context, g gateway.Gateway, _ gateway.EventType, _ int, ev gateway.EventData) {
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

		go a.storeSession(ctx, g)
	case gateway.EventResumed:
		slog.Info(
			"Discord shard RESUMED",
			slog.String("group_id", a.model.GroupID),
			slog.String("app_id", a.model.ID.String()),
			slog.Int("shard_id", g.ShardID()),
			slog.String("display_name", a.model.DisplayName),
		)

		go a.storeSession(ctx, g)
	case gateway.EventRateLimited:
		slog.Info(
			"Discord shard RATE_LIMITED",
			slog.String("group_id", a.model.GroupID),
			slog.String("app_id", a.model.ID.String()),
			slog.Int("shard_id", g.ShardID()),
			slog.String("display_name", a.model.DisplayName),
		)
	case gateway.EventHeartbeatAck:
		go a.storeSession(ctx, g)
	}
}
