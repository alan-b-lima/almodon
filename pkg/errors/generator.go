// Copyright (C) 2025 Alan Barbosa Lima.
//
// Almodon is licensed under the GNU General Public License
// version 3. You should have received a copy of the
// license, located in LICENSE, at the root of the source
// tree. If not, see <https://www.gnu.org/licenses/>.

package errors

import "fmt"

type gen struct {
	kind  Kind
	title string
}

// Gen creates a new error generator. The returned generator
// can be used to create multiple errors with the same
// kind and title, but different messages and causes.
func Gen(kind Kind, title string) *gen {
	return &gen{
		kind:  kind,
		title: title,
	}
}

// New creates a new error with the given message and cause.
func (gen *gen) New(message string, cause error) error {
	return &Error{
		Kind:    gen.kind,
		Title:   gen.title,
		Message: message,
		Cause:   cause,
	}
}

type imp struct {
	kind    Kind
	title   string
	message string
}

// Imp creates a new error implementer. The returned implementer
// can be used to create multiple errors with the same
// kind, title, and message, but different causes.
func Imp(kind Kind, title, message string) *imp {
	return &imp{
		kind:    kind,
		title:   title,
		message: message,
	}
}

// New creates a new error with the given cause.
func (gen *imp) New(cause error) error {
	return &Error{
		Kind:    gen.kind,
		Title:   gen.title,
		Message: gen.message,
		Cause:   cause,
	}
}

type fmt_ struct {
	kind   Kind
	title  string
	format string
}

// Fmt creates a new error formatter. The returned formatter
// can be used to create multiple errors with the same
// kind and title, but different formatted messages.
//
// The format follows the same rules as [fmt.Errorf].
//
// The cause is always nil.
func Fmt(kind Kind, title, format string) *fmt_ {
	return &fmt_{
		kind:   kind,
		title:  title,
		format: format,
	}
}

// New creates a new error with the given formatted message.
func (gen *fmt_) New(v ...any) error {
	err := fmt.Errorf(gen.format, v...)

	return &Error{
		Kind:    gen.kind,
		Title:   gen.title,
		Message: err.Error(),
		Cause:   nil,
	}
}
