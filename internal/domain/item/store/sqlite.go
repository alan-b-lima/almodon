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
	return store.List(ctx, s.db, scan, list)
}

func (s *SQLDB) ListByMaterial(ctx context.Context, material uuid.UUID) ([]item.Record, error) {
	return store.List(ctx, s.db, scan, list_by_material, material)
}

func (s *SQLDB) ListByECampus(ctx context.Context, ecampus int) ([]item.Record, error) {
	return store.List(ctx, s.db, scan, list_by_ecampus, ecampus)
}

func (s *SQLDB) ListByCATMAT(ctx context.Context, catmat int) ([]item.Record, error) {
	return store.List(ctx, s.db, scan, list_by_catmat, catmat)
}

func (s *SQLDB) ListBySIADS(ctx context.Context, siads int) ([]item.Record, error) {
	return store.List(ctx, s.db, scan, list_by_siads, siads)
}

func (s *SQLDB) ListByLot(ctx context.Context, lot uuid.UUID) ([]item.Record, error) {
	return store.List(ctx, s.db, scan, list_by_lot, lot)
}

func (s *SQLDB) Get(ctx context.Context, uuid uuid.UUID) (item.Record, error) {
	rec, err := store.Get(ctx, s.db, scan, get, uuid)
	if err != nil {
		if err == sql.ErrNoRows {
			return item.Record{}, item.ErrNotFound
		}

		return item.Record{}, err
	}

	return rec, nil
}

func (s *SQLDB) Create(ctx context.Context, ent item.Entity) error {
	return store.Create(ctx, s.db, create,
		ent.UUID,
		ent.Material,
		ent.Lot,
		ent.Available,
		ent.UnitCost.Cents(),
		ent.Expires,
		ent.Created,
		ent.Updated,
	)
}

func (s *SQLDB) Update(ctx context.Context, uuid uuid.UUID, ent item.UpdateEntity) error {
	return store.Update(ctx, s.db, update, ent.Amount, ent.Updated, uuid)
}

func (s *SQLDB) Patch(ctx context.Context, uuid uuid.UUID, ent item.PatchEntity) error {
	return store.Update(ctx, s.db, patch, ent.UnitCost, ent.Expires, ent.Updated, uuid)
}

func (s *SQLDB) Effective(ctx context.Context, ent item.EffectiveLot) error {
	return store.Update(ctx, s.db, effective, ent.Updated, ent.Lot)
}

func (s *SQLDB) Delete(ctx context.Context, uuid uuid.UUID) error {
	return store.Delete(ctx, s.db, delete, uuid)
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

const (
	list             = `select uuid, material, name, ecampus, catmat, siads, unit, lot, available, unit_cost, expires, created, updated from Items_View`
	list_by_material = list + ` where material = ?`
	list_by_ecampus  = list + ` where ecampus = ?`
	list_by_catmat   = list + ` where catmat = ?`
	list_by_siads    = list + ` where siads = ?`
	list_by_lot      = list + ` where lot = ?`

	get = list + ` where uuid = ?`

	create = `insert into Items (uuid, material, lot, available, unit_cost, expires, created, updated) values (?, ?, ?, ?, ?, ?, ?, ?)`

	update    = `update Items set amount = ?, updated = ? where uuid = ?`
	patch     = `update Items set unit_cost = coalesce(?, unit_cost), expires = coalesce(?, expires), updated = ? where uuid = ?`
	effective = `update Items set effective = true, updated = ? where lot = ?`

	delete = `delete from Items where uuid = ?`
)
