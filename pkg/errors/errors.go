// Copyright (C) 2025 Alan Barbosa Lima.
//
// Almodon is licensed under the GNU General Public License
// version 3. You should have received a copy of the
// license, located in LICENSE, at the root of the source
// tree. If not, see <https://www.gnu.org/licenses/>.

// Package errors implements especialized functionalities
// on the domain of Go's error handling. This package
// constructs over the foundation of the [error] interface
// and [errors] package.
package errors

import (
	"encoding/json"
	"fmt"
	"strconv"
)

// Error is an structured error type.
type Error struct {
	Kind    Kind
	Title   string
	Message string
	Cause   error
}

// New create a new error. Each new call generated a
// different error, regardless of the parameters.
func New(kind Kind, title, message string, cause error) error {
	return &Error{
		Kind:    kind,
		Title:   title,
		Message: message,
		Cause:   cause,
	}
}

// Error implements the [error] interface.
func (err *Error) Error() string {
	if err.Cause != nil {
		return err.Message + `: ` + err.Cause.Error()
	}

	return err.Message
}

// Unwrap returns the cause of the error, the cause might
// be nil.
func (err *Error) Unwrap() error {
	return err.Cause
}

// IsCLient identifies whether the error falls under the
// client category, see [Kind].
func (err *Error) IsClient() bool {
	return err.Kind.IsClient()
}

// IsCLient identifies whether the error falls under the
// internal category, see [Kind].
func (err *Error) IsInternal() bool {
	return err.Kind.IsInternal()
}

// MarshalJSON implements the [json.Marshaler] interface on
// the type, the cause is ommited if nil.
func (err Error) MarshalJSON() ([]byte, error) {
	var efj errorForJSON
	efj = errorForJSON(err)

	if efj.Cause != nil {
		efj.Cause = &wrapped{efj.Cause}
	}

	return json.Marshal(efj)
}

// UnmarshalJSON implements the [json.Unmarshaler]
// interface on the type, might fail if the cause cannot be
// unmarshalled into a sensible error type.
func (err *Error) UnmarshalJSON(buf []byte) error {
	var efj errorFromJSON
	if err := json.Unmarshal(buf, &efj); err != nil {
		return err
	}

	*err = Error{
		Kind:    efj.Kind,
		Title:   efj.Title,
		Message: efj.Message,
		Cause:   &efj.Cause,
	}
	return nil
}

// Kind in a enumeration of kinds of errors.
type Kind int

// client_errors_start is a sentinel value used to mark the start of the
// client-facing error kinds. It has no semantic meaning itself and should not
// be used as an actual error kind in returned errors; it exists only for grouping.

const (
	client_errors_start Kind = iota // This exists only for grouping.

	// InvalidInput signals that the request contains malformed or invalid
	// data. Use this for validation failures or when required fields are
	// missing.
	//
	// Rough HTTP equivalent: 400 Bad Request.
	InvalidInput

	// Unauthorized indicates that authentication is required or credentials
	// are missing/invalid. Use this when the client must authenticate before
	// accessing the resource.
	//
	// Rough HTTP equivalent: 401 Unauthorized.
	Unauthorized

	// Forbidden means the client is authenticated but does not have permission
	// to perform the requested operation. Use this for authorization failures.
	//
	// Rough HTTP equivalent: 403 Forbidden.
	Forbidden

	// PreconditionFailed denotes that a specified precondition (e.g. an
	// expected resource version, conditional header or media type) was not
	// met.
	//
	// Rough HTTP equivalent: 406 Not Acceptable, 412 Precondition Failed and
	// 415 Unsupported Media Type.
	PreconditionFailed

	// NotFound indicates that the requested resource does not exist. Use this
	// when a lookup by identifier yields no result.
	//
	// Rough HTTP equivalent: 404 Not Found.
	NotFound

	// Conflict represents a request that could not be completed because it
	// would cause a conflict with the current state of the target resource
	// (e.g. unique constraint violations, version conflicts).
	//
	// Rough HTTP equivalent: 409 Conflict.
	Conflict

	// Timeout denotes that the operation timed out before completion. Use
	// this when a request exceeds a predefined time limit.
	//
	// Rough HTTP equivalent: 408 Request Timeout.
	Timeout

	client_errors_end     // This exists only for grouping.
	internal_errors_start // This exists only for grouping.

	// Internal signals that an unexpected condition was encountered on the
	// server side. Use this for generic, non-client-related errors.
	//
	// Rough HTTP equivalent: 500 Internal Server Error.
	Internal

	// Unavailable indicates that the service is currently unavailable. Use
	// this when the server is overloaded or down for maintenance.
	//
	// Rough HTTP equivalent: 503 Service Unavailable.
	Unavailable

	// BadGateway indicates that an intermediary service received an invalid
	// response from an upstream service. Use this when acting as a proxy or
	// gateway.
	//
	// Rough HTTP equivalent: 502 Bad Gateway.
	BadGateway

	internal_errors_end // This exists only for grouping.
)

var kindStrings = map[Kind]string{
	InvalidInput:       "invalid input",
	Unauthorized:       "unauthorized",
	Forbidden:          "forbidden",
	PreconditionFailed: "precondition failed",
	NotFound:           "not found",
	Conflict:           "conflict",
	Timeout:            "timeout",

	Internal:    "internal error",
	Unavailable: "unavailable",
	BadGateway:  "bad gateway",
}

var stringKinds = invert(kindStrings)

// IsClient identifies whether the kind falls under the
// client category.
func (k Kind) IsClient() bool {
	return client_errors_start < k && k < client_errors_end
}

// IsInternal identifies whether the kind falls under the
// internal category.
func (k Kind) IsInternal() bool {
	return internal_errors_start < k && k < internal_errors_end
}

// String implements the [fmt.Stringer] interface on
// the type.
func (k Kind) String() string {
	return kindStrings[k]
}

// MarshalJSON implements the [json.Marshaler] interface on
// the type.
func (k Kind) MarshalJSON() ([]byte, error) {
	quoted := strconv.Quote(kindStrings[k])
	return []byte(quoted), nil
}

// UnmarshalJSON implements the [json.Unmarshaler] interface on
// the type.
func (k *Kind) UnmarshalJSON(buf []byte) error {
	unquoted, err := strconv.Unquote(string(buf))
	if err != nil {
		return fmt.Errorf("kind must be a valid JSON string: %w", err)
	}

	*k = stringKinds[unquoted]
	return nil
}

type errorFromJSON struct {
	Kind    Kind    `json:"kind"`
	Title   string  `json:"title"`
	Message string  `json:"message"`
	Cause   wrapped `json:"cause"`
}

type errorForJSON struct {
	Kind    Kind   `json:"kind"`
	Title   string `json:"title"`
	Message string `json:"message"`
	Cause   error  `json:"cause,omitempty"`
}

func invert[K, V comparable](m map[K]V) map[V]K {
	nm := make(map[V]K, len(m))
	for k, v := range m {
		nm[v] = k
	}

	return nm
}
