// Copyright (C) 2025 Alan Barbosa Lima.
//
// Almodon is licensed under the GNU General Public License
// version 3. You should have received a copy of the
// license, located in LICENSE, at the root of the source
// tree. If not, see <https://www.gnu.org/licenses/>.

package errors

import (
	"encoding/json"
	"strings"
)

// Multi is an error that represents multiple errors.
type Multi struct{ errs []error }

// Join combines multiple errors into a single error. Nil errors are
// ignored. If all errors are nil, Join returns nil.
func Join(errs ...error) error {
	var n int
	for _, err := range errs {
		if err != nil {
			n++
		}
	}
	if n == 0 {
		return nil
	}

	merr := Multi{errs: make([]error, 0, n)}
	for _, err := range errs {
		if err != nil {
			merr.errs = append(merr.errs, err)
		}
	}

	return &merr
}

// Error implements the [error] interface.
func (e *Multi) Error() string {
	if e == nil || len(e.errs) == 0 {
		return "<nil>"
	}

	if len(e.errs) == 1 {
		return e.errs[0].Error()
	}

	var b strings.Builder
	b.WriteString(e.errs[0].Error())
	for _, err := range e.errs[1:] {
		b.WriteByte('\n')
		b.WriteString(err.Error())
	}

	return b.String()
}

// MarshalJSON implements the [json.Marshaler] interface on the type.
func (e Multi) MarshalJSON() ([]byte, error) {
	return json.Marshal(e.errs)
}

// UnmarshalJSON implements the [json.Unmarshaler] interface on the
// type.
func (e *Multi) UnmarshalJSON(buf []byte) error {
	var errs []error
	if err := json.Unmarshal(buf, &errs); err != nil {
		return err
	}

	e.errs = errs
	return nil
}
