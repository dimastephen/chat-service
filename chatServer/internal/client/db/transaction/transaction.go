package transaction

import (
	"context"
	"github.com/dimastephen/chatServer/internal/client/db"
	"github.com/dimastephen/chatServer/internal/client/db/pg"
	"github.com/jackc/pgx/v4"
	"github.com/pkg/errors"
)

type manager struct {
	transactor db.Transactor
}

func NewTransactionManager(transactor db.Transactor) db.TxManager {
	return &manager{
		transactor: transactor,
	}
}

func (m *manager) transaction(ctx context.Context, opts pgx.TxOptions, fn db.Handler) (err error) {
	tx, ok := ctx.Value(pg.TxKey).(pgx.Tx)
	if ok {
		return fn(ctx)
	}
	tx, err = m.transactor.BeginTx(ctx, opts)
	if err != nil {
		errors.Wrap(err, "Cant begin transaction")
	}
	ctx = pg.MakeContextPg(ctx, opts)

	defer func() {
		if r := recover(); r != nil {
			err = errors.Errorf("panic recovered %v", r)
		}

		if err != nil {
			if errRollback := tx.Rollback(ctx); errRollback != nil {
				err = errors.Wrapf(err, "ErrRollback: %v", errRollback)
			}
			return
		}

		if nil == err {
			err = tx.Commit(ctx)
			if err != nil {
				err = errors.Wrapf(err, "tx commit failed")
			}
		}

	}()

	if err = fn(ctx); err != nil {
		err = errors.Wrap(err, "failed executing code inside transaction")
	}

	return err
}

func (m *manager) ReadCommited(ctx context.Context, f db.Handler) error {
	opts := pgx.TxOptions{IsoLevel: pgx.ReadCommitted}
	return m.transaction(ctx, opts, f)
}
