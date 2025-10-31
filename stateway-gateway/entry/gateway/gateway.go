package gateway

import (
	"context"
	"fmt"

	"github.com/merlinfuchs/stateway/stateway-gateway/app"
	"github.com/merlinfuchs/stateway/stateway-gateway/config"
	"github.com/merlinfuchs/stateway/stateway-gateway/db/postgres"
	"github.com/merlinfuchs/stateway/stateway-lib/event"
)

type eventHandler struct {
}

func (h *eventHandler) HandleEvent(event event.Event) {
	fmt.Println(event.EventType())
}

func Run(ctx context.Context, pg *postgres.Client, cfg *config.Config) error {
	appManager := app.NewAppManager(pg, &eventHandler{})

	appManager.Run(ctx)
	return nil
}
