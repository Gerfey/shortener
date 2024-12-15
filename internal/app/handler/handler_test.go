package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/Gerfey/shortener/internal/app/service"
	"github.com/Gerfey/shortener/internal/app/settings"
	"github.com/Gerfey/shortener/internal/models"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"
	"sort"
	"testing"
)

type MockRepository struct {
	urls     map[string]string
	userURLs map[string][]models.URLPair
}

func NewMockRepository() *MockRepository {
	return &MockRepository{
		urls:     make(map[string]string),
		userURLs: make(map[string][]models.URLPair),
	}
}

func (m *MockRepository) Save(shortURL, originalURL, userID string) (string, error) {
	m.urls[shortURL] = originalURL
	return shortURL, nil
}

func (m *MockRepository) Find(shortURL string) (string, bool) {
	originalURL, exists := m.urls[shortURL]
	return originalURL, exists
}

func (m *MockRepository) All() map[string]string {
	return m.urls
}

func (m *MockRepository) GetUserURLs(userID string) ([]models.URLPair, error) {
	urls, exists := m.userURLs[userID]
	if !exists {
		return []models.URLPair{}, nil
	}
	return urls, nil
}

func (m *MockRepository) SaveBatch(urls map[string]string, userID string) error {
	for shortURL, originalURL := range urls {
		m.urls[shortURL] = originalURL
	}
	return nil
}

func (m *MockRepository) FindShortURL(originalURL string) (string, error) {
	for shortURL, origURL := range m.urls {
		if origURL == originalURL {
			return shortURL, nil
		}
	}
	return "", nil
}

func (m *MockRepository) Ping() error {
	return nil
}

func TestURLHandler_GetUserURLsHandler(t *testing.T) {
	mockRepo := NewMockRepository()
	shortenerService := service.NewShortenerService(mockRepo)
	appSettings := settings.NewSettings(settings.ServerSettings{
		ServerRunAddress:       "localhost:8080",
		ServerShortenerAddress: "http://localhost:8080",
	})
	urlService := service.NewURLService(appSettings)
	handler := NewURLHandler(shortenerService, urlService, appSettings, mockRepo)

	userID := "test_user"
	testURLs := []models.URLPair{
		{ShortURL: "abc123", OriginalURL: "http://example.com/1"},
		{ShortURL: "def456", OriginalURL: "http://example.com/2"},
	}
	mockRepo.userURLs[userID] = testURLs

	req := httptest.NewRequest(http.MethodGet, "/api/user/urls", nil)
	req.AddCookie(&http.Cookie{Name: "user_id", Value: userID})

	rr := httptest.NewRecorder()

	handler.GetUserURLsHandler(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	body, err := io.ReadAll(rr.Body)
	require.NoError(t, err)

	var response []models.URLPair
	err = json.Unmarshal(body, &response)
	require.NoError(t, err)

	sort.Slice(testURLs, func(i, j int) bool {
		return testURLs[i].ShortURL < testURLs[j].ShortURL
	})
	sort.Slice(response, func(i, j int) bool {
		return response[i].ShortURL < response[j].ShortURL
	})

	assert.Equal(t, testURLs, response)
}

func TestURLHandler_HandleGetURL(t *testing.T) {
	mockRepo := NewMockRepository()
	shortenerService := service.NewShortenerService(mockRepo)
	appSettings := settings.NewSettings(settings.ServerSettings{
		ServerRunAddress:       "localhost:8080",
		ServerShortenerAddress: "http://localhost:8080",
	})
	urlService := service.NewURLService(appSettings)
	handler := NewURLHandler(shortenerService, urlService, appSettings, mockRepo)

	shortURL := "abc123"
	originalURL := "http://example.com"
	_, err := mockRepo.Save(shortURL, originalURL, "test_user")
	require.NoError(t, err)

	req := httptest.NewRequest(http.MethodGet, "/"+shortURL, nil)
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", shortURL)
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	rr := httptest.NewRecorder()

	handler.RedirectURLHandler(rr, req)

	assert.Equal(t, http.StatusTemporaryRedirect, rr.Code)

	location := rr.Header().Get("Location")
	assert.Equal(t, originalURL, location)
}

func TestURLHandler_ShortenHandler(t *testing.T) {
	mockRepo := NewMockRepository()
	shortenerService := service.NewShortenerService(mockRepo)
	appSettings := settings.NewSettings(settings.ServerSettings{
		ServerRunAddress:       "localhost:8080",
		ServerShortenerAddress: "http://localhost:8080",
	})
	urlService := service.NewURLService(appSettings)
	handler := NewURLHandler(shortenerService, urlService, appSettings, mockRepo)

	tests := []struct {
		name           string
		requestURL     string
		expectedStatus int
		expectedError  bool
	}{
		{
			name:           "Valid URL",
			requestURL:     "http://example.com",
			expectedStatus: http.StatusCreated,
			expectedError:  false,
		},
		{
			name:           "Empty URL",
			requestURL:     "",
			expectedStatus: http.StatusBadRequest,
			expectedError:  true,
		},
		{
			name:           "Invalid URL",
			requestURL:     "not-a-url",
			expectedStatus: http.StatusBadRequest,
			expectedError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reqBody := []byte(tt.requestURL)
			req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(reqBody))
			req.AddCookie(&http.Cookie{Name: "user_id", Value: "test_user"})

			rr := httptest.NewRecorder()

			handler.ShortenHandler(rr, req)

			assert.Equal(t, tt.expectedStatus, rr.Code)

			if !tt.expectedError {
				assert.NotEmpty(t, rr.Body.String())
			}
		})
	}
}

