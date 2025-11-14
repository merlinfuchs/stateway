package gateway

import (
	"context"
	"encoding/json"

	"github.com/disgoorg/snowflake/v2"
	"github.com/merlinfuchs/stateway/stateway-lib/broker"
	"github.com/merlinfuchs/stateway/stateway-lib/service"
)

var _ Gateway = (*GatewayClient)(nil)

type GatewayClient struct {
	b broker.Broker
}

func NewGatewayClient(b broker.Broker) *GatewayClient {
	return &GatewayClient{b: b}
}

func (c *GatewayClient) GetApp(ctx context.Context, appID snowflake.ID, withSecrets bool) (*App, error) {
	return gatewayRequest[*App](ctx, c.b, GatewayMethodAppGet, GetAppRequest{AppID: appID, WithSecrets: withSecrets})
}

func (c *GatewayClient) GetApps(ctx context.Context, params ListAppsRequest) ([]*App, error) {
	return gatewayRequest[[]*App](ctx, c.b, GatewayMethodAppList, params)
}

func (c *GatewayClient) UpsertApp(ctx context.Context, app UpsertAppRequest) (*App, error) {
	return gatewayRequest[*App](ctx, c.b, GatewayMethodAppUpsert, app)
}

func (c *GatewayClient) DisableApp(ctx context.Context, appID snowflake.ID) error {
	_, err := gatewayRequest[struct{}](ctx, c.b, GatewayMethodAppDisable, DisableAppRequest{AppID: appID})
	return err
}

func (c *GatewayClient) DeleteApp(ctx context.Context, appID snowflake.ID) error {
	_, err := gatewayRequest[struct{}](ctx, c.b, GatewayMethodAppDelete, DeleteAppRequest{AppID: appID})
	return err
}

func (c *GatewayClient) GetGroup(ctx context.Context, groupID string) (*Group, error) {
	return gatewayRequest[*Group](ctx, c.b, GatewayMethodGroupGet, GetGroupRequest{GroupID: groupID})
}

func (c *GatewayClient) GetGroups(ctx context.Context) ([]*Group, error) {
	return gatewayRequest[[]*Group](ctx, c.b, GatewayMethodGroupList, ListGroupsRequest{})
}

func (c *GatewayClient) UpsertGroup(ctx context.Context, group UpsertGroupRequest) (*Group, error) {
	return gatewayRequest[*Group](ctx, c.b, GatewayMethodGroupUpsert, group)
}

func (c *GatewayClient) DeleteGroup(ctx context.Context, groupID string) error {
	_, err := gatewayRequest[struct{}](ctx, c.b, GatewayMethodGroupDelete, DeleteGroupRequest{GroupID: groupID})
	return err
}

func gatewayRequest[R any](ctx context.Context, b broker.Broker, method GatewayMethod, request GatewayRequest) (R, error) {
	var r R

	response, err := b.Request(ctx, service.ServiceTypeGateway, string(method), request)
	if err != nil {
		return r, err
	}

	if !response.Success {
		if response.Error != nil {
			return r, response.Error
		}
		return r, service.ErrUnknown("unknown error")
	}

	err = json.Unmarshal(response.Data, &r)
	if err != nil {
		return r, err
	}

	return r, nil
}
