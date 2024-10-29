package repository

import (
	"sync"
)

type Repository interface {
	All() map[string]string
	Find(key string) (string, bool)
	Save(key, value string) error
}

type URLMemoryRepository struct {
	data map[string]string
	sync.RWMutex
}

func NewURLMemoryRepository() *URLMemoryRepository {
	return &URLMemoryRepository{
		data: make(map[string]string),
	}
}

func (r *URLMemoryRepository) All() map[string]string {
	return r.data
}

func (r *URLMemoryRepository) Find(key string) (string, bool) {
	r.RLock()
	defer r.RUnlock()
	value, exists := r.data[key]
	return value, exists
}

func (r *URLMemoryRepository) Save(key, value string) error {
	r.Lock()
	defer r.Unlock()
	r.data[key] = value
	return nil
}
