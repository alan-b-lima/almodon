package migration

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha512"
	"encoding/hex"
	"slices"
	"strconv"
)

type Migration struct {
	hash    Hash
	scripts []Hash
}

type Hash [32]byte

type Rel uint8

const (
	Identical Rel = iota
	Upgradable
	Incompatible
)

// New creates a new migration from the given sources. The migration's hash is
// computed from the sources' names and contents, and its scripts are the
// hashes of the sources' contents, sorted in ascending order. As of the Go
// spec, this algorithm does not depend on the order of iteration of the
// given sources' map.
//
// Migrations created from the same input are semantically and deeply equal,
// although not comparable by Go semantics. See [Migration.Compare] for more
// advanced usage.
func New(sources map[string]string) *Migration {
	scripts := make([]Hash, 0, len(sources))

	for name, source := range sources {
		h := hmac.New(sha512.New512_256, []byte(name))
		h.Write([]byte(source))

		script := Hash(h.Sum(nil))
		scripts = append(scripts, script)
	}

	slices.SortFunc(scripts, func(a, b Hash) int {
		return bytes.Compare(a[:], b[:])
	})

	h := sha512.New512_256()
	for _, hash := range scripts {
		h.Write(hash[:])
	}

	return &Migration{
		hash:    Hash(h.Sum(nil)),
		scripts: scripts,
	}
}

// IsValid reports whether the migration is valid, its hash is consistent with
// its scripts and its scripts are sorted in ascending order.
func (m *Migration) IsValid() bool {
	h := sha512.New512_256()

	for i, hash := range m.scripts {
		if i > 0 && bytes.Compare(m.scripts[i-1][:], hash[:]) < 0 {
			return false
		}

		h.Write(hash[:])
	}

	return m.hash == Hash(h.Sum(nil))
}

// Compare compares whether migrations are the same, a proper subset of one
// another or incompatible.
//
// Migrations are the same if, and only if, they have the same hash.
// [Identical] is returned.
//
// A previous migration is a proper subset of the next if all scripts in the
// previous migration are present in the next and they have the same hash. In
// this case, the next migration can be applied on top of the previous one, and
// [Upgradable] is returned.
//
// Otherwise, they are incompatible, and [Incompatible] is returned.
func (prev *Migration) Compare(next *Migration) Rel {
	if prev.hash == next.hash {
		return Identical
	}

	if len(prev.scripts) >= len(next.scripts) {
		// prev cannot be a subset of a next set bigger than it, and if they
		// are equal size, then their hashes should have been equal.
		return Incompatible
	}

	for i, j := 0, 0; i < len(prev.scripts); i, j = i+1, j+1 {
		iscript, jscript := prev.scripts[i], next.scripts[j]

		switch bytes.Compare(iscript[:], jscript[:]) {
		case 0:
			continue

		case 1:
			i-- // effectively advance only the next's script.

		case -1:
			// if a greater script is found in the next migration, then the
			// current prev's script is missing.
			return Incompatible
		}
	}

	return Upgradable
}

func (m Migration) Hash() Hash      { return m.hash }
func (m Migration) Scripts() []Hash { return slices.Clone(m.scripts) }

func (m Migration) String() string {
	return m.hash.String() + "(" + strconv.Itoa(len(m.scripts)) + ")"
}

func (h Hash) String() string {
	return hex.EncodeToString(h[:])
}

var rels = [...]string{
	Identical:    "identical",
	Upgradable:   "upgradable",
	Incompatible: "incompatible",
}

func (r Rel) String() string {
	if r < Rel(len(rels)) {
		return rels[r]
	}

	return "rel(" + strconv.Itoa(int(r)) + ")"
}
