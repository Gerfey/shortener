package usecase

import (
	"context"
	"errors"

	"github.com/Gerfey/shortener/internal/app/service"
	"github.com/Gerfey/shortener/internal/app/settings"
	"github.com/Gerfey/shortener/internal/models"
)

type ShortenUseCase struct {
	shortener *service.ShortenerService
	settings  *settings.Settings
}

func NewShortenUseCase(shortener *service.ShortenerService, settings *settings.Settings) *ShortenUseCase {
	return &ShortenUseCase{
		shortener: shortener,
		settings:  settings,
	}
}

type ShortenResult struct {
	ShortURL      string
	FullShortURL  string
	AlreadyExists bool
}

func (uc *ShortenUseCase) ShortenURL(ctx context.Context, originalURL string, userID string) (ShortenResult, error) {
	if originalURL == "" {
		return ShortenResult{}, errors.New("URL не может быть пустым")
	}

	shortID, err := uc.shortener.ShortenID(ctx, originalURL, userID)
	if err != nil && !errors.Is(err, models.ErrURLExists) {
		return ShortenResult{}, err
	}

	baseURL := uc.settings.BaseURL()
	fullShortURL := baseURL + "/" + shortID

	return ShortenResult{
		ShortURL:      shortID,
		FullShortURL:  fullShortURL,
		AlreadyExists: errors.Is(err, models.ErrURLExists),
	}, nil
}

type BatchItem struct {
	CorrelationID string
	OriginalURL   string
	ShortURL      string
	FullShortURL  string
}

func (uc *ShortenUseCase) ShortenBatch(ctx context.Context, items []models.BatchRequestItem, userID string) ([]BatchItem, error) {
	if len(items) == 0 {
		return nil, errors.New("пакет не может быть пустым")
	}

	urlMap := make(map[string]string)
	itemMap := make(map[string]models.BatchRequestItem)

	for _, item := range items {
		urlMap[item.CorrelationID] = item.OriginalURL
		itemMap[item.CorrelationID] = item
	}

	if err := uc.shortener.SaveBatch(ctx, urlMap, userID); err != nil {
		return nil, err
	}

	var result []BatchItem
	baseURL := uc.settings.BaseURL()

	for correlationID, originalURL := range urlMap {
		shortURL, err := uc.shortener.GetShortURL(ctx, originalURL)
		if err != nil {
			return nil, err
		}

		result = append(result, BatchItem{
			CorrelationID: correlationID,
			OriginalURL:   originalURL,
			ShortURL:      shortURL,
			FullShortURL:  baseURL + "/" + shortURL,
		})
	}

	return result, nil
}
