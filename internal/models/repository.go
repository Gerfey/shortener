package models

type Repository interface {
	All() map[string]string
	Find(key string) (string, bool)
	FindShortURL(originalURL string) (string, error)
	Save(key, value string) (string, error)
	SaveBatch(urls map[string]string) error
}

type StorageStrategy interface {
	Initialize() (Repository, error)
	Close() error
}
