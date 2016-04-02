package util

import (
	"io"
	"net/http"
)

type httpResponseWriter struct {
	io.Writer
}

func (httpResponseWriter) Header() http.Header {
	return make(map[string][]string)
}
func (httpResponseWriter) WriteHeader(int) {}

// 将 io.Writer 从语义上实现 http.ResponseWriter.
func HttpResponseWriter(w io.Writer) http.ResponseWriter {
	if rw, ok := w.(http.ResponseWriter); ok {
		return rw
	}
	return httpResponseWriter{Writer: w}
}
