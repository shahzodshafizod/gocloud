package orders

import (
	"context"
	"encoding/json"
	"net"
	"os"

	"github.com/pkg/errors"
	"github.com/shahzodshafizod/gocloud/pkg"
	"go.uber.org/fx"
	"google.golang.org/grpc"
)

type handler struct {
	server  *grpc.Server
	service Service
	queue   pkg.Subscriber
	tracer  pkg.Tracer
}

func NewHandler(lifecycle fx.Lifecycle, service Service, queue pkg.Queue, postgres pkg.Postgres, tracer pkg.Tracer) error {
	handler := &handler{
		server:  grpc.NewServer(),
		service: service,
		queue:   queue,
		tracer:  tracer,
	}
	RegisterOrdersServer(handler.server, handler)
	lis, err := net.Listen("tcp", os.Getenv("SERVICE_ADDRESS"))
	if err != nil {
		return errors.Wrap(err, "net.Listen")
	}
	lifecycle.Append(fx.Hook{
		OnStart: func(context.Context) error {
			handler.queue.Subscribe(context.Background(), &pkg.Subscribe{
				Topic:    "orders.ready",
				Callback: handler.pickupOrder,
			})
			go handler.server.Serve(lis)
			return nil
		},
		OnStop: func(ctx context.Context) error {
			tracer.Shutdown(ctx)
			queue.Close()
			postgres.Close(ctx)
			handler.server.GracefulStop()
			return lis.Close()
		},
	})
	return nil
}

func (h *handler) CreateOrder(ctx context.Context, order *Order) (*CreateResponse, error) {
	ctx, span := h.tracer.StartFromContext(ctx, "handler.CreateOrder")
	defer span.End()
	resp, err := h.service.createOrder(ctx, order)
	if err != nil {
		err = errors.Wrap(err, "h.service.createOrder")
		span.RecordError(err)
		return nil, err
	}
	return resp, nil
}

func (h *handler) PayOrder(ctx context.Context, req *PayRequest) (*PayResponse, error) {
	ctx, span := h.tracer.StartFromContext(ctx, "handler.PayOrder")
	defer span.End()
	resp, err := h.service.payOrder(ctx, req)
	if err != nil {
		err = errors.Wrap(err, "h.service.payOrder")
		span.RecordError(err)
		return nil, err
	}
	return resp, nil
}

func (h *handler) AssignOrder(ctx context.Context, req *AssignRequest) (*AssignResponse, error) {
	ctx, span := h.tracer.StartFromContext(ctx, "handler.AssignOrder")
	defer span.End()
	resp, err := h.service.assignOrder(ctx, req)
	if err != nil {
		err = errors.Wrap(err, "h.service.assignOrder")
		span.RecordError(err)
		return nil, err
	}
	return resp, nil
}

func (h *handler) pickupOrder(ctx context.Context, span pkg.Span, msg pkg.Message) error {
	ctx, span = h.tracer.StartFromSpan(ctx, span, "handler.pickupOrder")
	defer span.End()
	var req = &pickupRequest{}
	err := json.Unmarshal(msg.Body(), req)
	if err != nil {
		err = errors.Wrap(err, "json.Unmarshal")
		span.RecordError(err)
		return err
	}
	err = h.service.pickupOrder(ctx, req)
	if err != nil {
		err = errors.Wrap(err, "h.service.pickupOrder")
		span.RecordError(err)
		return err
	}
	return nil
}

func (h *handler) mustEmbedUnimplementedOrdersServer() {}
