package stem

import (
	"strings"
	"unicode"
)

func ProcessTitle(title string) (string, error) {
	title = strings.TrimSpace(title)

	if title == "" {
		return "", ErrTitleInvalid
	}

	for _, rune := range title {
		if !unicode.IsLetter(rune) && !unicode.IsDigit(rune) && rune != '-' && rune != '/' {
			return "", ErrTitleInvalid
		}
	}

	return title, nil
}
