//go:generate mockgen -source=email.go -package=mocks -destination=mocks/email.go
package pkg

import "context"

type Email interface {
	// Sends an email to a specified recipient (`to`) with a given `subject` and `body`.
	// It returns an error if the email fails to send.
	Send(ctx context.Context, to string, subject string, body string) error
}
