//go:generate mockgen -source=tracer.go -package=mocks -destination=mocks/tracer.go
package pkg

import "context"

type Tracer interface {
	// Starts a new span from an existing context, returning the updated context and the new span.
	StartFromContext(ctx context.Context, name string) (context.Context, Span)
	// Starts a new child span from an existing span, returning the updated context and the new span.
	StartFromSpan(ctx context.Context, span Span, name string) (context.Context, Span)
	// Gracefully shuts down the tracer, returning an error if the operation fails.
	Shutdown(ctx context.Context) error
	// Injects span context into a carrier (e.g., HTTP headers) for propagation across services.
	Inject(ctx context.Context, carrier TextMapCarrier)
	// Extracts span context from a carrier and returns the corresponding span.
	Extract(ctx context.Context, carrier TextMapCarrier) Span
}

type Span interface {
	// Marks the end of the span's execution.
	End()
	// Records an error associated with the span.
	RecordError(err error)
	// Adds key-value pairs as metadata to the span.
	SetAttributes(key string, value string)
	// Returns the context associated with the span.
	Context(ctx context.Context) context.Context
}

type TextMapCarrier interface {
	Get(key string) string
	Set(key string, value string)
	Keys() []string
}
