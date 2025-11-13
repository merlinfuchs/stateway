package broker

import (
	"context"

	"github.com/merlinfuchs/stateway/stateway-lib/event"
	"github.com/merlinfuchs/stateway/stateway-lib/service"
)

type Broker interface {
	Publish(ctx context.Context, event event.Event) error
	PublishComplete(ctx context.Context) error
	Listen(ctx context.Context, listener GenericListener) error
	Request(ctx context.Context, serviceType service.ServiceType, method string, request any, opts ...RequestOption) (service.Response, error)
	Provide(ctx context.Context, svc GenericBrokerService) error
	Close(ctx context.Context) error
}
