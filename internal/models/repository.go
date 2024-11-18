package models

type Repository interface {
	All() map[string]string
	Find(key string) (string, bool)
	Save(key, value string) error
}

type StorageStrategy interface {
	Initialize() (Repository, error)
	Close() error
}
