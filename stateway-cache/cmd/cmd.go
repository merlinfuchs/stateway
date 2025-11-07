package cmd

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

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
			Action: func(c *cli.Context) error {
				ctx, cancel := signal.NotifyContext(c.Context, syscall.SIGINT, syscall.SIGTERM)
				defer cancel()

				env, err := setupEnv(ctx)
				if err != nil {
					return fmt.Errorf("failed to setup environment: %w", err)
				}

				err = server.Run(ctx, env.cfg)
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
	cfg *config.CacheConfig
}

func setupEnv(ctx context.Context) (*env, error) {
	cfg, err := config.LoadConfig[*config.CacheConfig]()
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	logging.SetupLogger(logging.LoggerConfig(cfg.Logging))

	return &env{
		cfg: cfg,
	}, nil
}
