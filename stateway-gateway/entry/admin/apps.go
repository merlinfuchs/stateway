package admin

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"strconv"
	"time"

	"github.com/disgoorg/disgo/rest"
	"github.com/disgoorg/snowflake/v2"
	"github.com/merlinfuchs/stateway/stateway-gateway/db/postgres"
	"github.com/merlinfuchs/stateway/stateway-gateway/model"
	"github.com/merlinfuchs/stateway/stateway-gateway/store"
	"github.com/merlinfuchs/stateway/stateway-lib/config"
	"github.com/olekukonko/tablewriter"
	"gopkg.in/guregu/null.v4"
)

func ListApps(ctx context.Context, appStore store.AppStore, enabledOnly bool) error {
	var apps []*model.App
	var err error

	if enabledOnly {
		apps, err = appStore.GetEnabledApps(ctx, store.GetEnabledAppsParams{})
	} else {
		apps, err = appStore.GetApps(ctx)
	}
	if err != nil {
		return fmt.Errorf("failed to get apps: %w", err)
	}

	err = renderAppsTable(apps)
	if err != nil {
		return fmt.Errorf("failed to render apps table: %w", err)
	}

	return nil
}

func GetApp(ctx context.Context, appStore store.AppStore, id snowflake.ID) error {
	app, err := appStore.GetApp(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to get app: %w", err)
	}

	err = renderAppsTable([]*model.App{app})
	if err != nil {
		return fmt.Errorf("failed to render app table: %w", err)
	}

	return nil
}

func CreateApp(
	ctx context.Context,
	appStore store.AppStore,
	groupID string,
	token string,
	clientSecret string,
	config model.AppConfig,
) error {
	client := rest.New(rest.NewClient(token))

	discordApp, err := client.GetCurrentApplication(rest.WithCtx(ctx))
	if err != nil {
		return fmt.Errorf("failed to get current app: %w", err)
	}

	discordGateway, err := client.GetGatewayBot(rest.WithCtx(ctx))
	if err != nil {
		return fmt.Errorf("failed to get gateway bot: %w", err)
	}

	if groupID == "" {
		groupID = "default"
	}

	if !config.ShardConcurrency.Valid {
		config.ShardConcurrency = null.IntFrom(int64(discordGateway.SessionStartLimit.MaxConcurrency))
	}

	app, err := appStore.CreateApp(ctx, store.CreateAppParams{
		ID:                  discordApp.ID,
		GroupID:             groupID,
		DisplayName:         discordApp.Name,
		DiscordClientID:     discordApp.ID,
		DiscordBotToken:     token,
		DiscordPublicKey:    discordApp.VerifyKey,
		DiscordClientSecret: null.NewString(clientSecret, clientSecret != ""),
		ShardCount:          int(discordGateway.Shards),
		Config:              config,
		CreatedAt:           time.Now().UTC(),
		UpdatedAt:           time.Now().UTC(),
	})
	if err != nil {
		return fmt.Errorf("failed to create app: %w", err)
	}

	err = renderAppsTable([]*model.App{app})
	if err != nil {
		return fmt.Errorf("failed to render app table: %w", err)
	}

	return nil
}

func DeleteApp(ctx context.Context, appStore store.AppStore, id snowflake.ID) error {
	err := appStore.DeleteApp(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to delete app: %w", err)
	}
	return nil
}

func DisableApp(ctx context.Context, appStore store.AppStore, id snowflake.ID, code model.AppDisabledCode, message string) error {
	app, err := appStore.DisableApp(ctx, store.DisableAppParams{
		ID:              id,
		DisabledCode:    code,
		DisabledMessage: null.NewString(message, message != ""),
	})
	if err != nil {
		return fmt.Errorf("failed to disable app: %w", err)
	}
	err = renderAppsTable([]*model.App{app})
	if err != nil {
		return fmt.Errorf("failed to render app table: %w", err)
	}
	return nil
}

func InitializeApps(ctx context.Context, pg *postgres.Client, cfg *config.RootGatewayConfig) error {
	slog.Info("Initializing apps from config", slog.Int("app_count", len(cfg.Gateway.Apps)))

	for _, appCfg := range cfg.Gateway.Apps {
		client := rest.New(rest.NewClient(appCfg.Token))

		discordApp, err := client.GetCurrentApplication(rest.WithCtx(ctx))
		if err != nil {
			return fmt.Errorf("failed to get current app: %w", err)
		}

		config := model.AppConfig{
			Intents:          null.NewInt(appCfg.Intents, appCfg.Intents != 0),
			ShardConcurrency: null.NewInt(int64(appCfg.ShardConcurrency), appCfg.ShardConcurrency != 0),
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

		// TODO: Only update when it actually changed?
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
			"Initialized app",
			slog.String("app_id", discordApp.ID.String()),
			slog.String("group_id", appCfg.GroupID),
			slog.String("display_name", discordApp.Name),
			slog.Int("shard_count", appCfg.ShardCount),
			slog.Int64("intents", appCfg.Intents),
		)
	}
	return nil
}

func renderAppsTable(apps []*model.App) error {
	table := tablewriter.NewWriter(os.Stdout)
	table.Header([]string{"ID", "Group", "Name", "Shard Count", "Disabled", "Created At", "Updated At"})
	for _, app := range apps {
		err := table.Append([]string{
			app.ID.String(),
			app.GroupID,
			app.DisplayName,
			strconv.Itoa(app.ShardCount),
			strconv.FormatBool(app.Disabled),
			app.CreatedAt.Format(time.RFC3339),
			app.UpdatedAt.Format(time.RFC3339),
		})
		if err != nil {
			return fmt.Errorf("failed to append app to table: %w", err)
		}
	}
	return table.Render()
}
