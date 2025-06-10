package usecase

import (
	"context"
	"errors"
	"testing"

	"github.com/Gerfey/shortener/internal/app/settings"
	"github.com/Gerfey/shortener/internal/mock"
	"github.com/Gerfey/shortener/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestUserURLsUseCase_GetUserURLs(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock.NewMockRepository(ctrl)
	s := settings.NewSettings(settings.ServerSettings{
		ServerShortenerAddress: "http://localhost:8080",
		ShutdownTimeout:        0,
	})

	urls := []models.URLPair{
		{ShortURL: "abc123", OriginalURL: "http://example.com"},
		{ShortURL: "def456", OriginalURL: "http://example.org"},
	}

	mockRepo.EXPECT().
		GetUserURLs(gomock.Any(), "user123").
		Return(urls, nil)

	useCase := NewUserURLsUseCase(mockRepo, s)

	result, err := useCase.GetUserURLs(context.Background(), "user123")

	require.NoError(t, err)
	assert.Len(t, result, 2)
	assert.Equal(t, "http://localhost:8080/abc123", result[0].ShortURL)
	assert.Equal(t, "http://example.com", result[0].OriginalURL)
	assert.Equal(t, "http://localhost:8080/def456", result[1].ShortURL)
	assert.Equal(t, "http://example.org", result[1].OriginalURL)
}

func TestUserURLsUseCase_GetUserURLs_EmptyUserID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock.NewMockRepository(ctrl)
	s := settings.NewSettings(settings.ServerSettings{
		ServerShortenerAddress: "http://localhost:8080",
		ShutdownTimeout:        0,
	})

	useCase := NewUserURLsUseCase(mockRepo, s)

	_, err := useCase.GetUserURLs(context.Background(), "")

	require.Error(t, err)
	assert.Contains(t, err.Error(), "ID пользователя не может быть пустым")
}

func TestUserURLsUseCase_GetUserURLs_RepositoryError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock.NewMockRepository(ctrl)
	s := settings.NewSettings(settings.ServerSettings{
		ServerShortenerAddress: "http://localhost:8080",
		ShutdownTimeout:        0,
	})

	mockRepo.EXPECT().
		GetUserURLs(gomock.Any(), "user123").
		Return(nil, errors.New("database error"))

	useCase := NewUserURLsUseCase(mockRepo, s)

	_, err := useCase.GetUserURLs(context.Background(), "user123")

	require.Error(t, err)
	assert.Contains(t, err.Error(), "database error")
}

func TestUserURLsUseCase_DeleteUserURLs(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock.NewMockRepository(ctrl)
	s := settings.NewSettings(settings.ServerSettings{
		ServerShortenerAddress: "http://localhost:8080",
		ShutdownTimeout:        0,
	})

	mockRepo.EXPECT().
		DeleteUserURLsBatch(gomock.Any(), []string{"abc123", "def456"}, "user123").
		Return(nil)

	useCase := NewUserURLsUseCase(mockRepo, s)

	err := useCase.DeleteUserURLs(context.Background(), "user123", []string{"abc123", "def456"})

	require.NoError(t, err)
}

func TestUserURLsUseCase_DeleteUserURLs_EmptyUserID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock.NewMockRepository(ctrl)
	s := settings.NewSettings(settings.ServerSettings{
		ServerShortenerAddress: "http://localhost:8080",
		ShutdownTimeout:        0,
	})

	useCase := NewUserURLsUseCase(mockRepo, s)

	err := useCase.DeleteUserURLs(context.Background(), "", []string{"abc123", "def456"})

	require.Error(t, err)
	assert.Contains(t, err.Error(), "ID пользователя не может быть пустым")
}

func TestUserURLsUseCase_DeleteUserURLs_EmptyURLs(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock.NewMockRepository(ctrl)
	s := settings.NewSettings(settings.ServerSettings{
		ServerShortenerAddress: "http://localhost:8080",
		ShutdownTimeout:        0,
	})

	useCase := NewUserURLsUseCase(mockRepo, s)

	err := useCase.DeleteUserURLs(context.Background(), "user123", []string{})

	require.Error(t, err)
	assert.Contains(t, err.Error(), "список URL не может быть пустым")
}

func TestUserURLsUseCase_DeleteUserURLs_RepositoryError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock.NewMockRepository(ctrl)
	s := settings.NewSettings(settings.ServerSettings{
		ServerShortenerAddress: "http://localhost:8080",
		ShutdownTimeout:        0,
	})

	mockRepo.EXPECT().
		DeleteUserURLsBatch(gomock.Any(), []string{"abc123", "def456"}, "user123").
		Return(errors.New("database error"))

	useCase := NewUserURLsUseCase(mockRepo, s)

	err := useCase.DeleteUserURLs(context.Background(), "user123", []string{"abc123", "def456"})

	require.Error(t, err)
	assert.Contains(t, err.Error(), "database error")
}
