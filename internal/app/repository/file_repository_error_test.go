package repository

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/Gerfey/shortener/internal/models"
	"github.com/stretchr/testify/assert"
)

func TestFileRepository_InitializeErrors(t *testing.T) {
	tempDir := t.TempDir()

	t.Run("Failed to create directory", func(t *testing.T) {
		invalidDirPath := filepath.Join(tempDir, "invalid_dir")
		file, err := os.Create(invalidDirPath)
		assert.NoError(t, err)
		closeErr := file.Close()
		assert.NoError(t, closeErr)

		repo := &FileRepository{
			Path: filepath.Join(invalidDirPath, "urls.json"),
		}

		err = repo.Initialize()
		assert.Error(t, err)
	})

	t.Run("Failed to open file for reading", func(t *testing.T) {
		noReadDir := filepath.Join(tempDir, "no_read_dir")
		err := os.Mkdir(noReadDir, 0000)
		assert.NoError(t, err)
		defer func() {
			chmodErr := os.Chmod(noReadDir, 0755)
			assert.NoError(t, chmodErr)
		}()

		repo := &FileRepository{
			Path: filepath.Join(noReadDir, "urls.json"),
		}

		err = repo.Initialize()
		assert.Error(t, err)
	})
}

func TestFileRepository_GetUserURLsErrors(t *testing.T) {
	tempDir := t.TempDir()
	filePath := filepath.Join(tempDir, "urls.json")

	repo := &FileRepository{
		Path: filePath,
		data: make(map[string]models.URLInfo),
	}

	err := repo.Initialize()
	assert.NoError(t, err)

	t.Run("No URLs for user", func(t *testing.T) {
		urls, err := repo.GetUserURLs(context.Background(), "non_existent_user")
		assert.NoError(t, err)
		assert.Empty(t, urls)
	})
}
