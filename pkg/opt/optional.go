// Copyright (C) 2025 Alan Barbosa Lima.
//
// Almodon is licensed under the GNU General Public License
// version 3. You should have received a copy of the
// license, located in LICENSE, at the root of the source
// tree. If not, see <https://www.gnu.org/licenses/>.

// Package opt implements a optional type, it may or may not contatin
// a value. It can also be marshalled and unmarshalled with json,
// using the undelying marshelers, if such exist.
package opt

import (
	"bytes"
	"encoding/json"
	"fmt"
)

// Opt is an optional type, it may or may not contatin a value of a
// certain type, given by its type parameter. Its zero value is
// called None, when printed: <none>, it indicates the absence of
// value.
//
// As it stands, an arbitrary [Opt] can be marshalled and
// unmarshalled with JSON, with JSON's null being considered the zero
// value, regardless of the type param.
type Opt[T any] struct {
	val  T
	some bool
}

const _JSONNull = `null`

var _JSONNullBytes = []byte(_JSONNull)

// Some creates a new Opt value with a value.
func Some[T any](val T) Opt[T] {
	return Opt[T]{val, true}
}

// None creates a new Opt value with no value, ie, None.
func None[T any]() Opt[T] {
	return Opt[T]{}
}

// Unwrap unpacks the Opt struct and returns its components, a common
// assertion should be done, it may be done in the following way:
//
//	val, ok := opt.Unwrap()
//	if !ok {
//		// handle noneness
//	}
//
// Since this a non-build-in function, the ok return cannot be
// ommited.
func (o Opt[T]) Unwrap() (T, bool) {
	return o.val, o.some
}

// MarshalJSON implements the [json.Marshaler] interface, it marshals
// the Opt struct, if it is None, it returns JSON's null, otherwise
// it tries to marshal the underlying value, if it fails, the error
// is returned.
func (o Opt[T]) MarshalJSON() ([]byte, error) {
	if !o.some {
		return []byte(_JSONNull), nil
	}

	return json.Marshal(o.val)
}

// UnmarshalJSON implements the [json.Unmarshaler] interface, it
// unmarshals the value into the Opt struct, if the input is JSON's
// null, the Opt is set to None, otherwise it tries to unmarshal the
// value into the underlying type, if it fails, the Opt is set to
// None and the error is returned.
func (o *Opt[T]) UnmarshalJSON(b []byte) error {
	o.some = !bytes.Equal(_JSONNullBytes, b)
	if !o.some {
		return nil
	}

	if err := json.Unmarshal(b, &o.val); err != nil {
		o.some = false
		return err
	}

	return nil
}

// String implements the [fmt.Stringer] interface, it returns the
// string "<none>" if the Opt is None, otherwise it returns the
// string representation of the underlying value, using [fmt.Sprint].
func (o Opt[T]) String() string {
	if !o.some {
		return "<none>"
	}

	return fmt.Sprint(o.val)
}
