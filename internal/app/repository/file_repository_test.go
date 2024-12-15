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

	url, exists, isDeleted := repo.Find(shortID)
	assert.True(t, exists)
	assert.False(t, isDeleted)
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
		savedURL, exists, isDeleted := repo.Find(shortID)
		assert.True(t, exists)
		assert.False(t, isDeleted)
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

func TestFileRepository_DeleteUserURLsBatch(t *testing.T) {
	tmpFile := t.TempDir() + "/url_store.json"

	repo := NewFileRepository(tmpFile)
	err := repo.Initialize()
	assert.NoError(t, err)

	urls := []struct {
		shortURL    string
		originalURL string
		userID      string
	}{
		{"abc123", "http://example1.com", "user1"},
		{"def456", "http://example2.com", "user1"},
		{"ghi789", "http://example3.com", "user2"},
	}

	for _, u := range urls {
		_, err := repo.Save(u.shortURL, u.originalURL, u.userID)
		assert.NoError(t, err)
	}

	err = repo.DeleteUserURLsBatch([]string{"abc123", "def456"}, "user1")
	assert.NoError(t, err)

	_, _, isDeleted1 := repo.Find("abc123")
	assert.True(t, isDeleted1, "URL abc123 should be marked as deleted")
	_, _, isDeleted2 := repo.Find("def456")
	assert.True(t, isDeleted2, "URL def456 should be marked as deleted")

	err = repo.DeleteUserURLsBatch([]string{"ghi789"}, "user1")
	assert.NoError(t, err)

	_, _, isDeleted3 := repo.Find("ghi789")
	assert.False(t, isDeleted3, "URL ghi789 should not be marked as deleted")

	err = repo.DeleteUserURLsBatch([]string{"nonexistent"}, "user1")
	assert.NoError(t, err)

	repo2 := NewFileRepository(tmpFile)
	err = repo2.Initialize()
	assert.NoError(t, err)

	_, _, isDeleted1Again := repo2.Find("abc123")
	assert.True(t, isDeleted1Again, "URL abc123 should still be marked as deleted after reload")
	_, _, isDeleted2Again := repo2.Find("def456")
	assert.True(t, isDeleted2Again, "URL def456 should still be marked as deleted after reload")
	_, _, isDeleted3Again := repo2.Find("ghi789")
	assert.False(t, isDeleted3Again, "URL ghi789 should still not be marked as deleted after reload")
}

func TestFileRepository_Find_WithDeletedURLs(t *testing.T) {
	tmpFile := t.TempDir() + "/url_store.json"

	repo := NewFileRepository(tmpFile)
	err := repo.Initialize()
	assert.NoError(t, err)

	_, err = repo.Save("test123", "http://example.com", "user1")
	assert.NoError(t, err)

	originalURL, exists, isDeleted := repo.Find("test123")
	assert.True(t, exists)
	assert.False(t, isDeleted)
	assert.Equal(t, "http://example.com", originalURL)

	err = repo.DeleteUserURLsBatch([]string{"test123"}, "user1")
	assert.NoError(t, err)

	originalURL, exists, isDeleted = repo.Find("test123")
	assert.True(t, exists)
	assert.True(t, isDeleted)
	assert.Equal(t, "http://example.com", originalURL)

	repo2 := NewFileRepository(tmpFile)
	err = repo2.Initialize()
	assert.NoError(t, err)

	originalURL, exists, isDeleted = repo2.Find("test123")
	assert.True(t, exists)
	assert.True(t, isDeleted)
	assert.Equal(t, "http://example.com", originalURL)
}
