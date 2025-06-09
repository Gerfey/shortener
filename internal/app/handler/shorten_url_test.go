package handler

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Gerfey/shortener/internal/app/service"
	"github.com/Gerfey/shortener/internal/app/settings"
	"github.com/Gerfey/shortener/internal/mock"
	"github.com/Gerfey/shortener/internal/models"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestURLHandler_ShortenURLHandler(t *testing.T) {
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
		method        string
		url           string
		userID        string
		expectedCode  int
		expectedBody  string
		mockSetup     func()
	}{
		{
			name:          "Success",
			method:        http.MethodPost,
			url:           "https://example.com",
			userID:        "user123",
			expectedCode:  http.StatusCreated,
			expectedBody:  "http://localhost:8080/abc123",
			mockSetup: func() {
				mockRepo.EXPECT().
					FindShortURL(gomock.Any(), "https://example.com").
					Return("", models.ErrURLNotFound)
				mockRepo.EXPECT().
					Save(gomock.Any(), gomock.Any(), "https://example.com", "user123").
					Return("abc123", nil)
			},
		},
		{
			name:          "URL Already Exists",
			method:        http.MethodPost,
			url:           "https://example.com",
			userID:        "user123",
			expectedCode:  http.StatusConflict,
			expectedBody:  "http://localhost:8080/abc123",
			mockSetup: func() {
				mockRepo.EXPECT().
					FindShortURL(gomock.Any(), "https://example.com").
					Return("abc123", nil)
			},
		},
		{
			name:          "Method Not Allowed",
			method:        http.MethodGet,
			url:           "https://example.com",
			userID:        "user123",
			expectedCode:  http.StatusMethodNotAllowed,
			expectedBody:  "Method not allowed\n",
			mockSetup:     func() {},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()

			req := httptest.NewRequest(tt.method, "/", bytes.NewBufferString(tt.url))
			req.AddCookie(&http.Cookie{Name: UserIDCookieName, Value: tt.userID})
			w := httptest.NewRecorder()

			handler.ShortenURLHandler(w, req)

			assert.Equal(t, tt.expectedCode, w.Code)
			assert.Equal(t, tt.expectedBody, w.Body.String())
		})
	}
}
