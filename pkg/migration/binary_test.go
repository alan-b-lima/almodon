package migration_test

import (
	crand "crypto/rand"
	"io"
	"math/rand/v2"
	"strings"
	"testing"

	. "github.com/alan-b-lima/almodon/pkg/migration"
)

func FuzzInversabilityOfBinaryEncoding(f *testing.F) {
	const numEntries, nameLen, entryLen int64 = 4, 10, 62
	f.Add(numEntries, nameLen, entryLen)

	f.Fuzz(func(t *testing.T, e, nl, el int64) {
		if e < 0 || nl <= 0 || el <= 0 {
			return
		}

		sources := make(map[string]string, e)
		for range e {
			var b1, b2 strings.Builder

			io.CopyN(&b1, crand.Reader, rand.N(nl))
			io.CopyN(&b2, crand.Reader, rand.N(el))

			sources[b1.String()] = b2.String()
		}

		mIn := New(sources)

		b, err := mIn.MarshalBinary()
		if err != nil {
			t.Error(err)
		}

		var mOut Migration
		if err := mOut.UnmarshalBinary(b); err != nil {
			t.Error(err)
		}

		if cmp := mIn.Compare(&mOut); cmp != Identical {
			t.Errorf("expected to be equal, got %v", cmp)
		}
	})
}
