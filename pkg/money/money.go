package money

import (
	"errors"
	"math"
)

var (
	ErrSyntax = errors.New("invalid money syntax")
	ErrRange  = errors.New("value is out of range")
)

// Min is the highest value allowed by the money type. It represents the
// following:
//
//	9_999_999_999_999_999.99
//
// That it, 10¹⁶ - 0.01.
const Max = Money(9_999_999_999_999_999_99)

// Min is the lowest value allowed by the money type. It represents the
// following:
//
//	-9_999_999_999_999_999.99
//
// That it, -10¹⁶ + 0.01.
const Min = -Money(9_999_999_999_999_999_99)

// Money is a fixed-point decimal, with precision of 16 decimal integer places
// and 2 decimal fractional places.
//
// Even though this type can represent a little beyond the following, this package
// limits Money to be between [Min] and [Max].
//
// The money representation, as an int64, is the number or cents, addition and
// scalar multiplication are well defined, but do not account for overflow.
type Money int64

// FromFloat parses a money value as the number of monetary units, not cents.
// This means that the float literal 1.24 is equivalent 124 cents, i.e,
// [Money](124).
//
// Since the money type has 2 decimal places of precision on the fractional
// part, it may incur in loss of precision.
//
// Special cases are:
//
//	Round(±0) = 0
//	Round(-Inf) = Min
//	Round(+Inf) = Max
//	Round(NaN) = 0
//
// The output is under- and uppercapped be [Min] and [Max].
func FromFloat(amount float64) Money {
	switch {
	case math.IsInf(amount, -1):
		return Min
	case math.IsInf(amount, 1):
		return Max
	case math.IsNaN(amount):
		return 0
	}

	return clamp(Min, Money(math.Round(amount*100)), Max)
}

func FromCents(cents int64) Money {
	return clamp(Min, Money(cents), Max)
}

func FromString(str string) (Money, error) {
	var m Money
	err := m.UnmarshalText([]byte(str))
	return m, err
}

func (m Money) Cents() int64 {
	return int64(m)
}

func (m Money) InRange() bool {
	return Min <= m && m <= Max
}

func (m Money) String() string {
	b, _ := m.MarshalText()
	return string(b)
}

// AppendText stringifies the money and appends it to a byte slice.
// See [Money.MarshalText] for more details on format.
//
// AppendText always returns a nil error.
func (m Money) AppendText(b []byte) ([]byte, error) {
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

	b = append(b, buf[i+1:]...)
	return b, nil
}

// MarshalText stringifies the money. The money string format produced is
// defined, in a [variant] of Extended Backus-Naur Form (EBNF), as follows:
//
//	Money    = [ "-" ] integer "." fraction
//	integer  = "0" | "1" ... "9" { digit }
//	fraction = digit digit
//	digit    = "0" ... "9"
//
// MarshalText always returns a nil error.
func (m Money) MarshalText() ([]byte, error) {
	return m.AppendText(nil)
}

var def = [20]byte{
	'0', '0', '0', '0', '0',
	'0', '0', '0', '0', '0',
	'0', '0', '0', '0', '0',
	'0', '0', '0', '0', '0',
}

// UnmarshalText parses a money string. A money string is defined, in a
// [variant] of Extended Backus-Naur Form (EBNF), as follows:
//
//	Money  = [ "-" ] number
//	number = digits [ "." [ digits ] ] | "." digits
//	digits = digit { digit }
//	digit  = "0" ... "9"
//
// Note that if the fraction part exceeds two digits, it will be rounded to the
// closest two digit fractional part, which may incur loss of precision.
//
// [variant]: https://en.wikipedia.org/wiki/Wirth_syntax_notation
func (m *Money) UnmarshalText(buf []byte) error {
	if len(buf) == 0 {
		return ErrSyntax
	}

	var neg bool
	if buf[0] == '-' {
		neg = true
		buf = buf[1:]
	}

	int, isz, err := integer(buf)
	if err != nil {
		return err
	}
	buf = buf[isz:]

	if len(buf) == 0 {
		if isz == 0 {
			return ErrSyntax
		}

		return bind(m, neg, int, 0)
	}

	if buf[0] != '.' {
		return ErrSyntax
	}
	buf = buf[1:]

	frac, fsz, err := fraction(buf)
	if err != nil {
		return err
	}
	buf = buf[fsz:]

	if isz == 0 && fsz == 0 {
		return ErrSyntax
	}

	for _, c := range buf {
		if !is_digit(c) {
			return ErrSyntax
		}
	}

	return bind(m, neg, int, frac)
}

func integer(buf []byte) (Money, int, error) {
	var m Money
	var i int

	for i = 0; i < len(buf); i++ {
		c := buf[i]
		if !is_digit(c) {
			break
		}

		d := Money(c - '0')
		if m > (Max-d)/10 {
			return 0, 0, ErrRange
		}

		m = 10*m + d
	}

	return m, i, nil
}

func fraction(buf []byte) (Money, int, error) {
	const prec = 2

	var m Money
	var i int

	for i = 0; i < len(buf) && i < prec+1; i++ {
		c := buf[i]
		if !is_digit(c) {
			return 0, 0, ErrSyntax
		}

		m = 10*m + Money(c-'0')
	}

	if i > prec {
		var p Money
		p, m = m%10, m/10

		if p >= 5 {
			m++
		}
		i = prec
	}

	return m, i, nil
}

func is_digit(c byte) bool {
	return '0' <= c && c <= '9'
}

func bind(m *Money, neg bool, int, frac Money) error {
	if frac > 99 {
		return ErrRange
	}

	if int > (Max-frac)/100 {
		return ErrRange
	}

	amount := 100*int + frac
	if neg {
		amount = -amount
	}

	*m = amount
	return nil
}

func clamp(minv, v, maxv Money) Money {
	if minv >= v {
		return minv
	}
	if maxv <= v {
		return maxv
	}
	return v
}
