package cache

import (
	"context"
	"fmt"

	"github.com/merlinfuchs/stateway/stateway-lib/service"
)

type CacheService struct {
	caches Cache
}

func NewCacheService(caches Cache) *CacheService {
	return &CacheService{caches: caches}
}

func (s *CacheService) ServiceType() service.ServiceType {
	return service.ServiceTypeCache
}

func (s *CacheService) HandleRequest(ctx context.Context, method CacheMethod, request CacheRequest) (any, error) {
	switch req := request.(type) {
	case GuildGetRequest:
		return s.caches.GetGuild(ctx, req.GuildID, req.Options.Destructure()...)
	case GuildListRequest:
		return s.caches.GetGuilds(ctx, req.Options.Destructure()...)
	case GuildSearchRequest:
		return s.caches.SearchGuilds(ctx, req.Data, req.Options.Destructure()...)
	case ChannelGetRequest:
		return s.caches.GetChannel(ctx, req.GuildID, req.ChannelID, req.Options.Destructure()...)
	case ChannelListRequest:
		return s.caches.GetChannels(ctx, req.GuildID, req.Options.Destructure()...)
	case ChannelSearchRequest:
		return s.caches.SearchChannels(ctx, req.Data, req.Options.Destructure()...)
	case RoleGetRequest:
		return s.caches.GetRole(ctx, req.GuildID, req.RoleID, req.Options.Destructure()...)
	case RoleListRequest:
		return s.caches.GetRoles(ctx, req.GuildID, req.Options.Destructure()...)
	case RoleSearchRequest:
		return s.caches.SearchRoles(ctx, req.Data, req.Options.Destructure()...)
	}
	return nil, fmt.Errorf("unknown request type: %T", request)
}
