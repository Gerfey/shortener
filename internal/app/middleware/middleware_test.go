package middleware

import (
	log "github.com/sirupsen/logrus"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type TestHook struct {
	Entries []*log.Entry
}

func (hook *TestHook) Levels() []log.Level {
	return log.AllLevels
}

func (hook *TestHook) Fire(entry *log.Entry) error {
	hook.Entries = append(hook.Entries, entry)
	return nil
}

func TestLoggingMiddleware(t *testing.T) {
	hook := &TestHook{}
	log.AddHook(hook)

	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(10 * time.Millisecond)
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("Hello, World!"))
	})

	handler := LoggingMiddleware(testHandler)

	req := httptest.NewRequest("GET", "/test-uri", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)
	defer func() { _ = resp.Body.Close() }()

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "Hello, World!", string(body))

	assert.Equal(t, 1, len(hook.Entries), "Expected one log entry")

	entry := hook.Entries[0]

	assert.Equal(t, log.InfoLevel, entry.Level, "Expected log level to be Info")

	assert.Equal(t, "Handled request", entry.Message, "Unexpected log message")

	assert.Equal(t, "/test-uri", entry.Data["uri"], "Unexpected URI in log")
	assert.Equal(t, "GET", entry.Data["method"], "Unexpected method in log")
	assert.Equal(t, http.StatusOK, entry.Data["status"], "Unexpected status code in log")

	duration, ok := entry.Data["duration"].(time.Duration)
	assert.True(t, ok, "Duration should be of type time.Duration")
	assert.GreaterOrEqual(t, duration, 10*time.Millisecond, "Duration should be at least 10ms")
}

func TestLoggingMiddlewareDifferentStatuses(t *testing.T) {
	hook := &TestHook{}
	log.AddHook(hook)

	notFoundHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte("Not Found"))
	})

	handler := LoggingMiddleware(notFoundHandler)

	req := httptest.NewRequest("GET", "/not-found", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)
	defer func() { _ = resp.Body.Close() }()

	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	assert.Equal(t, "Not Found", string(body))
	assert.Equal(t, 1, len(hook.Entries), "Expected one log entry")

	entry := hook.Entries[0]
	assert.Equal(t, log.InfoLevel, entry.Level)
	assert.Equal(t, "Handled request", entry.Message)
	assert.Equal(t, "/not-found", entry.Data["uri"])
	assert.Equal(t, "GET", entry.Data["method"])
	assert.Equal(t, http.StatusNotFound, entry.Data["status"])

	duration, ok := entry.Data["duration"].(time.Duration)
	assert.True(t, ok, "Duration should be of type time.Duration")
	assert.GreaterOrEqual(t, duration, 0*time.Millisecond, "Duration should be non-negative")
}
