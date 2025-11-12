package resource

import (
	"net/http"

	"github.com/alan-b-lima/almodon/internal/auth"
	"github.com/alan-b-lima/almodon/internal/domain/user"
	"github.com/alan-b-lima/almodon/internal/xerrors"
	"github.com/alan-b-lima/almodon/pkg/errors"
	uuidpkg "github.com/alan-b-lima/almodon/pkg/uuid"
)

const SessionCookieName = "session"

func SessionCookie(cookie string, r *http.Request) (uuidpkg.UUID, error) {
	s, err := r.Cookie(cookie)
	if err != nil {
		return uuidpkg.UUID{}, xerrors.ErrUnauthenticatedUser.New(nil)
	}

	uuid, err := uuidpkg.FromString(s.Value)
	if err != nil {
		return uuidpkg.UUID{}, xerrors.ErrBadUUID
	}

	return uuid, nil
}

type actoer interface {
	Actor(user.ActorRequest) (auth.Actor, error)
}

func Session(rc actoer, r *http.Request) (auth.Actor, error) {
	session, err := SessionCookie(SessionCookieName, r)
	if err != nil {
		return auth.NewUnlogged(), nil
	}

	act, err := rc.Actor(user.ActorRequest{Session: session})
	if err, ok := errors.AsType[*errors.Error](err); ok && err.Kind.IsClient() {
		return auth.NewUnlogged(), nil
	}
	if err != nil {
		return auth.NewUnlogged(), err
	}

	return act, err
}
