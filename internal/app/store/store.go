package store

import "sync"

type Stored interface {
	Get(key string) (string, bool)
	Set(key, value string)
	Delete(key string)
}

type Store struct {
	sync.RWMutex
	data map[string]string
}

func NewStore() *Store {
	return &Store{
		data: make(map[string]string),
	}
}

func (s *Store) Get(key string) (string, bool) {
	s.RLock()
	defer s.RUnlock()
	value, exists := s.data[key]
	return value, exists
}

func (s *Store) Set(key, value string) {
	s.Lock()
	defer s.Unlock()
	s.data[key] = value
}

func (s *Store) Delete(key string) {
	s.Lock()
	defer s.Unlock()
	delete(s.data, key)
}
