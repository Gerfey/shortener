package strategy

import (
	"github.com/Gerfey/shortener/internal/app/repository"
	"github.com/Gerfey/shortener/internal/models"
)

// FileStrategy стратегия файлового хранилища
type FileStrategy struct {
	filePath string
	fileRepo *repository.FileRepository
}

// NewFileStrategy создает новую стратегию
func NewFileStrategy(filePath string) *FileStrategy {
	return &FileStrategy{
		filePath: filePath,
	}
}

// Initialize инициализирует хранилище
func (s *FileStrategy) Initialize() (models.Repository, error) {
	fileRepository := repository.NewFileRepository(s.filePath)
	if err := fileRepository.Initialize(); err != nil {
		return nil, err
	}
	s.fileRepo = fileRepository
	return fileRepository, nil
}

// Close закрывает хранилище
func (s *FileStrategy) Close() error {
	if s.fileRepo != nil {
		return s.fileRepo.Close()
	}
	return nil
}
