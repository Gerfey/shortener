package repository

import (
	"errors"
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/Gerfey/shortener/internal/models"
	"github.com/stretchr/testify/assert"
)

// ErrorFile имитирует ошибку при закрытии файла
type ErrorFile struct {
	*os.File
}

func (e *ErrorFile) Close() error {
	return errors.New("close error")
}

func TestFileCloseError(t *testing.T) {
	tempDir := t.TempDir()
	filePath := filepath.Join(tempDir, "urls.json")

	file, err := os.Create(filePath)
	assert.NoError(t, err)
	err = file.Close()
	assert.NoError(t, err)

	repo := &FileRepository{
		Path: filePath,
		data: make(map[string]models.URLInfo),
	}

	err = repo.Initialize()
	assert.NoError(t, err)

	assert.NotNil(t, repo.data)
}

func TestFileStatError(t *testing.T) {
	tempDir := t.TempDir()
	filePath := filepath.Join(tempDir, "non_existent_dir", "urls.json")

	repo := &FileRepository{
		Path: filePath,
		data: make(map[string]models.URLInfo),
	}

	err := repo.Initialize()
	assert.Error(t, err)
}

// MockReadCloser имитирует ошибку при чтении файла
type MockReadCloser struct {
	io.Reader
	io.Closer
}

func (m *MockReadCloser) Read(p []byte) (n int, err error) {
	return 0, errors.New("read error")
}

func (m *MockReadCloser) Close() error {
	return nil
}

func TestFileDecodeError(t *testing.T) {
	tempDir := t.TempDir()
	filePath := filepath.Join(tempDir, "urls.json")

	file, err := os.Create(filePath)
	assert.NoError(t, err)
	_, err = file.WriteString("{invalid json")
	assert.NoError(t, err)
	err = file.Close()
	assert.NoError(t, err)

	repo := &FileRepository{
		Path: filePath,
		data: make(map[string]models.URLInfo),
	}

	err = repo.Initialize()
	assert.Error(t, err)
}
