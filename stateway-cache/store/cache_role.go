package store

import (
	"context"
	"encoding/json"
	"time"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/snowflake/v2"
	"github.com/merlinfuchs/stateway/stateway-cache/model"
)

type UpsertRoleParams struct {
	AppID     snowflake.ID
	GuildID   snowflake.ID
	RoleID    snowflake.ID
	Data      discord.Role
	CreatedAt time.Time
	UpdatedAt time.Time
}

type SearchRolesParams struct {
	AppID   snowflake.ID
	GuildID snowflake.ID
	Limit   int
	Offset  int
	Data    json.RawMessage
}

type CacheRoleStore interface {
	GetRole(ctx context.Context, appID snowflake.ID, guildID snowflake.ID, roleID snowflake.ID) (*model.Role, error)
	GetRoles(ctx context.Context, appID snowflake.ID, guildID snowflake.ID, limit int, offset int) ([]*model.Role, error)
	SearchRoles(ctx context.Context, params SearchRolesParams) ([]*model.Role, error)
	UpsertRoles(ctx context.Context, roles ...UpsertRoleParams) error
	DeleteRole(ctx context.Context, appID snowflake.ID, guildID snowflake.ID, roleID snowflake.ID) error
}
