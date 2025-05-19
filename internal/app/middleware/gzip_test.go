package middleware

import (
	"bytes"
	"compress/gzip"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGzipMiddlewareNoCompression(t *testing.T) {
	next := http.HandlerFunc(testHandler)
	handler := GzipMiddleware(next)

	req := httptest.NewRequest("GET", "http://example.com/test", strings.NewReader("Hello!"))
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	assert.Empty(t, rr.Header().Get("Content-Encoding"), "Content-Encoding header should be empty")
	assert.Equal(t, http.StatusOK, rr.Code, "Status code should be 200 OK")
	assert.Equal(t, "Hello!", rr.Body.String(), "Response body should match the original data")
}

func TestGzipMiddlewareWithCompression(t *testing.T) {
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("Hello!"))
	})
	handler := GzipMiddleware(next)

	req := httptest.NewRequest("GET", "http://example.com/test", nil)
	req.Header.Set("Accept-Encoding", "gzip")

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	assert.Equal(t, "gzip", rr.Header().Get("Content-Encoding"), "Content-Encoding header should be gzip")
	assert.Equal(t, http.StatusOK, rr.Code, "Status code should be 200 OK")

	reader, err := gzip.NewReader(rr.Body)
	assert.NoError(t, err)
	defer func() { _ = reader.Close() }()

	content, err := io.ReadAll(reader)
	assert.NoError(t, err)
	assert.Equal(t, "Hello!", string(content))
}

func TestGzipMiddlewareWithGzipRequest(t *testing.T) {
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Failed to read body", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(body)
	})
	handler := GzipMiddleware(next)

	var buf bytes.Buffer
	gz := gzip.NewWriter(&buf)
	_, _ = gz.Write([]byte("Hello!"))
	_ = gz.Close()

	req := httptest.NewRequest("POST", "http://example.com/test", &buf)
	req.Header.Set("Content-Encoding", "gzip")

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code, "Status code should be 200 OK")
	assert.Equal(t, "Hello!", rr.Body.String(), "Response body should match the original data")
}

func testHandler(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read body", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(body)
}
