package usecase

import (
	"context"
	"testing"

	"github.com/Gerfey/shortener/internal/app/settings"
	"github.com/Gerfey/shortener/internal/mock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestStatsUseCase_GetStats(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock.NewMockRepository(ctrl)
	s := settings.NewSettings(settings.ServerSettings{
		TrustedSubnet: "192.168.1.0/24",
	})

	urls := map[string]string{
		"abc123": "http://example.com",
		"def456": "http://example.org",
	}

	mockRepo.EXPECT().
		All(gomock.Any()).
		Return(urls)

	useCase := NewStatsUseCase(mockRepo, s)

	result, err := useCase.GetStats(context.Background(), "192.168.1.5")

	require.NoError(t, err)
	assert.Equal(t, 2, result.URLs)
	assert.Equal(t, 1, result.Users)
}

func TestStatsUseCase_GetStats_EmptyIP(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock.NewMockRepository(ctrl)
	s := settings.NewSettings(settings.ServerSettings{
		TrustedSubnet: "192.168.1.0/24",
	})

	useCase := NewStatsUseCase(mockRepo, s)

	_, err := useCase.GetStats(context.Background(), "")

	require.Error(t, err)
	assert.Contains(t, err.Error(), "IP-адрес не может быть пустым")
}

func TestStatsUseCase_GetStats_NoTrustedSubnet(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock.NewMockRepository(ctrl)
	s := settings.NewSettings(settings.ServerSettings{})

	useCase := NewStatsUseCase(mockRepo, s)

	_, err := useCase.GetStats(context.Background(), "192.168.1.5")

	require.Error(t, err)
	assert.Contains(t, err.Error(), "доверенная подсеть не настроена")
}

func TestStatsUseCase_GetStats_IPNotInCIDR(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock.NewMockRepository(ctrl)
	s := settings.NewSettings(settings.ServerSettings{
		TrustedSubnet: "192.168.1.0/24",
	})

	useCase := NewStatsUseCase(mockRepo, s)

	_, err := useCase.GetStats(context.Background(), "10.0.0.1")

	require.Error(t, err)
	assert.Contains(t, err.Error(), "доступ запрещен")
}

func TestStatsUseCase_GetStats_InvalidIP(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock.NewMockRepository(ctrl)
	s := settings.NewSettings(settings.ServerSettings{
		TrustedSubnet: "192.168.1.0/24",
	})

	useCase := NewStatsUseCase(mockRepo, s)

	_, err := useCase.GetStats(context.Background(), "invalid-ip")

	require.Error(t, err)
	assert.Contains(t, err.Error(), "доступ запрещен")
}

func TestStatsUseCase_GetStats_InvalidCIDR(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock.NewMockRepository(ctrl)
	s := settings.NewSettings(settings.ServerSettings{
		TrustedSubnet: "invalid-cidr",
	})

	useCase := NewStatsUseCase(mockRepo, s)

	_, err := useCase.GetStats(context.Background(), "192.168.1.5")

	require.Error(t, err)
	assert.Contains(t, err.Error(), "доступ запрещен")
}
