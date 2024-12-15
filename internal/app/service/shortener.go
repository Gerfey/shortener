package service

import (
	"fmt"
	"github.com/Gerfey/shortener/internal/models"
	"math/rand"
)

const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
const lenShortID = 8

type ShortenerService struct {
	repository models.Repository
}

func NewShortenerService(r models.Repository) *ShortenerService {
	return &ShortenerService{repository: r}
}

func (s *ShortenerService) SaveBatch(urls map[string]string) error {
	return s.repository.SaveBatch(urls)
}

func (s *ShortenerService) GetShortURL(originalURL string) (string, error) {
	shortURL, err := s.repository.FindShortURL(originalURL)
	if err != nil {
		return "", fmt.Errorf("failed to find short URL: %w", err)
	}
	return shortURL, nil
}

func (s *ShortenerService) ShortenID(url string) (string, error) {
	shortID := generateShortID(lenShortID)

	shortID, err := s.repository.Save(shortID, url)
	if err != nil {
		return shortID, err
	}

	return shortID, err
}

func (s *ShortenerService) FindURL(code string) (string, error) {
	url, exists := s.repository.Find(code)
	if !exists {
		return "", fmt.Errorf("ничего не найдено по значению %v", code)
	}
	return url, nil
}

func generateShortID(length int) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}
