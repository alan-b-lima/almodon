package entrystore

import (
	"context"
	"database/sql"
	_ "embed"

	"github.com/alan-b-lima/almodon/internal/domain/entry"
	"github.com/alan-b-lima/almodon/internal/support/store"
	"github.com/alan-b-lima/almodon/pkg/uuid"
)

//go:embed sqlite.sql
var Script string

type SQLDB struct {
	db store.DBTx
}

var _ entry.Store = (*SQLDB)(nil)

func New(db *sql.DB) *SQLDB {
	return &SQLDB{db: db}
}

func (s *SQLDB) List(ctx context.Context, order uuid.UUID) ([]entry.Record, error) {
	rows, err := s.db.QueryContext(ctx, list, order)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var res []entry.Record
	for rows.Next() {
		rec, err := scan(rows)
		if err != nil {
			return nil, store.ErrQuery.Cause(err).Make()
		}

		res = append(res, rec)
	}
	if err := rows.Err(); err != nil {
		return nil, store.ErrQuery.Cause(err).Make()
	}

	return res, nil
}

func (s *SQLDB) Get(ctx context.Context, order uuid.UUID, item uuid.UUID) (entry.Record, error) {
	rec, err := scan(s.db.QueryRowContext(ctx, get, order, item))
	if err != nil {
		if err == sql.ErrNoRows {
			return entry.Record{}, entry.ErrNotFound
		}

		return entry.Record{}, store.ErrQuery.Cause(err).Make()
	}

	return rec, nil
}

func (s *SQLDB) Create(ctx context.Context, ent entry.Entity) error {
	_, err := s.db.ExecContext(ctx, create, ent.Order, ent.Item, ent.Amount)
	if err != nil {
		return store.ErrQuery.Cause(err).Make()
	}

	return nil
}

func (s *SQLDB) Update(ctx context.Context, order uuid.UUID, item uuid.UUID, amount int64) error {
	res, err := s.db.ExecContext(ctx, update, amount, order, item)
	if err != nil {
		return store.ErrQuery.Cause(err).Make()
	}

	if count, err := res.RowsAffected(); err == nil && count == 0 {
		return entry.ErrNotFound
	}

	return nil
}

func (s *SQLDB) Delete(ctx context.Context, order uuid.UUID, item uuid.UUID) error {
	_, err := s.db.ExecContext(ctx, delete, order, item)
	if err != nil {
		return store.ErrQuery.Cause(err).Make()
	}

	return nil
}

func scan(row store.Scanner) (entry.Record, error) {
	var rec entry.Record
	if err := row.Scan(
		&rec.Order,
		&rec.Item,
		&rec.Amount,
	); err != nil {
		return entry.Record{}, nil
	}

	return rec, nil
}

const (
	list = "select `order`, `item`, `amount` from Entries where `order` = ?"
	get  = "select `order`, `item`, `amount` from Entries where `order` = ? and `item` = ?"

	create = "insert into Entries (`order`, `item`, `amount`) values (?, ?, ?)"
	update = "update Entries set `amount` = ? where `order` = ? and `item` = ?"
	delete = "delete from Entries where `order` = ? and `item` = ?"
)
