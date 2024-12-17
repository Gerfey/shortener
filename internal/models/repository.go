package models

import "context"

type Repository interface {
	All(ctx context.Context) map[string]string
	Find(ctx context.Context, key string) (string, bool, bool) // returns originalURL, exists, isDeleted
	FindShortURL(ctx context.Context, originalURL string) (string, error)
	Save(ctx context.Context, key, value string, userID string) (string, error)
	SaveBatch(ctx context.Context, urls map[string]string, userID string) error
	GetUserURLs(ctx context.Context, userID string) ([]URLPair, error)
	DeleteUserURLsBatch(ctx context.Context, shortURLs []string, userID string) error
	Ping(ctx context.Context) error
}

type URLPair struct {
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}

type URLInfo struct {
	UUID        string `json:"uuid"`
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
	UserID      string `json:"user_id"`
	IsDeleted   bool   `json:"is_deleted"`
}

type StorageStrategy interface {
	Initialize() (Repository, error)
	Close() error
}
