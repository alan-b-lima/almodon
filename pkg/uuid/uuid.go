// Copyright (C) 2025-2026 Alan Barbosa Lima.
//
// Almodon is licensed under the GNU General Public License
// version 3. You should have received a copy of the
// license, located in LICENSE, at the root of the source
// tree. If not, see <https://www.gnu.org/licenses/>.

// Package uuid implements an UUID version 7 generator, defined in
// [RFC9562 section 5.7], as well as parsing and scanning of UUIDs. This
// package is complient to [RFC9562].
//
// [RFC9562]: https://www.rfc-editor.org/rfc/rfc9562
// [RFC9562 section 5.7]: https://www.rfc-editor.org/rfc/rfc9562#name-uuid-version-7
package uuid

import (
	"bytes"
	"crypto/rand"
	"database/sql/driver"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"sync"
	"time"
)

// UUID is the set of all 128-bit universally unique identifiers.
//
// Elements of this type have the should be compared using the == operator.
//
// The zero value of the UUID type it the Nil UUID.
type UUID [16]byte

var (
	// Nil is the zero value of UUID, 00000000-0000-0000-0000-000000000000.
	//
	// The Nil UUID is defined in [RFC9562 section 5.9]. Do not confuse this
	// with the Go value nil.
	//
	// [RFC9562 section 5.9]: https://www.rfc-editor.org/rfc/rfc9562#section-5.9
	Nil UUID

	// Max is the maximum value of UUID, FFFFFFFF-FFFF-FFFF-FFFF-FFFFFFFFFFFF.
	//
	// The Max UUID is defined in [RFC9562 section 5.10]. This may de used as a
	// sentinel value, in a implementation-specific way.
	//
	// [RFC9562 section 5.10]: https://www.rfc-editor.org/rfc/rfc9562#section-5.10
	Max = UUID{
		0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF,
		0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF,
	}
)

var errBadSource = errors.New("uuid: source cannot be converted to an UUID")

// NewUUIDv7 generates a new UUID accourding to version 7. This function is
// safe for concurrent use by multiple gorotines.
//
// The memory layout of a UUIDv7, as defined in [RFC9562 section 5.7], is as
// follows:
//
//   - Unix Timestamp: 48-bit big-endian unsigned number of the Unix
//     Epoch timestamp in milliseconds. Occupies bits 0 through 47,
//     octets 0 through 5.
//
//   - Version: 4-bit version field, set to 0b0111 (7). Occupies bits
//     48 through 51, octet 6.
//
//   - Random A: 12-bit pseudo-random data to provide uniqueness.
//     Occupies bits 52 through 63, octects 6 through 7.
//
//   - Variant: 2-bit variant field, set to 0b10. Occupies bits 64
//     through 65, octet 8.
//
//   - Random B: 62-bit pseudo-random data to provide uniqueness.
//     Occupies bits 66 through 127, octets 8 through 15.
//
// [RFC9562 section 5.7]: https://www.rfc-editor.org/rfc/rfc9562#name-uuid-version-7
func NewUUIDv7() UUID {
	const (
		mask_48bit = (1 << 48) - 1

		version = 7
		variant = 0b10
	)

	now := time.Now()
	timestamp := uint64(now.UnixMilli() & mask_48bit)
	randA, randB := next()

	var uuid UUID

	binary.BigEndian.PutUint64(uuid[:8], timestamp<<16|randA)
	binary.BigEndian.PutUint64(uuid[8:], randB)

	uuid[6] |= version << 4
	uuid[8] |= variant << 6

	return uuid
}

// IsNil reports whether the given UUID is the Nil UUID.
func (uuid UUID) IsNil() bool {
	return uuid == Nil
}

// IsMax reports whether the given UUID is the Max UUID.
func (uuid UUID) IsMax() bool {
	return uuid == Max
}

// FromBytes converts an UUID from a byte slice. The given byte slice
// should be of length 16.
//
// Due to this, this function is NOT interchangeable with
// [FromString], a byte slice that is the string representation of an
// UUID should be converted with [FromString].
func FromBytes(bytes []byte) (UUID, error) {
	if len(bytes) != 16 {
		return UUID{}, errBadSource
	}

	return UUID(bytes), nil
}

// FromString converts an UUID from the string format, in a
// case-insensitive manner.
//
// Note that this function is NOT interchangeable with [FromBytes],
// see [FromBytes] for more detail.
func FromString(str string) (UUID, error) {
	var uuid UUID
	err := uuid.UnmarshalText([]byte(str))
	return uuid, err
}

