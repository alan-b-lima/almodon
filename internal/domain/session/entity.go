package session

import (
	"time"

	"github.com/alan-b-lima/almodon/internal/xerrors"
	"github.com/alan-b-lima/almodon/pkg/errors"
	"github.com/alan-b-lima/almodon/pkg/uuid"
)

const _MaxMaxAge = 7 * 24 * time.Hour

type Session struct {
	uuid    uuid.UUID
	user    uuid.UUID
	expires time.Time
}

func New(user uuid.UUID, maxAge time.Duration) (Session, error) {
	session := Session{}

	err := errors.Join(
		session.setUser(user),
		session.SetMaxAge(maxAge),
	)
	if err != nil {
		return Session{}, err
	}

	session.uuid = uuid.NewUUIDv7()
	return session, nil
}

func (s *Session) UUID() uuid.UUID    { return s.uuid }
func (s *Session) User() uuid.UUID    { return s.user }
func (s *Session) Expires() time.Time { return s.expires }

func (s *Session) setUser(uuid uuid.UUID) error {
	s.user = uuid
	return nil
}

func (s *Session) SetMaxAge(maxAge time.Duration) error {
	if maxAge > _MaxMaxAge {
		return xerrors.ErrSessionTooLong.New(_MaxMaxAge)
	}

	s.expires = time.Now().Add(maxAge)
	return nil
}
