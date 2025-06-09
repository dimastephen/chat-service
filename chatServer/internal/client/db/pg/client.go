package pg

import (
	"context"
	"github.com/dimastephen/chatServer/internal/client/db"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/pkg/errors"
)

type DBClient struct {
	masterDBC db.DB
}

func New(ctx context.Context, dsn string) (db.Client, error) {
	dbc, err := pgxpool.Connect(ctx, dsn)
	if err != nil {
		return nil, errors.Errorf("Failed to connect db, %v", err)
	}
	return &DBClient{
		masterDBC: &pg{
			dbc: dbc,
		},
	}, nil
}

func (d *DBClient) DB() db.DB { return d.masterDBC }

func (d *DBClient) Close() {}
