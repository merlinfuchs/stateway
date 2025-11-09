package broker

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/merlinfuchs/stateway/stateway-lib/event"
	"github.com/merlinfuchs/stateway/stateway-lib/service"
)

type Broker interface {
	Publish(ctx context.Context, event event.Event) error
	PublishComplete(ctx context.Context) error
	Listen(ctx context.Context, listener GenericListener) error
	Request(ctx context.Context, service service.ServiceType, method string, request any, opts ...RequestOption) (Response, error)
	Provide(ctx context.Context, service GenericBrokerService) error
}

type Response struct {
	Success bool
	Error   *Error
	Data    json.RawMessage
}

type Error struct {
	Message string
	Code    string
}

func (e *Error) Error() string {
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}
