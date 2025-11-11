package gateway

import (
	"github.com/disgoorg/snowflake/v2"
	"gopkg.in/guregu/null.v4"
)

type GatewayMethod string

const (
	GatewayMethodAppGet      GatewayMethod = "app.get"
	GatewayMethodAppList     GatewayMethod = "app.list"
	GatewayMethodAppUpsert   GatewayMethod = "app.upsert"
	GatewayMethodAppDisable  GatewayMethod = "app.disable"
	GatewayMethodAppDelete   GatewayMethod = "app.delete"
	GatewayMethodGroupGet    GatewayMethod = "group.get"
	GatewayMethodGroupList   GatewayMethod = "group.list"
	GatewayMethodGroupUpsert GatewayMethod = "group.upsert"
	GatewayMethodGroupDelete GatewayMethod = "group.delete"
)

func (m GatewayMethod) RequestType() GatewayRequest {
	switch m {
	case GatewayMethodAppGet:
		return GetAppRequest{}
	case GatewayMethodAppList:
		return ListAppsRequest{}
	case GatewayMethodAppUpsert:
		return UpsertAppRequest{}
	case GatewayMethodAppDisable:
		return DisableAppRequest{}
	case GatewayMethodAppDelete:
		return UpsertAppRequest{}
	case GatewayMethodGroupGet:
		return GetGroupRequest{}
	case GatewayMethodGroupList:
		return ListGroupsRequest{}
	case GatewayMethodGroupUpsert:
		return UpsertGroupRequest{}
	case GatewayMethodGroupDelete:
		return DeleteGroupRequest{}
	default:
		return nil
	}
}

type GatewayRequest interface {
	gatewayRequest()
}

type GetAppRequest struct {
	AppID snowflake.ID `json:"app_id"`
}

func (r GetAppRequest) gatewayRequest() {}

type ListAppsRequest struct {
	GroupID null.String `json:"group_id,omitempty"`
	Limit   int         `json:"limit"`
	Offset  int         `json:"offset"`
}

func (r ListAppsRequest) gatewayRequest() {}

type UpsertAppRequest struct {
	ID                  snowflake.ID   `json:"id"`
	GroupID             string         `json:"group_id"`
	DisplayName         string         `json:"display_name"`
	DiscordClientID     snowflake.ID   `json:"discord_client_id"`
	DiscordBotToken     string         `json:"discord_bot_token"`
	DiscordPublicKey    string         `json:"discord_public_key"`
	DiscordClientSecret null.String    `json:"discord_client_secret"`
	ShardCount          int            `json:"shard_count"`
	Constraints         AppConstraints `json:"constraints"`
	Config              AppConfig      `json:"config"`
}

func (r UpsertAppRequest) gatewayRequest() {}

type DisableAppRequest struct {
	AppID snowflake.ID `json:"app_id"`
}

func (r DisableAppRequest) gatewayRequest() {}

type DeleteAppRequest struct {
	AppID snowflake.ID `json:"app_id"`
}

func (r DeleteAppRequest) gatewayRequest() {}

type GetGroupRequest struct {
	GroupID string `json:"group_id"`
}

func (r GetGroupRequest) gatewayRequest() {}

type ListGroupsRequest struct {
}

func (r ListGroupsRequest) gatewayRequest() {}

type UpsertGroupRequest struct {
	GroupID            string         `json:"group_id"`
	DisplayName        string         `json:"display_name"`
	DefaultConstraints AppConstraints `json:"default_constraints"`
	DefaultConfig      AppConfig      `json:"default_config"`
}

func (r UpsertGroupRequest) gatewayRequest() {}

type DeleteGroupRequest struct {
	GroupID string `json:"group_id"`
}

func (r DeleteGroupRequest) gatewayRequest() {}
