package usecase

import (
	"context"
	"errors"

	"github.com/Gerfey/shortener/internal/app/settings"
	"github.com/Gerfey/shortener/internal/models"
)

type UserURLsUseCase struct {
	repository models.Repository
	settings   *settings.Settings
}

func NewUserURLsUseCase(repository models.Repository, settings *settings.Settings) *UserURLsUseCase {
	return &UserURLsUseCase{
		repository: repository,
		settings:   settings,
	}
}

func (uc *UserURLsUseCase) GetUserURLs(ctx context.Context, userID string) ([]models.URLPair, error) {
	if userID == "" {
		return nil, errors.New("ID пользователя не может быть пустым")
	}

	urls, err := uc.repository.GetUserURLs(ctx, userID)
	if err != nil {
		return nil, err
	}

	baseURL := uc.settings.BaseURL()
	for i := range urls {
		urls[i].ShortURL = baseURL + "/" + urls[i].ShortURL
	}

	return urls, nil
}

func (uc *UserURLsUseCase) DeleteUserURLs(ctx context.Context, userID string, shortURLs []string) error {
	if userID == "" {
		return errors.New("ID пользователя не может быть пустым")
	}

	if len(shortURLs) == 0 {
		return errors.New("список URL не может быть пустым")
	}

	return uc.repository.DeleteUserURLsBatch(ctx, shortURLs, userID)
}
