package middleware

import (
	"errors"
	"net/http"
)

type ResponseWriter struct {
	w     http.ResponseWriter
	status int
}

func NewW(w http.ResponseWriter) *ResponseWriter {
	return &ResponseWriter{w: w, status: 200}
}

func (rw *ResponseWriter) WriteHeader(status int) {
	rw.w.WriteHeader(status)
	rw.status = status
}

func (rw *ResponseWriter) Header() http.Header           { return rw.w.Header() }
func (rw *ResponseWriter) Write(buf []byte) (int, error) { return rw.w.Write(buf) }
func (rw *ResponseWriter) StatusCode() int               { return rw.status }

type ResponseWriterFlusher struct {
	wf interface {
		http.ResponseWriter
		http.Flusher
	}
	status int
}

var errNoFlush = errors.New("flushing is not supported")

func NewWF(w http.ResponseWriter) (*ResponseWriterFlusher, error) {
	wf, ok := w.(interface {
		http.ResponseWriter
		http.Flusher
	})
	if !ok {
		return nil, errNoFlush
	}

	return &ResponseWriterFlusher{wf: wf, status: 200}, nil
}


func (rw *ResponseWriterFlusher) WriteHeader(status int) {
	rw.wf.WriteHeader(status)
	rw.status = status
}

func (rw *ResponseWriterFlusher) Header() http.Header           { return rw.wf.Header() }
func (rw *ResponseWriterFlusher) Write(buf []byte) (int, error) { return rw.wf.Write(buf) }
func (rw *ResponseWriterFlusher) Flush()                        { rw.wf.Flush() }
