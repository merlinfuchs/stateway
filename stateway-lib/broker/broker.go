package broker

import (
	"encoding/json"
	"fmt"
)

type Broker interface {
	Request(service ServiceType, method string, request any) (Response, error)
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
