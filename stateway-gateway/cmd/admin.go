package cmd

import (
	"fmt"
	"os/signal"
	"syscall"

	"github.com/disgoorg/snowflake/v2"
	"github.com/merlinfuchs/stateway/stateway-gateway/entry/admin"
	"github.com/merlinfuchs/stateway/stateway-gateway/model"
	"github.com/urfave/cli/v2"
	"gopkg.in/guregu/null.v4"
)

var adminCMD = cli.Command{
	Name:  "admin",
	Usage: "Manage admin tasks.",
	Subcommands: []*cli.Command{
		{
			Name:  "apps",
			Usage: "Manage apps.",
			Subcommands: []*cli.Command{
				{
					Name:  "list",
					Usage: "List all apps.",
					Flags: []cli.Flag{
						&cli.BoolFlag{
							Name:  "enabled",
							Usage: "Only list enabled apps.",
						},
					},
					Action: func(c *cli.Context) error {
						ctx, cancel := signal.NotifyContext(c.Context, syscall.SIGINT, syscall.SIGTERM)
						defer cancel()

						env, err := setupEnv(ctx)
						if err != nil {
							return fmt.Errorf("failed to setup environment: %w", err)
						}

						err = admin.ListApps(ctx, env.pg, c.Bool("enabled"))
						if err != nil {
							return fmt.Errorf("failed to list apps: %w", err)
						}
						return nil
					},
				},
				{
					Name:  "get",
					Usage: "Get an app.",
					Flags: []cli.Flag{
						&cli.StringFlag{
							Name:     "group",
							Usage:    "The group of the app to get.",
							Required: false,
						},
						&cli.StringFlag{
							Name:     "id",
							Usage:    "The ID of the app to get.",
							Required: true,
						},
					},
					Action: func(c *cli.Context) error {
						ctx, cancel := signal.NotifyContext(c.Context, syscall.SIGINT, syscall.SIGTERM)
						defer cancel()

						env, err := setupEnv(ctx)
						if err != nil {
							return fmt.Errorf("failed to setup environment: %w", err)
						}

						appID, err := snowflake.Parse(c.String("id"))
						if err != nil {
							return fmt.Errorf("failed to parse app ID: %w", err)
						}

						err = admin.GetApp(ctx, env.pg, appID)
						if err != nil {
							return fmt.Errorf("failed to get app: %w", err)
						}
						return nil
					},
				},
				{
					Name:  "create",
					Usage: "Create a new app.",
					Flags: []cli.Flag{
						&cli.StringFlag{
							Name:     "group",
							Usage:    "The group of the app to create.",
							Required: false,
						},
						&cli.StringFlag{
							Name:     "token",
							Usage:    "The token of the app to create.",
							Required: true,
						},
						&cli.StringFlag{
							Name:  "client-secret",
							Usage: "The client secret of the app to create.",
						},
						&cli.Int64Flag{
							Name:  "intents",
							Usage: "The intents of the app to create.",
						},
					},
					Action: func(c *cli.Context) error {
						ctx, cancel := signal.NotifyContext(c.Context, syscall.SIGINT, syscall.SIGTERM)
						defer cancel()

						env, err := setupEnv(ctx)
						if err != nil {
							return fmt.Errorf("failed to setup environment: %w", err)
						}

						var config model.AppConfig
						if c.IsSet("intents") {
							config.Intents = null.NewInt(c.Int64("intents"), true)
						}

						err = admin.CreateApp(
							ctx,
							env.pg,
							c.String("group"),
							c.String("token"),
							c.String("client-secret"),
							config,
						)
						if err != nil {
							return fmt.Errorf("failed to create app: %w", err)
						}
						return nil
					},
				},
				{
					Name:  "delete",
					Usage: "Delete an app.",
					Flags: []cli.Flag{
						&cli.StringFlag{
							Name:     "group",
							Usage:    "The group of the app to delete.",
							Required: false,
						},
						&cli.StringFlag{
							Name:     "id",
							Usage:    "The ID of the app to delete.",
							Required: true,
						},
						&cli.BoolFlag{
							Name:     "danger",
							Usage:    "Confirm that you want to delete the app.",
							Required: true,
						},
					},
					Action: func(c *cli.Context) error {
						ctx, cancel := signal.NotifyContext(c.Context, syscall.SIGINT, syscall.SIGTERM)
						defer cancel()

						env, err := setupEnv(ctx)
						if err != nil {
							return fmt.Errorf("failed to setup environment: %w", err)
						}

						appID, err := snowflake.Parse(c.String("id"))
						if err != nil {
							return fmt.Errorf("failed to parse app ID: %w", err)
						}

						err = admin.DeleteApp(ctx, env.pg, appID)
						if err != nil {
							return fmt.Errorf("failed to delete app: %w", err)
						}
						return nil
					},
				},
				{
					Name:  "disable",
					Usage: "Disable an app.",
					Flags: []cli.Flag{
						&cli.StringFlag{
							Name:     "group",
							Usage:    "The group of the app to disable.",
							Required: false,
						},
						&cli.StringFlag{
							Name:     "id",
							Usage:    "The Discord client ID of the app to disable.",
							Required: true,
						},
						&cli.StringFlag{
							Name:     "code",
							Usage:    "The code of the app to disable.",
							Required: true,
						},
						&cli.StringFlag{
							Name:     "message",
							Usage:    "The message of the app to disable.",
							Required: true,
						},
					},
					Action: func(c *cli.Context) error {
						ctx, cancel := signal.NotifyContext(c.Context, syscall.SIGINT, syscall.SIGTERM)
						defer cancel()

						env, err := setupEnv(ctx)
						if err != nil {
							return fmt.Errorf("failed to setup environment: %w", err)
						}

						appID, err := snowflake.Parse(c.String("id"))
						if err != nil {
							return fmt.Errorf("failed to parse app ID: %w", err)
						}

						err = admin.DisableApp(ctx, env.pg, appID, model.AppDisabledCode(c.String("code")), c.String("message"))
						if err != nil {
							return fmt.Errorf("failed to disable app: %w", err)
						}
						return nil
					},
				},
			},
		},
	},
}
