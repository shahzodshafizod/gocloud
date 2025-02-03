//go:generate mockgen -source=queue.go -package=mocks -destination=mocks/queue.go
package pkg

import "context"

type Queue interface {
	Subscriber
	Publisher
	// Closes the connection to the queue, returning an error if the operation fails.
	Close() error
}

type Subscriber interface {
	// Subscribes to a topic or queue for receiving messages, returning an error if the subscription fails.
	Subscribe(ctx context.Context, subscr *Subscribe) error
}

type Publisher interface {
	// Publishes a message to a specified topic, returning an error if the operation fails.
	Publish(ctx context.Context, topic string, request []byte) error
	// Sends a request to a specified topic and waits for a response, returning the response data or an error.
	Request(ctx context.Context, topic string, request []byte) ([]byte, error)
}

type Subscribe struct {
	Topic    string
	Callback func(context.Context, Span, Message) error
}

type Message interface {
	Body() []byte
	ReplyTo() string
}
