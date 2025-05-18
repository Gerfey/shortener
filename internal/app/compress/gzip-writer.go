package compress

import (
	"compress/gzip"
	"io"
	"net/http"
)

// GzipWriter - обертка над http.ResponseWriter для сжатия ответов с использованием gzip
type GzipWriter struct {
	w  http.ResponseWriter
	zw *gzip.Writer
}

// NewGzipWriter создает новый GzipWriter
func NewGzipWriter(w http.ResponseWriter) *GzipWriter {
	return &GzipWriter{
		w:  w,
		zw: gzip.NewWriter(w),
	}
}

// Header возвращает заголовки HTTP-ответа
func (c *GzipWriter) Header() http.Header {
	return c.w.Header()
}

// Write записывает сжатые данные в ответ
func (c *GzipWriter) Write(p []byte) (int, error) {
	return c.zw.Write(p)
}

// WriteHeader устанавливает код статуса HTTP-ответа и добавляет заголовок Content-Encoding: gzip
func (c *GzipWriter) WriteHeader(statusCode int) {
	if statusCode < 300 {
		c.w.Header().Set("Content-Encoding", "gzip")
	}
	c.w.WriteHeader(statusCode)
}

// Close закрывает gzip writer
func (c *GzipWriter) Close() error {
	return c.zw.Close()
}

// GzipReader обертка для чтения сжатых данных
type GzipReader struct {
	r  io.ReadCloser
	zr *gzip.Reader
}

func NewGzipReader(r io.ReadCloser) (*GzipReader, error) {
	zr, err := gzip.NewReader(r)
	if err != nil {
		return nil, err
	}

	return &GzipReader{
		r:  r,
		zr: zr,
	}, nil
}

// Read читает распакованные данные
func (c GzipReader) Read(p []byte) (n int, err error) {
	return c.zr.Read(p)
}

// Close закрывает gzip читатель и оригинальный ReadCloser
func (c *GzipReader) Close() error {
	if err := c.r.Close(); err != nil {
		return err
	}
	return c.zr.Close()
}
