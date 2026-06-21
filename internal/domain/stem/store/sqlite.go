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

func (s *SQLDB) Get(ctx context.Context) (stem.Record, error) {
	rec, err := store.Get(ctx, s.db, scan, get)
	if err == store.ErrEmpty {
		return stem.Record{}, stem.ErrNotFound
	}

	return rec, nil
}

func (s *SQLDB) Create(ctx context.Context, ent stem.Entity) error {
	return store.Create(ctx, s.db, create, ent.UUID, nil, ent.Created)
}

func (s *SQLDB) Upgrade(ctx context.Context, order uuid.UUID) error {
	err := store.Update(ctx, s.db, upgrade, order)
	if err == store.ErrEmpty {
		return stem.ErrNotFound
	}

	return err
}

func (s *SQLDB) Delete(ctx context.Context) error {
	return store.Delete(ctx, s.db, delete)
}

func (s *SQLDB) Tx() store.DBTx {
	return s.db
}

func (s *SQLDB) RunTx(ctx context.Context, proc func(context.Context, stem.Store) error) error {
	return store.WithTx(ctx, s.db, func(tx store.DBTx) error {
		return proc(ctx, &SQLDB{db: tx})
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
		&rec.Version,
		&rec.Created,
		&rec.Updated,
	); err != nil {
		return stem.Record{}, err
	}

	return rec, nil
}

const (
	get     = `select uuid, bloom, version, created, updated from Stem_View`
	create  = `insert into Stem (uuid, bloom, created) values (?, ?, ?)`
	upgrade = `update Stem set bloom = ?`
	delete  = `delete from Stem where`
)
