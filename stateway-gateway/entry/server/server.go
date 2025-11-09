package server

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/disgoorg/disgo/rest"
	"github.com/merlinfuchs/stateway/stateway-gateway/app"
	"github.com/merlinfuchs/stateway/stateway-gateway/db/postgres"
	"github.com/merlinfuchs/stateway/stateway-gateway/model"
	"github.com/merlinfuchs/stateway/stateway-gateway/store"
	"github.com/merlinfuchs/stateway/stateway-lib/broker"
	"github.com/merlinfuchs/stateway/stateway-lib/config"
	"gopkg.in/guregu/null.v4"
)

func Run(ctx context.Context, pg *postgres.Client, cfg *config.RootGatewayConfig) error {
	err := createInitialApps(ctx, pg, cfg)
	if err != nil {
		return fmt.Errorf("failed to create initial apps: %w", err)
	}

	broker, err := broker.NewNATSBroker(cfg.Broker.NATS.URL)
	if err != nil {
		return fmt.Errorf("failed to create NATS broker: %w", err)
	}

	err = broker.CreateGatewayStream(ctx)
	if err != nil {
		return fmt.Errorf("failed to create gateway stream: %w", err)
	}

	eventHandler := NewEventHandler(broker)
	go eventHandler.Run(ctx)

	appManager := app.NewAppManager(
		app.AppManagerConfig{
			GatewayCount: cfg.Gateway.GatewayCount,
			GatewayID:    cfg.Gateway.GatewayID,
		},
		pg,
		eventHandler,
	)

	appManager.Run(ctx)
	return nil
}

func createInitialApps(ctx context.Context, pg *postgres.Client, cfg *config.RootGatewayConfig) error {
	for _, appCfg := range cfg.Gateway.Apps {
		client := rest.New(rest.NewClient(appCfg.Token))

		discordApp, err := client.GetCurrentApplication(rest.WithCtx(ctx))
		if err != nil {
			return fmt.Errorf("failed to get current app: %w", err)
		}

		config := model.AppConfig{
			Intents: null.NewInt(appCfg.Intents, appCfg.Intents != 0),
		}
		if appCfg.Presence != nil {
			config.Presence = &model.AppPresenceConfig{
				Status: null.NewString(appCfg.Presence.Status, appCfg.Presence.Status != ""),
			}
			if appCfg.Presence.Activity != nil {
				config.Presence.Activity = &model.AppPresenceActivityConfig{
					Name:  appCfg.Presence.Activity.Name,
					State: appCfg.Presence.Activity.State,
					Type:  appCfg.Presence.Activity.Type,
					URL:   appCfg.Presence.Activity.URL,
				}
			}
		}

		err = pg.UpsertApp(ctx, store.UpsertAppParams{
			ID:               discordApp.ID,
			GroupID:          appCfg.GroupID,
			DisplayName:      discordApp.Name,
			DiscordClientID:  discordApp.ID,
			DiscordBotToken:  appCfg.Token,
			DiscordPublicKey: discordApp.VerifyKey,
			ShardCount:       appCfg.ShardCount,
			Constraints:      model.AppConstraints{},
			Config:           config,
			CreatedAt:        time.Now().UTC(),
			UpdatedAt:        time.Now().UTC(),
		})
		if err != nil {
			return fmt.Errorf("failed to upsert app: %w", err)
		}

		slog.Info(
			"Created initial app",
			slog.String("app_id", discordApp.ID.String()),
			slog.String("group_id", appCfg.GroupID),
			slog.String("display_name", discordApp.Name),
			slog.Int("shard_count", appCfg.ShardCount),
			slog.Int64("intents", appCfg.Intents),
		)
	}
	return nil
}
