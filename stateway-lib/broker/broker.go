package broker

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/merlinfuchs/stateway/stateway-lib/event"
)

type Broker interface {
	PublishEvent(event event.Event) error
	Request(service ServiceType, method string, request any, opts ...RequestOption) (Response, error)
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
