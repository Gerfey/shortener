package usecase

import (
	"context"
	"errors"
	"net"

	"github.com/Gerfey/shortener/internal/app/settings"
	"github.com/Gerfey/shortener/internal/models"
)

type StatsUseCase struct {
	repository models.Repository
	settings   *settings.Settings
}

type StatsResult struct {
	URLs  int
	Users int
}

func NewStatsUseCase(repository models.Repository, settings *settings.Settings) *StatsUseCase {
	return &StatsUseCase{
		repository: repository,
		settings:   settings,
	}
}

func (uc *StatsUseCase) GetStats(ctx context.Context, clientIP string) (StatsResult, error) {
	if clientIP == "" {
		return StatsResult{}, errors.New("IP-адрес не может быть пустым")
	}

	trustedSubnet := uc.settings.TrustedSubnet()
	if trustedSubnet == "" {
		return StatsResult{}, errors.New("доверенная подсеть не настроена")
	}

	if !isIPInCIDR(clientIP, trustedSubnet) {
		return StatsResult{}, errors.New("доступ запрещен")
	}

	urls := uc.repository.All(ctx)
	
	return StatsResult{
		URLs:  len(urls),
		Users: 1, // Минимум один пользователь
	}, nil
}

func isIPInCIDR(ipStr, cidrStr string) bool {
	ip := net.ParseIP(ipStr)
	if ip == nil {
		return false
	}

	_, ipNet, err := net.ParseCIDR(cidrStr)
	if err != nil {
		return false
	}

	return ipNet.Contains(ip)
}
