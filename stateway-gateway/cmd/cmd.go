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
					Name:  "gateway-count",
					Usage: "The number of gateways to run and balance the apps across.",
				},
				&cli.IntFlag{
					Name:  "gateway-id",
					Usage: "The ID of the gateway to run (0-based index).",
				},
				&cli.BoolFlag{
					Name:  "no-resume",
					Usage: "Disable resuming of shard sessions.",
				},
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

				if c.IsSet("gateway-count") {
					env.cfg.Gateway.GatewayCount = c.Int("gateway-count")
				}

				if c.IsSet("gateway-id") {
					env.cfg.Gateway.GatewayID = c.Int("gateway-id")
				}

				if c.IsSet("no-resume") {
					env.cfg.Gateway.NoResume = c.Bool("no-resume")
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

func setupEnv(ctx context.Context, debug bool) (*env, error) {
	cfg, err := config.LoadConfig[*config.RootGatewayConfig]()
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
