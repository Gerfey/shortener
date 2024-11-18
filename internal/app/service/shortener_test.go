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

	mockRepo.EXPECT().Save(gomock.Any(), originalURL).Return(shortID, nil).Times(1)

	id, err := shortener.ShortenID(originalURL)
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

	shortID := "s65fg"
	originalURL := "https://example.com"

	mockRepo.EXPECT().Save(gomock.Any(), originalURL).Return(shortID, errors.New("database error"))

	_, err := shortener.ShortenID("https://example.com")
	assert.Error(t, err)
	assert.Equal(t, "database error", err.Error())
}
