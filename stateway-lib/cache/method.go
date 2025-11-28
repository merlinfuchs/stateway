package cache

import (
	"encoding/json"
	"fmt"

	"github.com/disgoorg/snowflake/v2"
)

type CacheMethod string

const (
	CacheMethodGetGuild                    CacheMethod = "guild.get"
	CacheMethodGetGuildWithPermissions     CacheMethod = "guild.get_with_permissions"
	CacheMethodListGuilds                  CacheMethod = "guild.list"
	CacheMethodSearchGuilds                CacheMethod = "guild.search"
	CacheMethodCountGuilds                 CacheMethod = "guild.count"
	CacheMethodGetChannel                  CacheMethod = "channel.get"
	CacheMethodListChannels                CacheMethod = "channel.list"
	CacheMethodListChannelsWithPermissions CacheMethod = "channel.list_with_permissions"
	CacheMethodSearchChannels              CacheMethod = "channel.search"
	CacheMethodCountChannels               CacheMethod = "channel.count"
	CacheMethodGetRole                     CacheMethod = "role.get"
	CacheMethodListRoles                   CacheMethod = "role.list"
	CacheMethodSearchRoles                 CacheMethod = "role.search"
	CacheMethodCountRoles                  CacheMethod = "role.count"
	CacheMethodComputePermissions          CacheMethod = "permissions.compute"
	CacheMethodMassComputePermissions      CacheMethod = "permissions.mass_compute"
)

func (m CacheMethod) UnmarshalRequest(data json.RawMessage) (CacheRequest, error) {
	switch m {
	case CacheMethodGetGuild:
		var req GuildGetRequest
		err := json.Unmarshal(data, &req)
		return req, err
	case CacheMethodGetGuildWithPermissions:
		var req GuildGetWithPermissionsRequest
		err := json.Unmarshal(data, &req)
		return req, err
	case CacheMethodListGuilds:
		var req GuildListRequest
		err := json.Unmarshal(data, &req)
		return req, err
	case CacheMethodSearchGuilds:
		var req GuildSearchRequest
		err := json.Unmarshal(data, &req)
		return req, err
	case CacheMethodGetChannel:
		var req ChannelGetRequest
		err := json.Unmarshal(data, &req)
		return req, err
	case CacheMethodListChannels:
		var req ChannelListRequest
		err := json.Unmarshal(data, &req)
		return req, err
	case CacheMethodListChannelsWithPermissions:
		var req ChannelListWithPermissionsRequest
		err := json.Unmarshal(data, &req)
		return req, err
	case CacheMethodSearchChannels:
		var req ChannelSearchRequest
		err := json.Unmarshal(data, &req)
		return req, err
	case CacheMethodGetRole:
		var req RoleGetRequest
		err := json.Unmarshal(data, &req)
		return req, err
	case CacheMethodListRoles:
		var req RoleListRequest
		err := json.Unmarshal(data, &req)
		return req, err
	case CacheMethodSearchRoles:
		var req RoleSearchRequest
		err := json.Unmarshal(data, &req)
		return req, err
	case CacheMethodComputePermissions:
		var req PermissionsComputeRequest
		err := json.Unmarshal(data, &req)
		return req, err
	case CacheMethodMassComputePermissions:
		var req MassComputePermissionsRequest
		err := json.Unmarshal(data, &req)
		return req, err
	default:
		return nil, fmt.Errorf("unknown cache method: %v", m)
	}
}

type CacheRequest interface {
	cacheRequest()
}

type GuildGetRequest struct {
	GuildID snowflake.ID `json:"guild_id"`
	Options CacheOptions `json:"options,omitempty"`
}

func (r GuildGetRequest) cacheRequest() {}

type GuildGetWithPermissionsRequest struct {
	GuildID snowflake.ID   `json:"guild_id"`
	UserID  snowflake.ID   `json:"user_id"`
	RoleIDs []snowflake.ID `json:"role_ids"`
	Options CacheOptions   `json:"options,omitempty"`
}

