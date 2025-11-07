package cache

import "github.com/disgoorg/snowflake/v2"

type Caches interface {
	GuildCache
}

type CacheOptions struct {
	GroupID  string       `json:"group_id"`
	ClientID snowflake.ID `json:"client_id"`
}

func ResolveOptions(opts ...CacheOption) CacheOptions {
	options := CacheOptions{}
	for _, opt := range opts {
		opt(&options)
	}
	return options
}

func (o CacheOptions) Destructure() []CacheOption {
	res := []CacheOption{}
	if o.GroupID != "" {
		res = append(res, WithGroupID(o.GroupID))
	}
	if o.ClientID != 0 {
		res = append(res, WithClientID(o.ClientID))
	}
	return res
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
