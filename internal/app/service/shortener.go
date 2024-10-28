package service

import (
	"fmt"
	"github.com/google/uuid"
	"math/rand"

	"github.com/Gerfey/shortener/internal/app/repository"
)

const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
const lenShortID = 8

type ShortenerService struct {
	repository  repository.Repository
	fileStorage *FileStorage
}

func NewShortenerService(r repository.Repository, fileStorage *FileStorage) *ShortenerService {
	urlInfos, _ := fileStorage.Load()

	for _, urlInfo := range urlInfos {
		_ = r.Save(urlInfo.ShortURL, urlInfo.OriginalURL)
	}

	return &ShortenerService{repository: r, fileStorage: fileStorage}
}

func (s *ShortenerService) ShortenID(url string) (string, error) {
	shortID := generateShortID(lenShortID)

	err := s.repository.Save(shortID, url)
	if err != nil {
		return "", err
	}

	urlInfo := URLInfo{
		UUID:        uuid.New().String(),
		ShortURL:    shortID,
		OriginalURL: url,
	}

	err = s.fileStorage.Save(urlInfo)

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
