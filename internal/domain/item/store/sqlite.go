package itemstore

import (
	"context"
	"database/sql"
	_ "embed"

	"github.com/alan-b-lima/almodon/internal/domain/item"
	"github.com/alan-b-lima/almodon/internal/support/store"
	"github.com/alan-b-lima/almodon/pkg/uuid"
)

//go:embed sqlite.sql
var Script string

type SQLDB struct {
	db store.DBTx
}

var _ item.Store = (*SQLDB)(nil)

func New(db *sql.DB) *SQLDB {
	return &SQLDB{db: db}
}

func (s *SQLDB) List(ctx context.Context) ([]item.Record, error) {
	rows, err := s.db.QueryContext(ctx, list)
	if err != nil {
		return nil, store.ErrQuery.Cause(err).Make()
	}
	defer rows.Close()

	return scan_list(rows)
}

func (s *SQLDB) ListByMaterial(ctx context.Context, uuid uuid.UUID) ([]item.Record, error) {
	rows, err := s.db.QueryContext(ctx, list_by_material, uuid)
	if err != nil {
		return nil, store.ErrQuery.Cause(err).Make()
	}
	defer rows.Close()

	return scan_list(rows)
}

func (s *SQLDB) ListByECampus(ctx context.Context, ecampus int) ([]item.Record, error) {
	rows, err := s.db.QueryContext(ctx, list_by_ecampus, ecampus)
	if err != nil {
		return nil, store.ErrQuery.Cause(err).Make()
	}
	defer rows.Close()

	return scan_list(rows)
}

func (s *SQLDB) ListByCATMAT(ctx context.Context, catmat int) ([]item.Record, error) {
	rows, err := s.db.QueryContext(ctx, list_by_catmat, catmat)
	if err != nil {
		return nil, store.ErrQuery.Cause(err).Make()
	}
	defer rows.Close()

	return scan_list(rows)
}

func (s *SQLDB) ListBySIADS(ctx context.Context, siads int) ([]item.Record, error) {
	rows, err := s.db.QueryContext(ctx, list_by_siads, siads)
	if err != nil {
		return nil, store.ErrQuery.Cause(err).Make()
	}
	defer rows.Close()

	return scan_list(rows)
}

func (s *SQLDB) ListByLot(ctx context.Context, lot uuid.UUID) ([]item.Record, error) {
	rows, err := s.db.QueryContext(ctx, list_by_lot, lot)
	if err != nil {
		return nil, store.ErrQuery.Cause(err).Make()
	}
	defer rows.Close()

	return scan_list(rows)
}

func (s *SQLDB) Get(ctx context.Context, uuid uuid.UUID) (item.Record, error) {
	row := s.db.QueryRowContext(ctx, get, uuid)

	ent, err := scan(row)
	if err != nil {
		if err == sql.ErrNoRows {
			return item.Record{}, item.ErrNotFound
		}

		return item.Record{}, store.ErrQuery.Cause(err).Make()
	}

	return ent, nil
}

func (s *SQLDB) Create(ctx context.Context, ent item.Entity) error {
	_, err := s.db.ExecContext(ctx, create,
		ent.UUID,
		ent.Material,
		ent.Lot,
		ent.Available,
		ent.UnitCost.Cents(),
		ent.Expires,
		ent.Created,
		ent.Updated,
	)
	if err != nil {
		return store.ErrExec.Cause(err).Make()
	}

	return nil
}

func (s *SQLDB) Update(ctx context.Context, uuid uuid.UUID, ent item.UpdateEntity) error {
	res, err := s.db.ExecContext(ctx, update, ent.Amount, ent.Updated, uuid)
	if err != nil {
		return store.ErrExec.Cause(err).Make()
	}

	if count, err := res.RowsAffected(); err == nil && count == 0 {
		return item.ErrNotFound
	}

	return nil
}

func (s *SQLDB) Patch(ctx context.Context, uuid uuid.UUID, ent item.PatchEntity) error {
	res, err := s.db.ExecContext(ctx, patch,
		ent.UnitCost,
		ent.Expires,
		ent.Updated,
		uuid,
	)
	if err != nil {
		return store.ErrExec.Cause(err).Make()
	}

	if count, err := res.RowsAffected(); err == nil && count == 0 {
		return item.ErrNotFound
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

func (s *SQLDB) RunTx(ctx context.Context, proc func(item.Store) error) error {
	return store.WithTx(ctx, s.db, func(tx store.DBTx) error {
		return proc(&SQLDB{db: tx})
	})
}

func (s *SQLDB) JoinTx(other store.Store) (item.Store, error) {
	tx, err := store.JoinTx(other, s.db)
	if err != nil {
		return nil, err
	}

	return &SQLDB{db: tx}, nil
}

func scan(scanner store.Scanner) (item.Record, error) {
	var rec item.Record
	if err := scanner.Scan(
		&rec.UUID,
		&rec.Material,
		&rec.Name,
		&rec.ECampus,
		&rec.CATMAT,
		&rec.SIADS,
		&rec.Unit,
		&rec.Lot,
		&rec.Available,
		(*int64)(&rec.UnitCost),
		&rec.Expires,
		&rec.Created,
		&rec.Updated,
	); err != nil {
		return item.Record{}, err
	}

	return rec, nil
}

func scan_list(rows *sql.Rows) ([]item.Record, error) {
	var recs []item.Record
	for rows.Next() {
		rec, err := scan(rows)
		if err != nil {
			return nil, store.ErrQuery.Cause(err).Make()
		}

		recs = append(recs, rec)
	}
	if err := rows.Err(); err != nil {
		return nil, store.ErrQuery.Cause(err).Make()
	}

	return recs, nil
}

const (
	list             = `select uuid, material, name, ecampus, catmat, siads, unit, lot, available, unit_cost, expires, created, updated from Items_View`
	list_by_material = list + ` where material = ?`
	list_by_ecampus  = list + ` where ecampus = ?`
	list_by_catmat   = list + ` where catmat = ?`
	list_by_siads    = list + ` where siads = ?`
	list_by_lot      = list + ` where lot = ?`

	get = list + ` where uuid = ?`

	create = `insert into Items (uuid, material, lot, available, unit_cost, expires, created, updated) values (?, ?, ?, ?, ?, ?, ?, ?)`

	update = `update Items set amount = ?, updated = ? where uuid = ?`
	patch  = `update Items set unit_cost = coalesce(?, unit_cost), expires = coalesce(?, expires), updated = ? where uuid = ?`

	delete = `delete from Items where uuid = ?`
)
