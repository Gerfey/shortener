package usecase

import (
	"context"

	"github.com/Gerfey/shortener/internal/models"
)

type PingUseCase struct {
	repository models.Repository
}

func NewPingUseCase(repository models.Repository) *PingUseCase {
	return &PingUseCase{
		repository: repository,
	}
}

func (uc *PingUseCase) Ping(ctx context.Context) error {
	return uc.repository.Ping(ctx)
}