// Bytes returns the byte slice representation of the UUID. Changing the
// returned byte slice is safe, this won't change the original UUID.
func (uuid UUID) Bytes() []byte {
	return uuid[:]
}

// String implements the interface [fmt.Stringer] on the UUID type.
// An UUID is formated as:
//
//	xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx
//
// Each x is a lowercase hexadecimal digit.
func (uuid UUID) String() string {
	buf, _ := uuid.MarshalText()
	return string(buf)
}

// MarshalText implements the interface [encoding.TextMarshaler] on the
// UUID type. The output is in the format defined in [UUID.String].
func (uuid UUID) MarshalText() ([]byte, error) {
	buf := []byte("00000000-0000-0000-0000-000000000000")

	hex.Encode(buf[0:8], uuid[0:4])
	hex.Encode(buf[9:13], uuid[4:6])
	hex.Encode(buf[14:18], uuid[6:8])
	hex.Encode(buf[19:23], uuid[8:10])
	hex.Encode(buf[24:36], uuid[10:16])

	return buf, nil
}

// UnmarshalText implements the interface [encoding.TextUnmarshaler] on the
// UUID type.
//
// The byte slice should be in one of the following formats:
//
//   - xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx
//   - urn:uuid:xxxxxxxx-xxxx-xxxx-xxxx-xxxx-xxxxxxxx
//   - xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
//
// Where each x is a hexadecimal digit in any case.
func (uuid *UUID) UnmarshalText(buf []byte) error {
	var dst UUID

	switch len(buf) {
	case len("urn:uuid:") + 36:
		buf = bytes.TrimPrefix(buf, []byte("urn:uuid:"))

	case 32:
		if _, err := hex.Decode(dst[:], buf); err != nil {
			return err
		}

		*uuid = dst
		return nil
	}

	if len(buf) != 36 {
		return errBadSource
	}
	if buf[8] != '-' || buf[13] != '-' || buf[18] != '-' || buf[23] != '-' {
		return errBadSource
	}

	_, err0 := hex.Decode(dst[0:4], buf[0:8])
	_, err1 := hex.Decode(dst[4:6], buf[9:13])
	_, err2 := hex.Decode(dst[6:8], buf[14:18])
	_, err3 := hex.Decode(dst[8:10], buf[19:23])
	_, err4 := hex.Decode(dst[10:16], buf[24:36])

	if err := or(err0, err1, err2, err3, err4); err != nil {
		return err
	}

	*uuid = dst
	return nil
}

// Value implements SQL [driver.Valuer] interface, it returns the byte
// slice representation of the UUID.
func (uuid UUID) Value() (driver.Value, error) {
	return uuid.Bytes(), nil
}

// Scan implements SQL [driver.Scanner] interface, it expects the
// source to be:
//   - a byte slice, to be converted with [FromBytes].
//   - a string, to be converted with [FromString].
//   - nil, to be converted to the [Nil] UUID.
func (uuid *UUID) Scan(src any) error {
	switch src := src.(type) {
	case nil:
		*uuid = Nil

	case []byte:
		u, err := FromBytes(src)
		if err != nil {
			return err
		}
		*uuid = u

	case string:
		u, err := FromString(src)
		if err != nil {
			return err
		}
		*uuid = u

	default:
		return errBadSource
	}

	return nil
}

var (
	pool   [10 * 256]byte // pool of pseudo-random numbers.
	offset = len(pool)    // offset into the pool, its initialization makes the first call fill the pool.

	mu sync.Mutex
)

// next generates a 12bit and a 62bit cryptographically secure pseudo-random
// number, respectively.
func next() (uint64, uint64) {
	const (
		mask_62bit = (1 << 62) - 1
		mask_12bit = (1 << 12) - 1
	)

	mu.Lock()
	defer mu.Unlock()

	if offset >= len(pool) {
		rand.Read(pool[:])
		offset = 0
	}

	randA := binary.BigEndian.Uint16(pool[offset:]) & mask_12bit
	offset += 2

	randB := binary.BigEndian.Uint64(pool[offset:]) & mask_62bit
	offset += 8

	return uint64(randA), randB
}

func or(errs ...error) error {
	for _, err := range errs {
		if err != nil {
			return err
		}
	}
	return nil
}
