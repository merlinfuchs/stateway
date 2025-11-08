package admin

import (
	"context"
	"errors"
	"fmt"

	"github.com/merlinfuchs/stateway/stateway-lib/broker"
	"github.com/merlinfuchs/stateway/stateway-lib/config"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
)

func DeleteGatewayStream(ctx context.Context, config *config.RootGatewayConfig) error {
	nats, err := nats.Connect(config.Broker.NATS.URL)
	if err != nil {
		return fmt.Errorf("failed to connect to NATS: %w", err)
	}

	js, err := jetstream.New(nats)
	if err != nil {
		return fmt.Errorf("failed to create JetStream context: %w", err)
	}

	err = js.DeleteStream(ctx, broker.GatewayStreamName)
	if err != nil && !errors.Is(err, jetstream.ErrStreamNotFound) {
		return fmt.Errorf("failed to delete gateway stream: %w", err)
	}

	return nil
}
