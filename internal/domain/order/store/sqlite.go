package orderstore

import (
	"context"
	"database/sql"
	_ "embed"

	"github.com/alan-b-lima/almodon/internal/domain/order"
	"github.com/alan-b-lima/almodon/internal/support/store"
	"github.com/alan-b-lima/almodon/pkg/uuid"
)

//go:embed sqlite.sql
var Script string

type SQLDB struct {
	db store.DBTx
}

var _ order.Store = (*SQLDB)(nil)

func New(db *sql.DB) *SQLDB {
	return &SQLDB{db: db}
}

func (s *SQLDB) ListBlooms(ctx context.Context) ([]order.Record, error) {
	return store.List(ctx, s.db, scan, list_blooms)
}

func (s *SQLDB) ListByStem(ctx context.Context, stem uuid.UUID) ([]order.Record, error) {
	return store.List(ctx, s.db, scan, list_by_stem, stem)
}

func (s *SQLDB) Get(ctx context.Context, uuid uuid.UUID) (order.Record, error) {
	rec, err := store.Get(ctx, s.db, scan, get, uuid)
	if err == store.ErrEmpty {
		return order.Record{}, order.ErrNotFound
	}

	return rec, nil
}

func (s *SQLDB) GetBloom(ctx context.Context, stem uuid.UUID) (order.Record, error) {
	rec, err := store.Get(ctx, s.db, scan, get_bloom, stem)
	if err == store.ErrEmpty {
		return order.Record{}, order.ErrNotFound
	}

	return rec, nil
}

func (s *SQLDB) Create(ctx context.Context, ent order.Entity) error {
	return store.Create(ctx, s.db, create,
		ent.UUID,
		ent.Version,
		ent.Stem,
		ent.Author,
		ent.Created,
	)
}

func (s *SQLDB) RunTx(ctx context.Context, proc func(order.Store) error) error {
	return store.WithTx(ctx, s.db, func(tx store.DBTx) error {
		return proc(&SQLDB{db: tx})
	})
}

func scan(row store.Scanner) (order.Record, error) {
	var rec order.Record
	if err := row.Scan(
		&rec.UUID,
		&rec.Version,
		&rec.Stem,
		&rec.Author,
		&rec.Name,
		&rec.SIAPE,
		&rec.Created,
	); err != nil {
		return order.Record{}, nil
	}

	return rec, nil
}

const (
	list_blooms  = "select `uuid`, `version`, `stem`, `author`, `name`, `siape`, `created` from `Orders_Bloom_View`"
	list_by_stem = "select `uuid`, `version`, `stem`, `author`, `name`, `siape`, `created` from `Orders_View` where `stem` = ? order by `version` desc"
	get          = "select `uuid`, `version`, `stem`, `author`, `name`, `siape`, `created` from `Orders_View` where `uuid` = ?"
	get_bloom    = "select `uuid`, `version`, `stem`, `author`, `name`, `siape`, `created` from `Orders_Bloom_View` where `stem` = ?"

	create = "insert into `Orders` (`uuid`, `version`, `stem`, `author`, `created`) values (?, ?, ?, ?, ?)"
)
