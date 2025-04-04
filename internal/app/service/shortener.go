package service

import (
	"context"
	"fmt"
	"math/rand"

	"github.com/Gerfey/shortener/internal/models"
)

const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
const lenShortID = 8

type ShortenerService struct {
	repository models.Repository
}

func NewShortenerService(r models.Repository) *ShortenerService {
	return &ShortenerService{repository: r}
}

func (s *ShortenerService) SaveBatch(ctx context.Context, urls map[string]string, userID string) error {
	return s.repository.SaveBatch(ctx, urls, userID)
}

func (s *ShortenerService) GetShortURL(ctx context.Context, originalURL string) (string, error) {
	shortURL, err := s.repository.FindShortURL(ctx, originalURL)
	if err != nil {
		return "", fmt.Errorf("failed to find short URL: %w", err)
	}
	return shortURL, nil
}

func (s *ShortenerService) ShortenID(ctx context.Context, url string, userID string) (string, error) {
	existingShortURL, err := s.repository.FindShortURL(ctx, url)
	if err == nil {
		return existingShortURL, models.ErrURLExists
	}

	shortID := generateShortID(lenShortID)
	shortID, err = s.repository.Save(ctx, shortID, url, userID)
	if err != nil {
		return shortID, err
	}

	return shortID, nil
}

func (s *ShortenerService) FindURL(ctx context.Context, code string) (string, error) {
	url, exists, _ := s.repository.Find(ctx, code)
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
