package service

import (
	"github.com/alan-b-lima/almodon/internal/auth"
	"github.com/alan-b-lima/almodon/internal/xerrors"
)

func Authorize(auth auth.Permission, act auth.Actor) error {
	if role := act.Role(); !auth.Authorize(role) {
		return xerrors.ErrUnauthorizedUser.New(role, auth)
	}

	return nil
}
