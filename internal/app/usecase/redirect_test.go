package usecase

import (
	"context"
	"testing"

	"github.com/Gerfey/shortener/internal/mock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestRedirectUseCase_GetOriginalURL(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock.NewMockRepository(ctrl)

	mockRepo.EXPECT().
		Find(gomock.Any(), "abc123").
		Return("http://example.com", true, false)

	useCase := NewRedirectUseCase(mockRepo)

	result, err := useCase.GetOriginalURL(context.Background(), "abc123")

	require.NoError(t, err)
	assert.Equal(t, "http://example.com", result.OriginalURL)
	assert.False(t, result.IsDeleted)
}

func TestRedirectUseCase_GetOriginalURL_EmptyID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock.NewMockRepository(ctrl)
	useCase := NewRedirectUseCase(mockRepo)

	_, err := useCase.GetOriginalURL(context.Background(), "")

	require.Error(t, err)
	assert.Contains(t, err.Error(), "ID не может быть пустым")
}

func TestRedirectUseCase_GetOriginalURL_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock.NewMockRepository(ctrl)

	mockRepo.EXPECT().
		Find(gomock.Any(), "abc123").
		Return("", false, false)

	useCase := NewRedirectUseCase(mockRepo)

	_, err := useCase.GetOriginalURL(context.Background(), "abc123")

	require.Error(t, err)
	assert.Contains(t, err.Error(), "URL не найден")
}

func TestRedirectUseCase_GetOriginalURL_Deleted(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock.NewMockRepository(ctrl)

	mockRepo.EXPECT().
		Find(gomock.Any(), "abc123").
		Return("http://example.com", true, true)

	useCase := NewRedirectUseCase(mockRepo)

	result, err := useCase.GetOriginalURL(context.Background(), "abc123")

	require.NoError(t, err)
	assert.Equal(t, "http://example.com", result.OriginalURL)
	assert.True(t, result.IsDeleted)
}
