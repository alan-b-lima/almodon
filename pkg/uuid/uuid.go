// Copyright (C) 2025-2026 Alan Barbosa Lima.
//
// Almodon is licensed under the GNU General Public License
// version 3. You should have received a copy of the
// license, located in LICENSE, at the root of the source
// tree. If not, see <https://www.gnu.org/licenses/>.

// Package uuid implements an UUID version 7 generator, and parsing
// functionalities regarding UUID. This package is complient to [RFC9562].
//
// [RFC9562]: https://www.rfc-editor.org/rfc/rfc9562
package uuid

import (
	"crypto/rand"
	"database/sql/driver"
	"errors"
	"fmt"
	"sync"
	"time"
)

// UUID is the set of all 128-bit universal unique identifiers.
//
// Elements of this type have the == operator (equality operator)
// well defined, clients can and should use == to compare UUID's.
//
// The zero value of the UUID type it the Nil UUID.
type UUID [16]byte

// Nil is the zero value of UUID, according to RFC9562 section 5.9.
var Nil UUID

var (
	ErrBadSliceLength = errors.New("uuid: slice does not has 16 bytes")
	ErrBadString      = errors.New("uuid: string could not be parsed correctly")
	ErrBadJSONString  = errors.New("uuid: slice is a malformed JSON string")
	ErrBadSQLSource   = errors.New("uuid: source cannot be converted to an UUID")
)

var _Format = "%02x%02x%02x%02x-%02x%02x-%02x%02x-%02x%02x-%02x%02x%02x%02x%02x%02x"

// NewUUIDv7 generates a new UUID accourding to version 7. It's safe
// to call this function from multiple goroutines.
//
// The memory layout of a UUIDv7, as defined in [RFC9562], is as follows:
//
//   - Unix Timestamp: 48-bit big-endian unsigned number of the Unix
//     Epoch timestamp in milliseconds. Occupies bits 0 through 47,
//     octets 1 through 5.
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
// [RFC9562]: https://www.rfc-editor.org/rfc/rfc9562
func NewUUIDv7() UUID {
	const (
		mask_48bit = (1 << 48) - 1

		version = 0b0111
		variant = 0b10
	)

	unixTimestamp := uint64(time.Now().UnixMilli() & mask_48bit)
	randA, randB := next()

	return UUID{
		0x0: byte(unixTimestamp >> 0x28), 0x1: byte(unixTimestamp >> 0x20),
		0x2: byte(unixTimestamp >> 0x18), 0x3: byte(unixTimestamp >> 0x10),
		0x4: byte(unixTimestamp >> 0x08), 0x5: byte(unixTimestamp >> 0x00),

		0x6: version<<4 | byte(randA>>8),
		0x7: byte(randA),

		0x8: variant<<6 | byte(randB>>0x38),
		0x9: byte(randB >> 0x30), 0xA: byte(randB >> 0x28),
		0xB: byte(randB >> 0x20), 0xC: byte(randB >> 0x18),
		0xD: byte(randB >> 0x10), 0xE: byte(randB >> 0x08),
		0xF: byte(randB >> 0x00),
	}
}

// IsZero reports whether the given UUID is the Nil UUID.
func (uuid UUID) IsZero() bool {
	return uuid == Nil
}

// FromBytes converts an UUID from a byte slice. The given byte slice
// should be of length 16, otherwise an error will be returned.
//
// Due to this, this function is NOT interchangeable with
// [FromString], a byte slice that is the string representation of an
// UUID should be converted with [FromString].
func FromBytes(bytes []byte) (UUID, error) {
	if len(bytes) != 16 {
		return UUID{}, ErrBadSliceLength
	}

	return UUID(bytes), nil
}

// FromString converts an UUID from the string format, in a
// case-insensitive manner.
//
// Note that this function is NOT interchangeable with [FromBytes],
// see [FromBytes] for more detail.
func FromString(str string) (UUID, error) {
	if len(str) != 36 {
		return UUID{}, ErrBadString
	}

	var uuid UUID
	n, err := fmt.Sscanf(str, _Format,
		&uuid[0], &uuid[1], &uuid[2], &uuid[3],
		&uuid[4], &uuid[5],
		&uuid[6], &uuid[7],
		&uuid[8], &uuid[9],
		&uuid[10], &uuid[11], &uuid[12], &uuid[13], &uuid[14], &uuid[15],
	)
	if err != nil {
		return UUID{}, err
	}
	if n != 16 {
		return UUID{}, ErrBadString
	}

	return uuid, nil
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
	return fmt.Sprintf(_Format,
		uuid[0], uuid[1], uuid[2], uuid[3],
		uuid[4], uuid[5],
		uuid[6], uuid[7],
		uuid[8], uuid[9],
		uuid[10], uuid[11], uuid[12], uuid[13], uuid[14], uuid[15],
	)
}

// MarshalJSON implements the interface [json.Marshaler] on the UUID
// type.
func (uuid UUID) MarshalJSON() ([]byte, error) {
	return []byte(`"` + uuid.String() + `"`), nil
}

// UnmarshalJSON implements the interface [json.Unmarshaler] on the
// UUID type. The given byte slice should be a valid JSON string
// literal.
func (uuid *UUID) UnmarshalJSON(buf []byte) error {
	if len(buf) >= 2 && (buf[0] != '"' || buf[len(buf)-1] != '"') {
		return ErrBadJSONString
	}

	decoded, err := FromString(string(buf[1 : len(buf)-1]))
	if err != nil {
		return err
	}

	*uuid = decoded
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
//   - nil, to be converted to the Nil UUID.
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
		return ErrBadSQLSource
	}

	return nil
}

var (
	pool   [10 * 256]byte // pool of pseudo-random numbers.
	offset = len(pool)    // offset into the pool, its initialization makes the first call fill the pool.

	mu sync.Mutex
)

// next generates a 12bit and a 62bit pseudo-random number, respectively.
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

	var randA uint64
	randA |= uint64(pool[offset+0]) << 0x8
	randA |= uint64(pool[offset+1]) << 0x0

	randA &= mask_12bit
	offset += 2

	var randB uint64
	randB |= uint64(pool[offset+0]) << 0x38
	randB |= uint64(pool[offset+1]) << 0x30
	randB |= uint64(pool[offset+2]) << 0x28
	randB |= uint64(pool[offset+3]) << 0x20
	randB |= uint64(pool[offset+4]) << 0x18
	randB |= uint64(pool[offset+5]) << 0x10
	randB |= uint64(pool[offset+6]) << 0x08
	randB |= uint64(pool[offset+7]) << 0x00

	randB &= mask_62bit
	offset += 8

	return randA, randB
}
