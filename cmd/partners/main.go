package main

import (
	"os"

	"github.com/shahzodshafizod/gocloud/internal/migration"
	"github.com/shahzodshafizod/gocloud/internal/partners"
	"github.com/shahzodshafizod/gocloud/internal/products"
	"github.com/shahzodshafizod/gocloud/pkg"
	"github.com/shahzodshafizod/gocloud/pkg/http"
	"github.com/shahzodshafizod/gocloud/pkg/onprem"
	"github.com/shahzodshafizod/gocloud/pkg/onprem/postgres/pgx"
	"github.com/shahzodshafizod/gocloud/pkg/onprem/queue/rabbitmq"
	"go.uber.org/fx"
)

func main() {
	fx.New(
		fx.Invoke(func() error {
			return migration.Migrate(os.Getenv("MIGRATION_DIR"))
		}),

		fx.Provide(pgx.NewPostgres),
		fx.Provide(http.NewHTTPClient),
		fx.Provide(func() (pkg.Tracer, error) {
			return onprem.NewTracer(os.Getenv("SERVICE_NAME"))
		}),
		fx.Provide(func(tracer pkg.Tracer) (pkg.Queue, error) {
			return rabbitmq.NewQueue(os.Getenv("SERVICE_NAME"), tracer)
		}),

		fx.Provide(products.NewRepository),
		fx.Provide(products.NewService),
		fx.Provide(partners.NewRepository),
		fx.Provide(partners.NewService),
		fx.Invoke(partners.NewHandler),

		fx.NopLogger,
	).Run()
}
