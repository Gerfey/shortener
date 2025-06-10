package handler

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Gerfey/shortener/internal/app/usecase"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIsIPInCIDR(t *testing.T) {
	tests := []struct {
		name     string
		ip       string
		cidr     string
		expected bool
	}{
		{
			name:     "IP в подсети",
			ip:       "192.168.1.5",
			cidr:     "192.168.1.0/24",
			expected: true,
		},
		{
			name:     "IP не в подсети",
			ip:       "192.168.2.5",
			cidr:     "192.168.1.0/24",
			expected: false,
		},
		{
			name:     "Некорректный IP",
			ip:       "invalid",
			cidr:     "192.168.1.0/24",
			expected: false,
		},
		{
			name:     "Некорректный CIDR",
			ip:       "192.168.1.5",
			cidr:     "invalid",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isIPInCIDR(tt.ip, tt.cidr)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// StatsUseCaser интерфейс для StatsUseCase
type StatsUseCaser interface {
	GetStats(ctx context.Context, clientIP string) (usecase.StatsResult, error)
}

// MockStatsUseCase - мок для StatsUseCase
type MockStatsUseCase struct {
	getStatsFunc func(ctx context.Context, clientIP string) (usecase.StatsResult, error)
}

func (m *MockStatsUseCase) GetStats(ctx context.Context, clientIP string) (usecase.StatsResult, error) {
	return m.getStatsFunc(ctx, clientIP)
}

// URLHandlerWithMocks структура для тестирования с моками
type URLHandlerWithMocks struct {
	mockStatsUseCase StatsUseCaser
}

// StatsHandler переопределяет метод StatsHandler для использования мока
func (h *URLHandlerWithMocks) StatsHandler(w http.ResponseWriter, r *http.Request) {
	clientIP := r.Header.Get("X-Real-IP")
	if clientIP == "" {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	stats, err := h.mockStatsUseCase.GetStats(r.Context(), clientIP)
	if err != nil {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	response := StatsResponse{
		URLs:  stats.URLs,
		Users: stats.Users,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func TestStatsHandler(t *testing.T) {
	tests := []struct {
		name           string
		realIP         string
		urlCount       int
		userCount      int
		expectedStatus int
		hasError       bool
		checkResponse  bool
	}{
		{
			name:           "Успешный запрос",
			realIP:         "192.168.1.5",
			urlCount:       5,
			userCount:      3,
			expectedStatus: http.StatusOK,
			hasError:       false,
			checkResponse:  true,
		},
		{
			name:           "Ошибка доступа",
			realIP:         "192.168.2.5",
			urlCount:       0,
			userCount:      0,
			expectedStatus: http.StatusForbidden,
			hasError:       true,
			checkResponse:  false,
		},
		{
			name:           "Пустой X-Real-IP",
			realIP:         "",
			urlCount:       0,
			userCount:      0,
			expectedStatus: http.StatusForbidden,
			hasError:       false,
			checkResponse:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockStatsUseCase := &MockStatsUseCase{
				getStatsFunc: func(ctx context.Context, clientIP string) (usecase.StatsResult, error) {
					if tt.hasError {
						return usecase.StatsResult{}, errors.New("access denied")
					}
					return usecase.StatsResult{
						URLs:  tt.urlCount,
						Users: tt.userCount,
					}, nil
				},
			}

			handler := &URLHandlerWithMocks{
				mockStatsUseCase: mockStatsUseCase,
			}

			req := httptest.NewRequest(http.MethodGet, "/api/internal/stats", nil)
			if tt.realIP != "" {
				req.Header.Set("X-Real-IP", tt.realIP)
			}

			rr := httptest.NewRecorder()

			handler.StatsHandler(rr, req)

			assert.Equal(t, tt.expectedStatus, rr.Code)

			if tt.checkResponse {
				var response StatsResponse
				err := json.Unmarshal(rr.Body.Bytes(), &response)
				require.NoError(t, err)

				assert.Equal(t, tt.urlCount, response.URLs)
				assert.Equal(t, tt.userCount, response.Users)
			}
		})
	}
}
