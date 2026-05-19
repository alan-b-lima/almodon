package store

import (
	"context"
	"database/sql"
)

type Store any

type DBTx interface {
	ExecContext(context.Context, string, ...any) (sql.Result, error)

	QueryContext(context.Context, string, ...any) (*sql.Rows, error)
	QueryRowContext(context.Context, string, ...any) *sql.Row
}

type Txer interface {
	Tx() DBTx
}

type Scanner interface {
	Scan(...any) error
}

type txconn struct {
	*sql.Tx
	pool *sql.DB
}

func WithTx(ctx context.Context, dbtx DBTx, proc func(DBTx) error) error {
	db, ok := dbtx.(*sql.DB)
	if !ok {
		return ErrNestedTx
	}

	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return ErrTx.Cause(err).Make()
	}

	if err := proc(&txconn{Tx: tx, pool: db}); err != nil {
		if err := tx.Rollback(); err != nil {
			return ErrTx.Cause(err).Make()
		}
		return err
	}

	if err := tx.Commit(); err != nil {
		return ErrTx.Cause(err).Make()
	}
	return nil
}

func JoinTx(txed any, joiner DBTx) (DBTx, error) {
	txer, ok := txed.(Txer)
	if !ok {
		return nil, ErrNotExtendable
	}

	tx, ok := txer.Tx().(*txconn)
	if !ok {
		return nil, ErrNotInTx
	}

	db, ok := joiner.(*sql.DB)
	if !ok {
		return nil, ErrNestedTx
	}

	if db != tx.pool {
		return nil, ErrIllegalJoin
	}

	return tx, nil
}
