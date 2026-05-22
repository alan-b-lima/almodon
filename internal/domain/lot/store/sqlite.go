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
	rows, err := s.db.QueryContext(ctx, list)
	if err != nil {
		return nil, store.ErrQuery.Cause(err).Make()
	}
	defer rows.Close()

	return scan_list(rows)
}

func (s *SQLDB) ListByState(ctx context.Context, state lot.State) ([]lot.Record, error) {
	var query string
	switch state {
	case lot.Open:
		query = list_open
	case lot.Closed:
		query = list_closed
	default:
		return nil, support.ErrTODO
	}

	rows, err := s.db.QueryContext(ctx, query)
	if err != nil {
		return nil, store.ErrQuery.Cause(err).Make()
	}
	defer rows.Close()

	return scan_list(rows)
}

func (s *SQLDB) Get(ctx context.Context, uuid uuid.UUID) (lot.Record, error) {
	row := s.db.QueryRowContext(ctx, get, uuid)

	ent, err := scan(row)
	if err != nil {
		if err == sql.ErrNoRows {
			return lot.Record{}, lot.ErrNotFound
		}

		return lot.Record{}, store.ErrQuery.Cause(err).Make()
	}

	return ent, nil
}

func (s *SQLDB) Create(ctx context.Context, ent lot.Entity) error {
	_, err := s.db.ExecContext(ctx, create,
		ent.UUID,
		ent.Supplier,
		ent.Arrival,
		ent.Note,
		ent.Created,
		ent.Updated,
	)
	if err != nil {
		return store.ErrExec.Cause(err).Make()
	}

	return nil
}

func (s *SQLDB) Patch(ctx context.Context, uuid uuid.UUID, ent lot.PatchEntity) error {
	res, err := s.db.ExecContext(ctx, patch,
		ent.Supplier.Interface(),
		ent.Arrival.Interface(),
		ent.Note.Interface(),
		ent.Updated,
		uuid,
	)
	if err != nil {
		return store.ErrExec.Cause(err).Make()
	}

	if count, err := res.RowsAffected(); err == nil && count == 0 {
		return lot.ErrNotFound
	}

	return nil
}

func (s *SQLDB) Sign(ctx context.Context, uuid uuid.UUID, order uuid.UUID) error {
	res, err := s.db.ExecContext(ctx, sign, order, uuid)
	if err != nil {
		return store.ErrExec.Cause(err).Make()
	}

	if count, err := res.RowsAffected(); err == nil && count == 0 {
		return lot.ErrNotFound
	}

	return nil
}

func (s *SQLDB) Delete(ctx context.Context, uuid uuid.UUID) error {
	_, err := s.db.ExecContext(ctx, delete, uuid)
	if err != nil {
		return store.ErrExec.Cause(err).Make()
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
	if err := scanner.Scan(); err != nil {
		return lot.Record{}, err
	}

	return rec, nil
}

func scan_list(rows *sql.Rows) ([]lot.Record, error) {
	var recs []lot.Record
	for rows.Next() {
		ent, err := scan(rows)
		if err != nil {
			return nil, store.ErrQuery.Cause(err).Make()
		}

		recs = append(recs, ent)
	}
	if err := rows.Err(); err != nil {
		return nil, store.ErrQuery.Cause(err).Make()
	}

	return recs, nil
}

const (
	list        = "select `uuid`, `order`, `supplier`, `arrival`, `note`, `created`, `updated` from `Lots`"
	list_open   = "select `uuid`, `order`, `supplier`, `arrival`, `note`, `created`, `updated` from `Lots` where `order` is null"
	list_closed = "select `uuid`, `order`, `supplier`, `arrival`, `note`, `created`, `updated` from `Lots` where `order` is not null"

	get = "select `uuid`, `order`, `supplier`, `arrival`, `note`, `created`, `updated` from `Lots` where `uuid` = ?"

	create = "insert into `Lots` (`uuid`, `order`, `supplier`, `arrival`, `note`, `created`, `updated`) values (?, null, ?, ?, ?, ?, ?)"

	patch = "update `Lots` set `supplier` = coalesce(?, `supplier`), `arrival` = coalesce(?, `arrival`), `note` = coalesce(?, `note`), `updated` = ? where `uuid` = ?"
	sign  = "update `Lots` set `order` = ? where `uuid` = ? and `order` is null"

	delete = "delete from `Lots` where `uuid` = ?"
)
