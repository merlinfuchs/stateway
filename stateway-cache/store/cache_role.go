package store

import (
	"context"
	"encoding/json"
	"time"

	"github.com/disgoorg/snowflake/v2"
)

type UpsertRoleParams struct {
	AppID     snowflake.ID
	GuildID   snowflake.ID
	RoleID    snowflake.ID
	Data      json.RawMessage
	CreatedAt time.Time
	UpdatedAt time.Time
}

type RoleIdentifier struct {
	AppID   snowflake.ID
	GuildID snowflake.ID
	RoleID  snowflake.ID
}

type CacheRoleStore interface {
	UpsertRoles(ctx context.Context, roles ...UpsertRoleParams) error
	DeleteRole(ctx context.Context, params RoleIdentifier) error
}
