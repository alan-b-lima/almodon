package resource

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"regexp"
	"sync"

	"github.com/alan-b-lima/pkg/problem"
)

type ResponseWriter struct {
	http.ResponseWriter

	status int
	err    error
}

func NewResponseWriter(w http.ResponseWriter) *ResponseWriter {
	return &ResponseWriter{ResponseWriter: w, status: 200, err: nil}
}

var _ http.ResponseWriter = (*ResponseWriter)(nil)

func (rw *ResponseWriter) WriteHeader(status int) {
	rw.ResponseWriter.WriteHeader(status)
	rw.status = status
}

func (rw *ResponseWriter) StatusCode() int { return rw.status }
func (rw *ResponseWriter) Error() error    { return rw.err }

func (rw *ResponseWriter) Unwrap() http.ResponseWriter {
	return rw.ResponseWriter
}

func WriteError(w http.ResponseWriter, err error) {
	if err == nil {
		return
	}

	if err, ok := errors.AsType[*problem.Error](err); ok {
		writeErrorJson(w, err, int(err.Kind))
		return
	}

	writeErrorJson(w, err, http.StatusInternalServerError)
}

func writeErrorJson(w http.ResponseWriter, err error, status int) {
	b := buffers.Get().(*bytes.Buffer)
	defer buffers.Put(b)
	b.Reset()

	if rw, ok := w.(*ResponseWriter); ok {
		rw.err = err
	}

	if err := json.NewEncoder(b).Encode(err); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	w.Write(b.Bytes())
}

var (
	reContentTypeApplicationJson = regexp.MustCompile(`^\s*(\*/\*|application/(json|\*))\s*(;.*)?\s*$`)
	reAcceptApplicationJson      = regexp.MustCompile(`(^|.*,)\s*(\*/\*|application/(json|\*))\s*(;.*)?\s*($|,.*)`)
)

func DecodeJSON(r *http.Request, req any) error {
	contentType := r.Header.Get("Content-Type")
	if contentType == "" {
		return ErrNoContentType
	}

	if !reContentTypeApplicationJson.MatchString(contentType) {
		return ErrUnsupportedContentType.Make(contentType)
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		if err == io.EOF {
			return ErrJSON.Message("unexpected end of input").Make()
		}

		return ErrJSON.Cause(err).Make()
	}

	return nil
}

var buffers = sync.Pool{New: func() any { return new(bytes.Buffer) }}

func EncodeJSON(w http.ResponseWriter, r *http.Request, res any, status int) error {
	accept := r.Header.Get("Accept")
	if !reAcceptApplicationJson.MatchString(accept) {
		return ErrNotAcceptable.Make("application/json")
	}

	b := buffers.Get().(*bytes.Buffer)
	defer buffers.Put(b)
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

func NotFound(w http.ResponseWriter, r *http.Request) {
	WriteError(w, ErrNotFound.Make(r.URL.Path))
}
