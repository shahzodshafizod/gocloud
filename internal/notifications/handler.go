package notifications

import (
	"context"
	"encoding/json"

	"github.com/pkg/errors"
	"github.com/shahzodshafizod/gocloud/pkg"
	"go.uber.org/fx"
)

type handler struct {
	service Service
	queue   pkg.Subscriber
	tracer  pkg.Tracer
}

func NewHandler(lifecycle fx.Lifecycle, service Service, queue pkg.Queue, postgres pkg.Postgres, tracer pkg.Tracer) {
	handler := &handler{
		service: service,
		queue:   queue,
		tracer:  tracer,
	}

	lifecycle.Append(fx.Hook{
		OnStart: func(context.Context) error {
			handler.queue.Subscribe(context.Background(), &pkg.Subscribe{
				Topic:    "orders.delivering",
				Callback: handler.delivering,
			})
			return nil
		},
		OnStop: func(ctx context.Context) error {
			postgres.Close(ctx)
			queue.Close()
			tracer.Shutdown(ctx)
			return nil
		},
	})
}

func (h *handler) delivering(ctx context.Context, span pkg.Span, msg pkg.Message) error {
	ctx, span = h.tracer.StartFromSpan(ctx, span, "handler.delivering")
	defer span.End()
	var req = &Message{}
	err := json.Unmarshal(msg.Body(), req)
	if err != nil {
		err = errors.Wrap(err, "json.Unmarshal")
		span.RecordError(err)
		return err
	}
	err = h.service.sendNotification(ctx, req)
	if err != nil {
		err = errors.Wrap(err, "h.service.sendNotification")
		span.RecordError(err)
		return err
	}
	return nil
}
