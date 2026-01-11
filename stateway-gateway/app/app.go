package app

import (
	"context"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"time"

	disgateway "github.com/disgoorg/disgo/gateway"
	"github.com/disgoorg/disgo/rest"
	"github.com/disgoorg/disgo/sharding"
	"github.com/disgoorg/snowflake/v2"
	"github.com/merlinfuchs/stateway/stateway-gateway/model"
	"github.com/merlinfuchs/stateway/stateway-gateway/store"
	"github.com/merlinfuchs/stateway/stateway-lib/event"
	"github.com/merlinfuchs/stateway/stateway-lib/gateway"
)

type AppConfig struct {
	GatewayCount int
	GatewayID    int
	NoResume     bool
}

type App struct {
	cfg   AppConfig
	model *model.App
	group *model.Group

	appStore               store.AppStore
	shardSessionStore      store.ShardSessionStore
	identifyRateLimitStore store.IdentifyRateLimitStore
	eventHandler           event.EventHandler

	shardManager sharding.ShardManager
}

func NewApp(
	cfg AppConfig,
	model *model.App,
	group *model.Group,
	appStore store.AppStore,
	shardSessionStore store.ShardSessionStore,
	identifyRateLimitStore store.IdentifyRateLimitStore,
	eventHandler event.EventHandler,
) *App {
	return &App{
		cfg:                    cfg,
		model:                  model,
		group:                  group,
		appStore:               appStore,
		shardSessionStore:      shardSessionStore,
		identifyRateLimitStore: identifyRateLimitStore,
		eventHandler:           eventHandler,
	}
}

func (a *App) Run(ctx context.Context) {
	constraints := a.resolveConstraints()
	config := a.resolveConfig()

	intents := intentsFromConfig(config)
	presenceOpts := presenceOptsFromConfig(config)

	shardCount, shardConcurrency, shards, err := a.shardsFromApp(ctx, a.cfg.GatewayCount, a.cfg.GatewayID, a.cfg.NoResume)
	if err != nil {
		slog.Error("Failed to get shards", slog.Any("error", err))
		return
	}

	if constraints.MaxShards.Valid && int64(shardCount) > constraints.MaxShards.Int64 {
		a.disable(ctx, gateway.AppDisabledConstraintExceeded, "Max shards constraint exceeded")
		return
	}

	logger := slog.Default().With("group_id", a.model.GroupID, "app_name", a.model.DisplayName, "app_id", a.model.ID.String())

	// Check if the bot token is valid
	restClient := rest.New(rest.NewClient(a.model.DiscordBotToken))
	_, err = restClient.GetCurrentApplication(rest.WithCtx(ctx))
	if err != nil {
		var restErr *rest.Error
		if errors.As(err, &restErr) {
			if restErr.Response != nil && restErr.Response.StatusCode == http.StatusUnauthorized {
				a.disable(ctx, gateway.AppDisabledCodeInvalidToken, "Invalid authentication token")
				return
			}
		}
		slog.Error("Failed to get current application", slog.Any("error", err))
		return
	}

	shardManager := sharding.New(
		a.model.DiscordBotToken,
		func(g disgateway.Gateway, eventType disgateway.EventType, sequenceNumber int, ev disgateway.EventData) {
			a.handleEvent(ctx, g, eventType, sequenceNumber, ev)
		},
		sharding.WithAutoScaling(false),
		sharding.WithShardCount(shardCount),
		sharding.WithIdentifyRateLimiter(
			NewIdentifyRateLimiter(a.identifyRateLimitStore, a.model.ID, shardConcurrency),
		),
		sharding.WithShardIDsWithStates(shards),
		sharding.WithLogger(logger),
		sharding.WithGatewayConfigOpts(
			disgateway.WithIntents(intents),
			disgateway.WithEnableRawEvents(true),
			disgateway.WithPresenceOpts(presenceOpts...),
			disgateway.WithLogger(logger),
			disgateway.WithAutoReconnect(true),
		),
		sharding.WithCloseHandler(func(gateway disgateway.Gateway, err error, reconnect bool) {
			a.handleClose(ctx, gateway, err, reconnect)
		}),
	)

	a.shardManager = shardManager
	shardManager.Open(ctx)
}

func (a *App) Close(ctx context.Context) {
	if a.shardManager != nil {
		a.shardManager.Close(ctx)
		a.shardManager = nil
	}
}

func (a *App) Update(ctx context.Context, model *model.App, group *model.Group) {
	a.Close(ctx)
	// TODO: Re-use client if possible
	a.model = model
	a.group = group
	go a.Run(ctx)
}

func (a *App) handleClose(ctx context.Context, g disgateway.Gateway, err error, reconnect bool) {
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
}

func (a *App) handleEvent(ctx context.Context, g disgateway.Gateway, _ disgateway.EventType, _ int, ev disgateway.EventData) {
	switch e := ev.(type) {
	case disgateway.EventRaw:
		data, err := io.ReadAll(e.Payload)
		if err != nil {
			slog.Error("Failed to read event payload", slog.Any("error", err))
			return
		}

		a.eventHandler.HandleEvent(&event.GatewayEvent{
			ID:        snowflake.New(time.Now().UTC()),
			GatewayID: a.cfg.GatewayID,
			AppID:     a.model.ID,
			// TODO: Set GuildID if available
			GroupID: a.model.GroupID,
			ShardID: g.ShardID(),
			Type:    string(e.EventType),
			Data:    data,
		})
	case disgateway.EventReady:
		slog.Info(
			"Discord shard READY",
			slog.String("group_id", a.model.GroupID),
			slog.String("app_id", a.model.ID.String()),
			slog.Int("shard_id", e.Shard[0]),
			slog.Int("shard_count", e.Shard[1]),
			slog.String("display_name", a.model.DisplayName),
		)

		go a.storeSession(ctx, g)
	case disgateway.EventResumed:
		slog.Info(
			"Discord shard RESUMED",
			slog.String("group_id", a.model.GroupID),
			slog.String("app_id", a.model.ID.String()),
			slog.Int("shard_id", g.ShardID()),
			slog.String("display_name", a.model.DisplayName),
		)

		go a.storeSession(ctx, g)
	case disgateway.EventRateLimited:
		slog.Info(
			"Discord shard RATE_LIMITED",
			slog.String("group_id", a.model.GroupID),
			slog.String("app_id", a.model.ID.String()),
			slog.Int("shard_id", g.ShardID()),
			slog.String("display_name", a.model.DisplayName),
		)
	case disgateway.EventHeartbeatAck:
		go a.storeSession(ctx, g)
	}
}

func (a *App) resolveConstraints() gateway.AppConstraints {
	return a.group.DefaultConstraints.Merge(a.model.Constraints)
}

func (a *App) resolveConfig() gateway.AppConfig {
	return a.group.DefaultConfig.Merge(a.model.Config)
}
