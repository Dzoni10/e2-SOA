package events

import (
	"encoding/json"

	"github.com/nats-io/nats.go"
)

// Event struct
type BlogCreatedEvent struct {
	BlogID    string `json:"blog_id"`
	AuthorID  string `json:"author_id"`
	Title     string `json:"title"`
	Content   string `json:"content"`
	CreatedAt string `json:"created_at"`
}

// Publisher helper
type Publisher struct {
	Conn *nats.Conn
}

func ConnectNATS(url string) (*nats.Conn, error) {
	return nats.Connect(url)
}

func (p *Publisher) Publish(subject string, v interface{}) error {
	data, err := json.Marshal(v)
	if err != nil {
		return err
	}
	return p.Conn.Publish(subject, data)
}
