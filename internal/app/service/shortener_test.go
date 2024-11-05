package service

import (
	"github.com/Gerfey/shortener/internal/app/repository"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestShortenSuccess(t *testing.T) {
	url := "https://example.com"

	memoryRepository := repository.NewURLMemoryRepository()

	s := NewShortenerService(memoryRepository)

	shortURL, err := s.ShortenID(url)
	findURL, _ := memoryRepository.Find(shortURL)

	assert.NoError(t, err)
	assert.NotEmpty(t, shortURL)
	assert.Equal(t, findURL, url)
}

func TestFindURLSuccess(t *testing.T) {
	url := "https://example.com"
	shortURL := "s65fg"

	memoryRepository := repository.NewURLMemoryRepository()
	_ = memoryRepository.Save(shortURL, url)

	s := NewShortenerService(memoryRepository)

	findURL, err := s.FindURL(shortURL)

	assert.NoError(t, err)
	assert.Equal(t, findURL, url)
}

func TestNotFound(t *testing.T) {
	shortURL := "s65fg"

	memoryRepository := repository.NewURLMemoryRepository()

	s := NewShortenerService(memoryRepository)

	_, err := s.FindURL(shortURL)

	assert.Error(t, err)
}
