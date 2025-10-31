package admin

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/disgoorg/disgo/rest"
	"github.com/disgoorg/snowflake/v2"
	"github.com/merlinfuchs/stateway/stateway-gateway/model"
	"github.com/merlinfuchs/stateway/stateway-gateway/store"
	"github.com/olekukonko/tablewriter"
	"gopkg.in/guregu/null.v4"
)

func ListApps(ctx context.Context, appStore store.AppStore, enabledOnly bool) error {
	var apps []*model.App
	var err error

	if enabledOnly {
		apps, err = appStore.GetEnabledApps(ctx)
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

func CreateApp(ctx context.Context, appStore store.AppStore, token string, clientSecret string) error {
	client := rest.New(rest.NewClient(token))

	discordApp, err := client.GetCurrentApplication(rest.WithCtx(ctx))
	if err != nil {
		return fmt.Errorf("failed to get current app: %w", err)
	}

	app, err := appStore.CreateApp(ctx, store.CreateAppParams{
		ID:                  snowflake.ID(discordApp.ID),
		DisplayName:         discordApp.Name,
		DiscordClientID:     snowflake.ID(discordApp.ID),
		DiscordBotToken:     token,
		DiscordPublicKey:    discordApp.VerifyKey,
		DiscordClientSecret: null.NewString(clientSecret, clientSecret != ""),
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

func renderAppsTable(apps []*model.App) error {
	table := tablewriter.NewWriter(os.Stdout)
	table.Header([]string{"ID", "Display Name", "Disabled", "Disabled Code", "Disabled Message", "Created At", "Updated At"})
	for _, app := range apps {
		err := table.Append([]string{
			app.ID.String(),
			app.DisplayName,
			app.DiscordClientID.String(),
			strconv.FormatBool(app.Disabled),
			string(app.DisabledCode),
			app.DisabledMessage.String,
			app.CreatedAt.Format(time.RFC3339),
			app.UpdatedAt.Format(time.RFC3339),
		})
		if err != nil {
			return fmt.Errorf("failed to append app to table: %w", err)
		}
	}
	return table.Render()
}
