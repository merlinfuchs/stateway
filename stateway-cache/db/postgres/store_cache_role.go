package postgres

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/snowflake/v2"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/merlinfuchs/stateway/stateway-cache/db/postgres/pgmodel"
	"github.com/merlinfuchs/stateway/stateway-cache/model"
	"github.com/merlinfuchs/stateway/stateway-cache/store"
)

func (c *Client) GetRole(ctx context.Context, appID snowflake.ID, guildID snowflake.ID, roleID snowflake.ID) (*model.Role, error) {
	row, err := c.Q.GetRole(ctx, pgmodel.GetRoleParams{
		AppID:   int64(appID),
		GuildID: int64(guildID),
		RoleID:  int64(roleID),
	})
	if err != nil {
		return nil, err
	}
	return rowToRole(row)
}

func (c *Client) GetRoles(ctx context.Context, appID snowflake.ID, guildID snowflake.ID, limit int, offset int) ([]*model.Role, error) {
	rows, err := c.Q.GetRoles(ctx, pgmodel.GetRolesParams{
		AppID:   int64(appID),
		GuildID: int64(guildID),
		Limit:   int32(limit),
		Offset:  int32(offset),
	})
	if err != nil {
		return nil, err
	}

	roles := make([]*model.Role, len(rows))
	for i, row := range rows {
		role, err := rowToRole(row)
		if err != nil {
			return nil, err
		}
		roles[i] = role
	}
	return roles, nil
}

func (c *Client) SearchRoles(ctx context.Context, params store.SearchRolesParams) ([]*model.Role, error) {
	rows, err := c.Q.SearchRoles(ctx, pgmodel.SearchRolesParams{
		AppID:   int64(params.AppID),
		GuildID: int64(params.GuildID),
		Data:    params.Data,
		Limit:   int32(params.Limit),
		Offset:  int32(params.Offset),
	})
	if err != nil {
		return nil, err
	}

	roles := make([]*model.Role, len(rows))
	for i, row := range rows {
		role, err := rowToRole(row)
		if err != nil {
			return nil, err
		}
		roles[i] = role
	}
	return roles, nil
}

func (c *Client) UpsertRoles(ctx context.Context, roles ...store.UpsertRoleParams) error {
	if len(roles) == 0 {
		return nil
	}

	params := make([]pgmodel.UpsertRolesParams, len(roles))
	for i, role := range roles {
		data, err := json.Marshal(role.Data)
		if err != nil {
			return fmt.Errorf("failed to marshal role data: %w", err)
		}

		params[i] = pgmodel.UpsertRolesParams{
			AppID:   int64(role.AppID),
			GuildID: int64(role.GuildID),
			RoleID:  int64(role.RoleID),
			Data:    data,
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

func (c *Client) DeleteRole(ctx context.Context, appID snowflake.ID, guildID snowflake.ID, roleID snowflake.ID) error {
	return c.Q.DeleteRole(ctx, pgmodel.DeleteRoleParams{
		AppID:   int64(appID),
		GuildID: int64(guildID),
		RoleID:  int64(roleID),
	})
}

func rowToRole(row pgmodel.CacheRole) (*model.Role, error) {
	var data discord.Role
	err := json.Unmarshal(row.Data, &data)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal role data: %w", err)
	}

	return &model.Role{
		AppID:     snowflake.ID(row.AppID),
		GuildID:   snowflake.ID(row.GuildID),
		RoleID:    snowflake.ID(row.RoleID),
		Data:      data,
		Tainted:   row.Tainted,
		CreatedAt: row.CreatedAt.Time,
		UpdatedAt: row.UpdatedAt.Time,
	}, nil
}
