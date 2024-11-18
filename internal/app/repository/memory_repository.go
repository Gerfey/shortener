package repository

import (
	"sync"
)

type MemoryRepository struct {
	data map[string]string
	sync.RWMutex
}

func NewMemoryRepository() *MemoryRepository {
	return &MemoryRepository{
		data: make(map[string]string),
	}
}

func (r *MemoryRepository) All() map[string]string {
	return r.data
}

func (r *MemoryRepository) Find(key string) (string, bool) {
	r.RLock()
	defer r.RUnlock()
	value, exists := r.data[key]
	return value, exists
}

func (r *MemoryRepository) Save(key, value string) error {
	r.Lock()
	defer r.Unlock()
	r.data[key] = value
	return nil
}

func (r *MemoryRepository) SaveBatch(urls map[string]string) error {
	r.Lock()
	defer r.Unlock()

	for shortURL, originalURL := range urls {
		r.data[shortURL] = originalURL
	}

	return nil
}
