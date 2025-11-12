package resource

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"regexp"
	"sync"

	"github.com/alan-b-lima/almodon/internal/xerrors"
	"github.com/alan-b-lima/almodon/pkg/errors"
)

func WriteJsonError(w http.ResponseWriter, err error) {
	if err == nil {
		return
	}

	if err, ok := errors.AsType[*errors.Error](err); ok {
		writeJsonError(w, err, toHTTPStatus(err.Kind))
		return
	}

	writeJsonError(w, err, http.StatusInternalServerError)
}

func writeJsonError(w http.ResponseWriter, err error, status int) {
	body, e := json.Marshal(err)
	if e != nil {
		http.Error(w, e.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	w.Write(body)
}

var statusCodes = map[errors.Kind]int{
	errors.InvalidInput:       http.StatusBadRequest,
	errors.Unauthorized:       http.StatusUnauthorized,
	errors.Forbidden:          http.StatusForbidden,
	errors.PreconditionFailed: http.StatusPreconditionFailed,
	errors.NotFound:           http.StatusNotFound,
	errors.Conflict:           http.StatusConflict,
	errors.Timeout:            http.StatusRequestTimeout,

	errors.Internal:    http.StatusInternalServerError,
	errors.Unavailable: http.StatusServiceUnavailable,
	errors.BadGateway:  http.StatusBadGateway,
}

func toHTTPStatus(kind errors.Kind) int {
	if status, in := statusCodes[kind]; in {
		return status
	}

	return http.StatusInternalServerError
}

var (
	reContentTypeApplicationJson = regexp.MustCompile(`^\s*(\*/\*|application/(json|\*))\s*(;.*)?\s*$`)
	reAcceptApplicationJson      = regexp.MustCompile(`(^|.*,)\s*(\*/\*|application/(json|\*))\s*(;.*)?\s*($|,.*)`)
)

func DecodeJSON(req any, r *http.Request) error {
	contentType := r.Header.Get("Content-Type")
	if contentType == "" {
		return xerrors.ErrNoContentType
	}

	if !reContentTypeApplicationJson.MatchString(contentType) {
		return xerrors.ErrNoContentType
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		switch err := err.(type) {
		case *json.SyntaxError:
			return xerrors.ErrJsonSyntax.New(err.Offset)

		case *json.UnmarshalTypeError:
			return xerrors.ErrJsonType.New(err.Offset, err.Type.Kind(), err.Value)

		}

		return err
	}

	return nil
}

var bufPool = sync.Pool{New: func() any { return new(bytes.Buffer) }}

func EncodeJSON(res any, status int, w http.ResponseWriter, r *http.Request) error {
	accept := r.Header.Get("Accept")
	if !reAcceptApplicationJson.MatchString(accept) {
		return xerrors.ErrNotAcceptableJson
	}

	b := bufPool.Get().(*bytes.Buffer)
	defer bufPool.Put(b)
	b.Reset()

	if err := json.NewEncoder(b).Encode(res); err != nil {
		return err
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)

	if _, err := io.Copy(w, b); err != nil {
		return err
	}

	return nil
}
