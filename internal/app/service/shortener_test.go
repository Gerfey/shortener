package service

import (
	"github.com/Gerfey/shortener/internal/app/repository"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestShortenSuccess(t *testing.T) {
	path := "test.json"
	url := "https://example.com"

	fileStorageService := NewFileStorage(path)
	memoryRepository := repository.NewURLMemoryRepository()

	s := NewShortenerService(memoryRepository, fileStorageService)

	shortURL, err := s.ShortenID(url)
	findURL, _ := memoryRepository.Find(shortURL)

	assert.NoError(t, err)
	assert.NotEmpty(t, shortURL)
	assert.Equal(t, findURL, url)
}

func TestFindURLSuccess(t *testing.T) {
	path := "test.json"
	url := "https://example.com"
	shortURL := "s65fg"

	fileStorageService := NewFileStorage(path)
	memoryRepository := repository.NewURLMemoryRepository()
	_ = memoryRepository.Save(shortURL, url)

	s := NewShortenerService(memoryRepository, fileStorageService)

	findURL, err := s.FindURL(shortURL)

	assert.NoError(t, err)
	assert.Equal(t, findURL, url)
}

func TestNotFound(t *testing.T) {
	path := "test.json"
	shortURL := "s65fg"

	fileStorageService := NewFileStorage(path)
	memoryRepository := repository.NewURLMemoryRepository()

	s := NewShortenerService(memoryRepository, fileStorageService)

	_, err := s.FindURL(shortURL)

	assert.Error(t, err)
}
