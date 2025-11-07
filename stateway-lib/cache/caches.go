package cache

import "github.com/disgoorg/snowflake/v2"

type Caches interface {
	GuildCache
}

type CacheOptions struct {
	GroupID  string       `json:"group_id"`
	ClientID snowflake.ID `json:"client_id"`
}

type CacheOption func(*CacheOptions)

func WithGroupID(groupID string) CacheOption {
	return func(o *CacheOptions) {
		o.GroupID = groupID
	}
}

func WithClientID(clientID snowflake.ID) CacheOption {
	return func(o *CacheOptions) {
		o.ClientID = clientID
	}
}
