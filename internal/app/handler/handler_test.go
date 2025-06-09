package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Gerfey/shortener/internal/app/service"
	"github.com/Gerfey/shortener/internal/app/settings"
	"github.com/Gerfey/shortener/internal/mock"
	"github.com/Gerfey/shortener/internal/models"
	chi "github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestURLHandler_ShortenURL(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock.NewMockRepository(ctrl)
	shortener := service.NewShortenerService(mockRepo)
	appSettings := settings.NewSettings(settings.ServerSettings{
		ServerRunAddress:       "localhost:8080",
		ServerShortenerAddress: "http://localhost:8080",
	})

	handler := NewURLHandler(shortener, appSettings, mockRepo)

	tests := []struct {
		name          string
		url           string
		expectedCode  int
		mockSetup     func()
		expectedShort string
	}{
		{
			name:          "Success",
			url:           "https://example.com",
			expectedCode:  http.StatusCreated,
			expectedShort: "http://localhost:8080/abc123",
			mockSetup: func() {
				mockRepo.EXPECT().
					FindShortURL(gomock.Any(), "https://example.com").
					Return("", models.ErrURLNotFound)
				mockRepo.EXPECT().
					Save(gomock.Any(), gomock.Any(), "https://example.com", gomock.Any()).
					Return("abc123", nil)
			},
		},
		{
			name:          "URL Already Exists",
			url:           "https://example.com",
			expectedCode:  http.StatusConflict,
			expectedShort: "http://localhost:8080/existing123",
			mockSetup: func() {
				mockRepo.EXPECT().
					FindShortURL(gomock.Any(), "https://example.com").
					Return("existing123", nil)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()

			body := bytes.NewBufferString(tt.url)
			req := httptest.NewRequest(http.MethodPost, "/", body)
			req.AddCookie(&http.Cookie{Name: UserIDCookieName, Value: "test-user"})
			w := httptest.NewRecorder()

			handler.ShortenHandler(w, req)

			assert.Equal(t, tt.expectedCode, w.Code)
			assert.Equal(t, tt.expectedShort, w.Body.String())
		})
	}
}

func TestURLHandler_GetOriginalURL(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock.NewMockRepository(ctrl)
	shortener := service.NewShortenerService(mockRepo)
	appSettings := settings.NewSettings(settings.ServerSettings{
		ServerRunAddress:       "localhost:8080",
		ServerShortenerAddress: "http://localhost:8080",
	})

	handler := NewURLHandler(shortener, appSettings, mockRepo)

	tests := []struct {
		name         string
		id           string
		expectedCode int
		mockSetup    func()
		expectedURL  string
		isDeleted    bool
	}{
		{
			name:         "Success",
			id:           "abc123",
			expectedCode: http.StatusTemporaryRedirect,
			expectedURL:  "https://example.com",
			mockSetup: func() {
				mockRepo.EXPECT().
					Find(gomock.Any(), "abc123").
					Return("https://example.com", true, false)
			},
		},
		{
			name:         "Not Found",
			id:           "notfound",
			expectedCode: http.StatusNotFound,
			mockSetup: func() {
				mockRepo.EXPECT().
					Find(gomock.Any(), "notfound").
					Return("", false, false)
			},
		},
		{
			name:         "Deleted URL",
			id:           "deleted123",
			expectedCode: http.StatusGone,
			mockSetup: func() {
				mockRepo.EXPECT().
					Find(gomock.Any(), "deleted123").
					Return("", true, true)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()

			req := httptest.NewRequest(http.MethodGet, "/"+tt.id, nil)
			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("id", tt.id)
			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

			w := httptest.NewRecorder()
			handler.RedirectURLHandler(w, req)

			assert.Equal(t, tt.expectedCode, w.Code)
			if tt.expectedCode == http.StatusTemporaryRedirect {
				assert.Equal(t, tt.expectedURL, w.Header().Get("Location"))
			}
		})
	}
}

func TestURLHandler_ShortenURLJSON(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock.NewMockRepository(ctrl)
	shortener := service.NewShortenerService(mockRepo)
	appSettings := settings.NewSettings(settings.ServerSettings{
		ServerRunAddress:       "localhost:8080",
		ServerShortenerAddress: "http://localhost:8080",
	})

	handler := NewURLHandler(shortener, appSettings, mockRepo)

	tests := []struct {
		name           string
		request        models.ShortenRequest
		expectedCode   int
		mockSetup      func()
		expectedResult models.ShortenResponse
	}{
		{
			name: "Success",
			request: models.ShortenRequest{
				URL: "https://example.com",
			},
			expectedCode: http.StatusCreated,
			mockSetup: func() {
				mockRepo.EXPECT().
					FindShortURL(gomock.Any(), "https://example.com").
					Return("", models.ErrURLNotFound)
				mockRepo.EXPECT().
					Save(gomock.Any(), gomock.Any(), "https://example.com", gomock.Any()).
					Return("abc123", nil)
			},
			expectedResult: models.ShortenResponse{
				Result: "http://localhost:8080/abc123",
			},
		},
		{
			name: "URL Already Exists",
			request: models.ShortenRequest{
				URL: "https://example.com",
			},
			expectedCode: http.StatusConflict,
			mockSetup: func() {
				mockRepo.EXPECT().
					FindShortURL(gomock.Any(), "https://example.com").
					Return("existing123", nil)
			},
			expectedResult: models.ShortenResponse{
				Result: "http://localhost:8080/existing123",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()

			reqBody, _ := json.Marshal(tt.request)
			req := httptest.NewRequest(http.MethodPost, "/api/shorten", bytes.NewBuffer(reqBody))
			req.Header.Set("Content-Type", "application/json")
			req.AddCookie(&http.Cookie{Name: UserIDCookieName, Value: "test-user"})
			w := httptest.NewRecorder()

			handler.ShortenJSONHandler(w, req)
			assert.Equal(t, tt.expectedCode, w.Code)

			var response models.ShortenResponse
			err := json.NewDecoder(w.Body).Decode(&response)
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedResult, response)
		})
	}
}

