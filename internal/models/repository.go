package models

import "context"

// Repository определяет интерфейс для работы с хранилищем URL
type Repository interface {
	// All возвращает все сохраненные URL в виде карты короткий->оригинальный
	All(ctx context.Context) map[string]string
	// Find ищет URL по короткому идентификатору и возвращает оригинальный URL, флаг существования и флаг удаления
	Find(ctx context.Context, key string) (string, bool, bool)
	// FindShortURL ищет короткий URL по оригинальному URL
	FindShortURL(ctx context.Context, originalURL string) (string, error)
	// Save сохраняет пару короткий->оригинальный URL с привязкой к пользователю
	Save(ctx context.Context, key, value string, userID string) (string, error)
	// SaveBatch сохраняет несколько пар короткий->оригинальный URL с привязкой к пользователю
	SaveBatch(ctx context.Context, urls map[string]string, userID string) error
	// GetUserURLs возвращает все URL, принадлежащие пользователю
	GetUserURLs(ctx context.Context, userID string) ([]URLPair, error)
	// DeleteUserURLsBatch помечает указанные URL пользователя как удаленные
	DeleteUserURLsBatch(ctx context.Context, shortURLs []string, userID string) error
	// Ping проверяет доступность хранилища
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

// StorageStrategy определяет интерфейс для стратегии хранения данных
type StorageStrategy interface {
	// Initialize инициализирует хранилище и возвращает репозиторий
	Initialize() (Repository, error)
	// Close закрывает соединение с хранилищем
	Close() error
}
