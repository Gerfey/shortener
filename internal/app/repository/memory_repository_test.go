package repository

import (
	"context"
	"github.com/Gerfey/shortener/internal/models"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMemoryRepository_Find(t *testing.T) {
	repo := NewMemoryRepository()

	urls := map[string]models.URLInfo{
		"abc123": {
			OriginalURL: "https://example.com",
			UserID:      "user1",
			IsDeleted:   false,
		},
		"def456": {
			OriginalURL: "https://google.com",
			UserID:      "user1",
			IsDeleted:   false,
		},
	}
	repo.urls = urls

	url, exists, isDeleted := repo.Find(context.Background(), "abc123")
	assert.True(t, exists)
	assert.False(t, isDeleted)
	assert.Equal(t, "https://example.com", url)

	url, exists, isDeleted = repo.Find(context.Background(), "nonexistent")
	assert.False(t, exists)
	assert.False(t, isDeleted)
	assert.Empty(t, url)
}

func TestMemoryRepository_Save(t *testing.T) {
	repo := NewMemoryRepository()

	shortID := "abc123"
	originalURL := "https://example.com"
	userID := "user1"

	savedID, err := repo.Save(context.Background(), shortID, originalURL, userID)
	assert.NoError(t, err)
	assert.Equal(t, shortID, savedID)

	info, ok := repo.urls[shortID]
	assert.True(t, ok)
	assert.Equal(t, originalURL, info.OriginalURL)
	assert.Equal(t, userID, info.UserID)
}

func TestMemoryRepository_All(t *testing.T) {
	repo := NewMemoryRepository()

	urls := map[string]models.URLInfo{
		"abc123": {
			OriginalURL: "https://example.com",
			UserID:      "user1",
		},
		"def456": {
			OriginalURL: "https://google.com",
			UserID:      "user1",
		},
	}
	repo.urls = urls

	allURLs := repo.All(context.Background())
	assert.Equal(t, 2, len(allURLs))
	assert.Equal(t, "https://example.com", allURLs["abc123"])
	assert.Equal(t, "https://google.com", allURLs["def456"])
}

func TestMemoryRepository_GetUserURLs(t *testing.T) {
	repo := NewMemoryRepository()

	urls := map[string]models.URLInfo{
		"abc123": {
			OriginalURL: "https://example.com",
			UserID:      "user1",
		},
		"def456": {
			OriginalURL: "https://google.com",
			UserID:      "user1",
		},
		"ghi789": {
			OriginalURL: "https://github.com",
			UserID:      "user2",
		},
	}
	repo.urls = urls

	userURLs, err := repo.GetUserURLs(context.Background(), "user1")
	assert.NoError(t, err)
	assert.Equal(t, 2, len(userURLs))

	userURLs, err = repo.GetUserURLs(context.Background(), "nonexistent")
	assert.NoError(t, err)
	assert.Empty(t, userURLs)
}

func TestMemoryRepository_FindShortURL(t *testing.T) {
	repo := NewMemoryRepository()

	urls := map[string]models.URLInfo{
		"abc123": {
			OriginalURL: "https://example.com",
			UserID:      "user1",
		},
		"def456": {
			OriginalURL: "https://google.com",
			UserID:      "user1",
		},
	}
	repo.urls = urls

	shortURL, err := repo.FindShortURL(context.Background(), "https://example.com")
	assert.NoError(t, err)
	assert.Equal(t, "abc123", shortURL)

	shortURL, err = repo.FindShortURL(context.Background(), "https://nonexistent.com")
	assert.Error(t, err)
	assert.Empty(t, shortURL)
}

func TestMemoryRepository_SaveBatch(t *testing.T) {
	repo := NewMemoryRepository()

	urls := map[string]string{
		"abc123": "https://example.com",
		"def456": "https://google.com",
	}
	userID := "user1"

	err := repo.SaveBatch(context.Background(), urls, userID)
	assert.NoError(t, err)

	for shortID, originalURL := range urls {
		info, ok := repo.urls[shortID]
		assert.True(t, ok)
		assert.Equal(t, originalURL, info.OriginalURL)
		assert.Equal(t, userID, info.UserID)
	}

	err = repo.SaveBatch(context.Background(), map[string]string{}, userID)
	assert.NoError(t, err)
}

func TestMemoryRepository_Ping(t *testing.T) {
	repo := NewMemoryRepository()

	err := repo.Ping(context.Background())
	assert.NoError(t, err)
}

func TestMemoryRepository_DeleteUserURLsBatch(t *testing.T) {
	repo := NewMemoryRepository()

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
		_, err := repo.Save(context.Background(), u.shortURL, u.originalURL, u.userID)
		assert.NoError(t, err)
	}

	err := repo.DeleteUserURLsBatch(context.Background(), []string{"abc123", "def456"}, "user1")
	assert.NoError(t, err)

	_, _, isDeleted1 := repo.Find(context.Background(), "abc123")
	assert.True(t, isDeleted1, "URL abc123 should be marked as deleted")
	_, _, isDeleted2 := repo.Find(context.Background(), "def456")
	assert.True(t, isDeleted2, "URL def456 should be marked as deleted")

	err = repo.DeleteUserURLsBatch(context.Background(), []string{"ghi789"}, "user1")
	assert.NoError(t, err)

	_, _, isDeleted3 := repo.Find(context.Background(), "ghi789")
	assert.False(t, isDeleted3, "URL ghi789 should not be marked as deleted")

	err = repo.DeleteUserURLsBatch(context.Background(), []string{"nonexistent"}, "user1")
	assert.NoError(t, err)
}

func TestMemoryRepository_Find_WithDeletedURLs(t *testing.T) {
	repo := NewMemoryRepository()

	_, err := repo.Save(context.Background(), "test123", "http://example.com", "user1")
	assert.NoError(t, err)

	originalURL, exists, isDeleted := repo.Find(context.Background(), "test123")
	assert.True(t, exists)
	assert.False(t, isDeleted)
	assert.Equal(t, "http://example.com", originalURL)

	err = repo.DeleteUserURLsBatch(context.Background(), []string{"test123"}, "user1")
	assert.NoError(t, err)

	originalURL, exists, isDeleted = repo.Find(context.Background(), "test123")
	assert.True(t, exists)
	assert.True(t, isDeleted)
	assert.Equal(t, "http://example.com", originalURL)
}
