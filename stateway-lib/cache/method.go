package cache

import (
	"encoding/json"
	"fmt"

	"github.com/disgoorg/snowflake/v2"
)

type CacheMethod string

const (
	CacheMethodGetGuild       CacheMethod = "guild.get"
	CacheMethodListGuilds     CacheMethod = "guild.list"
	CacheMethodSearchGuilds   CacheMethod = "guild.search"
	CacheMethodGetChannel     CacheMethod = "channel.get"
	CacheMethodListChannels   CacheMethod = "channel.list"
	CacheMethodSearchChannels CacheMethod = "channel.search"
	CacheMethodGetRole        CacheMethod = "role.get"
	CacheMethodListRoles      CacheMethod = "role.list"
	CacheMethodSearchRoles    CacheMethod = "role.search"
)

func (m CacheMethod) UnmarshalRequest(data json.RawMessage) (CacheRequest, error) {
	switch m {
	case CacheMethodGetGuild:
		var req GuildGetRequest
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

type GuildListRequest struct {
	Options CacheOptions `json:"options,omitempty"`
}

func (r GuildListRequest) cacheRequest() {}

type GuildSearchRequest struct {
	Data    json.RawMessage `json:"data,omitempty"`
	Options CacheOptions    `json:"options,omitempty"`
}

func (r GuildSearchRequest) cacheRequest() {}

type ChannelGetRequest struct {
	GuildID   snowflake.ID `json:"guild_id"`
	ChannelID snowflake.ID `json:"channel_id"`
	Options   CacheOptions `json:"options,omitempty"`
}

func (r ChannelGetRequest) cacheRequest() {}

type ChannelListRequest struct {
	GuildID snowflake.ID `json:"guild_id"`
	Options CacheOptions `json:"options,omitempty"`
}

func (r ChannelListRequest) cacheRequest() {}

type ChannelSearchRequest struct {
	Data    json.RawMessage `json:"data,omitempty"`
	Options CacheOptions    `json:"options,omitempty"`
}

func (r ChannelSearchRequest) cacheRequest() {}

type RoleGetRequest struct {
	GuildID snowflake.ID `json:"guild_id"`
	RoleID  snowflake.ID `json:"role_id"`
	Options CacheOptions `json:"options,omitempty"`
}

func (r RoleGetRequest) cacheRequest() {}

type RoleListRequest struct {
	GuildID snowflake.ID `json:"guild_id"`
	Options CacheOptions `json:"options,omitempty"`
}

func (r RoleListRequest) cacheRequest() {}

type RoleSearchRequest struct {
	Data    json.RawMessage `json:"data,omitempty"`
	Options CacheOptions    `json:"options,omitempty"`
}

func (r RoleSearchRequest) cacheRequest() {}
