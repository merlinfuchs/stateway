package gateway

import (
	"encoding/json"
	"fmt"

	"github.com/disgoorg/snowflake/v2"
	validation "github.com/go-ozzo/ozzo-validation/v4"
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

func (m GatewayMethod) UnmarshalRequest(data json.RawMessage) (GatewayRequest, error) {
	switch m {
	case GatewayMethodAppGet:
		var req GetAppRequest
		err := json.Unmarshal(data, &req)
		return req, err
	case GatewayMethodAppList:
		var req ListAppsRequest
		err := json.Unmarshal(data, &req)
		return req, err
	case GatewayMethodAppUpsert:
		var req UpsertAppRequest
		err := json.Unmarshal(data, &req)
		return req, err
	case GatewayMethodAppDisable:
		var req DisableAppRequest
		err := json.Unmarshal(data, &req)
		return req, err
	case GatewayMethodAppDelete:
		var req UpsertAppRequest
		err := json.Unmarshal(data, &req)
		return req, err
	case GatewayMethodGroupGet:
		var req GetGroupRequest
		err := json.Unmarshal(data, &req)
		return req, err
	case GatewayMethodGroupList:
		var req ListGroupsRequest
		err := json.Unmarshal(data, &req)
		return req, err
	case GatewayMethodGroupUpsert:
		var req UpsertGroupRequest
		err := json.Unmarshal(data, &req)
		return req, err
	case GatewayMethodGroupDelete:
		var req DeleteGroupRequest
		err := json.Unmarshal(data, &req)
		return req, err
	default:
		return nil, fmt.Errorf("unknown gateway method: %v", m)
	}
}

type GatewayRequest interface {
	gatewayRequest()
}

type GetAppRequest struct {
	AppID       snowflake.ID `json:"app_id"`
	WithSecrets bool         `json:"with_secrets,omitempty"`
}

func (r GetAppRequest) gatewayRequest() {}

func (r GetAppRequest) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.AppID, validation.Required),
	)
}

type ListAppsRequest struct {
	GroupID     null.String `json:"group_id,omitempty"`
	Limit       null.Int    `json:"limit,omitempty"`
	Offset      null.Int    `json:"offset,omitempty"`
	WithSecrets bool        `json:"with_secrets,omitempty"`
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

func (r UpsertAppRequest) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.ID, validation.Required),
		validation.Field(&r.GroupID, validation.Required),
		validation.Field(&r.DisplayName, validation.Required),
		validation.Field(&r.DiscordClientID, validation.Required),
		validation.Field(&r.DiscordBotToken, validation.Required),
		validation.Field(&r.ShardCount, validation.Required, validation.Min(1)),
	)
}

type DisableAppRequest struct {
	AppID snowflake.ID `json:"app_id"`
}

func (r DisableAppRequest) gatewayRequest() {}

func (r DisableAppRequest) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.AppID, validation.Required),
	)
}

type DeleteAppRequest struct {
	AppID snowflake.ID `json:"app_id"`
}

func (r DeleteAppRequest) gatewayRequest() {}

func (r DeleteAppRequest) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.AppID, validation.Required),
	)
}

type GetGroupRequest struct {
	GroupID string `json:"group_id"`
}

func (r GetGroupRequest) gatewayRequest() {}

func (r GetGroupRequest) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.GroupID, validation.Required),
	)
}

type ListGroupsRequest struct {
}

func (r ListGroupsRequest) gatewayRequest() {}

type UpsertGroupRequest struct {
	ID                 string         `json:"id"`
	DisplayName        string         `json:"display_name"`
	DefaultConstraints AppConstraints `json:"default_constraints"`
	DefaultConfig      AppConfig      `json:"default_config"`
}

func (r UpsertGroupRequest) gatewayRequest() {}

func (r UpsertGroupRequest) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.ID, validation.Required),
		validation.Field(&r.DisplayName, validation.Required),
	)
}

type DeleteGroupRequest struct {
	GroupID string `json:"group_id"`
}

func (r DeleteGroupRequest) gatewayRequest() {}

func (r DeleteGroupRequest) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.GroupID, validation.Required),
	)
}
