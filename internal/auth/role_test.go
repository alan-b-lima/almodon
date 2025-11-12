package auth_test

import (
	"testing"

	. "github.com/alan-b-lima/almodon/internal/auth"
)

func TestDefaultHierarchy(t *testing.T) {
	type Tests struct {
		x, y     Role
		inherits bool
	}

	tests := []Tests{
		{Chief, Chief, true},
		{Chief, Promoted, false},
		{Chief, Admin, false},
		{Chief, User, false},
		{Chief, Unlogged, false},

		{Promoted, Chief, true},
		{Promoted, Promoted, true},
		{Promoted, Admin, false},
		{Promoted, User, false},
		{Promoted, Unlogged, false},

		{Admin, Chief, true},
		{Admin, Promoted, true},
		{Admin, Admin, true},
		{Admin, User, false},
		{Admin, Unlogged, false},

		{User, Chief, true},
		{User, Promoted, true},
		{User, Admin, true},
		{User, User, true},
		{User, Unlogged, false},

		{Unlogged, Chief, true},
		{Unlogged, Promoted, true},
		{Unlogged, Admin, true},
		{Unlogged, User, true},
		{Unlogged, Unlogged, true},
	}

	for _, test := range tests {
		if DefaultHierarchy(test.x, test.y) != test.inherits {
			if test.inherits {
				t.Errorf("%v does inherits %v", test.x, test.y)
			} else {
				t.Errorf("%v does not inherits %v", test.x, test.y)
			}
		}
	}
}
