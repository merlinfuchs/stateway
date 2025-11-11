package gateway

import (
	"context"
	"fmt"

	"github.com/merlinfuchs/stateway/stateway-lib/service"
)

type GatewayService struct {
	gateway Gateway
}

func NewGatewayService(gateway Gateway) *GatewayService {
	return &GatewayService{gateway: gateway}
}

func (s *GatewayService) ServiceType() service.ServiceType {
	return service.ServiceTypeGateway
}

func (s *GatewayService) HandleRequest(ctx context.Context, method GatewayMethod, request GatewayRequest) (any, error) {
	switch req := request.(type) {
	case GetAppRequest:
		return s.gateway.GetApp(ctx, req.AppID, req.WithSecrets)
	case ListAppsRequest:
		return s.gateway.GetApps(ctx, req)
	case UpsertAppRequest:
		return s.gateway.UpsertApp(ctx, req)
	case DisableAppRequest:
		err := s.gateway.DisableApp(ctx, req.AppID)
		if err != nil {
			return nil, err
		}
		return struct{}{}, nil
	case DeleteAppRequest:
		err := s.gateway.DeleteApp(ctx, req.AppID)
		if err != nil {
			return nil, err
		}
		return struct{}{}, nil
	case GetGroupRequest:
		return s.gateway.GetGroup(ctx, req.GroupID)
	case ListGroupsRequest:
		return s.gateway.GetGroups(ctx)
	case UpsertGroupRequest:
		return s.gateway.UpsertGroup(ctx, req)
	case DeleteGroupRequest:
		err := s.gateway.DeleteGroup(ctx, req.GroupID)
		if err != nil {
			return nil, err
		}
		return struct{}{}, nil
	}
	return nil, fmt.Errorf("unknown method: %s", method)
}
