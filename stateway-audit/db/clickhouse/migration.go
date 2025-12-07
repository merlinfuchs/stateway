package clickhouse

import (
	"embed"
	"fmt"
	"net"
	"strconv"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/clickhouse" // For NewMigrator().
	"github.com/golang-migrate/migrate/v4/source/iofs"
)

//go:embed migrations/*.sql
var migrationsFS embed.FS

// Migrator is a wrapper around a `migrate.Migrate` struct, with the migrations embedded in the binary.
type Migrator struct {
	m     *migrate.Migrate
	close func() error
}

// Up migrates up to the latest version.
func (mig *Migrator) Up() error {
	err := mig.m.Up()
	if err != nil {
		return fmt.Errorf("clickhousedb failed to migrate up: %w", err)
	}
	return nil
}

// Down migrates down to the lowest version.
func (mig *Migrator) Down() error {
	err := mig.m.Down()
	if err != nil {
		return fmt.Errorf("clickhousedb failed to migrate down: %w", err)
	}
	return nil
}

// Version returns the current migration version.
func (mig *Migrator) Version() (uint, bool, error) {
	v, dirty, err := mig.m.Version()
	if err != nil {
		return 0, false, fmt.Errorf("clickhousedb failed to get migration version: %w", err)
	}
	return v, dirty, nil
}

// To migrates to the given version.
func (mig *Migrator) To(version uint) error {
	err := mig.m.Migrate(version)
	if err != nil {
		return fmt.Errorf("clickhousedb failed to migrate to version %d: %w", version, err)
	}
	return nil
}

// Force forces the migration to the given version.
func (mig *Migrator) Force(version int) error {
	err := mig.m.Force(version)
	if err != nil {
		return fmt.Errorf("clickhousedb failed to force migrate to version %d: %w", version, err)
	}
	return nil
}

// List returns a list of all migration filenames.
func (mig *Migrator) List() ([]string, error) {
	dirEntries, err := migrationsFS.ReadDir("migrations")
	if err != nil {
		return nil, fmt.Errorf("clickhousedb failed to read migrations dir: %w", err)
	}

	migrationFiles := make([]string, 0)
	for _, entry := range dirEntries {
		migrationFiles = append(migrationFiles, entry.Name())
	}
	return migrationFiles, nil
}

// Close closes the migrator.
func (mig *Migrator) Close() error {
	err := mig.close()
	if err != nil {
		return fmt.Errorf("clickhousedb failed to close migrator: %w", err)
	}
	return nil
}

// SetLogger sets the logger for the wrapped `*migrate.Migrate` library.
func (mig *Migrator) SetLogger(logger migrate.Logger) {
	mig.m.Log = logger
}

// NewMigrator returns a new Migrator, which can be used to run migrations.
func NewMigrator(host string, port int, username, password, database string) (*Migrator, error) {
	d, err := iofs.New(migrationsFS, "migrations")
	if err != nil {
		return nil, fmt.Errorf("failed to open Clickhouse migrations iofs: %w", err)
	}

	params := fmt.Sprintf("username=%s&password=%s&x-multi-statement=true&x-migrations-table-engine=MergeTree",
		username, password)
	dsn := fmt.Sprintf("clickhouse://%s/%s?%s",
		net.JoinHostPort(host, strconv.Itoa(port)), database, params)

	m, err := migrate.NewWithSourceInstance("iofs", d, dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to create clickhouse migrate instance: %w", err)
	}

	closeFunc := func() error {
		err1, err2 := m.Close()
		if err1 != nil || err2 != nil {
			return fmt.Errorf("source close error: %w, driver close error: %w", err1, err2)
		}
		return nil
	}

	return &Migrator{
		m:     m,
		close: closeFunc,
	}, nil
}
