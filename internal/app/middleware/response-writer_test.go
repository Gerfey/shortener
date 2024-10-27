package middleware

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestWriteHeader(t *testing.T) {
	tests := []struct {
		name       string
		headerCode int
		expectCode int
	}{
		{
			name:       "OK Status",
			headerCode: http.StatusOK,
			expectCode: http.StatusOK,
		},
		{
			name:       "Not Found Status",
			headerCode: http.StatusNotFound,
			expectCode: http.StatusNotFound,
		},
		{
			name:       "Server Error Status",
			headerCode: http.StatusInternalServerError,
			expectCode: http.StatusInternalServerError,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			recorder := httptest.NewRecorder()
			writer := &responseWriter{
				ResponseWriter: recorder,
			}

			writer.WriteHeader(test.headerCode)

			if writer.statusCode != test.expectCode {
				t.Errorf("expected statusCode to be %d, but got %d", test.expectCode, writer.statusCode)
			}

			if result := recorder.Result(); result.StatusCode != test.expectCode {
				t.Errorf("expected ResponseWriter StatusCode to be %d, but got %d", test.expectCode, result.StatusCode)
			}
		})
	}
}

func TestWrite(t *testing.T) {
	tests := []struct {
		name     string
		input    []byte
		expect   []byte
		hasError bool
	}{
		{
			name:     "Write some bytes",
			input:    []byte("Some data here"),
			expect:   []byte("Some data here"),
			hasError: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			recorder := httptest.NewRecorder()
			writer := &responseWriter{
				ResponseWriter: recorder,
			}

			_, err := writer.Write(test.input)

			if (!test.hasError && err != nil) || (test.hasError && err == nil) {
				t.Errorf("Expected error to be %v but got %v", test.hasError, err != nil)
			}

			if got := recorder.Body.Bytes(); !bytes.Equal(got, test.expect) {
				t.Errorf("Expected written bytes to be %s but got %s", string(test.expect), string(got))
			}
		})
	}
}
