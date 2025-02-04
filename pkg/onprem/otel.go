package onprem

// https://opentelemetry.io/docs/languages/go/

import (
	"context"
	"os"
	"time"

	"github.com/pkg/errors"
	"github.com/shahzodshafizod/gocloud/pkg"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
	"go.opentelemetry.io/otel/trace"
)

type tracer struct {
	provider *sdktrace.TracerProvider
	tracer   trace.Tracer
}

func NewTracer(serviceName string) (pkg.Tracer, error) {
	grpcExporter, err := otlptrace.New(
		context.Background(),
		otlptracegrpc.NewClient(
			otlptracegrpc.WithInsecure(),
			otlptracegrpc.WithEndpoint(os.Getenv("JAEGER_EXPORTER_GRPC_ENDPOINT")),
			// otlptracehttp.WithHeaders(map[string]string{"content-type": "application/json"}),
		),
	)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to create an exporter")
	}

	// Create the resource to be traced
	resource, err := resource.Merge(
		resource.Default(),
		resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceName(serviceName),
			semconv.ServiceVersion("v0.0.1"),
			attribute.String("library.language", "go"),
		),
	)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to create a resource")
	}

	var provider = sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(
			grpcExporter,
			sdktrace.WithMaxExportBatchSize(sdktrace.DefaultMaxExportBatchSize),
			sdktrace.WithBatchTimeout(sdktrace.DefaultScheduleDelay*time.Millisecond),
		),
		sdktrace.WithResource(resource),
	)

	otel.SetTracerProvider(provider)

	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	))

	return &tracer{
		provider: provider,
		tracer:   provider.Tracer(serviceName),
	}, nil
}

func (t *tracer) Shutdown(ctx context.Context) error {
	return t.provider.Shutdown(ctx)
}

func (t *tracer) StartFromContext(ctx context.Context, name string) (context.Context, pkg.Span) {
	ctx, sp := t.tracer.Start(ctx, name)
	return ctx, &span{sp}
}

func (t *tracer) StartFromSpan(ctx context.Context, span pkg.Span, name string) (context.Context, pkg.Span) {
	return t.StartFromContext(span.Context(ctx), name)
}

// Inject context to request headers when you are sending a request
func (t *tracer) Inject(ctx context.Context, carrier pkg.TextMapCarrier) {
	otel.GetTextMapPropagator().Inject(ctx, carrier)
}

// Extract context on the receiver side, happening in middleware in this case
func (t *tracer) Extract(ctx context.Context, carrier pkg.TextMapCarrier) pkg.Span {
	return &span{trace.SpanFromContext(otel.GetTextMapPropagator().Extract(
		ctx,
		carrier,
	))}
}

type span struct {
	span trace.Span
}

func (s *span) End()                  { s.span.End() }
func (s *span) RecordError(err error) { s.span.RecordError(err) }
func (s *span) SetAttributes(key string, value string) {
	s.span.SetAttributes(attribute.String(key, value))
}
func (s *span) Context(ctx context.Context) context.Context {
	return trace.ContextWithSpan(ctx, s.span)
}
