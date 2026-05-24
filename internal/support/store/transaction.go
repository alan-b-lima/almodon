package store

import (
	"context"
	"database/sql"

	"github.com/alan-b-lima/almodon/internal/support"
	"github.com/alan-b-lima/pkg/problem"
)

var (
	ErrTx       = problem.Imp(problem.UnexpectedError, "transaction-error").Message("unexpected error while processing transaction")
	ErrNestedTx = problem.New(problem.UnexpectedError, "open-transaction-inside-transaction", "cannot open a transaction inside another transaction", nil, nil)

	ErrNotExtendable = problem.New(problem.UnexpectedError, "not-extendable", "store does not support transaction extension", nil, nil)
	ErrNotInTx       = problem.New(problem.UnexpectedError, "not-transaction", "store is not in a transaction", nil, nil)
	ErrIllegalJoin   = problem.New(problem.UnexpectedError, "illegal-join", "cannot join transactions from different pools", nil, nil)
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

type txconn struct {
	*sql.Tx
	pool *sql.DB
}

func WithTx(ctx context.Context, dbtx DBTx, proc func(DBTx) error) error {
	db, ok := dbtx.(*sql.DB)
	if !ok {
		if tx, ok := dbtx.(*txconn); ok {
			return proc(tx)
		}

		return support.ErrUnreacheable
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

func JoinTx(txed Store, joiner DBTx) (DBTx, error) {
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
