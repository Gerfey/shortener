package compress

import (
	"bytes"
	"compress/gzip"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGzipWriter(t *testing.T) {
	writer := httptest.NewRecorder()

	gw := NewGzipWriter(writer)

	testData := []byte("Hello, Gzip!")

	n, err := gw.Write(testData)
	assert.NoError(t, err, "Write should not return an error")
	assert.Equal(t, len(testData), n, "Write should return the number of bytes written")

	err = gw.Close()
	assert.NoError(t, err, "Close should not return an error")

	compressedData := writer.Body.Bytes()

	gr, err := gzip.NewReader(bytes.NewReader(compressedData))
	assert.NoError(t, err, "Failed to create gzip reader")

	defer gr.Close()

	decompressedData, err := io.ReadAll(gr)
	assert.NoError(t, err, "Failed to read decompressed data")

	assert.Equal(t, testData, decompressedData, "Decompressed data should match original data")
}

func TestGzipWriter_WriteHeader(t *testing.T) {
	writer := httptest.NewRecorder()

	gw := NewGzipWriter(writer)

	statusCode := http.StatusOK
	gw.WriteHeader(statusCode)

	assert.Equal(t, "gzip", writer.Header().Get("Content-Encoding"), "Content-Encoding header should be 'gzip'")

	assert.Equal(t, statusCode, writer.Code, "Status code should match the one set")
}

func TestGzipReader(t *testing.T) {
	originalData := []byte("Hello, Gzip Reader!")

	var compressedBuf bytes.Buffer
	gw := gzip.NewWriter(&compressedBuf)
	_, err := gw.Write(originalData)
	assert.NoError(t, err, "Failed to write to gzip.Writer")

	err = gw.Close()
	assert.NoError(t, err, "Failed to close gzip.Writer")

	compressedReader := io.NopCloser(bytes.NewReader(compressedBuf.Bytes()))

	gr, err := NewGzipReader(compressedReader)
	assert.NoError(t, err, "NewGzipReader should not return an error")

	decompressedData, err := io.ReadAll(gr)
	assert.NoError(t, err, "Read should not return an error")

	assert.Equal(t, originalData, decompressedData, "Decompressed data should match original data")

	err = gr.Close()
	assert.NoError(t, err, "Close should not return an error")
}

func TestGzipWriter_WriteHeader_NoGzip(t *testing.T) {
	writer := httptest.NewRecorder()
	gw := NewGzipWriter(writer)

	statusCode := http.StatusNotFound
	gw.WriteHeader(statusCode)

	assert.Empty(t, writer.Header().Get("Content-Encoding"), "Content-Encoding header should be empty")

	assert.Equal(t, statusCode, writer.Code, "Status code should match the one set")
}

func TestGzipReader_InvalidData(t *testing.T) {
	invalidCompressedData := []byte("invalid gzip data")

	compressedReader := io.NopCloser(bytes.NewReader(invalidCompressedData))

	gr, err := NewGzipReader(compressedReader)
	assert.Error(t, err, "NewGzipReader should return an error with invalid data")
	assert.Nil(t, gr, "GzipReader should be nil when creation fails")
}
