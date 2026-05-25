package uuid_test

import (
	"crypto/rand"
	"sync"
	"testing"

	. "github.com/alan-b-lima/almodon/pkg/uuid"
)

func TestInvariant(t *testing.T) {
	const (
		mask_4bit = (1 << 4) - 1
	)

	const numTests = 100

	for range numTests {
		uuid := NewUUIDv7()

		version := uint64(uuid[6]) >> 4
		variant := uint64(uuid[8]) >> 6

		if version != 7 {
			t.Errorf("unexpected version, expected 7, got %d", version)
		}

		if variant != 0b10 {
			t.Errorf("unexpected version, expected 10, got %02b", variant)
		}
	}
}

func TestConcurrentUUIDGeneration(t *testing.T) {
	const batches, size = 123, 1999
	const limit = batches * size

	res := make([]UUID, limit)
	var wg sync.WaitGroup

	wg.Add(batches)
	for range batches {
		var r []UUID
		r, res = res[:size], res[size:]

		go func() {
			for i := range size {
				uuid := NewUUIDv7()
				r[i] = uuid
			}

			wg.Done()
		}()
	}

	wg.Wait()

	set := make(map[UUID]struct{}, limit)
	for _, v := range res {
		set[v] = struct{}{}
	}

	if len(set) < limit {
		t.Error("equal UUIDs have been generated")
	}
}

func TestInversabilityOfStringAndFromString(t *testing.T) {
	const numTests = 1000

	for range numTests {
		var uuid UUID
		rand.Read(uuid[:])

		str := uuid.String()
		if uuid2, err := FromString(str); err != nil {
			t.Error(err)
		} else if uuid != uuid2 {
			t.Errorf("%x and %x should be equal", uuid, uuid2)
		}
	}
}

func TestInversabilityOfMarshalTextAndUnmarshalText(t *testing.T) {
	const numTests = 1000

	for range numTests {
		var uuid UUID
		rand.Read(uuid[:])

		str, _ := uuid.MarshalText()
		var uuid2 UUID

		if err := uuid2.UnmarshalText(str); err != nil {
			t.Error(err)
		} else if uuid != uuid2 {
			t.Errorf("%x and %x should be equal", uuid, uuid2)
		}
	}
}
