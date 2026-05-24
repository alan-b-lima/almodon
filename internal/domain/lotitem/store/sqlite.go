package lotitemstore

import (
	"context"
	"database/sql"
	_ "embed"

	"github.com/alan-b-lima/almodon/internal/domain/lotitem"
	"github.com/alan-b-lima/almodon/internal/support/store"
	"github.com/alan-b-lima/almodon/pkg/money"
	"github.com/alan-b-lima/almodon/pkg/uuid"
)

//go:embed sqlite.sql
var Script string

type SQLDB struct {
	db store.DBTx
}

var _ lotitem.Store = (*SQLDB)(nil)

func New(db *sql.DB) *SQLDB {
	return &SQLDB{db: db}
}

func (s *SQLDB) List(ctx context.Context, lot uuid.UUID) ([]lotitem.Record, error) {
	return store.List(ctx, s.db, scan, list)
}

func (s *SQLDB) Get(ctx context.Context, uuid uuid.UUID) (lotitem.Record, error) {
	rec, err := store.Get(ctx, s.db, scan, get, uuid)
	switch err {
	case store.ErrEmpty:
		return lotitem.Record{}, lotitem.ErrNotFound
	}

	return rec, err
}

func (s *SQLDB) Create(ctx context.Context, ent lotitem.Entity) error {
	err := store.Create(ctx, s.db, create,
		ent.UUID,
		ent.Material,
		ent.UnitCost.Cents(),
		ent.Expires,
		ent.Created,
		ent.Updated,
	)

	return err
}

func (s *SQLDB) Patch(ctx context.Context, uuid uuid.UUID, ent lotitem.PatchEntity) error {
	err := store.Update(ctx, s.db, patch, uuid, ent.UnitCost, ent.Expires, ent.Updated, uuid)
	switch err {
	case store.ErrEmpty:
		return lotitem.ErrNotFound
	}

	return err
}

func (s *SQLDB) Delete(ctx context.Context, uuid uuid.UUID) error {
	return store.Delete(ctx, s.db, delete, uuid)
}

func (s *SQLDB) Tx() store.DBTx {
	return s.db
}

func (s *SQLDB) RunTx(ctx context.Context, proc func(lotitem.Store) error) error {
	return store.WithTx(ctx, s.db, func(tx store.DBTx) error {
		return proc(&SQLDB{db: tx})
	})
}

func (s *SQLDB) JoinTx(other store.Store) (lotitem.Store, error) {
	tx, err := store.JoinTx(other, s.db)
	if err != nil {
		return nil, err
	}

	return &SQLDB{db: tx}, nil
}

func scan(scanner store.Scanner) (lotitem.Record, error) {
	var cents int64

	var rec lotitem.Record
	if err := scanner.Scan(
		&rec.UUID,
		&rec.Lot,
		&rec.Material,
		&rec.Name,
		&rec.Unit,
		&rec.Amount,
		&cents,
		&rec.Expires,
		&rec.Created,
		&rec.Updated,
	); err != nil {
		return lotitem.Record{}, err
	}

	rec.UnitCost = money.FromCents(cents)
	return rec, nil
}

const (
	list = `select uuid, lot, material, name, unit, amount, unit_cost, expires, created, updated from LotItems_View where lot = ?`
	get  = `select uuid, lot, material, name, unit, amount, unit_cost, expires, created, updated from LotItems_View where uuid = ?`

	create = `insert into LotItems (uuid, lot, material, amount, unit_cost, expires, created, updated) values (?, ?, ?, ?, ?, ?, ?, ?)`
	patch  = `update LotItems set material = coalesce(?, material), amount = coalesce(?, amount), unit_cost = coalesce(?, unit_cost), expires = coalesce(?, expires), updated = ? where uuid = ?`
	delete = `delete from LotItems where uuid = ?`
)
