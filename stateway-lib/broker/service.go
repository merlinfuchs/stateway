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

type genericBrokerService[REQUEST any, RESPONSE any, METHOD ServiceMethod[REQUEST]] struct {
	inner BrokerService[REQUEST, RESPONSE, METHOD]
}

func (s *genericBrokerService[REQUEST, RESPONSE, METHOD]) ServiceType() service.ServiceType {
	return s.inner.ServiceType()
}

func (s *genericBrokerService[REQUEST, RESPONSE, METHOD]) HandleRequest(ctx context.Context, rawMethod string, request json.RawMessage) (any, error) {
	method := METHOD(rawMethod)

	// Let the method handle unmarshaling directly - it knows the concrete type
	req, err := method.UnmarshalRequest(request)
	if err != nil {
		return nil, err
	}

	if vReq, ok := any(req).(service.RequestValidate); ok {
		if err := vReq.Validate(); err != nil {
			return nil, service.ErrInvalidRequest("request validation failed", err)
		}
	}

	response, err := s.inner.HandleRequest(ctx, METHOD(method), req)
	if err != nil {
		return nil, err
	}

	return response, nil
}

type BrokerService[REQUEST any, RESPONSE any, METHOD ServiceMethod[REQUEST]] interface {
	ServiceType() service.ServiceType
	HandleRequest(ctx context.Context, method METHOD, request REQUEST) (RESPONSE, error)
}

func Provide[REQUEST any, RESPONSE any, METHOD ServiceMethod[REQUEST]](ctx context.Context, b Broker, server BrokerService[REQUEST, RESPONSE, METHOD]) error {
	return b.Provide(ctx, &genericBrokerService[REQUEST, RESPONSE, METHOD]{inner: server})
}

type ServiceMethod[REQUEST any] interface {
	~string

	UnmarshalRequest(data json.RawMessage) (REQUEST, error)
}
