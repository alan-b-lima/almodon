package session_test

import (
	"testing"

	. "github.com/alan-b-lima/almodon/internal/domain/session"
)

func Test_InversabilityBetweenStringAndParseString(t *testing.T) {
	const num_tests = 1000

	for range num_tests {
		token := NewToken()

		str := token.String()
		if back, err := ParseString(str); err != nil {
			t.Error(err)
		} else if token != back {
			t.Errorf("%x and %x should be equal", token, back)
		}
	}
}
