package handlers

import (
	"compress/gzip"
	"net/http"
	"strings"
)

type GzipHandler struct {
}

func (g *GzipHandler) GzipMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		if strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			wrw := NewWrapperResponseWriter(rw)
			wrw.Header().Set("Content-Encoding", "gzip")
			next.ServeHTTP(wrw, r)
			defer wrw.Flush()
			return
		}
		// handle normal
		next.ServeHTTP(rw, r)
	})
}

type WrapperResponseWriter struct {
	rw http.ResponseWriter
	gw *gzip.Writer
}

func NewWrapperResponseWriter(rw http.ResponseWriter) *WrapperResponseWriter {
	gw := gzip.NewWriter(rw)
	return &WrapperResponseWriter{
		rw: rw,
		gw: gw,
	}
}

func (w *WrapperResponseWriter) Header() http.Header {
	return w.rw.Header()
}

func (w *WrapperResponseWriter) Write(b []byte) (int, error) {
	return w.gw.Write(b)
}

func (w *WrapperResponseWriter) WriteHeader(code int) {
	w.rw.WriteHeader(code)
}

func (w *WrapperResponseWriter) Flush() {
	w.gw.Flush()
	w.gw.Close()
}
