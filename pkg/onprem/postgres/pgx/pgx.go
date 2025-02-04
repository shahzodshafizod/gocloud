package pgx

import (
	"context"
	"os"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pkg/errors"
	"github.com/shahzodshafizod/gocloud/pkg"
)

type postgres struct {
	pool *pgxpool.Pool
}

func NewPostgres() (pkg.Postgres, error) {
	poolConfig, err := pgxpool.ParseConfig(os.Getenv("POSTGRES_DSN"))
	if err != nil {
		return nil, errors.Wrap(err, "pgxpool.ParseConfig")
	}

	poolConfig.MaxConns = 10

	pool, err := pgxpool.NewWithConfig(context.Background(), poolConfig)
	if err != nil {
		return nil, errors.Wrap(err, "pgxpool.NewWithConfig")
	}

	return &postgres{pool: pool}, nil
}

func (p *postgres) Close(context.Context) error {
	p.pool.Close()
	return nil
}

func (p *postgres) Begin(ctx context.Context) (pkg.Tx, error) {
	t, err := p.pool.Begin(ctx)
	return &tx{tx: t}, err
}

func (p *postgres) Exec(ctx context.Context, sql string, arguments ...any) error {
	res, err := p.pool.Exec(ctx, sql, arguments...)
	if err != nil {
		return p.ParseError(err)
	}
	if res.RowsAffected() == 0 {
		return pkg.ErrNoRowsAffected
	}
	return nil
}

func (p *postgres) QueryRow(ctx context.Context, sql string, args ...any) pkg.Row {
	return p.pool.QueryRow(ctx, sql, args...)
}

func (p *postgres) Query(ctx context.Context, sql string, args ...any) (pkg.Rows, error) {
	return p.pool.Query(ctx, sql, args...)
}

func (p *postgres) ParseError(err error) error {
	if pgErr, ok := err.(*pgconn.PgError); ok && pgErr.Code == "23505" {
		return pkg.ErrDuplicate
	} else if err == pgx.ErrNoRows {
		return pkg.ErrNoRows
	}
	return err
}

type tx struct {
	tx pgx.Tx
}

func (t *tx) Exec(ctx context.Context, sql string, args ...any) error {
	_, err := t.tx.Exec(ctx, sql, args...)
	return err
}

func (t *tx) QueryRow(ctx context.Context, sql string, args ...any) pkg.Row {
	return t.tx.QueryRow(ctx, sql, args...)
}

func (t *tx) Commit(ctx context.Context) error {
	return t.tx.Commit(ctx)
}

func (t *tx) Rollback(ctx context.Context) error {
	return t.tx.Rollback(ctx)
}
