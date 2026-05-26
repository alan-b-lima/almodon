package user

import (
	"regexp"
	"slices"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/alan-b-lima/almodon/internal/domain/auth"
	"github.com/alan-b-lima/almodon/internal/support/service"

	"golang.org/x/crypto/bcrypt"
)

const (
	NameMaxLen = 128

	PasswordMinLen = 8
	PasswordMaxLen = 64
)

var (
	re_siape = regexp.MustCompile(`^\d{7}$`)
	re_email = regexp.MustCompile(`^[0-9A-Za-z_%+-]+(\.[0-9A-Za-z_%+-]+)*@[0-9A-Za-z-]+(\.[0-9A-Za-zA-Z-]+)*\.[A-Za-z]{2,}$`)

	accept_roles = [...]auth.Role{auth.User, auth.Admin, auth.Chief, auth.Maintainer}
)

const Root = "0000000"

func ProcessSIAPE(siape string) (string, error) {
	if !re_siape.MatchString(siape) {
		return "", ErrSiapeInvalid
	}

	return siape, nil
}

func ProcessName(name string) (string, error) {
	name = strings.TrimSpace(name)

	if name == "" {
		return "", ErrNameEmpty
	}

	if service.StringLenMax(name, NameMaxLen) {
		return "", ErrNameTooLong
	}

	return name, nil
}

func ProcessEmail(email string) (string, error) {
	email = strings.TrimSpace(email)

	if !re_email.MatchString(email) {
		return "", ErrEmailInvalid
	}

	return email, nil
}

func ProcessPassword(password string) ([]byte, error) {
	if len(password) < PasswordMinLen {
		return nil, ErrPasswordTooShort
	}

	if len(password) > PasswordMaxLen {
		return nil, ErrPasswordTooLong
	}

	if unicode.IsSpace(rune(password[0])) {
		return nil, ErrPasswordLeadOrTrailWhitespace
	}

	if unicode.IsSpace(rune(password[len(password)-1])) {
		return nil, ErrPasswordLeadOrTrailWhitespace
	}

	for _, rune := range password {
		if rune < ' ' || !utf8.ValidRune(rune) {
			return nil, ErrPasswordIllegalCharacters
		}
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, ErrPasswordFailedToHash.Cause(err).Make()
	}

	return hash, nil
}

func ComparePassword(hash []byte, password string) error {
	if bcrypt.CompareHashAndPassword(hash, []byte(password)) == nil {
		return nil
	}

	return ErrPasswordIncorrect
}

func ProcessRole(role auth.Role) (auth.Role, error) {
	if !slices.Contains(accept_roles[:], role) {
		return 0, ErrRoleNotAccepted
	}

	return role, nil
}
