package service

import (
	"github.com/Gerfey/shortener/internal/app/repository"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestShortenSuccess(t *testing.T) {
	url := "https://example.com"

	repository := repository.NewURLMemoryRepository()

	s := NewShortenerService(repository)

	shortURL, err := s.ShortenID(url)
	findURL, _ := repository.Find(shortURL)

	assert.NoError(t, err)
	assert.NotEmpty(t, shortURL)
	assert.Equal(t, findURL, url)
}

func TestFindURLSuccess(t *testing.T) {
	url := "https://example.com"
	shortURL := "s65fg"

	repository := repository.NewURLMemoryRepository()
	_ = repository.Save(shortURL, url)

	s := NewShortenerService(repository)

	findURL, err := s.FindURL(shortURL)

	assert.NoError(t, err)
	assert.Equal(t, findURL, url)
}

func TestNotFound(t *testing.T) {
	shortURL := "s65fg"

	repository := repository.NewURLMemoryRepository()

	s := NewShortenerService(repository)

	_, err := s.FindURL(shortURL)

	assert.Error(t, err)
}