func TestURLHandler_GetUserURLs(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock.NewMockRepository(ctrl)
	shortener := service.NewShortenerService(mockRepo)
	appSettings := settings.NewSettings(settings.ServerSettings{
		ServerRunAddress:       "localhost:8080",
		ServerShortenerAddress: "http://localhost:8080",
	})

	handler := NewURLHandler(shortener, appSettings, mockRepo)

	tests := []struct {
		name         string
		userID       string
		expectedCode int
		mockSetup    func()
		withCookie   bool
		expectedURLs []models.URLPair
	}{
		{
			name:         "Success with URLs",
			userID:       "user123",
			expectedCode: http.StatusOK,
			withCookie:   true,
			mockSetup: func() {
				mockRepo.EXPECT().
					GetUserURLs(gomock.Any(), "user123").
					Return([]models.URLPair{
						{ShortURL: "abc123", OriginalURL: "https://example.com"},
						{ShortURL: "def456", OriginalURL: "https://example.org"},
					}, nil)
			},
			expectedURLs: []models.URLPair{
				{ShortURL: "http://localhost:8080/abc123", OriginalURL: "https://example.com"},
				{ShortURL: "http://localhost:8080/def456", OriginalURL: "https://example.org"},
			},
		},
		{
			name:         "No URLs Found",
			userID:       "user456",
			expectedCode: http.StatusNoContent,
			withCookie:   true,
			mockSetup: func() {
				mockRepo.EXPECT().
					GetUserURLs(gomock.Any(), "user456").
					Return([]models.URLPair{}, nil)
			},
			expectedURLs: nil,
		},
		{
			name:         "No Cookie",
			expectedCode: http.StatusNoContent,
			withCookie:   false,
			mockSetup:    func() {},
			expectedURLs: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()

			req := httptest.NewRequest(http.MethodGet, "/api/user/urls", nil)
			if tt.withCookie {
				req.AddCookie(&http.Cookie{Name: UserIDCookieName, Value: tt.userID})
			}
			w := httptest.NewRecorder()

			handler.GetUserURLsHandler(w, req)

			assert.Equal(t, tt.expectedCode, w.Code)

			if tt.expectedCode == http.StatusOK {
				var response []models.URLPair
				err := json.NewDecoder(w.Body).Decode(&response)
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedURLs, response)
			}
		})
	}
}

func TestURLHandler_DeleteUserURLs(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock.NewMockRepository(ctrl)
	shortener := service.NewShortenerService(mockRepo)
	appSettings := settings.NewSettings(settings.ServerSettings{
		ServerRunAddress:       "localhost:8080",
		ServerShortenerAddress: "http://localhost:8080",
	})

	handler := NewURLHandler(shortener, appSettings, mockRepo)

	tests := []struct {
		name         string
		userID       string
		urls         []string
		expectedCode int
		mockSetup    func()
	}{
		{
			name:         "Success",
			userID:       "user123",
			urls:         []string{"abc123", "def456"},
			expectedCode: http.StatusAccepted,
			mockSetup: func() {
				mockRepo.EXPECT().
					DeleteUserURLsBatch(gomock.Any(), []string{"abc123", "def456"}, "user123").
					Return(nil)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()

			reqBody, _ := json.Marshal(tt.urls)
			req := httptest.NewRequest(http.MethodDelete, "/api/user/urls", bytes.NewBuffer(reqBody))
			req.AddCookie(&http.Cookie{Name: UserIDCookieName, Value: tt.userID})
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			handler.DeleteUserURLsHandler(w, req)

			assert.Equal(t, tt.expectedCode, w.Code)
		})
	}
}

func TestURLHandler_Ping(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock.NewMockRepository(ctrl)
	shortener := service.NewShortenerService(mockRepo)
	appSettings := settings.NewSettings(settings.ServerSettings{
		ServerRunAddress:       "localhost:8080",
		ServerShortenerAddress: "http://localhost:8080",
	})

	handler := NewURLHandler(shortener, appSettings, mockRepo)

	tests := []struct {
		name         string
		expectedCode int
		mockSetup    func()
	}{
		{
			name:         "Success",
			expectedCode: http.StatusOK,
			mockSetup: func() {
				mockRepo.EXPECT().Ping(gomock.Any()).Return(nil)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()

			req := httptest.NewRequest(http.MethodGet, "/ping", nil)
			w := httptest.NewRecorder()

			handler.PingHandler(w, req)

			assert.Equal(t, tt.expectedCode, w.Code)
		})
	}
}
