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

var _ stem.Store = (*SQLDB)(nil)

func (s *SQLDB) List(ctx context.Context) ([]stem.Record, error) {
	rows, err := s.db.QueryContext(ctx, list)
	if err != nil {
		return nil, store.ErrQuery.Cause(err).Make()
	}
	defer rows.Close()

	var recs []stem.Record
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

func (s *SQLDB) Get(ctx context.Context, uuid uuid.UUID) (stem.Record, error) {
	row := s.db.QueryRowContext(ctx, get)

	rec, err := scan(row)
	if err != nil {
		if err == sql.ErrNoRows {
			return stem.Record{}, stem.ErrNotFound
		}

		return stem.Record{}, store.ErrQuery.Cause(err).Make()
	}

	return rec, nil
}

func (s *SQLDB) GetByName(ctx context.Context, name string) (stem.Record, error) {
	row := s.db.QueryRowContext(ctx, get_by_name, name)

	rec, err := scan(row)
	if err != nil {
		if err == sql.ErrNoRows {
			return stem.Record{}, stem.ErrNotFound
		}

		return stem.Record{}, store.ErrQuery.Cause(err).Make()
	}

	return rec, nil
}

func (s *SQLDB) Create(ctx context.Context, ent stem.Entity) error {
	if _, err := s.db.ExecContext(ctx, create, ent.UUID, nil, ent.Name, ent.Created); err != nil {
		return store.ErrExec.Cause(err).Make()
	}

	return nil
}

func (s *SQLDB) Rename(ctx context.Context, uuid uuid.UUID, title string) error {
	res, err := s.db.ExecContext(ctx, rename, title, uuid)
	if err != nil {
		return store.ErrExec.Cause(err).Make()
	}

	if count, err := res.RowsAffected(); err == nil && count == 0 {
		return stem.ErrNotFound
	}

	return nil
}

func (s *SQLDB) Upgrade(ctx context.Context, uuid uuid.UUID, order uuid.UUID) error {
	res, err := s.db.ExecContext(ctx, upgrade, order, uuid)
	if err != nil {
		return store.ErrExec.Cause(err).Make()
	}

	if count, err := res.RowsAffected(); err == nil && count == 0 {
		return stem.ErrNotFound
	}

	return nil
}

func (s *SQLDB) Delete(ctx context.Context, uuid uuid.UUID) error {
	if _, err := s.db.ExecContext(ctx, delete, uuid); err != nil {
		return store.ErrExec.Cause(err).Make()
	}

	return nil
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
	if err := rows.Scan(&rec.UUID, &rec.Bloom, &rec.Name, &rec.Version, &rec.Created, &rec.Updated); err != nil {
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
