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
		return s.gateway.GetApp(ctx, req.AppID)
	case ListAppsRequest:
		return s.gateway.GetApps(ctx, req.GroupID, req.Limit, req.Offset)
	case UpsertAppRequest:
		err := s.gateway.UpsertApp(ctx, req)
		if err != nil {
			return nil, err
		}
		return struct{}{}, nil
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
		err := s.gateway.UpsertGroup(ctx, req)
		if err != nil {
			return nil, err
		}
		return struct{}{}, nil
	case DeleteGroupRequest:
		err := s.gateway.DeleteGroup(ctx, req.GroupID)
		if err != nil {
			return nil, err
		}
		return struct{}{}, nil
	}
	return nil, fmt.Errorf("unknown method: %s", method)
}
