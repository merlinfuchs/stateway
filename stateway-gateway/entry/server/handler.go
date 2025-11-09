package server

import (
	"context"
	"log/slog"
	"time"

	"github.com/merlinfuchs/stateway/stateway-lib/broker"
	"github.com/merlinfuchs/stateway/stateway-lib/event"
)

const MaxQueueSize = 100_000

type EventHandler struct {
	broker broker.Broker

	queue chan event.Event
}

func NewEventHandler(broker broker.Broker) *EventHandler {
	return &EventHandler{
		broker: broker,
		queue:  make(chan event.Event, MaxQueueSize),
	}
}

func (h *EventHandler) HandleEvent(event event.Event) {
	select {
	case h.queue <- event:
	case <-time.After(10 * time.Second):
		slog.Error(
			"Queue is full, dropping event",
			slog.String("event_id", event.EventID().String()),
			slog.String("service_type", string(event.ServiceType())),
		)
	}
}

func (h *EventHandler) Run(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case event := <-h.queue:
			err := h.broker.Publish(ctx, event)
			if err != nil {
				slog.Error(
					"Failed to publish event",
					slog.String("event_id", event.EventID().String()),
					slog.String("service_type", string(event.ServiceType())),
					slog.String("error", err.Error()),
				)
			}
		}
	}
}
