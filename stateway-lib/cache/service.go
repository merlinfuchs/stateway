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
		if req.GuildID == nil {
			return s.caches.GetChannel(ctx, req.ChannelID, req.Options.Destructure()...)
		} else {
			return s.caches.GetGuildChannel(ctx, *req.GuildID, req.ChannelID, req.Options.Destructure()...)
		}
	case ChannelListRequest:
		if req.GuildID == nil {
			return s.caches.GetChannels(ctx, req.Options.Destructure()...)
		} else {
			return s.caches.GetGuildChannels(ctx, *req.GuildID, req.Options.Destructure()...)
		}
	case ChannelSearchRequest:
		if req.GuildID == nil {
			return s.caches.SearchChannels(ctx, req.Data, req.Options.Destructure()...)
		} else {
			return s.caches.SearchGuildChannels(ctx, *req.GuildID, req.Data, req.Options.Destructure()...)
		}
	case RoleGetRequest:
		if req.GuildID == nil {
			return s.caches.GetRole(ctx, req.RoleID, req.Options.Destructure()...)
		} else {
			return s.caches.GetGuildRole(ctx, *req.GuildID, req.RoleID, req.Options.Destructure()...)
		}
	case RoleListRequest:
		if req.GuildID == nil {
			return s.caches.GetRoles(ctx, req.Options.Destructure()...)
		} else {
			return s.caches.GetGuildRoles(ctx, *req.GuildID, req.Options.Destructure()...)
		}
	case RoleSearchRequest:
		if req.GuildID == nil {
			return s.caches.SearchRoles(ctx, req.Data, req.Options.Destructure()...)
		} else {
			return s.caches.SearchGuildRoles(ctx, *req.GuildID, req.Data, req.Options.Destructure()...)
		}
	}
	return nil, fmt.Errorf("unknown request type: %T", request)
}
