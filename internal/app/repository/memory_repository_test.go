package repository

import (
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
		},
		"def456": {
			OriginalURL: "https://google.com",
			UserID:      "user1",
		},
	}
	repo.urls = urls

	url, exists := repo.Find("abc123")
	assert.True(t, exists)
	assert.Equal(t, "https://example.com", url)

	url, exists = repo.Find("nonexistent")
	assert.False(t, exists)
	assert.Empty(t, url)
}

func TestMemoryRepository_Save(t *testing.T) {
	repo := NewMemoryRepository()

	shortID := "abc123"
	originalURL := "https://example.com"
	userID := "user1"

	savedID, err := repo.Save(shortID, originalURL, userID)
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

	allURLs := repo.All()
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

	userURLs, err := repo.GetUserURLs("user1")
	assert.NoError(t, err)
	assert.Equal(t, 2, len(userURLs))

	userURLs, err = repo.GetUserURLs("nonexistent")
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

	shortURL, err := repo.FindShortURL("https://example.com")
	assert.NoError(t, err)
	assert.Equal(t, "abc123", shortURL)

	shortURL, err = repo.FindShortURL("https://nonexistent.com")
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

	err := repo.SaveBatch(urls, userID)
	assert.NoError(t, err)

	for shortID, originalURL := range urls {
		info, ok := repo.urls[shortID]
		assert.True(t, ok)
		assert.Equal(t, originalURL, info.OriginalURL)
		assert.Equal(t, userID, info.UserID)
	}

	err = repo.SaveBatch(map[string]string{}, userID)
	assert.NoError(t, err)
}

func TestMemoryRepository_Ping(t *testing.T) {
	repo := NewMemoryRepository()

	err := repo.Ping()
	assert.NoError(t, err)
}
