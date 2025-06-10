package usecase

import (
	"context"
	"errors"
	"testing"

	"github.com/Gerfey/shortener/internal/app/service"
	"github.com/Gerfey/shortener/internal/app/settings"
	"github.com/Gerfey/shortener/internal/mock"
	"github.com/Gerfey/shortener/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestShortenUseCase_ShortenURL(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock.NewMockRepository(ctrl)
	s := settings.NewSettings(settings.ServerSettings{
		ServerShortenerAddress: "http://localhost:8080",
		ShutdownTimeout:        0,
	})

	mockRepo.EXPECT().
		FindShortURL(gomock.Any(), "http://example.com").
		Return("", errors.New("not found"))

	mockRepo.EXPECT().
		Save(gomock.Any(), gomock.Any(), "http://example.com", "user123").
		Return("abc123", nil)

	shortenerService := service.NewShortenerService(mockRepo)
	useCase := NewShortenUseCase(shortenerService, s)

	result, err := useCase.ShortenURL(context.Background(), "http://example.com", "user123")

	require.NoError(t, err)
	assert.NotEmpty(t, result.ShortURL)
	assert.Equal(t, "http://localhost:8080/"+result.ShortURL, result.FullShortURL)
	assert.False(t, result.AlreadyExists)
}

func TestShortenUseCase_ShortenURL_EmptyURL(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock.NewMockRepository(ctrl)
	s := settings.NewSettings(settings.ServerSettings{
		ServerShortenerAddress: "http://localhost:8080",
		ShutdownTimeout:        0,
	})

	shortenerService := service.NewShortenerService(mockRepo)
	useCase := NewShortenUseCase(shortenerService, s)

	_, err := useCase.ShortenURL(context.Background(), "", "user123")

	require.Error(t, err)
	assert.Contains(t, err.Error(), "URL не может быть пустым")
}

func TestShortenUseCase_ShortenURL_AlreadyExists(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock.NewMockRepository(ctrl)
	s := settings.NewSettings(settings.ServerSettings{
		ServerShortenerAddress: "http://localhost:8080",
		ShutdownTimeout:        0,
	})

	mockRepo.EXPECT().
		FindShortURL(gomock.Any(), "http://example.com").
		Return("", errors.New("not found"))

	mockRepo.EXPECT().
		Save(gomock.Any(), gomock.Any(), "http://example.com", "user123").
		Return("", models.ErrURLExists)

	shortenerService := service.NewShortenerService(mockRepo)
	useCase := NewShortenUseCase(shortenerService, s)

	result, err := useCase.ShortenURL(context.Background(), "http://example.com", "user123")

	require.NoError(t, err)
	assert.True(t, result.AlreadyExists)
}

func TestShortenUseCase_ShortenBatch(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock.NewMockRepository(ctrl)
	s := settings.NewSettings(settings.ServerSettings{
		ServerShortenerAddress: "http://localhost:8080",
		ShutdownTimeout:        0,
	})

	batchItems := []models.BatchRequestItem{
		{CorrelationID: "1", OriginalURL: "http://example.com"},
		{CorrelationID: "2", OriginalURL: "http://example.org"},
	}

	mockRepo.EXPECT().
		SaveBatch(gomock.Any(), gomock.Any(), "user123").
		Return(nil)

	mockRepo.EXPECT().
		FindShortURL(gomock.Any(), "http://example.com").
		Return("abc123", nil)

	mockRepo.EXPECT().
		FindShortURL(gomock.Any(), "http://example.org").
		Return("def456", nil)

	shortenerService := service.NewShortenerService(mockRepo)
	useCase := NewShortenUseCase(shortenerService, s)

	result, err := useCase.ShortenBatch(context.Background(), batchItems, "user123")

	require.NoError(t, err)
	assert.Len(t, result, 2)
	assert.Equal(t, "1", result[0].CorrelationID)
	assert.Equal(t, "http://example.com", result[0].OriginalURL)
	assert.Equal(t, "abc123", result[0].ShortURL)
	assert.Equal(t, "http://localhost:8080/abc123", result[0].FullShortURL)
	assert.Equal(t, "2", result[1].CorrelationID)
	assert.Equal(t, "http://example.org", result[1].OriginalURL)
	assert.Equal(t, "def456", result[1].ShortURL)
	assert.Equal(t, "http://localhost:8080/def456", result[1].FullShortURL)
}

func TestShortenUseCase_ShortenBatch_EmptyBatch(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock.NewMockRepository(ctrl)
	s := settings.NewSettings(settings.ServerSettings{
		ServerShortenerAddress: "http://localhost:8080",
		ShutdownTimeout:        0,
	})

	shortenerService := service.NewShortenerService(mockRepo)
	useCase := NewShortenUseCase(shortenerService, s)

	_, err := useCase.ShortenBatch(context.Background(), []models.BatchRequestItem{}, "user123")

	require.Error(t, err)
	assert.Contains(t, err.Error(), "пакет не может быть пустым")
}

func TestShortenUseCase_ShortenBatch_SaveBatchError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock.NewMockRepository(ctrl)
	s := settings.NewSettings(settings.ServerSettings{
		ServerShortenerAddress: "http://localhost:8080",
		ShutdownTimeout:        0,
	})

	batchItems := []models.BatchRequestItem{
		{CorrelationID: "1", OriginalURL: "http://example.com"},
	}

	mockRepo.EXPECT().
		SaveBatch(gomock.Any(), gomock.Any(), "user123").
		Return(errors.New("database error"))

	shortenerService := service.NewShortenerService(mockRepo)
	useCase := NewShortenUseCase(shortenerService, s)

	_, err := useCase.ShortenBatch(context.Background(), batchItems, "user123")

	require.Error(t, err)
	assert.Contains(t, err.Error(), "database error")
}

func TestShortenUseCase_ShortenBatch_GetShortURLError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock.NewMockRepository(ctrl)
	s := settings.NewSettings(settings.ServerSettings{
		ServerShortenerAddress: "http://localhost:8080",
		ShutdownTimeout:        0,
	})

	batchItems := []models.BatchRequestItem{
		{CorrelationID: "1", OriginalURL: "http://example.com"},
	}

	mockRepo.EXPECT().
		SaveBatch(gomock.Any(), gomock.Any(), "user123").
		Return(nil)

	mockRepo.EXPECT().
		FindShortURL(gomock.Any(), "http://example.com").
		Return("", errors.New("not found"))

	shortenerService := service.NewShortenerService(mockRepo)
	useCase := NewShortenUseCase(shortenerService, s)

	_, err := useCase.ShortenBatch(context.Background(), batchItems, "user123")

	require.Error(t, err)
	assert.Contains(t, err.Error(), "not found")
}
