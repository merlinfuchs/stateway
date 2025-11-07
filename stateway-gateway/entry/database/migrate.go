package database

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"github.com/golang-migrate/migrate/v4"
	"github.com/merlinfuchs/stateway/stateway-gateway/db/postgres"
	"github.com/merlinfuchs/stateway/stateway-gateway/store"
	"github.com/merlinfuchs/stateway/stateway-lib/config"
	"github.com/merlinfuchs/stateway/stateway-lib/logging"
)

func RunMigrations(ctx context.Context, db string, opts DatabaseMigrationOpts) error {
	cfg, err := config.LoadConfig[*config.GatewayConfig]()
	if err != nil {
		return fmt.Errorf("Failed to load server config: %v", err)
	}

	logging.SetupLogger(logging.LoggerConfig(cfg.Logging))

	l := slog.With("database", db).With("operation", opts.Operation)

	pg, err := postgres.New(ctx, postgres.ClientConfig(cfg.Database.Postgres))
	if err != nil {
		l.With("error", err).Error("Failed to create postgres client")
		os.Exit(1)
	}

	switch db {
	case "postgres":
		migrater, err := pg.GetMigrater()
		if err != nil {
			l.With("error", err).Error("Failed to get migrater")
			os.Exit(1)
		}
		defer migrater.Close()

		return runMigrationsAgainstMigrater(l, migrater, opts)
	default:
		return fmt.Errorf("invalid database: %s", db)
	}
}

func runMigrationsAgainstMigrater(l *slog.Logger, migrater store.Migrater, opts DatabaseMigrationOpts) error {
	migrater.SetLogger(databaseMigrationLogger{
		logger: l,
	})

	var err error
	switch opts.Operation {
	case "up":
		err = migrater.Up()
	case "down":
		err = migrater.Down()
	case "list":
		var migrations []string
		migrations, err = migrater.List()
		if err != nil {
			break
		}
		l.With("migrations", migrations).Info("")
	case "version":
		var version uint
		var dirty bool
		version, dirty, err = migrater.Version()
		if err != nil {
			break
		}
		l.With("version", version).With("dirty", dirty).Info("")

	case "force":
		l = l.With("target_version", opts.TargetVersion)
		err = migrater.Force(opts.TargetVersion)
		if err != nil {
			break
		}
	case "to":
		l = l.With("target_version", opts.TargetVersion)
		if opts.TargetVersion < 0 {
			l.With("error", err).Error("Invalid target version for migrate")
			return err
		}
		err = migrater.To(uint(opts.TargetVersion))
		if err != nil {
			break
		}
	}

	if err == migrate.ErrNoChange {
		l.Warn("Already at the correct version, migration was skipped")
	} else if err == migrate.ErrNilVersion {
		l.Warn("Migration is at nil version (no migrations have been performed)")
	} else if err != nil {
		l.With("error", err).Error("Migration operation failed")
		return err
	}

	l.Info("Migration completed")
	return nil
}

type DatabaseMigrationOpts struct {
	Operation     string
	TargetVersion int
}

type databaseMigrationLogger struct {
	logger  *slog.Logger
	verbose bool
}

// Printf is like fmt.Printf
func (ml databaseMigrationLogger) Printf(format string, v ...interface{}) {
	ml.logger.Info(fmt.Sprintf(format, v...))
}

// Verbose returns the verbose flag
func (ml databaseMigrationLogger) Verbose() bool {
	return ml.verbose
}
