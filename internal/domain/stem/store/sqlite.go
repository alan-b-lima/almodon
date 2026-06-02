package stemstore

import (
	"context"
	"database/sql"
	_ "embed"

	"github.com/alan-b-lima/almodon/internal/domain/stem"
	"github.com/alan-b-lima/almodon/internal/support/store"
	"github.com/alan-b-lima/almodon/pkg/uuid"
)

//go:embed sqlite.sql
var Script string

type SQLDB struct {
	db store.DBTx
}

func New(db *sql.DB) *SQLDB {
	return &SQLDB{db: db}
}

var _ stem.Store = (*SQLDB)(nil)

func (s *SQLDB) List(ctx context.Context) ([]stem.Record, error) {
	return store.List(ctx, s.db, scan, list)
}

func (s *SQLDB) Get(ctx context.Context, uuid uuid.UUID) (stem.Record, error) {
	rec, err := store.Get(ctx, s.db, scan, get)
	if err == store.ErrEmpty {
		return stem.Record{}, stem.ErrNotFound
	}

	return rec, nil
}

func (s *SQLDB) GetByName(ctx context.Context, name string) (stem.Record, error) {
	rec, err := store.Get(ctx, s.db, scan, get_by_name, name)
	if err == store.ErrEmpty {
		return stem.Record{}, stem.ErrNotFound
	}

	return rec, nil
}

func (s *SQLDB) Create(ctx context.Context, ent stem.Entity) error {
	return store.Create(ctx, s.db, create,
		ent.UUID,
		nil,
		ent.Name,
		ent.Created,
	)
}

func (s *SQLDB) Rename(ctx context.Context, uuid uuid.UUID, title string) error {
	err := store.Update(ctx, s.db, rename, title, uuid)
	if err == store.ErrEmpty {
		return stem.ErrNotFound
	}

	return err
}

func (s *SQLDB) Upgrade(ctx context.Context, uuid uuid.UUID, order uuid.UUID) error {
	err := store.Update(ctx, s.db, upgrade, order, uuid)
	if err == store.ErrEmpty {
		return stem.ErrNotFound
	}

	return err
}

func (s *SQLDB) Delete(ctx context.Context, uuid uuid.UUID) error {
	return store.Delete(ctx, s.db, delete, uuid)
}

func (s *SQLDB) Tx() store.DBTx {
	return s.db
}

func (s *SQLDB) RunTx(ctx context.Context, proc func(stem.Store) error) error {
	return store.WithTx(ctx, s.db, func(tx store.DBTx) error {
		return proc(&SQLDB{db: tx})
	})
}

func (s *SQLDB) JoinTx(other store.Store) (stem.Store, error) {
	tx, err := store.JoinTx(other, s.db)
	if err != nil {
		return nil, err
	}

	return &SQLDB{db: tx}, nil
}

func scan(rows store.Scanner) (stem.Record, error) {
	var rec stem.Record
	if err := rows.Scan(
		&rec.UUID,
		&rec.Bloom,
		&rec.Name,
		&rec.Version,
		&rec.Created,
		&rec.Updated,
	); err != nil {
		return stem.Record{}, err
	}

	return rec, nil
}

const (
	list        = `select uuid, bloom, name, version, created, updated from Stems_View`
	get         = list + ` where uuid = ?`
	get_by_name = list + ` where uuid = ?`

	create = `insert into Stems (uuid, bloom, name, created) values (?, ?, ?, ?)`

	rename  = `update Stems set name = ? where uuid = ?`
	upgrade = `update Stems set bloom = ? where uuid = ?`

	delete = `delete from Stems where uuid = ?`
)
