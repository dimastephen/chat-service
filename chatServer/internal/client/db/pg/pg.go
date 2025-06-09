package pg

import (
	"context"
	"fmt"
	"github.com/dimastephen/chatServer/internal/client/db"
	"github.com/dimastephen/chatServer/internal/client/db/decorate"
	"github.com/georgysavva/scany/pgxscan"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"log"
)

type key string

const TxKey key = "tx"

type pg struct {
	dbc *pgxpool.Pool
}

func NewDB(dbc *pgxpool.Pool) db.DB {
	return &pg{
		dbc: dbc,
	}
}

func (p *pg) ScanOneContext(ctx context.Context, dest interface{}, q db.Query, args ...interface{}) error {
	logQuery(ctx, q, args...)

	row, err := p.QueryContext(ctx, q, args...)
	if err != nil {
		return err
	}

	return pgxscan.ScanOne(&dest, row)
}

func (p *pg) ScanAllContext(ctx context.Context, dest interface{}, q db.Query, args ...interface{}) error {
	logQuery(ctx, q, args)

	rows, err := p.QueryContext(ctx, q, args...)
	if err != nil {
		return err
	}

	return pgxscan.ScanAll(dest, rows)
}

func (p *pg) ExecContext(ctx context.Context, q db.Query, args ...interface{}) (pgconn.CommandTag, error) {
	logQuery(ctx, q, args)

	if tx, ok := ctx.Value(TxKey).(pgx.Tx); ok {
		return tx.Exec(ctx, q.QueryRaw, args...)
	}

	return p.dbc.Exec(ctx, q.QueryRaw, args...)
}

func (p *pg) QueryContext(ctx context.Context, q db.Query, args ...interface{}) (pgx.Rows, error) {
	logQuery(ctx, q, args)

	if tx, ok := ctx.Value(TxKey).(pgx.Tx); ok {
		return tx.Query(ctx, q.QueryRaw, args...)
	}

	return p.dbc.Query(ctx, q.QueryRaw, args...)
}

func (p *pg) QueryRowContext(ctx context.Context, q db.Query, args ...interface{}) pgx.Row {
	logQuery(ctx, q, args...)

	if tx, ok := ctx.Value(TxKey).(pgx.Tx); ok {
		return tx.QueryRow(ctx, q.QueryRaw, args...)
	}

	return p.dbc.QueryRow(ctx, q.QueryRaw, args...)
}

func (p *pg) Ping(ctx context.Context) error {
	err := p.dbc.Ping(ctx)
	if err != nil {
		return err
	}
	return nil
}

func (p *pg) Close() {
	p.dbc.Close()
}

func MakeContextPg(ctx context.Context, opts pgx.TxOptions) context.Context {
	return context.WithValue(ctx, TxKey, "tx")
}

func (p *pg) BeginTx(ctx context.Context, options pgx.TxOptions) (pgx.Tx, error) {
	return p.dbc.BeginTx(ctx, options)
}

func logQuery(ctx context.Context, q db.Query, args ...interface{}) {
	decoratedQuery := decorate.Decorate(q.QueryRaw, decorate.PlaceholderDollar, args...)
	log.Println(ctx, fmt.Sprintf("sql: %s", q.Name), fmt.Sprintf("query: %s", decoratedQuery))
}
