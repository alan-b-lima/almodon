package session

import (
	"crypto/rand"
	"encoding/binary"
	"encoding/hex"
	"time"
)

type Token [24]byte

const TokenLen = len(Token{})

func NewToken() Token {
	var token Token

	nanos := time.Now().UnixNano()

	binary.BigEndian.PutUint64(token[:8], uint64(nanos))
	rand.Read(token[8:])

	return token
}

func (t *Token) Bytes() []byte {
	return t[:]
}

func (t Token) String() string {
	b, _ := t.MarshalText()
	return string(b)
}

func ParseString(str string) (Token, error) {
	var t Token
	err := t.UnmarshalJSON([]byte(str))
	return t, err
}

func (t Token) AppendText(b []byte) ([]byte, error) {
	return hex.AppendEncode(b, t[:]), nil
}

func (t Token) MarshalText() ([]byte, error) {
	return t.AppendText(nil)
}

func (t *Token) UnmarshalJSON(b []byte) error {
	if TokenLen != hex.DecodedLen(len(b)) {
		return ErrInvalidToken
	}

	var token Token
	n, err := hex.Decode(token[:], b)
	if err != nil || n != TokenLen {
		return ErrInvalidToken
	}

	*t = token
	return nil
}
