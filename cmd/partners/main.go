package main

import (
	"os"

	"github.com/shahzodshafizod/gocloud/internal/migration"
	"github.com/shahzodshafizod/gocloud/internal/partners"
	"github.com/shahzodshafizod/gocloud/internal/products"
	"go.uber.org/fx"
)

func main() {
	fx.New(
		fx.Invoke(func() error {
			return migration.Migrate(os.Getenv("MIGRATION_DIR"))
		}),

		// fx.Provide(NewPostgres),
		// fx.Provide(NewHTTPClient),
		// fx.Provide(func() (pkg.Tracer, error) {
		// 	return NewTracer(os.Getenv("SERVICE_NAME"))
		// }),
		// fx.Provide(func(tracer pkg.Tracer) (pkg.Queue, error) {
		// 	return NewQueue(os.Getenv("SERVICE_NAME"), tracer)
		// }),

		fx.Provide(products.NewRepository),
		fx.Provide(products.NewService),
		fx.Provide(partners.NewRepository),
		fx.Provide(partners.NewService),
		fx.Invoke(partners.NewHandler),

		// fx.NopLogger,
	).Run()
}
