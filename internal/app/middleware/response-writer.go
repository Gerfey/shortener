package middleware

import (
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"
)

// responseWriter - обертка над стандартным http.ResponseWriter для отслеживания статус-кода и размера ответа
type responseWriter struct {
	http.ResponseWriter
	statusCode int
	size       int
}

// WriteHeader устанавливает код статуса HTTP-ответа и вызывает оригинальный метод
func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

// Write записывает данные в ответ и отслеживает их размер
func (rw *responseWriter) Write(b []byte) (int, error) {
	if rw.statusCode == 0 {
		rw.statusCode = http.StatusOK
	}
	size, err := rw.ResponseWriter.Write(b)
	rw.size += size
	return size, err
}

// LoggingMiddleware - middleware для логирования HTTP-запросов
func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		rw := &responseWriter{ResponseWriter: w}

		next.ServeHTTP(rw, r)

		duration := time.Since(start)

		log.WithFields(log.Fields{
			"uri":      r.RequestURI,
			"method":   r.Method,
			"duration": duration,
			"status":   rw.statusCode,
		}).Info("Handled request")
	})
}
