package materialstore

import (
	"context"
	"database/sql"
	_ "embed"

	"github.com/alan-b-lima/almodon/internal/domain/material"
	"github.com/alan-b-lima/almodon/internal/support/service"
	"github.com/alan-b-lima/almodon/internal/support/store"
	"github.com/alan-b-lima/almodon/pkg/uuid"
)

//go:embed sqlite.sql
var Script string

type SQLDB struct {
	db store.DBTx
}

var _ material.Store = (*SQLDB)(nil)

func New(db store.DBTx) *SQLDB {
	return &SQLDB{db: db}
}

func (s *SQLDB) List(ctx context.Context) ([]material.Record, error) {
	return s.list(ctx, list)
}

func (s *SQLDB) ListByECampus(ctx context.Context, ecampus int) ([]material.Record, error) {
	return s.list(ctx, list_by_ecampus, ecampus)
}

func (s *SQLDB) ListByCATMAT(ctx context.Context, catmat int) ([]material.Record, error) {
	return s.list(ctx, list_by_catmat, catmat)
}

func (s *SQLDB) ListBySIADS(ctx context.Context, siads int) ([]material.Record, error) {
	return s.list(ctx, list_by_siads, siads)
}

func (s *SQLDB) list(ctx context.Context, query string, args ...any) ([]material.Record, error) {
	rows, err := s.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, store.ErrQuery.Cause(err).Make()
	}
	defer rows.Close()

	var ents []material.Record
	for rows.Next() {
		var ent material.Record
		if ok, err := scan(&ent, rows); err != nil {
			if ok {
				return nil, err
			}

			return nil, store.ErrQuery.Cause(err).Make()
		}

		ents = append(ents, ent)
	}
	if err := rows.Err(); err != nil {
		return nil, store.ErrQuery.Cause(err).Make()
	}

	return ents, nil
}

func (s *SQLDB) Get(ctx context.Context, uuid uuid.UUID) (material.Record, error) {
	row := s.db.QueryRowContext(ctx, get, uuid.Bytes())

	var ent material.Record
	if ok, err := scan(&ent, row); err != nil {
		if ok {
			return material.Record{}, err
		}

		if err == sql.ErrNoRows {
			return material.Record{}, material.ErrNotFound
		}
		return material.Record{}, store.ErrQuery.Cause(err).Make()
	}

	return ent, nil
}

func (s *SQLDB) Create(ctx context.Context, rec material.Entity) error {
	_, err := s.db.ExecContext(ctx, create,
		rec.UUID.Bytes(),
		rec.Name,
		rec.ECampus,
		rec.CATMAT,
		rec.SIADS,
		rec.Descript,
		rec.Unit,
		rec.Min,
		rec.Created,
		rec.Updated,
	)
	if err != nil {
		return store.ErrExec.Cause(err).Make()
	}

	return nil
}

func (s *SQLDB) Patch(ctx context.Context, uuid uuid.UUID, rec material.PatchEntity) error {
	_, err := s.db.ExecContext(ctx, patch,
		store.NoneNil(rec.Name),
		store.NoneNil(rec.ECampus),
		store.NoneNil(rec.CATMAT),
		store.NoneNil(rec.SIADS),
		store.NoneNil(rec.Descript),
		store.NoneNil(rec.Unit),
		store.NoneNil(rec.Min),
		rec.Updated,
		uuid.Bytes(),
	)
	if err != nil {
		return store.ErrExec.Cause(err).Make()
	}

	return nil
}

func (s *SQLDB) Delete(ctx context.Context, uuid uuid.UUID) error {
	_, err := s.db.ExecContext(ctx, delete, uuid.Bytes())
	if err != nil {
		return store.ErrExec.Cause(err).Make()
	}

	return nil
}

func scan(ent *material.Record, scanner interface{ Scan(...any) error }) (bool, error) {
	var bytes []byte

	if err := scanner.Scan(
		&bytes,
		&ent.Name,
		&ent.ECampus,
		&ent.CATMAT,
		&ent.SIADS,
		&ent.Descript,
		&ent.Unit,
		&ent.Min,
		&ent.Created,
		&ent.Updated,
	); err != nil {
		return false, err
	}

	if err := service.Set(&ent.UUID, bytes, uuid.FromBytes); err != nil {
		return true, store.ErrQuery.Cause(err).Make()
	}

	return true, nil
}

const (
	list            = `select uuid, name, ecampus, catmat, siads, descript, unit, min, created, updated from Materials`
	list_by_ecampus = `select uuid, name, ecampus, catmat, siads, descript, unit, min, created, updated from Materials where ecampus = ?`
	list_by_catmat  = `select uuid, name, ecampus, catmat, siads, descript, unit, min, created, updated from Materials where catmat = ?`
	list_by_siads   = `select uuid, name, ecampus, catmat, siads, descript, unit, min, created, updated from Materials where siads = ?`
	get             = `select uuid, name, ecampus, catmat, siads, descript, unit, min, created, updated from Materials where uuid = ?`
	create          = `insert into Materials (uuid, name, ecampus, catmat, siads, descript, unit, min, created, updated) values (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
	patch           = `update Materials set name = coalesce(?, name), ecampus = coalesce(?, ecampus), catmat = coalesce(?, catmat), siads = coalesce(?, siads), descript = coalesce(?, descript), unit = coalesce(?, unit), min = coalesce(?, min), updated = ? where uuid = ?`
	delete          = `delete from Materials where uuid = ?`
)
