package cache

import "github.com/disgoorg/snowflake/v2"

type Caches interface {
	GuildCache
}

type CacheOptions struct {
	AppID snowflake.ID `json:"app_id"`
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
