package db

import (
	"context"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
)

type Handler func(ctx context.Context) error

type Client interface {
	DB() DB
	Close()
}

type Transactor interface {
	BeginTx(ctx context.Context, txOptions pgx.TxOptions) (pgx.Tx, error)
}

type TxManager interface {
	ReadCommited(ctx context.Context, f Handler) error
}

type DB interface {
	Transactor
	SQLExecer
	Pinger
	Close()
}

type SQLExecer interface {
	NamedExecer
	QueryExecer
}

type Pinger interface {
	Ping(ctx context.Context) error
}

type NamedExecer interface {
	ScanOneContext(ctx context.Context, dest interface{}, q Query, args ...interface{}) error
	ScanAllContext(ctx context.Context, dest interface{}, q Query, args ...interface{}) error
}

type QueryExecer interface {
	ExecContext(ctx context.Context, q Query, args ...interface{}) (pgconn.CommandTag, error)
	QueryContext(ctx context.Context, q Query, args ...interface{}) (pgx.Rows, error)
	QueryRowContext(ctx context.Context, q Query, args ...interface{}) pgx.Row
}

type Query struct {
	Name     string
	QueryRaw string
}
