package broker

import (
	"context"
	"fmt"

	"github.com/merlinfuchs/stateway/stateway-lib/event"
	"github.com/merlinfuchs/stateway/stateway-lib/service"
)

type GenericListener interface {
	BalanceKey() string
	ServiceType() service.ServiceType
	EventFilters() []string
	HandleEvent(ctx context.Context, event event.Event) error
}

type genericEventListener[E event.Event] struct {
	serviceType service.ServiceType
	inner       EventListener[E]
}

func (l *genericEventListener[E]) BalanceKey() string {
	return l.inner.BalanceKey()
}

func (l *genericEventListener[E]) ServiceType() service.ServiceType {
	return l.serviceType
}

func (l *genericEventListener[E]) EventFilters() []string {
	return l.inner.EventFilters()
}

func (l *genericEventListener[E]) HandleEvent(ctx context.Context, event event.Event) error {
	e, ok := event.(E)
	if !ok {
		return fmt.Errorf("event is not of type %T", e)
	}

	return l.inner.HandleEvent(ctx, e)
}

type EventListener[E event.Event] interface {
	BalanceKey() string
	EventFilters() []string
	HandleEvent(ctx context.Context, event E) error
}

func Listen[E event.Event](ctx context.Context, b Broker, listener EventListener[E]) error {
	return b.Listen(ctx, &genericEventListener[E]{
		// Get the service type from the event type
		serviceType: (*new(E)).ServiceType(),
		inner:       listener,
	})
}

type EventFilter struct{}
