package usecase

import (
	"context"
	"errors"
	"testing"

	"github.com/Gerfey/shortener/internal/mock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestPingUseCase_Ping(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock.NewMockRepository(ctrl)

	mockRepo.EXPECT().
		Ping(gomock.Any()).
		Return(nil)

	useCase := NewPingUseCase(mockRepo)

	err := useCase.Ping(context.Background())

	require.NoError(t, err)
}

func TestPingUseCase_Ping_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock.NewMockRepository(ctrl)

	mockRepo.EXPECT().
		Ping(gomock.Any()).
		Return(errors.New("database error"))

	useCase := NewPingUseCase(mockRepo)

	err := useCase.Ping(context.Background())

	require.Error(t, err)
	assert.Contains(t, err.Error(), "database error")
}
