package xerrors

import "github.com/alan-b-lima/almodon/pkg/errors"

var (
	ErrSessionNotFound = errors.New(errors.NotFound, "session-not-found", "session not found", nil)
	ErrSessionTooLong  = errors.Fmt(errors.InvalidInput, "session-too-long", "session must not last longer than %v")
)

var (
	ErrUnauthenticatedUser = errors.Imp(errors.Unauthorized, "unauthenticated-user", "user is not logged in")
	ErrUnauthorizedUser    = errors.Fmt(errors.Forbidden, "unauthorized-user", "auth level %v does not match any criteria in %v")
)
