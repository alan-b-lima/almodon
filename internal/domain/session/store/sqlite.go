package sessionstore

import (
	"context"
	"database/sql"
	_ "embed"
	"time"

	"github.com/alan-b-lima/almodon/internal/domain/session"
	"github.com/alan-b-lima/almodon/internal/support/service"
	"github.com/alan-b-lima/almodon/internal/support/store"
	"github.com/alan-b-lima/almodon/pkg/uuid"
)

//go:embed sqlite.sql
var Script string

const (
	list = `select token, user, hard_deadline, idle_deadline from Sessions`
	get  = list + ` where token = ?`

	create = `insert into Sessions (token, user, hard_deadline, idle_deadline) values (?, ?, ?, ?, ?)`
	update = `update Sessions set idle_deadline = ? where token = ?`

	delete         = `delete from Sessions where token = ?`
	delete_by_user = `delete from Sessions where user = ?`
	delete_expired = `delete from Sessions where hard_deadline < ? or idle_deadline < ?`
)

type SQLDB struct {
	db store.DBTx
}

var _ session.Store = (*SQLDB)(nil)

func New(db store.DBTx) *SQLDB {
	return &SQLDB{db: db}
}

func (s *SQLDB) List(ctx context.Context) ([]session.Record, error) {
	rows, err := s.db.QueryContext(ctx, list)
	if err != nil {
		return nil, store.ErrQuery.Cause(err).Make()
	}

	var recs []session.Record
	for rows.Next() {
		var rec session.Record
		if err := scan(&rec, rows); err != nil {
			return nil, store.ErrQuery.Cause(err).Make()
		}

		recs = append(recs, rec)
	}
	if err := rows.Err(); err != nil {
		return nil, store.ErrQuery.Cause(err).Make()
	}

	return recs, nil
}

func (s *SQLDB) Get(ctx context.Context, token session.Token) (session.Record, error) {
	row := s.db.QueryRowContext(ctx, get, token.Bytes())

	var rec session.Record
	if err := scan(&rec, row); err != nil {
		if err == sql.ErrNoRows {
			return session.Record{}, session.ErrNotFound
		}

		return session.Record{}, store.ErrQuery.Cause(err).Make()
	}

	return rec, nil
}

func (s *SQLDB) Create(ctx context.Context, rec session.Entity) error {
	_, err := s.db.ExecContext(ctx, create, rec.Token.Bytes(), rec.User, rec.HardDeadline, rec.IdleDeadline)
	if err != nil {
		return store.ErrExec.Cause(err).Make()
	}

	return nil
}

func (s *SQLDB) Update(ctx context.Context, token session.Token, deadline time.Time) error {
	res, err := s.db.ExecContext(ctx, update, deadline, token.Bytes())
	if err != nil {
		return store.ErrExec.Cause(err).Make()
	}

	changed, err := res.RowsAffected()
	if err == nil && changed == 0 {
		return session.ErrNotFound
	}

	return nil
}

func (s *SQLDB) Delete(ctx context.Context, token session.Token) error {
	_, err := s.db.ExecContext(ctx, delete, token.Bytes())
	if err != nil {
		return store.ErrExec.Cause(err).Make()
	}

	return nil
}

func (s *SQLDB) DeleteByUser(ctx context.Context, user uuid.UUID) error {
	_, err := s.db.ExecContext(ctx, delete_by_user, user)
	if err != nil {
		return store.ErrExec.Cause(err).Make()
	}

	return nil
}

func (s *SQLDB) DeleteExpired(ctx context.Context, deadline time.Time) error {
	_, err := s.db.ExecContext(ctx, delete_expired, deadline, deadline)
	if err != nil {
		return store.ErrExec.Cause(err).Make()
	}

	return nil
}

func (s *SQLDB) RunTx(ctx context.Context, proc func(session.Store) error) error {
	return store.WithTx(ctx, s.db, func(tx store.DBTx) error {
		return proc(New(tx))
	})
}

func scan(ent *session.Record, scanner store.Scanner) error {
	var token []byte

	err := scanner.Scan(
		&token,
		&ent.User,
		&ent.HardDeadline,
		&ent.IdleDeadline,
	)
	if err != nil {
		return err
	}

	return service.Set(&ent.Token, token, token_from_bytes)
}

func token_from_bytes(bytes []byte) (session.Token, error) {
	if len(bytes) != session.TokenLen {
		return session.Token{}, session.ErrInvalidToken
	}

	return session.Token(bytes), nil
}
