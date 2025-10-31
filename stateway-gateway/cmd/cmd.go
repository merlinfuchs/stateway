package cmd

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/merlinfuchs/stateway/stateway-gateway/config"
	"github.com/merlinfuchs/stateway/stateway-gateway/db/postgres"
	"github.com/merlinfuchs/stateway/stateway-gateway/entry/gateway"
	"github.com/merlinfuchs/stateway/stateway-gateway/logging"
	"github.com/urfave/cli/v2"
)

var CLI = cli.App{
	Name:        "stateway-gateway",
	Description: "Stateway Gateway CLI",
	Commands: []*cli.Command{
		{
			Name:  "start",
			Usage: "Start the Stateway Gateway.",
			Action: func(c *cli.Context) error {
				ctx, cancel := signal.NotifyContext(c.Context, syscall.SIGINT, syscall.SIGTERM)
				defer cancel()

				env, err := setupEnv(ctx)
				if err != nil {
					return fmt.Errorf("failed to setup environment: %w", err)
				}

				err = gateway.Run(ctx, env.pg, env.cfg)
				if err != nil {
					return fmt.Errorf("failed to run gateway: %w", err)
				}
				return nil
			},
		},
		&adminCMD,
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
	cfg *config.Config
}

func setupEnv(ctx context.Context) (*env, error) {
	cfg, err := config.LoadConfig(config.ConfigFile)
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
