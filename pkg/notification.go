//go:generate mockgen -source=notification.go -package=mocks -destination=mocks/notification.go
package pkg

import "context"

// Push Notification
type Notification interface {
	// Sends a push notification to a single recipient (`to`) with a message body and priority level, returning a notification ID or error.
	SendPush(ctx context.Context, to string, body string, priority int) (string, error)
	// Sends a push notification to multiple recipients (`tos`) with a message body and priority level, returning a notification ID or error.
	SendToMany(ctx context.Context, tos []string, body string, priority int) (string, error)
	// Sends a push notification to all devices subscribed to a specific topic, including custom data, and returns a notification ID or error.
	SendToTopic(ctx context.Context, topic string, data map[string]string) (string, error)
}
