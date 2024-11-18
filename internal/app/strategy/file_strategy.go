package strategy

import (
	"github.com/Gerfey/shortener/internal/app/repository"
	"github.com/Gerfey/shortener/internal/models"
)

type FileStrategy struct {
	FilePath string
}

func NewFileStrategy(filePath string) *FileStrategy {
	return &FileStrategy{FilePath: filePath}
}

func (s *FileStrategy) Initialize() (models.Repository, error) {
	return repository.NewFileRepository(s.FilePath), nil
}

func (s *FileStrategy) Close() error {
	return nil
}
