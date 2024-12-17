package service

import (
	"context"
	"errors"
	"testing"

	"github.com/Gerfey/shortener/internal/mock"
	"github.com/Gerfey/shortener/internal/models"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestShortenSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock.NewMockRepository(ctrl)
	shortener := NewShortenerService(mockRepo)

	shortID := "s65fg"
	originalURL := "https://example.com"
	userID := "user123"
	ctx := context.Background()

	mockRepo.EXPECT().FindShortURL(ctx, originalURL).Return("", models.ErrURLNotFound)
	mockRepo.EXPECT().Save(ctx, gomock.Any(), originalURL, userID).Return(shortID, nil)

	id, err := shortener.ShortenID(ctx, originalURL, userID)
	assert.NoError(t, err)
	assert.Equal(t, len(id), 5)
}

func TestFindURLSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock.NewMockRepository(ctrl)
	shortener := NewShortenerService(mockRepo)

	shortID := "s65fg"
	originalURL := "https://example.com"
	ctx := context.Background()

	mockRepo.EXPECT().Find(ctx, shortID).Return(originalURL, true, false)

	url, err := shortener.FindURL(ctx, shortID)
	assert.NoError(t, err)
	assert.Equal(t, originalURL, url)

	mockRepo.EXPECT().Find(ctx, "notfound").Return("", false, false)

	_, err = shortener.FindURL(ctx, "notfound")
	assert.Error(t, err)
}

func TestShortenerService_ShortenID_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock.NewMockRepository(ctrl)
	shortener := NewShortenerService(mockRepo)

	originalURL := "https://example.com"
	userID := "user123"
	ctx := context.Background()

	expectedErr := errors.New("database error")
	mockRepo.EXPECT().FindShortURL(ctx, originalURL).Return("", models.ErrURLNotFound)
	mockRepo.EXPECT().Save(ctx, gomock.Any(), originalURL, userID).Return("", expectedErr)

	_, err := shortener.ShortenID(ctx, originalURL, userID)
	assert.Error(t, err)
	assert.Equal(t, expectedErr, err)
}

func TestShortenerService_ShortenID_URLExists(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock.NewMockRepository(ctrl)
	shortener := NewShortenerService(mockRepo)

	originalURL := "https://example.com"
	existingShortURL := "abc123"
	userID := "user123"
	ctx := context.Background()

	mockRepo.EXPECT().FindShortURL(ctx, originalURL).Return(existingShortURL, nil)

	shortURL, err := shortener.ShortenID(ctx, originalURL, userID)
	assert.NoError(t, err)
	assert.Equal(t, existingShortURL, shortURL)
}
