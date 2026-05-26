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
	return store.List(ctx, s.db, scan, list, order)
}

func (s *SQLDB) Get(ctx context.Context, order uuid.UUID, item uuid.UUID) (entry.Record, error) {
	rec, err := store.Get(ctx, s.db, scan, get, order, item)
	if err == store.ErrEmpty {
		return entry.Record{}, entry.ErrNotFound
	}

	return rec, err
}

func (s *SQLDB) Create(ctx context.Context, order uuid.UUID, ent entry.Entity) error {
	return store.Create(ctx, s.db, create, order, ent.Item, ent.Amount)
}

func (s *SQLDB) Tx() store.DBTx {
	return s.db
}

func (s *SQLDB) RunTx(ctx context.Context, proc func(entry.Store) error) error {
	return store.WithTx(ctx, s.db, func(tx store.DBTx) error {
		return proc(&SQLDB{db: tx})
	})
}

func (s *SQLDB) JoinTx(other store.Store) (entry.Store, error) {
	tx, err := store.JoinTx(other, s.db)
	if err != nil {
		return nil, err
	}

	return &SQLDB{db: tx}, nil
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
)
