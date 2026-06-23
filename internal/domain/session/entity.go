package session

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"strconv"
	"time"
)

const (
	IdleTimeout     = 30 * time.Minute
	HardTimeout     = 24 * time.Hour
	PasswordTimeout = 12 * time.Hour
)

func PasswordVerificationExpired(verifiedAt time.Time) bool {
	return time.Now().After(verifiedAt.Add(PasswordTimeout))
}

func ProcessIdleTimeout(idle_timeout time.Duration) (time.Time, error) {
	if idle_timeout > IdleTimeout {
		return time.Time{}, ErrTooLong
	}

	expires := time.Now().Add(idle_timeout)
	return expires, nil
}

const TokenLen = 24

type Token [TokenLen]byte

func NewToken() Token {
	var token Token
	rand.Read(token[:])

	return token
}

func (t *Token) Bytes() []byte {
	return t[:]
}

var _Format = `%0` + strconv.Itoa(2*TokenLen) + `x`

func (t Token) String() string {
	return fmt.Sprintf(_Format, t[:])
}

func FromString(string string) (Token, error) {
	if 2*TokenLen != len(string) {
		return Token{}, ErrInvalidToken
	}

	var token Token
	for i := range TokenLen {
		b, err := strconv.ParseUint(string[:2], 16, 8)
		if err != nil {
			return Token{}, ErrInvalidToken
		}

		string = string[2:]
		token[i] = byte(b)
	}

	return token, nil
}

func (t Token) MarshalJSON() ([]byte, error) {
	return []byte(`"` + t.String() + `"`), nil
}

func (t *Token) UnmarshalJSON(b []byte) error {
	var string string
	err := json.Unmarshal(b, &string)
	if err != nil {
		return err
	}

	token, err := FromString(string)
	if err != nil {
		return err
	}

	*t = token
	return nil
}

func Expired(hard_deadline, idle_deadline time.Time) bool {
	now := time.Now()
	return hard_deadline.Before(now) || idle_deadline.Before(now)
}
