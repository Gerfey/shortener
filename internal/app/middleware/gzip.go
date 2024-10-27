package middleware

import (
	"compress/gzip"
	"net/http"
	"strings"
)

type gzipResponseWriter struct {
	http.ResponseWriter
	Writer *gzip.Writer
	data   []byte
}

func GzipMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			next.ServeHTTP(w, r)
			return
		}

		gzipWriter := &gzipResponseWriter{ResponseWriter: w, Writer: nil}
		next.ServeHTTP(gzipWriter, r)

		contentType := w.Header().Get("Content-Type")
		if !strings.HasPrefix(contentType, "application/json") && !strings.HasPrefix(contentType, "text/html") {
			return
		}

		if gzipWriter.Writer == nil {
			w.Header().Set("Content-Encoding", "gzip")
			gz := gzip.NewWriter(w)
			gzipWriter.Writer = gz
			defer gz.Close()
		}

		gzipWriter.Writer.Write(gzipWriter.data)
	})
}

func (w *gzipResponseWriter) Write(b []byte) (int, error) {
	w.data = append(w.data, b...)
	return len(b), nil
}

func (w *gzipResponseWriter) WriteHeader(statusCode int) {
	w.ResponseWriter.Header().Del("Content-Length")
	w.ResponseWriter.WriteHeader(statusCode)
}
