package migration_test

import (
	"testing"

	. "github.com/alan-b-lima/almodon/pkg/migration"
)

func TestMigration(t *testing.T) {
	type Test struct {
		Prev map[string]string
		Next map[string]string
		Want Rel
	}

	tests := []Test{
		{
			Prev: map[string]string{"a.go": "Hello, World"},
			Next: map[string]string{"a.go": "Hello, World"},
			Want: Equal,
		},
		{
			Prev: map[string]string{"a.go": "Hello, World"},
			Next: map[string]string{"a.go": "Hello, World", "b.go": "Hello, World"},
			Want: Upgradable,
		},
		{
			Prev: map[string]string{"a.go": "Hello, World", "b.go": "Hello, World"},
			Next: map[string]string{"a.go": "Hello, World", "b2.go": "Hello, World"},
			Want: Incompatible,
		},
		{
			Prev: map[string]string{"a.go": "Hello, World"},
			Next: map[string]string{"a.go": "Goodbye, World"},
			Want: Incompatible,
		},
	}

	for _, test := range tests {
		prev := New(test.Prev)
		next := New(test.Next)

		if got := prev.Compare(next); got != test.Want {
			t.Errorf("want %v, got %v", test.Want, got)
		}
	}
}
