package broker

import (
	"context"
	"fmt"

	"github.com/disgoorg/snowflake/v2"
	"github.com/merlinfuchs/stateway/stateway-lib/event"
	"github.com/merlinfuchs/stateway/stateway-lib/service"
)

type GenericListener interface {
	BalanceKey() string
	ServiceType() service.ServiceType
	EventFilter() EventFilter
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

func (l *genericEventListener[E]) EventFilter() EventFilter {
	return l.inner.EventFilter()
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
	EventFilter() EventFilter
	HandleEvent(ctx context.Context, event E) error
}

func Listen[E event.Event](ctx context.Context, b Broker, listener EventListener[E]) error {
	return b.Listen(ctx, &genericEventListener[E]{
		// Get the service type from the event type
		serviceType: (*new(E)).ServiceType(),
		inner:       listener,
	})
}

type FuncListener[E event.Event] struct {
	balanceKey  string
	eventFilter EventFilter
	handleEvent func(ctx context.Context, event E) error
}

func NewFuncListener[E event.Event](
	balanceKey string,
	eventFilter EventFilter,
	handleEvent func(ctx context.Context, event E) error,
) *FuncListener[E] {
	return &FuncListener[E]{
		balanceKey:  balanceKey,
		eventFilter: eventFilter,
		handleEvent: handleEvent,
	}
}

func (l *FuncListener[E]) BalanceKey() string {
	return l.balanceKey
}

func (l *FuncListener[E]) EventFilter() EventFilter {
	return l.eventFilter
}

func (l *FuncListener[E]) HandleEvent(ctx context.Context, event E) error {
	return l.handleEvent(ctx, event)
}

type EventFilter struct {
	GatewayIDs []int
	GroupIDs   []string
	AppIDs     []snowflake.ID
	EventTypes []string
}

func (f EventFilter) Subjects() []string {
	subjects := []string{}

	gatewayIDs := make([]string, len(f.GatewayIDs))
	for i, gatewayID := range f.GatewayIDs {
		gatewayIDs[i] = fmt.Sprintf("%d", gatewayID)
	}
	if len(gatewayIDs) == 0 {
		gatewayIDs = []string{"*"}
	}

	groupIDs := f.GroupIDs
	if len(groupIDs) == 0 {
		groupIDs = []string{"*"}
	}

	appIDs := make([]string, len(f.AppIDs))
	for i, appID := range f.AppIDs {
		appIDs[i] = appID.String()
	}
	if len(appIDs) == 0 {
		appIDs = []string{"*"}
	}

	eventTypes := f.EventTypes
	if len(eventTypes) == 0 {
		eventTypes = []string{"*"}
	}

	for _, gatewayID := range gatewayIDs {
		for _, groupID := range groupIDs {
			for _, appID := range appIDs {
				for _, eventType := range eventTypes {
					subjects = append(
						subjects,
						fmt.Sprintf(
							"%s.%s.%s.%s",
							gatewayID, groupID, appID, eventType,
						),
					)
				}
			}
		}
	}
	return subjects
}
