package strategy

import (
	"context"
	"github.com/Gerfey/shortener/internal/app/repository"
	"github.com/Gerfey/shortener/internal/models"
)

type FileStrategy struct {
	filePath string
	fileRepo *repository.FileRepository
}

func NewFileStrategy(filePath string) *FileStrategy {
	return &FileStrategy{
		filePath: filePath,
	}
}

func (s *FileStrategy) Initialize(ctx context.Context) (models.Repository, error) {
	fileRepository := repository.NewFileRepository(s.filePath)
	if err := fileRepository.Initialize(); err != nil {
		return nil, err
	}
	s.fileRepo = fileRepository
	return fileRepository, nil
}

func (s *FileStrategy) Close() error {
	if s.fileRepo != nil {
		return s.fileRepo.Close()
	}
	return nil
}
