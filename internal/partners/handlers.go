package partners

import (
	"context"
	"encoding/json"
	"net"
	"os"

	"github.com/pkg/errors"
	"github.com/shahzodshafizod/gocloud/internal/orders"
	"github.com/shahzodshafizod/gocloud/internal/products"
	"github.com/shahzodshafizod/gocloud/pkg"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.uber.org/fx"
	"google.golang.org/grpc"
)

type handler struct {
	server  *grpc.Server
	service Service
	queue   pkg.Queue
	tracer  pkg.Tracer
}

func NewHandler(lifecycle fx.Lifecycle, service Service, queue pkg.Queue, postgres pkg.Postgres, tracer pkg.Tracer) error {
	var handler = &handler{
		server:  grpc.NewServer(grpc.StatsHandler(otelgrpc.NewServerHandler())),
		service: service,
		queue:   queue,
		tracer:  tracer,
	}
	RegisterPartnersServer(handler.server, handler)
	lis, err := net.Listen("tcp", os.Getenv("SERVICE_ADDRESS"))
	if err != nil {
		return errors.Wrap(err, "net.Listen")
	}
	lifecycle.Append(fx.Hook{
		OnStart: func(context.Context) error {
			err := queue.Subscribe(context.Background(), &pkg.Subscribe{
				Topic:    "orders.paid",
				Callback: handler.sendToPartner,
			})
			if err != nil {
				return nil
			}
			go handler.server.Serve(lis)
			return nil
		},
		OnStop: func(ctx context.Context) error {
			tracer.Shutdown(ctx)
			postgres.Close(ctx)
			queue.Close()
			handler.server.GracefulStop()
			return lis.Close()
		},
	})
	return nil
}

func (h *handler) GetPartnerProducts(ctx context.Context, req *products.GetAllRequest) (*products.GetAllResponse, error) {
	ctx, span := h.tracer.StartFromContext(ctx, "handler.GetPartnerProducts")
	defer span.End()
	resp, err := h.service.getPartnerProducts(ctx)
	if err != nil {
		err = errors.Wrap(err, "h.service.getPartnerProducts")
		span.RecordError(err)
		return nil, err
	}
	return resp, nil
}

func (h *handler) CheckPartnerProducts(ctx context.Context, req *CheckRequest) (*CheckResponse, error) {
	ctx, span := h.tracer.StartFromContext(ctx, "handler.CheckPartnerProducts")
	defer span.End()
	resp, err := h.service.checkPartnerProducts(ctx, req)
	if err != nil {
		err = errors.Wrap(err, "h.service.checkPartnerProducts")
		span.RecordError(err)
		return nil, err
	}
	return resp, nil
}

func (h *handler) sendToPartner(ctx context.Context, span pkg.Span, msg pkg.Message) error {
	ctx, span = h.tracer.StartFromSpan(ctx, span, "handler.sendToPartner")
	defer span.End()
	var order = &orders.PaidOrder{}
	err := json.Unmarshal(msg.Body(), order)
	if err != nil {
		err = errors.Wrap(err, "json.Unmarshal")
		span.RecordError(err)
		return err
	}
	err = h.service.sendToPartner(ctx, order)
	if err != nil {
		err = errors.Wrap(err, "h.service.sendToPartner")
		span.RecordError(err)
		return err
	}
	return nil
}

func (h *handler) mustEmbedUnimplementedPartnersServer() {}
