package service

import (
	"fmt"
	"math/rand"

	"github.com/Gerfey/shortener/internal/app/repository"
)

const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
const lenShortID = 8

type ShortenerService struct {
	repository repository.Repository
}

func NewShortenerService(r repository.Repository) *ShortenerService {
	return &ShortenerService{repository: r}
}

func (s *ShortenerService) ShortenID(url string) (string, error) {
	shortID := generateShortID(lenShortID)

	err := s.repository.Save(shortID, url)
	if err != nil {
		return "", err
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
