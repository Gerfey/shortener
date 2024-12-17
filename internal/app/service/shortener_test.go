package service

import (
	"context"
	"errors"
	"github.com/Gerfey/shortener/internal/mock"
	"github.com/Gerfey/shortener/internal/models"
	"go.uber.org/mock/gomock"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestShortenSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock.NewMockRepository(ctrl)
	shortener := NewShortenerService(mockRepo)

	shortID := "s65fg"
	originalURL := "https://example.com"
	userID := "user123"

	mockRepo.EXPECT().FindShortURL(originalURL).Return("", models.ErrURLNotFound)
	mockRepo.EXPECT().Save(gomock.Any(), originalURL, userID).Return(shortID, nil)

	id, err := shortener.ShortenID(context.Background(), originalURL, userID)
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

	mockRepo.EXPECT().Find(shortID).Return(originalURL, true, false)

	url, err := shortener.FindURL(context.Background(), shortID)
	assert.NoError(t, err)
	assert.Equal(t, originalURL, url)

	mockRepo.EXPECT().Find("notfound").Return("", false, false)

	_, err = shortener.FindURL(context.Background(), "notfound")
	assert.Error(t, err)
}

func TestShortenerService_ShortenID_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock.NewMockRepository(ctrl)
	shortener := NewShortenerService(mockRepo)

	originalURL := "https://example.com"
	userID := "user123"
	expectedErr := errors.New("repository error")

	mockRepo.EXPECT().FindShortURL(originalURL).Return("", models.ErrURLNotFound)
	mockRepo.EXPECT().Save(gomock.Any(), originalURL, userID).Return("", expectedErr)

	_, err := shortener.ShortenID(context.Background(), originalURL, userID)
	assert.Error(t, err)
	assert.Equal(t, expectedErr, err)
}

func TestShortenerService_ShortenID_URLExists(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock.NewMockRepository(ctrl)
	shortener := NewShortenerService(mockRepo)

	originalURL := "https://example.com"
	userID := "user123"
	existingShortURL := "abc123"

	mockRepo.EXPECT().FindShortURL(originalURL).Return(existingShortURL, nil)

	shortURL, err := shortener.ShortenID(context.Background(), originalURL, userID)
	assert.Equal(t, models.ErrURLExists, err)
	assert.Equal(t, existingShortURL, shortURL)
}