func (r GuildGetWithPermissionsRequest) cacheRequest() {}

type GuildListRequest struct {
	Options CacheOptions `json:"options,omitempty"`
}

func (r GuildListRequest) cacheRequest() {}

type GuildSearchRequest struct {
	Data    json.RawMessage `json:"data,omitempty"`
	Options CacheOptions    `json:"options,omitempty"`
}

func (r GuildSearchRequest) cacheRequest() {}

type GuildCountRequest struct {
	Options CacheOptions `json:"options,omitempty"`
}

func (r GuildCountRequest) cacheRequest() {}

type ChannelGetRequest struct {
	GuildID   *snowflake.ID `json:"guild_id,omitempty"`
	ChannelID snowflake.ID  `json:"channel_id"`
	Options   CacheOptions  `json:"options,omitempty"`
}

func (r ChannelGetRequest) cacheRequest() {}

type ChannelListRequest struct {
	GuildID *snowflake.ID `json:"guild_id,omitempty"`
	Options CacheOptions  `json:"options,omitempty"`
}

func (r ChannelListRequest) cacheRequest() {}

type ChannelListWithPermissionsRequest struct {
	GuildID snowflake.ID   `json:"guild_id"`
	UserID  snowflake.ID   `json:"user_id"`
	RoleIDs []snowflake.ID `json:"role_ids"`
	Options CacheOptions   `json:"options,omitempty"`
}

func (r ChannelListWithPermissionsRequest) cacheRequest() {}

type ChannelSearchRequest struct {
	GuildID *snowflake.ID   `json:"guild_id,omitempty"`
	Data    json.RawMessage `json:"data,omitempty"`
	Options CacheOptions    `json:"options,omitempty"`
}

func (r ChannelSearchRequest) cacheRequest() {}

type ChannelCountRequest struct {
	GuildID *snowflake.ID `json:"guild_id,omitempty"`
	Options CacheOptions  `json:"options,omitempty"`
}

func (r ChannelCountRequest) cacheRequest() {}

type RoleGetRequest struct {
	GuildID *snowflake.ID `json:"guild_id,omitempty"`
	RoleID  snowflake.ID  `json:"role_id"`
	Options CacheOptions  `json:"options,omitempty"`
}

func (r RoleGetRequest) cacheRequest() {}

type RoleListRequest struct {
	GuildID *snowflake.ID `json:"guild_id,omitempty"`
	Options CacheOptions  `json:"options,omitempty"`
}

func (r RoleListRequest) cacheRequest() {}

type RoleSearchRequest struct {
	GuildID *snowflake.ID   `json:"guild_id,omitempty"`
	Data    json.RawMessage `json:"data,omitempty"`
	Options CacheOptions    `json:"options,omitempty"`
}

func (r RoleSearchRequest) cacheRequest() {}

type RoleCountRequest struct {
	GuildID *snowflake.ID `json:"guild_id,omitempty"`
	Options CacheOptions  `json:"options,omitempty"`
}

func (r RoleCountRequest) cacheRequest() {}

type PermissionsComputeRequest struct {
	GuildID   *snowflake.ID  `json:"guild_id,omitempty"`
	ChannelID *snowflake.ID  `json:"channel_id,omitempty"`
	UserID    snowflake.ID   `json:"user_id"`
	RoleIDs   []snowflake.ID `json:"role_ids"`
	Options   CacheOptions   `json:"options,omitempty"`
}

func (r PermissionsComputeRequest) cacheRequest() {}

type MassComputePermissionsRequest struct {
	GuildID    snowflake.ID   `json:"guild_id,omitempty"`
	ChannelIDs []snowflake.ID `json:"channel_ids,omitempty"`
	UserID     snowflake.ID   `json:"user_id"`
	RoleIDs    []snowflake.ID `json:"role_ids"`
	Options    CacheOptions   `json:"options,omitempty"`
}

func (r MassComputePermissionsRequest) cacheRequest() {}
