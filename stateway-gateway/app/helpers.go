package app

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/gateway"
	"github.com/disgoorg/disgo/sharding"
	"github.com/gorilla/websocket"
	"github.com/merlinfuchs/stateway/stateway-gateway/model"
	"github.com/merlinfuchs/stateway/stateway-gateway/store"
	"gopkg.in/guregu/null.v4"
)

const resumeTimeout = time.Minute

func (a *App) shardsFromApp(ctx context.Context, gatewayCount int, gatewayID int) (int, int, map[int]sharding.ShardState, error) {
	shardCount := a.model.ShardCount
	if shardCount == 0 {
		shardCount = 1
	}

	shardConcurrency := 1
	if a.model.Config.ShardConcurrency.Valid {
		shardConcurrency = int(a.model.Config.ShardConcurrency.Int64)
	}

	shards := make(map[int]sharding.ShardState, shardCount)
	for shardID := 0; shardID < shardCount; shardID++ {
		// if shardCount == 1, this gateway is the only one that should run the app
		// otherwise, we are splitting the app shards across the gateways
		// shardID % gatewayCount gives us the index of the gateway that should run the shard
		if shardCount == 1 || shardID%gatewayCount == gatewayID {
			shardSession, err := a.shardSessionStore.GetLastShardSession(ctx, a.model.ID, shardID, shardCount)
			if err != nil && !errors.Is(err, store.ErrNotFound) {
				return 0, 0, nil, fmt.Errorf("failed to get last shard session: %w", err)
			}

			var state sharding.ShardState
			// We only try to resume the shard if it hasn't been invalidated and it's been updated in the last resumeTimeout
			if shardSession != nil && !shardSession.InvalidatedAt.Valid && shardSession.UpdatedAt.After(time.Now().UTC().Add(-resumeTimeout)) {
				state = sharding.ShardState{
					SessionID: shardSession.ID,
					ResumeURL: shardSession.ResumeURL,
					Sequence:  shardSession.LastSequence,
				}
			}

			shards[shardID] = state
		}
	}

	return shardCount, shardConcurrency, shards, nil
}

func intentsFromApp(app *model.App) gateway.Intents {
	intents := gateway.IntentsNonPrivileged
	if app.Config.Intents.Valid {
		intents = gateway.Intents(app.Config.Intents.Int64)
	}
	return intents
}

func presenceOptsFromApp(app *model.App) []gateway.PresenceOpt {
	res := []gateway.PresenceOpt{}
	if app.Config.Presence != nil {
		presence := app.Config.Presence
		if presence.Status.Valid {
			res = append(res, gateway.WithOnlineStatus(discord.OnlineStatus(presence.Status.String)))
		}

		if presence.Activity != nil {
			activityOpts := []gateway.ActivityOpt{}

			switch presence.Activity.Type {
			case "watching":
				res = append(
					res,
					gateway.WithWatchingActivity(presence.Activity.Name, activityOpts...),
				)
			case "listening":
				res = append(
					res,
					gateway.WithListeningActivity(presence.Activity.Name, activityOpts...),
				)
			case "competing":
				res = append(
					res,
					gateway.WithCompetingActivity(presence.Activity.Name, activityOpts...),
				)
			case "streaming":
				res = append(
					res,
					gateway.WithStreamingActivity(presence.Activity.Name, presence.Activity.URL, activityOpts...),
				)
			case "playing":
				res = append(
					res,
					gateway.WithPlayingActivity(presence.Activity.Name, activityOpts...),
				)
			default:
				res = append(
					res,
					gateway.WithCustomActivity(presence.Activity.Name, activityOpts...),
				)
			}
		}
	}

	return res
}

func (a *App) disableIfFatal(ctx context.Context, err error) {
	var wsError *websocket.CloseError
	if errors.As(err, &wsError) {
		switch wsError.Code {
		case 4004:
			a.disable(ctx, model.AppDisabledCodeInvalidToken, wsError.Text)
		case 4013:
			a.disable(ctx, model.AppDisabledCodeInvalidIntents, wsError.Text)
		case 4014:
			a.disable(ctx, model.AppDisabledCodeDisallowedIntents, wsError.Text)
		}
	}
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

func (a *App) storeSession(ctx context.Context, g gateway.Gateway) {
	sessionID := g.SessionID()
	resumeURL := g.ResumeURL()
	sequenceNumber := g.LastSequenceReceived()

	if sessionID == nil || resumeURL == nil || sequenceNumber == nil {
		return
	}

	err := a.shardSessionStore.UpsertShardSession(ctx, store.UpsertShardSessionParams{
		ID:           *sessionID,
		AppID:        a.model.ID,
		ShardID:      g.ShardID(),
		ShardCount:   g.ShardCount(),
		LastSequence: *sequenceNumber,
		ResumeURL:    *resumeURL,
		CreatedAt:    time.Now().UTC(),
		UpdatedAt:    time.Now().UTC(),
	})
	if err != nil {
		slog.Error(
			"Failed to upsert shard session",
			slog.String("app_id", a.model.ID.String()),
			slog.String("group_id", a.model.GroupID),
			slog.Int("shard_id", g.ShardID()),
			slog.String("display_name", a.model.DisplayName),
			slog.Any("error", err),
		)
	}
}

func (a *App) invalidateSession(ctx context.Context, g gateway.Gateway) {
	err := a.shardSessionStore.InvalidateShardSession(ctx, a.model.ID, g.ShardID(), g.ShardCount())
	if err != nil {
		slog.Error(
			"Failed to invalidate shard session",
			slog.String("app_id", a.model.ID.String()),
			slog.String("group_id", a.model.GroupID),
			slog.Int("shard_id", g.ShardID()),
			slog.String("display_name", a.model.DisplayName),
			slog.Any("error", err),
		)
	}
}
