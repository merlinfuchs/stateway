package cmd

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/merlinfuchs/stateway/stateway-gateway/db/postgres"
	"github.com/merlinfuchs/stateway/stateway-gateway/entry/server"
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
			Usage: "Start the Stateway Gateway Server.",
			Flags: []cli.Flag{
				&cli.IntFlag{
					Name:  "instance-count",
					Usage: "The number of instances to run the gateway server on.",
				},
				&cli.IntFlag{
					Name:  "instance-index",
					Usage: "The index of the instance to run the gateway server on.",
				},
			},
			Action: func(c *cli.Context) error {
				ctx, cancel := signal.NotifyContext(c.Context, syscall.SIGINT, syscall.SIGTERM)
				defer cancel()

				env, err := setupEnv(ctx)
				if err != nil {
					return fmt.Errorf("failed to setup environment: %w", err)
				}

				if c.IsSet("instance-count") {
					env.cfg.Gateway.InstanceCount = c.Int("instance-count")
				}

				if c.IsSet("instance-index") {
					env.cfg.Gateway.InstanceIndex = c.Int("instance-index")
				}

				err = server.Run(ctx, env.pg, env.cfg)
				if err != nil {
					return fmt.Errorf("failed to run gateway server: %w", err)
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
	cfg *config.RootGatewayConfig
}

func setupEnv(ctx context.Context) (*env, error) {
	cfg, err := config.LoadConfig[*config.RootGatewayConfig]()
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
