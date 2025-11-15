package postgres

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/snowflake/v2"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/merlinfuchs/stateway/stateway-cache/db/postgres/pgmodel"
	"github.com/merlinfuchs/stateway/stateway-cache/model"
	"github.com/merlinfuchs/stateway/stateway-cache/store"
)

func (c *Client) GetGuildRole(ctx context.Context, appID snowflake.ID, guildID snowflake.ID, roleID snowflake.ID) (*model.Role, error) {
	row, err := c.Q.GetGuildRole(ctx, pgmodel.GetGuildRoleParams{
		AppID:   int64(appID),
		GuildID: int64(guildID),
		RoleID:  int64(roleID),
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, store.ErrNotFound
		}
		return nil, err
	}
	return rowToRole(row)
}

func (c *Client) GetRole(ctx context.Context, appID snowflake.ID, roleID snowflake.ID) (*model.Role, error) {
	row, err := c.Q.GetRole(ctx, pgmodel.GetRoleParams{
		AppID:  int64(appID),
		RoleID: int64(roleID),
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, store.ErrNotFound
		}
		return nil, err
	}
	return rowToRole(row)
}

func (c *Client) GetGuildRoles(ctx context.Context, appID snowflake.ID, guildID snowflake.ID, limit int, offset int) ([]*model.Role, error) {
	rows, err := c.Q.GetGuildRoles(ctx, pgmodel.GetGuildRolesParams{
		AppID:   int64(appID),
		GuildID: int64(guildID),
		Limit: pgtype.Int4{
			Int32: int32(limit),
			Valid: limit != 0,
		},
		Offset: pgtype.Int4{
			Int32: int32(offset),
			Valid: offset != 0,
		},
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, store.ErrNotFound
		}
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

func (c *Client) GetGuildRolesByIDs(ctx context.Context, appID snowflake.ID, guildID snowflake.ID, roleIDs []snowflake.ID) ([]*model.Role, error) {
	roleIDInts := make([]int64, len(roleIDs))
	for i, roleID := range roleIDs {
		roleIDInts[i] = int64(roleID)
	}

	rows, err := c.Q.GetGuildRolesByIDs(ctx, pgmodel.GetGuildRolesByIDsParams{
		AppID:   int64(appID),
		GuildID: int64(guildID),
		RoleIds: roleIDInts,
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

func (c *Client) GetRoles(ctx context.Context, appID snowflake.ID, limit int, offset int) ([]*model.Role, error) {
	rows, err := c.Q.GetRoles(ctx, pgmodel.GetRolesParams{
		AppID: int64(appID),
		Limit: pgtype.Int4{
			Int32: int32(limit),
			Valid: limit != 0,
		},
		Offset: pgtype.Int4{
			Int32: int32(offset),
			Valid: offset != 0,
		},
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, store.ErrNotFound
		}
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

func (c *Client) SearchGuildRoles(ctx context.Context, params store.SearchGuildRolesParams) ([]*model.Role, error) {
	rows, err := c.Q.SearchGuildRoles(ctx, pgmodel.SearchGuildRolesParams{
		AppID:   int64(params.AppID),
		GuildID: int64(params.GuildID),
		Data:    params.Data,
		Limit: pgtype.Int4{
			Int32: int32(params.Limit),
			Valid: params.Limit != 0,
		},
		Offset: pgtype.Int4{
			Int32: int32(params.Offset),
			Valid: params.Offset != 0,
		},
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, store.ErrNotFound
		}
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
		AppID: int64(params.AppID),
		Data:  params.Data,
		Limit: pgtype.Int4{
			Int32: int32(params.Limit),
			Valid: params.Limit != 0,
		},
		Offset: pgtype.Int4{
			Int32: int32(params.Offset),
			Valid: params.Offset != 0,
		},
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, store.ErrNotFound
		}
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

func (c *Client) CountGuildRoles(ctx context.Context, appID snowflake.ID, guildID snowflake.ID) (int, error) {
	res, err := c.Q.CountGuildRoles(ctx, pgmodel.CountGuildRolesParams{
		AppID:   int64(appID),
		GuildID: int64(guildID),
	})
	if err != nil {
		return 0, err
	}
	return int(res), nil
}

func (c *Client) CountRoles(ctx context.Context, appID snowflake.ID) (int, error) {
	res, err := c.Q.CountRoles(ctx, int64(appID))
	if err != nil {
		return 0, err
	}
	return int(res), nil
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
