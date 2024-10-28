package middleware

import (
	"github.com/Gerfey/shortener/internal/app/compress"
	"net/http"
	"strings"
)

func GzipMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ow := w
		if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			next.ServeHTTP(w, r)
			return
		}

		contentType := r.Header.Get("Content-Type")
		if !strings.HasPrefix(contentType, "application/json") && !strings.HasPrefix(contentType, "text/html") {
			next.ServeHTTP(w, r)
			return
		}

		cw := compress.NewGzipWriter(w)
		ow = cw
		defer cw.Close()

		contentEncoding := r.Header.Get("Content-Encoding")
		if !strings.Contains(contentEncoding, "gzip") {
			cr, err := compress.NewGzipReader(r.Body)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			r.Body = cr
			defer cr.Close()
		}

		next.ServeHTTP(ow, r)
	})
}
