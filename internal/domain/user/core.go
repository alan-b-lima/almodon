package user

import (
	"time"

	"github.com/alan-b-lima/almodon/internal/auth"
	sessionpkg "github.com/alan-b-lima/almodon/internal/domain/session"
	"github.com/alan-b-lima/almodon/internal/xerrors"
	"github.com/alan-b-lima/almodon/pkg/hash"
	"github.com/alan-b-lima/almodon/pkg/opt"
	"github.com/alan-b-lima/almodon/pkg/uuid"
)

func List(users Lister, offset, limit int) (Entities, error) {
	return users.List(offset, limit)
}

func Get(users Getter, uuid uuid.UUID) (Entity, error) {
	return users.Get(uuid)
}

func GetBySIAPE(users GetterBySIAPE, siape int) (Entity, error) {
	return users.GetBySIAPE(siape)
}

func Create(users Creater, siape int, name, email, password string, role auth.Role) (Entity, error) {
	return users.Create(siape, name, email, password, role)
}

func Patch(users Patcher, uuid uuid.UUID, name, email, password opt.Opt[string], role opt.Opt[auth.Role]) (Entity, error) {
	return users.Patch(uuid, name, email, password, role)
}

func Delete(users Deleter, uuid uuid.UUID) error {
	return users.Delete(uuid)
}

func Authenticate(users GetterBySIAPE, sessions sessionpkg.Creater, siape int, password string) (AuthEntity, error) {
	res, err := users.GetBySIAPE(siape)
	if err != nil {
		return AuthEntity{}, err
	}

	if !hash.Compare(res.Password[:], []byte(password)) {
		return AuthEntity{}, xerrors.ErrIncorrectPassword
	}

	s, err := sessions.Create(res.UUID, 10*time.Minute)
	if err != nil {
		return AuthEntity{}, err
	}

	ares := AuthEntity{
		UUID:    s.UUID,
		User:    res.UUID,
		Expires: s.Expires,
	}
	return ares, nil
}

func Actor(users Getter, sessions sessionpkg.Getter, session uuid.UUID) (auth.Actor, error) {
	res, err := sessionpkg.Get(sessions, session)
	if err != nil {
		return auth.NewUnlogged(), err
	}

	ures, err := users.Get(res.User)
	if err != nil {
		return auth.NewUnlogged(), err
	}

	return auth.NewLogged(
		ures.UUID,
		ures.Role,
	), nil
}
