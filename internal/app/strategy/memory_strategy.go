package strategy

import (
	"context"
	"github.com/Gerfey/shortener/internal/app/repository"
	"github.com/Gerfey/shortener/internal/models"
)

type MemoryStrategy struct{}

func NewMemoryStrategy() *MemoryStrategy {
	return &MemoryStrategy{}
}

func (s *MemoryStrategy) Initialize(ctx context.Context) (models.Repository, error) {
	return repository.NewMemoryRepository(), nil
}

func (s *MemoryStrategy) Close() error {
	return nil
}
