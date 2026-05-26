package user

import "github.com/alan-b-lima/pkg/problem"

var (
	ErrCreate   = problem.Imp(problem.SemanticalError, "user-create").Message("could not create user")
	ErrUpdate   = problem.Imp(problem.SemanticalError, "user-update").Message("could not update user")
	ErrNotFound = problem.New(problem.NotFound, "user-not-found", "user not found", nil, nil)

	ErrNameEmpty   = problem.New(problem.SemanticalError, "user-name-empty", "name must not be empty", nil, nil)
	ErrNameTooLong = problem.Imp(problem.SemanticalError, "user-name-too-long").Format("name must be less than %d characters").Details(map[string]any{"max": NameMaxLen}).Make(NameMaxLen)

	ErrSiapeInvalid = problem.New(problem.Malformed, "user-siape-invalid", "siape must be a seven-character number", nil, nil)
	ErrSiapeTaken   = problem.New(problem.Conflict, "user-siape-in-use", "siape already taken", nil, nil)

	ErrEmailInvalid = problem.New(problem.Malformed, "user-email-invalid", "email invalid", nil, map[string]any{"pattern": re_email})

	ErrPasswordTooShort              = problem.New(problem.SemanticalError, "password-too-short", "password too short", nil, map[string]any{"max": PasswordMinLen})
	ErrPasswordTooLong               = problem.New(problem.SemanticalError, "password-too-long", "password too long", nil, map[string]any{"max": PasswordMaxLen})
	ErrPasswordLeadOrTrailWhitespace = problem.New(problem.SemanticalError, "password-edge-whitespace", "password must not start or end in whitespace", nil, nil)
	ErrPasswordIllegalCharacters     = problem.New(problem.Malformed, "password-illegal-chars", "password must not include unprintable or illegal UTF-8 characters", nil, nil)
	ErrPasswordFailedToHash          = problem.Imp(problem.UnexpectedError, "password-hash-failed").Message("unexpected error ocurred while hashing password")
	ErrPasswordIncorrect             = problem.New(problem.Unauthenticated, "password-incorrect", "incorrect password", nil, nil)

	ErrRoleInvalid     = problem.New(problem.Malformed, "user-role-invalid", "role was not identified", nil, nil)
	ErrRoleNotAccepted = problem.Imp(problem.SemanticalError, "user-role-not-accepted").Format("role must be one of %v").Details(map[string]any{"accept_roles": accept_roles}).Make(accept_roles)

	ErrDeleteRootUser  = problem.New(problem.Conflict, "delete-root-user", "the root user cannot be deleted", nil, nil)
	ErrNotEnoughChiefs = problem.New(problem.Conflict, "not-enough-chiefs", "there must be at least one chief at any moment", nil, nil)
)
