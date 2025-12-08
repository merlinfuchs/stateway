package postgres

import (
	"context"

	"github.com/disgoorg/snowflake/v2"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/merlinfuchs/stateway/stateway-audit/db/postgres/pgmodel"
	"github.com/merlinfuchs/stateway/stateway-audit/model"
)

func (c *Client) GetAuditConfig(ctx context.Context, appID snowflake.ID, guildID snowflake.ID) (*model.AuditConfig, error) {
	row, err := c.Q.GetAuditConfig(ctx, pgmodel.GetAuditConfigParams{
		AppID:   int64(appID),
		GuildID: int64(guildID),
	})
	if err != nil {
		return nil, err
	}
	return rowToAuditConfig(row)
}

func (c *Client) UpsertAuditConfig(ctx context.Context, auditConfig model.AuditConfig) error {
	return c.Q.UpsertAuditConfig(ctx, pgmodel.UpsertAuditConfigParams{
		AppID:     int64(auditConfig.AppID),
		GuildID:   int64(auditConfig.GuildID),
		Enabled:   auditConfig.Enabled,
		CreatedAt: pgtype.Timestamp{Time: auditConfig.CreatedAt, Valid: true},
		UpdatedAt: pgtype.Timestamp{Time: auditConfig.UpdatedAt, Valid: true},
	})
}

func rowToAuditConfig(row pgmodel.AuditConfig) (*model.AuditConfig, error) {
	return &model.AuditConfig{
		AppID:     snowflake.ID(row.AppID),
		GuildID:   snowflake.ID(row.GuildID),
		Enabled:   row.Enabled,
		CreatedAt: row.CreatedAt.Time,
		UpdatedAt: row.UpdatedAt.Time,
	}, nil
}
