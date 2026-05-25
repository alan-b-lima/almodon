package money_test

import (
	"testing"

	. "github.com/alan-b-lima/almodon/pkg/money"
)

func TestMarshalText(t *testing.T) {
	type Test struct {
		Got  Money
		Want string
	}

	tests := []Test{
		{Got: 0, Want: "0.00"},
		{Got: 1, Want: "0.01"},
		{Got: 10, Want: "0.10"},
		{Got: -10, Want: "-0.10"},
		{Got: 1234567890, Want: "12345678.90"},
		{Got: -1234567890, Want: "-12345678.90"},
	}

	for _, test := range tests {
		str, _ := test.Got.MarshalText()
		if string(str) != test.Want {
			t.Errorf("Expected %s, got %s", test.Want, test.Got.String())
		}
	}
}

func TestUnmarshalText(t *testing.T) {
	type Test struct {
		Got  string
		Want Money
		Err  error
	}

	tests := []Test{
		{Got: ".0", Want: 0},
		{Got: ".", Err: ErrSyntax},
		{Got: "-.", Err: ErrSyntax},
		{Got: "-a.", Err: ErrSyntax},
		{Got: "-123456789987654321.", Err: ErrRange},
		{Got: "-1234567899876543.", Want: -123456789987654300},
		{Got: ".00", Want: 0},
		{Got: "00.016", Want: 2},
		{Got: "0.01", Want: 1},
		{Got: "0.009999999999", Want: 1},
		{Got: "0.004999999999", Want: 0},
		{Got: "0.10", Want: 10},
		{Got: "-0.10", Want: -10},
		{Got: "12345678.90", Want: 1234567890},
		{Got: "-12345678.905", Want: -1234567891},
	}

	for i, test := range tests {
		var money Money
		err := money.UnmarshalText([]byte(test.Got))
		if err != test.Err {
			t.Errorf("#%d: exepected %v, got %v", i, test.Err, err)
			continue
		}

		if test.Err == nil && money != test.Want {
			t.Errorf("#%d: expected %v, got %v", i, test.Want, money)
		}
	}
}
