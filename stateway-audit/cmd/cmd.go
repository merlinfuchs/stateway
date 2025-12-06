package cmd

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/merlinfuchs/stateway/stateway-audit/db/postgres"
	"github.com/merlinfuchs/stateway/stateway-audit/entry/server"
	"github.com/merlinfuchs/stateway/stateway-lib/config"
	"github.com/merlinfuchs/stateway/stateway-lib/logging"
	"github.com/urfave/cli/v2"
)

var CLI = cli.App{
	Name:        "stateway-audit",
	Description: "Stateway Audit CLI",
	Commands: []*cli.Command{
		{
			Name:  "server",
			Usage: "Start the Stateway Audit Server.",
			Flags: []cli.Flag{
				&cli.BoolFlag{
					Name:  "debug",
					Usage: "Enable debug logging.",
				},
			},
			Action: func(c *cli.Context) error {
				ctx, cancel := signal.NotifyContext(c.Context, syscall.SIGINT, syscall.SIGTERM)
				defer cancel()

				env, err := setupEnv(ctx, c.Bool("debug"))
				if err != nil {
					return fmt.Errorf("failed to setup environment: %w", err)
				}

				err = server.Run(ctx, env.pg, env.cfg)
				if err != nil {
					return fmt.Errorf("failed to run cache server: %w", err)
				}
				return nil
			},
		},
	},
}

func Execute() {
	if err := CLI.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

type env struct {
	pg  *postgres.Client
	cfg *config.RootAuditConfig
}

func setupEnv(ctx context.Context, debug bool) (*env, error) {
	cfg, err := config.LoadConfig[*config.RootAuditConfig]()
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	loggingConfig := logging.LoggerConfig(cfg.Logging)
	if debug {
		loggingConfig.Debug = true
	}

	logging.SetupLogger(loggingConfig)

	pg, err := postgres.New(ctx, postgres.ClientConfig(cfg.Database.Postgres))
	if err != nil {
		return nil, fmt.Errorf("failed to create postgres client: %w", err)
	}

	return &env{
		pg:  pg,
		cfg: cfg,
	}, nil
}
