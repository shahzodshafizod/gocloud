package main

import (
	"os"

	"github.com/shahzodshafizod/gocloud/internal/migration"
	"github.com/shahzodshafizod/gocloud/internal/orders"
	"go.uber.org/fx"
)

func main() {
	fx.New(
		fx.Invoke(func() error {
			return migration.Migrate(os.Getenv("MIGRATION_DIR"))
		}),

		// fx.Provide(NewPostgres),
		// fx.Provide(NewNoSQL),
		// fx.Provide(func() (pkg.Tracer, error) {
		// 	return NewTracer(os.Getenv("SERVICE_NAME"))
		// }),
		// fx.Provide(func(tracer pkg.Tracer) (pkg.Queue, error) {
		// 	return NewQueue(os.Getenv("SERVICE_NAME"), tracer)
		// }),

		fx.Provide(orders.NewRepository),
		fx.Provide(orders.NewService),
		fx.Invoke(orders.NewHandler),

		// fx.NopLogger,
	).Run()
}
