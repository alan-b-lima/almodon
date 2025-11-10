// Copyright (C) 2025 Alan Barbosa Lima.
//
// Almodon is licensed under the GNU General Public License
// version 3. You should have received a copy of the
// license, located in LICENSE, at the root of the source
// tree. If not, see <https://www.gnu.org/licenses/>.

package errors

import (
	"encoding/json"
	"errors"
	"fmt"
)

// AsType is a generic function that replaces the common
// boilerplate for using the [errors.As] function.
func AsType[T error](err error) (T, bool) {
	var target T
	ok := errors.As(err, &target)
	return target, ok
}

// Wrap wraps an error into a type that implements JSON
// marshalling and unmarshalling.
func Wrap(err error) error {
	return &wrapped{err}
}

type wrapped struct{ error }

var errUnmarshal = errors.New("errors: failed to unmarsheled error into a sensible type")

// Wrap wraps an error into a type that implements JSON marshalling
// and unmarshalling.
func (e *wrapped) Error() string {
	return e.error.Error()
}

// String implements the [fmt.Stringer] interface on the type.
func (e *wrapped) String() string {
	return e.error.Error()
}

// Unwrap returns the wrapped error.
func (e *wrapped) Unwrap() error {
	return e.error
}

// MarshalJSON implements the [json.Marshaler] interface on the type.
func (e wrapped) MarshalJSON() ([]byte, error) {
	switch err := e.error.(type) {
	case *Multi:
		errs := make([]wrapped, 0, len(err.errs))
		for _, err := range err.errs {
			errs = append(errs, wrapped{err})
		}

		return json.Marshal(errs)

	case json.Marshaler:
		return err.MarshalJSON()

	case fmt.Stringer:
		return json.Marshal(err.String())
	}

	return json.Marshal(e.Error())
}

// UnmarshalJSON implements the [json.Unmarshaler] interface on the
// type. It tries to unmarshal the error into an [Error], if that
// fails, it tries to unmarshal it into a string, if that also fails,
// it returns an error.
func (e *wrapped) UnmarshalJSON(buf []byte) error {
	var err_ Error
	if err := json.Unmarshal(buf, &err_); err == nil {
		e.error = &err_
		return nil
	}

	var str_ string
	if err := json.Unmarshal(buf, &str_); err == nil {
		e.error = errors.New(str_)
		return nil
	}

	return errUnmarshal
}
