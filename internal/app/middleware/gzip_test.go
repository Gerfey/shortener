package middleware

import (
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

	responseData := rr.Body.Bytes()
	assert.Equal(t, "Hello!", string(responseData), "Response body should match the original data")
}

func testHandler(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	write, err := w.Write(body)
	if err != nil {
		return
	}
}
