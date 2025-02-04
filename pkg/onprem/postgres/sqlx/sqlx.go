package sqlx

import (
	"context"
	"database/sql"
	"os"

	"github.com/jackc/pgx"
	"github.com/jackc/pgx/v5/pgconn"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"github.com/shahzodshafizod/gocloud/pkg"
)

type postgres struct {
	client *sqlx.DB
}

func NewPostgres() (pkg.Postgres, error) {
	dbClient, err := sqlx.Open("pgx", os.Getenv("POSTGRES_DSN"))
	if err != nil {
		return nil, errors.Wrap(err, "sqlx.Open")
	}
	return &postgres{client: dbClient}, nil
}

func (p *postgres) Close(context.Context) error {
	return p.client.Close()
}

func (p *postgres) Begin(ctx context.Context) (pkg.Tx, error) {
	t, err := p.client.Begin()
	return &tx{t}, err
}

func (p *postgres) Exec(ctx context.Context, sqlQuery string, arguments ...any) error {
	res, err := p.client.Exec(sqlQuery, arguments...)
	if err != nil {
		return p.parseError(err)
	}
	rows, _ := res.RowsAffected()
	if rows == 0 {
		return pkg.ErrNoRowsAffected
	}
	return nil
}

func (p *postgres) QueryRow(ctx context.Context, sqlQuery string, args ...any) pkg.Row {
	return p.client.QueryRow(sqlQuery, args...)
}

type rows struct {
	rows *sql.Rows
}

func (r *rows) Close()                 { r.rows.Close() }
func (r *rows) Next() bool             { return r.rows.Next() }
func (r *rows) Scan(dest ...any) error { return r.rows.Scan(dest...) }

func (p *postgres) Query(ctx context.Context, sqlQuery string, args ...any) (pkg.Rows, error) {
	rs, err := p.client.Query(sqlQuery, args...)
	if err != nil {
		return nil, p.parseError(err)
	}
	return &rows{rs}, nil
}

func (p *postgres) parseError(err error) error {
	if pgErr, ok := err.(*pgconn.PgError); ok && pgErr.Code == "23505" {
		return pkg.ErrDuplicate
	} else if err == pgx.ErrNoRows {
		return pkg.ErrNoRows
	}
	return err
}

type tx struct {
	tx *sql.Tx
}

func (t *tx) Exec(ctx context.Context, sql string, args ...any) error {
	_, err := t.tx.Exec(sql, args...)
	return err
}

func (t *tx) QueryRow(ctx context.Context, sql string, args ...any) pkg.Row {
	return t.tx.QueryRow(sql, args...)
}

func (t *tx) Commit(ctx context.Context) error {
	return t.tx.Commit()
}

func (t *tx) Rollback(ctx context.Context) error {
	return t.tx.Rollback()
}
