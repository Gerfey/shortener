package repository

import (
	"sync"
)

type Repository interface {
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

func (s *URLMemoryRepository) Find(key string) (string, bool) {
	s.RLock()
	defer s.RUnlock()
	value, exists := s.data[key]
	return value, exists
}

func (s *URLMemoryRepository) Save(key, value string) error {
	s.Lock()
	defer s.Unlock()
	s.data[key] = value
	return nil
}
