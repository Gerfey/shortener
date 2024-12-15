package service

import (
	"errors"
	"github.com/Gerfey/shortener/internal/mock"
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

	mockRepo.EXPECT().Save(gomock.Any(), originalURL, userID).Return(shortID, nil).Times(1)

	id, err := shortener.ShortenID(originalURL, userID)
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

	mockRepo.EXPECT().Find(shortID).Return(originalURL, true).Times(1)

	url, err := shortener.FindURL(shortID)
	assert.NoError(t, err)
	assert.Equal(t, originalURL, url)

	mockRepo.EXPECT().Find("notfound").Return("", false).Times(1)

	_, err = shortener.FindURL("notfound")
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

	mockRepo.EXPECT().Save(gomock.Any(), originalURL, userID).Return("", expectedErr).Times(1)

	_, err := shortener.ShortenID(originalURL, userID)
	assert.Error(t, err)
	assert.Equal(t, expectedErr, err)
}