func TestURLHandler_ShortenJSONHandler(t *testing.T) {
	mockRepo := NewMockRepository()
	shortenerService := service.NewShortenerService(mockRepo)
	appSettings := settings.NewSettings(settings.ServerSettings{
		ServerRunAddress:       "localhost:8080",
		ServerShortenerAddress: "http://localhost:8080",
	})
	urlService := service.NewURLService(appSettings)
	handler := NewURLHandler(shortenerService, urlService, appSettings, mockRepo)

	tests := []struct {
		name           string
		requestBody    map[string]string
		expectedStatus int
		expectedError  bool
	}{
		{
			name: "Valid URL",
			requestBody: map[string]string{
				"url": "http://example.com",
			},
			expectedStatus: http.StatusCreated,
			expectedError:  false,
		},
		{
			name: "Empty URL",
			requestBody: map[string]string{
				"url": "",
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  true,
		},
		{
			name:           "Invalid JSON",
			requestBody:    nil,
			expectedStatus: http.StatusBadRequest,
			expectedError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var reqBody []byte
			var err error

			if tt.requestBody != nil {
				reqBody, err = json.Marshal(tt.requestBody)
				require.NoError(t, err)
			} else {
				reqBody = []byte("{invalid json}")
			}

			req := httptest.NewRequest(http.MethodPost, "/api/shorten", bytes.NewBuffer(reqBody))
			req.Header.Set("Content-Type", "application/json")
			req.AddCookie(&http.Cookie{Name: "user_id", Value: "test_user"})

			rr := httptest.NewRecorder()

			handler.ShortenJSONHandler(rr, req)

			assert.Equal(t, tt.expectedStatus, rr.Code)

			if !tt.expectedError {
				var response map[string]string
				err = json.NewDecoder(rr.Body).Decode(&response)
				require.NoError(t, err)
				assert.NotEmpty(t, response["result"])
			}
		})
	}
}

func TestURLHandler_ShortenBatchHandler(t *testing.T) {
	mockRepo := NewMockRepository()
	shortenerService := service.NewShortenerService(mockRepo)
	appSettings := settings.NewSettings(settings.ServerSettings{
		ServerRunAddress:       "localhost:8080",
		ServerShortenerAddress: "http://localhost:8080",
	})
	urlService := service.NewURLService(appSettings)
	handler := NewURLHandler(shortenerService, urlService, appSettings, mockRepo)

	tests := []struct {
		name           string
		requestBody    []models.BatchRequestItem
		expectedStatus int
		expectedError  bool
	}{
		{
			name: "Valid URLs",
			requestBody: []models.BatchRequestItem{
				{
					CorrelationID: "1",
					OriginalURL:   "http://example1.com",
				},
				{
					CorrelationID: "2",
					OriginalURL:   "http://example2.com",
				},
			},
			expectedStatus: http.StatusCreated,
			expectedError:  false,
		},
		{
			name:           "Empty request",
			requestBody:    []models.BatchRequestItem{},
			expectedStatus: http.StatusBadRequest,
			expectedError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reqBody, err := json.Marshal(tt.requestBody)
			require.NoError(t, err)

			req := httptest.NewRequest(http.MethodPost, "/api/shorten/batch", bytes.NewBuffer(reqBody))
			req.Header.Set("Content-Type", "application/json")
			req.AddCookie(&http.Cookie{Name: "user_id", Value: "test_user"})

			rr := httptest.NewRecorder()

			handler.ShortenBatchHandler(rr, req)

			assert.Equal(t, tt.expectedStatus, rr.Code)

			if !tt.expectedError {
				var response []models.BatchResponseItem
				err = json.NewDecoder(rr.Body).Decode(&response)
				require.NoError(t, err)
				assert.Equal(t, len(tt.requestBody), len(response))
				for _, resp := range response {
					assert.NotEmpty(t, resp.ShortURL)
				}
			}
		})
	}
}
