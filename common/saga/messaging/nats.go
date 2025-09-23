package messaging

import (
	"log"

	"github.com/nats-io/nats.go"
)

type NatsClient struct {
	Conn *nats.Conn
}

func NewNatsClient(url string) *NatsClient {
	nc, err := nats.Connect(url)
	if err != nil {
		log.Fatalf("Error connecting to NATS: %v", err)
	}
	return &NatsClient{Conn: nc}
}

func (c *NatsClient) Publish(subject string, data []byte) error {
	return c.Conn.Publish(subject, data)
}

func (c *NatsClient) Subscribe(subject string, handler func(msg *nats.Msg)) {
	_, err := c.Conn.Subscribe(subject, handler)
	if err != nil {
		log.Fatalf("Error subscribing to subject %s: %v", subject, err)
	}
}
