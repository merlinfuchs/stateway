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

func GetApp(ctx context.Context, appStore store.AppStore, groupID string, discordClientID snowflake.ID) error {
	if groupID == "" {
		groupID = "default"
	}

	app, err := appStore.GetApp(ctx, groupID, discordClientID)
	if err != nil {
		return fmt.Errorf("failed to get app: %w", err)
	}

	err = renderAppsTable([]*model.App{app})
	if err != nil {
		return fmt.Errorf("failed to render app table: %w", err)
	}

	return nil
}

func CreateApp(ctx context.Context, appStore store.AppStore, groupID string, token string, clientSecret string) error {
	client := rest.New(rest.NewClient(token))

	discordApp, err := client.GetCurrentApplication(rest.WithCtx(ctx))
	if err != nil {
		return fmt.Errorf("failed to get current app: %w", err)
	}

	if groupID == "" {
		groupID = "default"
	}

	app, err := appStore.CreateApp(ctx, store.CreateAppParams{
		ID:                  snowflake.New(time.Now().UTC()),
		GroupID:             groupID,
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

func DeleteApp(ctx context.Context, appStore store.AppStore, groupID string, discordClientID snowflake.ID) error {
	if groupID == "" {
		groupID = "default"
	}

	err := appStore.DeleteApp(ctx, groupID, discordClientID)
	if err != nil {
		return fmt.Errorf("failed to delete app: %w", err)
	}
	return nil
}

func DisableApp(ctx context.Context, appStore store.AppStore, groupID string, discordClientID snowflake.ID, code model.AppDisabledCode, message string) error {
	if groupID == "" {
		groupID = "default"
	}

	app, err := appStore.DisableApp(ctx, store.DisableAppParams{
		GroupID:         groupID,
		DiscordClientID: discordClientID,
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
	table.Header([]string{"Group", "Name", "Disabled", "Disabled Code", "Disabled Message", "Created At", "Updated At"})
	for _, app := range apps {
		err := table.Append([]string{
			app.GroupID,
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
