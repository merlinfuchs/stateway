package postgres

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/merlinfuchs/stateway/stateway-cache/db/postgres/pgmodel"
	"github.com/merlinfuchs/stateway/stateway-cache/store"
)

func (c *Client) UpsertRoles(ctx context.Context, roles ...store.UpsertRoleParams) error {
	if len(roles) == 0 {
		return nil
	}

	params := make([]pgmodel.UpsertRolesParams, len(roles))
	for i, role := range roles {
		params[i] = pgmodel.UpsertRolesParams{
			AppID:   int64(role.AppID),
			GuildID: int64(role.GuildID),
			RoleID:  int64(role.RoleID),
			Data:    role.Data,
			CreatedAt: pgtype.Timestamp{
				Time:  role.CreatedAt,
				Valid: true,
			},
			UpdatedAt: pgtype.Timestamp{
				Time:  role.UpdatedAt,
				Valid: true,
			},
		}
	}
	res := c.Q.UpsertRoles(ctx, params)
	return res.Close()
}

func (c *Client) DeleteRole(ctx context.Context, params store.RoleIdentifier) error {
	return c.Q.DeleteRole(ctx, pgmodel.DeleteRoleParams{
		AppID:   int64(params.AppID),
		GuildID: int64(params.GuildID),
		RoleID:  int64(params.RoleID),
	})
}
