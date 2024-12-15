package repository

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFileRepository_SaveAndFind(t *testing.T) {
	tmpDir := t.TempDir()
	tmpFile := filepath.Join(tmpDir, "test_urls.json")

	repo := NewFileRepository(tmpFile)

	shortID := "abc123"
	originalURL := "https://example.com"
	userID := "user1"

	savedID, err := repo.Save(shortID, originalURL, userID)
	assert.NoError(t, err)
	assert.Equal(t, shortID, savedID)

	url, exists := repo.Find(shortID)
	assert.True(t, exists)
	assert.Equal(t, originalURL, url)
}

func TestFileRepository_GetUserURLs(t *testing.T) {
	tmpDir := t.TempDir()
	tmpFile := filepath.Join(tmpDir, "test_urls.json")

	repo := NewFileRepository(tmpFile)

	userID := "user1"
	urls := map[string]string{
		"abc123": "https://example.com",
		"def456": "https://google.com",
	}

	for shortID, originalURL := range urls {
		_, err := repo.Save(shortID, originalURL, userID)
		assert.NoError(t, err)
	}

	userURLs, err := repo.GetUserURLs(userID)
	assert.NoError(t, err)
	assert.Equal(t, len(urls), len(userURLs))
}

func TestFileRepository_Initialize(t *testing.T) {
	tmpDir := t.TempDir()
	tmpFile := filepath.Join(tmpDir, "test_urls.json")

	repo := NewFileRepository(tmpFile)

	err := repo.Initialize()
	assert.NoError(t, err)

	_, err = os.Stat(tmpFile)
	assert.NoError(t, err)
}

func TestFileRepository_Close(t *testing.T) {
	tmpDir := t.TempDir()
	tmpFile := filepath.Join(tmpDir, "test_urls.json")

	repo := NewFileRepository(tmpFile)

	shortID := "abc123"
	originalURL := "https://example.com"
	userID := "user1"

	_, err := repo.Save(shortID, originalURL, userID)
	assert.NoError(t, err)

	err = repo.Close()
	assert.NoError(t, err)

	data, err := os.ReadFile(tmpFile)
	assert.NoError(t, err)
	assert.NotEmpty(t, data)
}

func TestFileRepository_All(t *testing.T) {
	tmpDir := t.TempDir()
	tmpFile := filepath.Join(tmpDir, "test_urls.json")

	repo := NewFileRepository(tmpFile)

	urls := map[string]string{
		"abc123": "https://example.com",
		"def456": "https://google.com",
	}

	for shortID, originalURL := range urls {
		_, err := repo.Save(shortID, originalURL, "user1")
		assert.NoError(t, err)
	}

	allURLs := repo.All()
	assert.Equal(t, len(urls), len(allURLs))
	for shortID, originalURL := range urls {
		assert.Equal(t, originalURL, allURLs[shortID])
	}
}

func TestFileRepository_FindShortURL(t *testing.T) {
	tmpFile := t.TempDir() + "/url_store.json"

	repo := NewFileRepository(tmpFile)
	err := repo.Initialize()
	assert.NoError(t, err)

	_, err = repo.Save("abc123", "https://example.com", "user1")
	assert.NoError(t, err)

	shortURL, err := repo.FindShortURL("https://example.com")
	assert.NoError(t, err)
	assert.Equal(t, "abc123", shortURL)

	shortURL, err = repo.FindShortURL("https://nonexistent.com")
	assert.Error(t, err)
	assert.Empty(t, shortURL)

	err = repo.Close()
	assert.NoError(t, err)
}

func TestFileRepository_SaveBatch(t *testing.T) {
	tmpFile := t.TempDir() + "/url_store.json"

	repo := NewFileRepository(tmpFile)
	err := repo.Initialize()
	assert.NoError(t, err)

	urls := map[string]string{
		"abc123": "https://example.com",
		"def456": "https://google.com",
	}
	userID := "user1"

	err = repo.SaveBatch(urls, userID)
	assert.NoError(t, err)

	for shortID, originalURL := range urls {
		savedURL, exists := repo.Find(shortID)
		assert.True(t, exists)
		assert.Equal(t, originalURL, savedURL)
	}

	err = repo.SaveBatch(map[string]string{}, userID)
	assert.NoError(t, err)

	err = repo.Close()
	assert.NoError(t, err)
}

func TestFileRepository_Ping(t *testing.T) {
	tmpFile := t.TempDir() + "/url_store.json"

	repo := NewFileRepository(tmpFile)
	err := repo.Initialize()
	assert.NoError(t, err)

	err = repo.Ping()
	assert.NoError(t, err)

	err = repo.Close()
	assert.NoError(t, err)
}
