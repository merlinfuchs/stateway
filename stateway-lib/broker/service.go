package broker

import (
	"context"
	"encoding/json"
)

type ServiceType string

const (
	SerivceTypeGateway ServiceType = "gateway"
	ServiceTypeCache   ServiceType = "cache"
)

type GenericBrokerService interface {
	ServiceType() ServiceType
	HandleRequest(method string, request json.RawMessage) (any, error)
}

type genericBrokerService[REQUEST any, RESPONSE any] struct {
	inner BrokerService[REQUEST, RESPONSE]
}

func (s *genericBrokerService[REQUEST, RESPONSE]) ServiceType() ServiceType {
	return s.inner.ServiceType()
}

func (s *genericBrokerService[REQUEST, RESPONSE]) HandleRequest(method string, request json.RawMessage) (any, error) {
	var req REQUEST
	err := json.Unmarshal(request, &req)
	if err != nil {
		return nil, err
	}

	response, err := s.inner.HandleRequest(method, req)
	if err != nil {
		return nil, err
	}

	return response, nil
}

type BrokerService[REQUEST any, RESPONSE any] interface {
	ServiceType() ServiceType
	HandleRequest(method string, request REQUEST) (RESPONSE, error)
}

func Provide[REQUEST any, RESPONSE any](b Broker, ctx context.Context, server BrokerService[REQUEST, RESPONSE]) error {
	return b.Provide(ctx, &genericBrokerService[REQUEST, RESPONSE]{inner: server})
}
