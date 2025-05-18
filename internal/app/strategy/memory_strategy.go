package strategy

import (
	"github.com/Gerfey/shortener/internal/app/repository"
	"github.com/Gerfey/shortener/internal/models"
)

// MemoryStrategy стратегия хранилища в памяти
type MemoryStrategy struct{}

// NewMemoryStrategy создает новую стратегию
func NewMemoryStrategy() *MemoryStrategy {
	return &MemoryStrategy{}
}

// Initialize инициализирует хранилище
func (s *MemoryStrategy) Initialize() (models.Repository, error) {
	return repository.NewMemoryRepository(), nil
}

// Close закрывает хранилище
func (s *MemoryStrategy) Close() error {
	return nil
}
