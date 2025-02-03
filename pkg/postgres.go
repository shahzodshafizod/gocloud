//go:generate mockgen -source=postgres.go -package=mocks -destination=mocks/postgres.go --build_flags=--mod=mod github.com/jackc/pgx/v5
package pkg

import (
	"context"
)

type Postgres interface {
	// Closes the connection to the PostgreSQL database, returning an error if the operation fails.
	Close(ctx context.Context) error
	// Starts a new database transaction, returning a `Tx` object or an error.
	Begin(ctx context.Context) (Tx, error)
	// Executes a SQL query (e.g., INSERT, UPDATE, DELETE) with optional arguments, returning an error if the operation fails.
	Exec(ctx context.Context, sql string, args ...any) error
	// Executes a SQL query expected to return a single row, returning a `Row` object for result scanning.
	QueryRow(ctx context.Context, sql string, args ...any) Row
	// Executes a SQL query expected to return multiple rows, returning a `Rows` object or an error.
	Query(ctx context.Context, sql string, args ...any) (Rows, error)
}

type Row interface {
	Scan(dest ...any) error
}

type Rows interface {
	Close()
	Next() bool
	Scan(dest ...any) error
}

type Tx interface {
	Exec(ctx context.Context, sql string, args ...any) error
	QueryRow(ctx context.Context, sql string, args ...any) Row
	Commit(ctx context.Context) error
	Rollback(ctx context.Context) error
}
