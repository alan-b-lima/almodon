// Copyright (C) 2025 Alan Barbosa Lima.
//
// Almodon is licensed under the GNU General Public License
// version 3. You should have received a copy of the
// license, located in LICENSE, at the root of the source
// tree. If not, see <https://www.gnu.org/licenses/>.

// This package implements a UUID generator and some other
// functionalities revolving UUIDs. Read in RFC9562 for more info
// about UUIDs.
package uuid

import (
	crand "crypto/rand"
	"errors"
	"fmt"
	"math/rand/v2"
	"sync"
	"time"
	"unsafe"
)

// This represents a 128bit UUID type. Elements of this type can, and
// should, be compared using the == operator. The zero value of UUID
// is classified as a Nil UUID according to RFC9562 section 5.9.
type UUID [16]byte

const (
	_62BitMask = (1 << 62) - 1
	_48BitMask = (1 << 48) - 1
	_12BitMask = (1 << 12) - 1
)

var (
	source rand.Source // The source for all pseudo-random number needed.
	mulock sync.Mutex  // A mutex for safe concurrent UUID generation.
)

func init() {
	var seed [2]uint64
	crand.Read(unsafe.Slice((*byte)(unsafe.Pointer(&seed)), 16))

	source = rand.NewPCG(seed[0], seed[1])
}

var (
	ErrBadSliceLength = errors.New("uuid: slice does not has 16 bytes")
	ErrBadString      = errors.New("uuid: string could not be parsed correctly")
	ErrBadJSONString  = errors.New("uuid: slice is a malformed JSON string")
)

var _UUIDFormat = "%02x%02x%02x%02x-%02x%02x-%02x%02x-%02x%02x-%02x%02x%02x%02x%02x%02x"

// Generates a new UUID accourding to version 7. It's safe to call
// this function from multiple goroutines.
//
// The memory layout of a UUIDv7 is as follows:
//
//   - Unix Timestamp: 48-bit big-endian unsigned number of the Unix
//     Epoch timestamp in milliseconds. Occupies bits 0 through 47,
//     octets 1 through 5.
//
//   - Version: 4-bit version field, set to 0b0111 (7). Occupies bits
//     48 through 51, octet 6.
//
//   - Random A: 12-bit pseudorandom data to provide uniqueness.
//     Occupies bits 52 through 63, octects 6 through 7.
//
//   - Variant: 2-bit variant field, set to 0b10. Occupies bits 64
//     through 65, octet 8.
//
//   - Random B: 62-bit pseudorandom data to provide uniqueness.
//     Occupies bits 66 through 127, octets 8 through 15.
func NewUUIDv7() UUID {
	const (
		version = 0b0111
		variant = 0b01
	)

	unixTimestamp := uint64(time.Now().UnixMilli() & _48BitMask)

	mulock.Lock()
	randA := source.Uint64() & _12BitMask
	randB := source.Uint64() & _62BitMask
	mulock.Unlock()

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

// Converts an UUID from a byte slice. The given byte slice should be
// of length 16, otherwise an error will be returned. Due to this,
// this function is NOT interchangeable with [FromString], a byte
// slice that represents the string representation of an UUID should
// be converted with [FromString].
func FromBytes(bytes []byte) (UUID, error) {
	if len(bytes) != 16 {
		return UUID{}, ErrBadSliceLength
	}

	return UUID(bytes), nil
}

// Converts an UUID from a string format. Note that this function is
// NOT interchangeable with [FromBytes], even considering convertion.
func FromString(str string) (UUID, error) {
	if len(str) != 36 {
		return UUID{}, ErrBadString
	}

	var uuid UUID
	n, err := fmt.Sscanf(str, _UUIDFormat,
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

// Verifies whether the given UUID is the Nil UUID. Not be confused
// with a nil pointer to an UUID. Equivalent to using == against the
// zero value of this type.
func (uuid UUID) IsNil() bool {
	return uuid == UUID{}
}

// Implements the interface [fmt.Stringer] on the UUID type.
func (uuid UUID) String() string {
	return fmt.Sprintf(_UUIDFormat,
		uuid[0], uuid[1], uuid[2], uuid[3],
		uuid[4], uuid[5],
		uuid[6], uuid[7],
		uuid[8], uuid[9],
		uuid[10], uuid[11], uuid[12], uuid[13], uuid[14], uuid[15],
	)
}

// Implements the interface [json.Marshaler] on the UUID type.
func (uuid UUID) MarshalJSON() ([]byte, error) {
	return []byte(`"` + uuid.String() + `"`), nil
}

// Implements the interface [json.Unmarshaler] on the UUID type. The
// given byte slice should be a valid JSON string literal.
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
