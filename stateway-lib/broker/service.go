package broker

import (
	"context"
	"encoding/json"

	"github.com/merlinfuchs/stateway/stateway-lib/service"
)

type GenericBrokerService interface {
	ServiceType() service.ServiceType
	HandleRequest(ctx context.Context, method string, request json.RawMessage) (any, error)
}

type genericBrokerService[REQUEST any, RESPONSE any] struct {
	inner BrokerService[REQUEST, RESPONSE]
}

func (s *genericBrokerService[REQUEST, RESPONSE]) ServiceType() service.ServiceType {
	return s.inner.ServiceType()
}

func (s *genericBrokerService[REQUEST, RESPONSE]) HandleRequest(ctx context.Context, method string, request json.RawMessage) (any, error) {
	var req REQUEST
	err := json.Unmarshal(request, &req)
	if err != nil {
		return nil, err
	}

	response, err := s.inner.HandleRequest(ctx, method, req)
	if err != nil {
		return nil, err
	}

	return response, nil
}

type BrokerService[REQUEST any, RESPONSE any] interface {
	ServiceType() service.ServiceType
	HandleRequest(ctx context.Context, method string, request REQUEST) (RESPONSE, error)
}

func Provide[REQUEST any, RESPONSE any](ctx context.Context, b Broker, server BrokerService[REQUEST, RESPONSE]) error {
	return b.Provide(ctx, &genericBrokerService[REQUEST, RESPONSE]{inner: server})
}
