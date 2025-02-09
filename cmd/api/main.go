package main

import (
	"os"

	"github.com/pkg/errors"
	"github.com/shahzodshafizod/gocloud/internal/gateway"
	"github.com/shahzodshafizod/gocloud/internal/orders"
	"github.com/shahzodshafizod/gocloud/internal/partners"
	"github.com/shahzodshafizod/gocloud/pkg"
	"github.com/shahzodshafizod/gocloud/pkg/onprem"
	"github.com/shahzodshafizod/gocloud/pkg/onprem/keycloak"
	"github.com/shahzodshafizod/gocloud/pkg/onprem/queue/rabbitmq"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.uber.org/fx"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// Delivery godoc
//
//	@title						Delivery API Gateway
//	@version					1.0
//	@description				Delivery Requests' Entry Point
//	@contact.name				Shahzod Shafizod
//	@contact.url				http://github.com/shahzodshafizod
//	@contact.email				shahzodshafizod@gmail.com
//	@license.name				© Shahzod Shafizod
//	@host						delivery.local
//	@schemes					http
//	@BasePath					/api/v1
//	@securityDefinitions.apikey	Authorization Token
//	@in							header
//	@name						Authorization
//	@securityDefinitions.apikey	Request Signature
//	@in							header
//	@name						Signature
func main() {
	fx.New(
		fx.Provide(func() (partners.PartnersClient, error) {
			conn, err := grpc.NewClient(
				os.Getenv("PARTNERS_SERVICE_ADDRESS"),
				grpc.WithTransportCredentials(insecure.NewCredentials()),
				grpc.WithStatsHandler(otelgrpc.NewClientHandler()),
			)
			if err != nil {
				return nil, errors.Wrap(err, "grpc.Dial")
			}
			return partners.NewPartnersClient(conn), nil
		}),
		fx.Provide(func() (orders.OrdersClient, error) {
			conn, err := grpc.NewClient(
				os.Getenv("ORDERS_SERVICE_ADDRESS"),
				grpc.WithTransportCredentials(insecure.NewCredentials()),
				grpc.WithStatsHandler(otelgrpc.NewClientHandler()),
			)
			if err != nil {
				return nil, errors.Wrap(err, "grpc.Dial")
			}
			return orders.NewOrdersClient(conn), nil
		}),

		fx.Provide(onprem.NewCache),
		fx.Provide(keycloak.NewAuth),
		fx.Provide(onprem.NewStorage),
		fx.Provide(func() (pkg.Tracer, error) {
			return onprem.NewTracer(os.Getenv("SERVICE_NAME"))
		}),
		fx.Provide(func(tracer pkg.Tracer) (pkg.Queue, error) {
			return rabbitmq.NewQueue(os.Getenv("SERVICE_NAME"), tracer)
		}),

		fx.Provide(gateway.NewService),
		fx.Invoke(gateway.NewHandler),

		fx.NopLogger,
	).Run()
}
