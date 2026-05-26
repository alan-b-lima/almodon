package lotstore

import (
	"context"
	"database/sql"
	_ "embed"

	"github.com/alan-b-lima/almodon/internal/domain/lot"
	"github.com/alan-b-lima/almodon/internal/support"
	"github.com/alan-b-lima/almodon/internal/support/store"
	"github.com/alan-b-lima/almodon/pkg/uuid"
)

//go:embed sqlite.sql
var Script string

type SQLDB struct {
	db store.DBTx
}

var _ lot.Store = (*SQLDB)(nil)

func New(db *sql.DB) *SQLDB {
	return &SQLDB{db: db}
}

func (s *SQLDB) List(ctx context.Context) ([]lot.Record, error) {
	return store.List(ctx, s.db, scan, list)
}

func (s *SQLDB) ListByState(ctx context.Context, state lot.State) ([]lot.Record, error) {
	switch state {
	case lot.Open:
		return store.List(ctx, s.db, scan, list_open)

	case lot.Closed:
		return store.List(ctx, s.db, scan, list_closed)
	}

	return nil, support.ErrTODO
}

func (s *SQLDB) Get(ctx context.Context, uuid uuid.UUID) (lot.Record, error) {
	return store.Get(ctx, s.db, scan, get, uuid)
}

func (s *SQLDB) Create(ctx context.Context, ent lot.Entity) error {
	return store.Create(ctx, s.db, create,
		ent.UUID,
		ent.Supplier,
		ent.Author,
		ent.Arrival,
		ent.Note,
		ent.Created,
		ent.Updated,
	)
}

func (s *SQLDB) Patch(ctx context.Context, uuid uuid.UUID, ent lot.PatchEntity) error {
	return store.Update(ctx, s.db, patch,
		ent.Supplier.Interface(),
		ent.Arrival.Interface(),
		ent.Note.Interface(),
		ent.Updated,
		uuid,
	)
}

func (s *SQLDB) Sign(ctx context.Context, uuid uuid.UUID, ent lot.SignEntity) error {
	return store.Update(ctx, s.db, sign, ent.Order, ent.Updated, uuid)
}

func (s *SQLDB) Delete(ctx context.Context, uuid uuid.UUID) error {
	return store.Delete(ctx, s.db, delete, uuid)
}

func (s *SQLDB) Mutable(ctx context.Context, uuid uuid.UUID) error {
	rec, err := s.Get(ctx, uuid)
	if err != nil {
		if err == lot.ErrNotFound {
			return nil
		}
		return err
	}

	if !rec.Order.IsNil() {
		return lot.ErrModifySigned
	}

	return nil
}

func (s *SQLDB) Tx() store.DBTx {
	return s.db
}

func (s *SQLDB) RunTx(ctx context.Context, proc func(lot.Store) error) error {
	return store.WithTx(ctx, s.db, func(tx store.DBTx) error {
		return proc(&SQLDB{db: tx})
	})
}

func (s *SQLDB) JoinTx(other store.Store) (lot.Store, error) {
	tx, err := store.JoinTx(other, s.db)
	if err != nil {
		return nil, err
	}

	return &SQLDB{db: tx}, nil
}

func scan(scanner store.Scanner) (lot.Record, error) {
	var rec lot.Record
	if err := scanner.Scan(
		rec.UUID,
		rec.Order,
		rec.Supplier,
		rec.Author,
		rec.Arrival,
		rec.Note,
		rec.Created,
		rec.Updated,
	); err != nil {
		return lot.Record{}, err
	}

	return rec, nil
}

const (
	list        = "select `uuid`, `order`, `supplier`, `author`, `arrival`, `note`, `created`, `updated` from `Lots`"
	list_open   = "select `uuid`, `order`, `supplier`, `author`, `arrival`, `note`, `created`, `updated` from `Lots` where `order` is null"
	list_closed = "select `uuid`, `order`, `supplier`, `author`, `arrival`, `note`, `created`, `updated` from `Lots` where `order` is not null"

	get = "select `uuid`, `order`, `supplier`, `author`, `arrival`, `note`, `created`, `updated` from `Lots` where `uuid` = ?"

	create = "insert into `Lots` (`uuid`, `order`, `supplier`, `author`, `arrival`, `note`, `created`, `updated`) values (?, null, ?, ?, ?, ?, ?)"

	patch = "update `Lots` set `supplier` = coalesce(?, `supplier`), `arrival` = coalesce(?, `arrival`), `note` = coalesce(?, `note`), `updated` = ? where `uuid` = ?"
	sign  = "update `Lots` set `order` = ?, `updated` = ? where `uuid` = ? and `order` is null"

	delete = "delete from `Lots` where `uuid` = ?"
)
