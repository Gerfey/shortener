package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Gerfey/shortener/internal/app/settings"
	"github.com/Gerfey/shortener/internal/mock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

type MockSettings struct {
	trustedSubnet string
}

func (s *MockSettings) ServerAddress() string           { return "" }
func (s *MockSettings) BaseURL() string                 { return "" }
func (s *MockSettings) FileStoragePath() string         { return "" }
func (s *MockSettings) DatabaseDSN() string             { return "" }
func (s *MockSettings) EnableHTTPS() bool               { return false }
func (s *MockSettings) TrustedSubnet() string           { return s.trustedSubnet }
func (s *MockSettings) Server() settings.ServerSettings { return settings.ServerSettings{} }

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

func TestStatsHandler(t *testing.T) {
	tests := []struct {
		name           string
		trustedSubnet  string
		realIP         string
		urlCount       int
		expectedStatus int
		checkResponse  bool
	}{
		{
			name:           "Успешный запрос из доверенной подсети",
			trustedSubnet:  "192.168.1.0/24",
			realIP:         "192.168.1.5",
			urlCount:       5,
			expectedStatus: http.StatusOK,
			checkResponse:  true,
		},
		{
			name:           "Запрос из недоверенной подсети",
			trustedSubnet:  "192.168.1.0/24",
			realIP:         "192.168.2.5",
			urlCount:       0,
			expectedStatus: http.StatusForbidden,
			checkResponse:  false,
		},
		{
			name:           "Пустой X-Real-IP",
			trustedSubnet:  "192.168.1.0/24",
			realIP:         "",
			urlCount:       0,
			expectedStatus: http.StatusForbidden,
			checkResponse:  false,
		},
		{
			name:           "Пустая доверенная подсеть",
			trustedSubnet:  "",
			realIP:         "192.168.1.5",
			urlCount:       0,
			expectedStatus: http.StatusForbidden,
			checkResponse:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			repo := mock.NewMockRepository(ctrl)

			urlMap := make(map[string]string)
			for i := 0; i < tt.urlCount; i++ {
				urlMap["short"+string(rune('a'+i))] = "original" + string(rune('a'+i))
			}
			repo.EXPECT().All(gomock.Any()).Return(urlMap).AnyTimes()

			h := &URLHandler{
				repository: repo,
			}

			h.settings = &settings.Settings{
				Server: settings.ServerSettings{
					TrustedSubnet: tt.trustedSubnet,
				},
			}

			req := httptest.NewRequest(http.MethodGet, "/api/internal/stats", nil)
			if tt.realIP != "" {
				req.Header.Set("X-Real-IP", tt.realIP)
			}

			rr := httptest.NewRecorder()

			h.StatsHandler(rr, req)

			assert.Equal(t, tt.expectedStatus, rr.Code)

			if tt.checkResponse {
				var response StatsResponse
				err := json.Unmarshal(rr.Body.Bytes(), &response)
				require.NoError(t, err)

				assert.Equal(t, tt.urlCount, response.URLs)
				assert.GreaterOrEqual(t, response.Users, 1)
			}
		})
	}
}
