package broker

import (
	"fmt"

	"github.com/nats-io/nats.go"
)

type NATSBroker struct {
	nc *nats.Conn
}

func (b *NATSBroker) Request(service ServiceType, method string, request any) (Response, error) {
	return Response{
		Success: true,
		Error:   nil,
		Data:    nil,
	}, nil
}

func NewNATSBroker(url string) (*NATSBroker, error) {
	if url == "" {
		url = nats.DefaultURL
	}

	nc, err := nats.Connect(url)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to NATS: %w", err)
	}

	return &NATSBroker{nc: nc}, nil
}
