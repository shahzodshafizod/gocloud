package main

import (
	"os"

	"github.com/shahzodshafizod/gocloud/internal/migration"
	"github.com/shahzodshafizod/gocloud/internal/notifications"
	"go.uber.org/fx"
)

func main() {
	fx.New(
		fx.Invoke(func() error {
			return migration.Migrate(os.Getenv("MIGRATION_DIR"))
		}),

		// fx.Provide(NewPostgres),
		// fx.Provide(NewNoSQL),
		// fx.Provide(NewNotification),
		// fx.Provide(func() (pkg.Tracer, error) {
		// 	return NewTracer(os.Getenv("SERVICE_NAME"))
		// }),
		// fx.Provide(func(tracer pkg.Tracer) (pkg.Queue, error) {
		// 	return NewQueue(os.Getenv("SERVICE_NAME"), tracer)
		// }),

		fx.Provide(notifications.NewRepository),
		fx.Provide(notifications.NewService),
		fx.Invoke(notifications.NewHandler),

		// fx.NopLogger,
	).Run()
}
