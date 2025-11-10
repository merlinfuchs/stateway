package cache

import "github.com/disgoorg/snowflake/v2"

type Caches interface {
	GuildCache
	ChannelCache
	RoleCache
}

type CacheOptions struct {
	AppID  snowflake.ID `json:"app_id"`
	Limit  int          `json:"limit"`
	Offset int          `json:"offset"`
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
	if o.AppID != 0 {
		res = append(res, WithAppID(o.AppID))
	}
	return res
}

type CacheOption func(*CacheOptions)

func WithAppID(appID snowflake.ID) CacheOption {
	return func(o *CacheOptions) {
		o.AppID = appID
	}
}

func WithLimit(limit int) CacheOption {
	return func(o *CacheOptions) {
		o.Limit = limit
	}
}

func WithOffset(offset int) CacheOption {
	return func(o *CacheOptions) {
		o.Offset = offset
	}
}
