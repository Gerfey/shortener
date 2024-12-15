package models

type Repository interface {
	All() map[string]string
	Find(key string) (string, bool)
	FindShortURL(originalURL string) (string, error)
	Save(key, value string, userID string) (string, error)
	SaveBatch(urls map[string]string, userID string) error
	GetUserURLs(userID string) ([]URLPair, error)
	Ping() error
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
}

type StorageStrategy interface {
	Initialize() (Repository, error)
	Close() error
}
