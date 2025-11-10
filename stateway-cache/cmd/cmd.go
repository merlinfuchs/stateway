package cmd

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/merlinfuchs/stateway/stateway-cache/db/postgres"
	"github.com/merlinfuchs/stateway/stateway-cache/entry/server"
	"github.com/merlinfuchs/stateway/stateway-lib/config"
	"github.com/merlinfuchs/stateway/stateway-lib/logging"
	"github.com/urfave/cli/v2"
)

var CLI = cli.App{
	Name:        "stateway-gateway",
	Description: "Stateway Gateway CLI",
	Commands: []*cli.Command{
		{
			Name:  "server",
			Usage: "Start the Stateway Cache Server.",
			Flags: []cli.Flag{
				&cli.IntSliceFlag{
					Name:  "gateway-ids",
					Usage: "The gateway IDs to process events from. Leave empty to process events from all gateways.",
				},
			},
			Action: func(c *cli.Context) error {
				ctx, cancel := signal.NotifyContext(c.Context, syscall.SIGINT, syscall.SIGTERM)
				defer cancel()

				env, err := setupEnv(ctx)
				if err != nil {
					return fmt.Errorf("failed to setup environment: %w", err)
				}

				if c.IsSet("gateway-ids") {
					env.cfg.Cache.GatewayIDs = c.IntSlice("gateway-ids")
				}

				err = server.Run(ctx, env.pg, env.cfg)
				if err != nil {
					return fmt.Errorf("failed to run cache server: %w", err)
				}
				return nil
			},
		},
		&databaseCMD,
	},
}

func Execute() {
	if err := CLI.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

type env struct {
	pg  *postgres.Client
	cfg *config.RootCacheConfig
}

func setupEnv(ctx context.Context) (*env, error) {
	cfg, err := config.LoadConfig[*config.RootCacheConfig]()
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	logging.SetupLogger(logging.LoggerConfig(cfg.Logging))

	pg, err := postgres.New(ctx, postgres.ClientConfig(cfg.Database.Postgres))
	if err != nil {
		return nil, fmt.Errorf("failed to create postgres client: %w", err)
	}

	return &env{
		pg:  pg,
		cfg: cfg,
	}, nil
}
