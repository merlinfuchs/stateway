package server

import (
	"context"

	"github.com/merlinfuchs/stateway/stateway-audit/db/postgres"
	"github.com/merlinfuchs/stateway/stateway-lib/config"
)

func Run(ctx context.Context, pg *postgres.Client, cfg *config.RootAuditConfig) error {
	return nil
}
