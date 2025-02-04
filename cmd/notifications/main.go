package main

import (
	"os"

	"github.com/shahzodshafizod/gocloud/internal/migration"
	"github.com/shahzodshafizod/gocloud/internal/notifications"
	"github.com/shahzodshafizod/gocloud/pkg"
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
		fx.Provide(onprem.NewNoSQL),
		fx.Provide(onprem.NewNotification),
		fx.Provide(func() (pkg.Tracer, error) {
			return onprem.NewTracer(os.Getenv("SERVICE_NAME"))
		}),
		fx.Provide(func(tracer pkg.Tracer) (pkg.Queue, error) {
			return rabbitmq.NewQueue(os.Getenv("SERVICE_NAME"), tracer)
		}),

		fx.Provide(notifications.NewRepository),
		fx.Provide(notifications.NewService),
		fx.Invoke(notifications.NewHandler),

		fx.NopLogger,
	).Run()
}
