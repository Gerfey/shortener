package usecase

import (
	"context"
	"errors"

	"github.com/Gerfey/shortener/internal/models"
)

type RedirectUseCase struct {
	repository models.Repository
}

type RedirectResult struct {
	OriginalURL string
	IsDeleted   bool
}

func NewRedirectUseCase(repository models.Repository) *RedirectUseCase {
	return &RedirectUseCase{
		repository: repository,
	}
}

func (uc *RedirectUseCase) GetOriginalURL(ctx context.Context, shortID string) (RedirectResult, error) {
	if shortID == "" {
		return RedirectResult{}, errors.New("ID не может быть пустым")
	}

	originalURL, exists, isDeleted := uc.repository.Find(ctx, shortID)
	if !exists {
		return RedirectResult{}, errors.New("URL не найден")
	}

	return RedirectResult{
		OriginalURL: originalURL,
		IsDeleted:   isDeleted,
	}, nil
}
