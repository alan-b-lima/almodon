package user

import (
	"regexp"
	"unicode/utf8"

	"github.com/alan-b-lima/almodon/internal/auth"
	"github.com/alan-b-lima/almodon/internal/xerrors"
	"github.com/alan-b-lima/almodon/pkg/errors"
	"github.com/alan-b-lima/almodon/pkg/hash"
	"github.com/alan-b-lima/almodon/pkg/uuid"
)

type User struct {
	uuid     uuid.UUID
	siape    int
	name     string
	email    string
	password [60]byte
	role     auth.Role
}

func New(siape int, name, email, password string, role auth.Role) (User, error) {
	var u User

	errpwd := u.SetPassword(password)
	if err, ok := errors.AsType[*errors.Error](errpwd); ok && err.IsInternal() {
		return User{}, err
	}

	err := errors.Join(
		u.SetSIAPE(siape),
		u.SetName(name),
		u.SetEmail(email),
		errpwd,
		u.SetRole(role),
	)
	if err != nil {
		return User{}, xerrors.ErrUserCreation.New(err)
	}

	u.uuid = uuid.NewUUIDv7()
	return u, nil
}

func (u *User) UUID() uuid.UUID    { return u.uuid }
func (u *User) SIAPE() int         { return u.siape }
func (u *User) Name() string       { return u.name }
func (u *User) Email() string      { return u.email }
func (u *User) Password() [60]byte { return u.password }
func (u *User) Role() auth.Role    { return u.role }

func (u *User) SetSIAPE(siape int) error          { return set(&u.siape, siape, ProcessSiape) }
func (u *User) SetName(name string) error         { return set(&u.name, name, ProcessName) }
func (u *User) SetEmail(email string) error       { return set(&u.email, email, ProcessEmail) }
func (u *User) SetPassword(password string) error { return set(&u.password, password, ProcessPassword) }
func (u *User) SetRole(role auth.Role) error      { return set(&u.role, role, ProcessRole) }

func ProcessSiape(siape int) (int, error) {
	return siape, nil
}

func ProcessName(name string) (string, error) {
	if name == "" {
		return "", xerrors.ErrNameEmpty
	}

	return name, nil
}

var reEmail = regexp.MustCompile(`^[0-9A-Za-z_%+-]+(\.[0-9A-Za-z_%+-]+)*@[0-9A-Za-z-]+(\.[0-9A-Za-zA-Z-]+)*\.[A-Za-z]{2,}$`)

func ProcessEmail(email string) (string, error) {
	if reEmail.MatchString(email) {
		return "", xerrors.ErrEmailInvalid
	}

	return email, nil
}

func ProcessPassword(password string) ([60]byte, error) {
	if len(password) < 8 {
		return [60]byte{}, xerrors.ErrPasswordTooShort
	}

	if len(password) > 64 {
		return [60]byte{}, xerrors.ErrPasswordTooLong
	}

	switch password[0] {
	case ' ', '\t', '\n', '\r':
		return [60]byte{}, xerrors.ErrPasswordLeadOrTrailWhitespace
	}

	switch password[len(password)-1] {
	case ' ', '\t', '\n', '\r':
		return [60]byte{}, xerrors.ErrPasswordLeadOrTrailWhitespace
	}

	for _, rune := range password {
		if rune < ' ' || !utf8.ValidRune(rune) {
			return [60]byte{}, xerrors.ErrPasswordIllegalCharacters
		}
	}

	hash, err := hash.Hash([]byte(password))
	if err != nil {
		return [60]byte{}, xerrors.ErrFailedToHashPassword.New(err)
	}

	return hash, nil
}

func ProcessRole(role auth.Role) (auth.Role, error) {
	if !role.IsValid() {
		return 0, nil
	}

	return role, nil
}

func set[D, S any](dst *D, src S, proc func(S) (D, error)) error {
	val, err := proc(src)
	if err != nil {
		return err
	}

	*dst = val
	return nil
}
