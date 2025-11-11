package gateway

import (
	"context"
	"encoding/json"

	"github.com/disgoorg/snowflake/v2"
	"github.com/merlinfuchs/stateway/stateway-lib/broker"
	"github.com/merlinfuchs/stateway/stateway-lib/service"
)

type GatewayClient struct {
	b broker.Broker
}

func NewGatewayClient(b broker.Broker) *GatewayClient {
	return &GatewayClient{b: b}
}

func (c *GatewayClient) GetApp(ctx context.Context, appID snowflake.ID) (*App, error) {
	return gatewayRequest[*App](ctx, c.b, GatewayMethodAppGet, GetAppRequest{AppID: appID})
}

func (c *GatewayClient) ListApps(ctx context.Context, params ListAppsRequest) ([]*App, error) {
	return gatewayRequest[[]*App](ctx, c.b, GatewayMethodAppList, params)
}

func (c *GatewayClient) UpsertApp(ctx context.Context, app UpsertAppRequest) error {
	_, err := gatewayRequest[struct{}](ctx, c.b, GatewayMethodAppUpsert, app)
	return err
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

func (c *GatewayClient) ListGroups(ctx context.Context) ([]*Group, error) {
	return gatewayRequest[[]*Group](ctx, c.b, GatewayMethodGroupList, ListGroupsRequest{})
}

func (c *GatewayClient) UpsertGroup(ctx context.Context, group UpsertGroupRequest) error {
	_, err := gatewayRequest[struct{}](ctx, c.b, GatewayMethodGroupUpsert, group)
	return err
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

	err = json.Unmarshal(response.Data, &r)
	if err != nil {
		return r, err
	}

	return r, nil
}
