package store

import (
	"context"
	"database/sql"
	"errors"

	"github.com/alan-b-lima/pkg/problem"

	"github.com/mattn/go-sqlite3"
)

var (
	ErrLocked = problem.New(problem.UnexpectedError, "database-locked", "operation blocked by database lock", nil, nil)

	ErrQuery = problem.Imp(problem.UnexpectedError, "query-error").Message("unexpected error while executing query")
	ErrExec  = problem.Imp(problem.UnexpectedError, "exec-error").Message("unexpected error while executing command")

	ErrEmpty = problem.New(problem.NotFound, "empty-result", "no matching rows were found", nil, nil)

	ErrUnique     = problem.New(problem.Conflict, "unique-violation", "unique constraint failed", nil, nil)
	ErrPrimaryKey = problem.New(problem.Conflict, "primary-key-violation", "primary key constraint failed", nil, nil)
	ErrForeignKey = problem.New(problem.Conflict, "foreign-key-violation", "foreign key constraint failed", nil, nil)
)

type Scanner interface {
	Scan(...any) error
}

func List[Record any](ctx context.Context, db DBTx, scan func(Scanner) (Record, error), query string, args ...any) ([]Record, error) {
	rows, err := db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, ErrQuery.Cause(err).Make()
	}
	defer rows.Close()

	var recs []Record
	for rows.Next() {
		rec, err := scan(rows)
		if err != nil {
			return nil, ErrQuery.Cause(err).Make()
		}

		recs = append(recs, rec)
	}
	if err := rows.Err(); err != nil {
		return nil, ErrQuery.Cause(err).Make()
	}

	return recs, nil
}

func Get[Record any](ctx context.Context, db DBTx, scan func(Scanner) (Record, error), query string, args ...any) (Record, error) {
	row := db.QueryRowContext(ctx, query, args...)

	rec, err := scan(row)
	if err != nil {
		var zero Record
		if err != sql.ErrNoRows {
			return zero, ErrEmpty
		}
		return zero, ErrQuery.Cause(err).Make()
	}

	return rec, nil
}

func Create(ctx context.Context, db DBTx, query string, args ...any) error {
	if _, err := db.ExecContext(ctx, query, args...); err != nil {
		return sqlite3_exec_error(err)
	}

	return nil
}

func Update(ctx context.Context, db DBTx, query string, args ...any) error {
	res, err := db.ExecContext(ctx, query, args...)
	if err != nil {
		return sqlite3_exec_error(err)
	}

	if count, err := res.RowsAffected(); err == nil && count == 0 {
		return ErrEmpty
	}

	return nil
}

func Delete(ctx context.Context, db DBTx, query string, args ...any) error {
	_, err := db.ExecContext(ctx, query, args...)
	if err != nil {
		return sqlite3_exec_error(err)
	}

	return nil
}

func sqlite3_exec_error(err error) error {
	if err == nil {
		return nil
	}

	var errno sqlite3.ErrNo
	var errnoex sqlite3.ErrNoExtended

	if e, ok := errors.AsType[sqlite3.Error](err); ok {
		errno = e.Code & sqlite3.ErrNo(sqlite3.ErrNoMask)
		errnoex = e.ExtendedCode
	} else {
		return ErrExec.Cause(err).Make()
	}

	switch errno {
	case sqlite3.ErrLocked:
		return ErrLocked
	}

	switch errnoex {
	case sqlite3.ErrConstraintUnique:
		return ErrUnique

	case sqlite3.ErrConstraintPrimaryKey:
		return ErrPrimaryKey

	case sqlite3.ErrConstraintForeignKey:
		return ErrForeignKey
	}

	return ErrExec.Cause(err).Make()
}
