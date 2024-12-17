package repository

import (
	"context"
	"fmt"
	"sync"

	"github.com/Gerfey/shortener/internal/models"
)

type MemoryRepository struct {
	urls map[string]models.URLInfo
	mu   sync.RWMutex
}

func NewMemoryRepository() *MemoryRepository {
	return &MemoryRepository{
		urls: make(map[string]models.URLInfo),
	}
}

func (r *MemoryRepository) All(ctx context.Context) map[string]string {
	r.mu.RLock()
	defer r.mu.RUnlock()

	result := make(map[string]string)
	for shortURL, urlInfo := range r.urls {
		result[shortURL] = urlInfo.OriginalURL
	}
	return result
}

func (r *MemoryRepository) Find(ctx context.Context, key string) (string, bool, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if urlInfo, ok := r.urls[key]; ok {
		return urlInfo.OriginalURL, true, urlInfo.IsDeleted
	}
	return "", false, false
}

func (r *MemoryRepository) FindShortURL(ctx context.Context, originalURL string) (string, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for shortURL, urlInfo := range r.urls {
		if urlInfo.OriginalURL == originalURL {
			return shortURL, nil
		}
	}
	return "", fmt.Errorf("original URL not found")
}

func (r *MemoryRepository) Save(ctx context.Context, key, value string, userID string) (string, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.urls[key] = models.URLInfo{
		ShortURL:    key,
		OriginalURL: value,
		UserID:      userID,
	}
	return key, nil
}

func (r *MemoryRepository) SaveBatch(ctx context.Context, urls map[string]string, userID string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	for shortURL, originalURL := range urls {
		r.urls[shortURL] = models.URLInfo{
			ShortURL:    shortURL,
			OriginalURL: originalURL,
			UserID:      userID,
		}
	}
	return nil
}

func (r *MemoryRepository) GetUserURLs(ctx context.Context, userID string) ([]models.URLPair, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var userURLs []models.URLPair
	for _, urlInfo := range r.urls {
		if urlInfo.UserID == userID {
			userURLs = append(userURLs, models.URLPair{
				ShortURL:    urlInfo.ShortURL,
				OriginalURL: urlInfo.OriginalURL,
			})
		}
	}
	return userURLs, nil
}

func (r *MemoryRepository) DeleteUserURLsBatch(ctx context.Context, shortURLs []string, userID string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	for _, shortURL := range shortURLs {
		if urlInfo, exists := r.urls[shortURL]; exists && urlInfo.UserID == userID {
			urlInfo.IsDeleted = true
			r.urls[shortURL] = urlInfo
		}
	}
	return nil
}

func (r *MemoryRepository) Ping(ctx context.Context) error {
	return nil
}
