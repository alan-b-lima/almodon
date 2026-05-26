package stem

import (
	"strings"
	"unicode"

	"github.com/alan-b-lima/almodon/internal/support/service"
)

const (
	NameMinLen = 3
	NameMaxLen = 32
)

func ProcessName(title string) (string, error) {
	title = strings.TrimSpace(title)

	if !service.StringLenMin(title, NameMinLen) {
		return "", ErrTitleTooShort
	}

	if !service.StringLenMax(title, NameMaxLen) {
		return "", ErrTitleTooLong
	}

	for _, rune := range title {
		if !unicode.IsLetter(rune) && !unicode.IsDigit(rune) && rune != '-' && rune != '/' {
			return "", ErrTitleInvalid
		}
	}

	return title, nil
}
