package money

import (
	"encoding/json"
	"math"
)

type Money int64

func FromInt(cents int) Money {
	return Money(cents)
}

func FromFloat(amount float64) Money {
	return Money(math.Round(amount * 100))
}

func FromCents(cents int64) Money {
	return Money(cents)
}

func (m Money) Cents() int64 {
	return int64(m)
}

func (m Money) String() string {
	buf := def
	i := len(buf) - 2

	var neg bool
	if m < 0 {
		neg = true
		m = -m
	}

	for m > 0 {
		rem := m % 10
		m = m / 10

		buf[i] = '0' + byte(rem)
		i--
	}

	buf[len(buf)-1] = buf[len(buf)-2]
	buf[len(buf)-2] = buf[len(buf)-3]
	buf[len(buf)-3] = '.'

	i = min(i, len(buf)-5)

	if neg {
		buf[i] = '-'
		i--
	}

	return string(buf[i+1:])
}

func (m Money) MarshalJSON() ([]byte, error) {
	return []byte(m.String()), nil
}

func (m *Money) UnmarshalJSON(b []byte) error {
	var amount float64
	if err := json.Unmarshal(b, &amount); err != nil {
		return err
	}

	*m = FromFloat(amount)
	return nil
}

var def = [20]byte{
	'0', '0', '0', '0', '0',
	'0', '0', '0', '0', '0',
	'0', '0', '0', '0', '0',
	'0', '0', '0', '0', '0',
}
