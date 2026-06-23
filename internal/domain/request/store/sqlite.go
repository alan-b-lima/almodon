package requeststore

import (
	"context"
	"database/sql"
	_ "embed"
	"errors"

	"github.com/alan-b-lima/almodon/internal/domain/request"
	"github.com/alan-b-lima/almodon/internal/support/store"
	"github.com/alan-b-lima/almodon/pkg/uuid"
)

//go:embed sqlite.sql
var Script string

const (
	list           = `select uuid, number, author, name, siape, title, memo, status, created, updated from Request_View order by created desc`
	list_by_status = `select uuid, number, author, name, siape, title, memo, status, created, updated from Request_View where status = ? order by created desc`

	get           = `select uuid, number, author, name, siape, title, memo, status, created, updated from Request_View where uuid = ?`
	get_by_number = `select uuid, number, author, name, siape, title, memo, status, created, updated from Request_View where number = ?`

	create = `insert into Request (uuid, number, author, title, memo, status, created, updated) values (?, ?, ?, ?, ?, ?, ?, ?)`
	delete = `delete from Request where uuid = ?`
)

type SQLDB struct {
	db store.DBTx
}

var _ request.Store = (*SQLDB)(nil)

func New(db *sql.DB) *SQLDB {
	return &SQLDB{db: db}
}

func (s *SQLDB) List(ctx context.Context) ([]request.Record, error) {
	return store.List(ctx, s.db, scan, list)
}

func (s *SQLDB) ListByStatus(ctx context.Context, status request.Status) ([]request.Record, error) {
	return store.List(ctx, s.db, scan, list_by_status, statusInt[status])
}

func (s *SQLDB) Get(ctx context.Context, uuid uuid.UUID) (request.Record, error) {
	return store.Get(ctx, s.db, scan, get, uuid)
}

func (s *SQLDB) GetByNumber(ctx context.Context, number int) (request.Record, error) {
	return store.Get(ctx, s.db, scan, get_by_number, number)
}

func (s *SQLDB) Create(ctx context.Context, ent request.Entity) error {
	return store.Create(ctx, s.db, create,
		ent.UUID,
		ent.Number,
		ent.Author,
		ent.Title,
		ent.Memo,
		ent.Status,
		ent.Created,
		ent.Updated,
	)
}

func (s *SQLDB) Delete(ctx context.Context, uuid uuid.UUID) error {
	return store.Delete(ctx, s.db, delete, uuid)
}

func scan(scanner store.Scanner) (request.Record, error) {
	var status int
	var rec request.Record

	if err := scanner.Scan(
		&rec.UUID,
		&rec.Number,
		&rec.Author,
		&rec.Name,
		&rec.SIAPE,
		&rec.Title,
		&rec.Memo,
		&status,
		&rec.Created,
		&rec.Updated,
	); err != nil {
		return request.Record{}, err
	}

	var in bool
	rec.Status, in = intStatus[status]
	if !in {
		return request.Record{}, errBadStatus
	}

	return rec, nil
}

var errBadStatus = errors.New("bad status")

var statusInt = map[request.Status]int{
	request.Open:     10,
	request.Closed:   20,
	request.Staged:   30,
	request.Accepted: 40,
}

var intStatus = mirror(statusInt)

func mirror[K, V comparable](m map[K]V) map[V]K {
	mn := make(map[V]K, len(m))
	for k, v := range m {
		if _, in := mn[v]; in {
			panic("mirror: duplicate value")
		}

		mn[v] = k
	}

	return mn
}
