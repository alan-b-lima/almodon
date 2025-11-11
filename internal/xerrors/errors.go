package xerrors

import "github.com/alan-b-lima/almodon/pkg/errors"

var (
	ErrSessionTooLong = errors.Fmt(errors.InvalidInput, "session-too-long", "session must not last longer than %v")

	ErrSessionNotFound = errors.New(errors.NotFound, "session-not-found", "session not found", nil)
)

var (
	ErrUserCreation = errors.Imp(errors.InvalidInput, "user-creation", "given data does not satisfy the user type")

	ErrNameEmpty                     = errors.New(errors.InvalidInput, "name-empty", "name cannot be empty", nil)
	ErrEmailInvalid                  = errors.New(errors.InvalidInput, "email-invalid", "email must be valid", nil)
	ErrPasswordTooShort              = errors.New(errors.InvalidInput, "password-too-short", "password must be at least 8 characters long", nil)
	ErrPasswordTooLong               = errors.New(errors.InvalidInput, "password-too-long", "password must be a maximum of 64 characters long", nil)
	ErrPasswordLeadOrTrailWhitespace = errors.New(errors.InvalidInput, "password-edge-whitespace", "password must not begin or end with whitespaces", nil)
	ErrPasswordIllegalCharacters     = errors.New(errors.InvalidInput, "password-illegal-chars", "password must not contain unprintable or invalid uft-8 characters", nil)

	ErrIncorrectPassword    = errors.New(errors.Unauthorized, "incorrect-password", "given password is incorrect", nil)
	ErrFailedToHashPassword = errors.Imp(errors.Internal, "hash-failure", "failed to hash the password")

	ErrUnauthenticatedUser = errors.Imp(errors.Unauthorized, "unauthenticated-user", "user is not logged in")
	ErrUnauthorizedUser    = errors.Fmt(errors.Forbidden, "unauthorized-user", "auth role %v does not match any criteria in %v")

	ErrUserNotFound = errors.New(errors.NotFound, "user-not-found", "user not found", nil)
)
