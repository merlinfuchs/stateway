package cache

import (
	"encoding/json"

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

func (m CacheMethod) RequestType() CacheRequest {
	switch m {
	case CacheMethodGetGuild:
		return GuildGetRequest{}
	case CacheMethodListGuilds:
		return GuildListRequest{}
	case CacheMethodSearchGuilds:
		return GuildSearchRequest{}
	case CacheMethodGetChannel:
		return ChannelGetRequest{}
	case CacheMethodListChannels:
		return ChannelListRequest{}
	case CacheMethodSearchChannels:
		return ChannelSearchRequest{}
	case CacheMethodGetRole:
		return RoleGetRequest{}
	case CacheMethodListRoles:
		return RoleListRequest{}
	case CacheMethodSearchRoles:
		return RoleSearchRequest{}
	default:
		return nil
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
