package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Gerfey/shortener/internal/app/service"
	"github.com/Gerfey/shortener/internal/app/settings"
	"github.com/Gerfey/shortener/internal/mock"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestURLHandler_DeleteUserURLsHandler(t *testing.T) {
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
		withCookie   bool
		requestBody  []string
		expectedCode int
		mockSetup    func()
	}{
		{
			name:         "Success",
			userID:       "user123",
			withCookie:   true,
			requestBody:  []string{"abc123", "def456"},
			expectedCode: http.StatusAccepted,
			mockSetup: func() {
				mockRepo.EXPECT().
					DeleteUserURLsBatch(gomock.Any(), []string{"abc123", "def456"}, "user123").
					Return(nil)
			},
		},
		{
			name:         "No Cookie",
			withCookie:   false,
			requestBody:  []string{"abc123", "def456"},
			expectedCode: http.StatusUnauthorized,
			mockSetup:    func() {},
		},
		{
			name:         "Invalid JSON",
			userID:       "user123",
			withCookie:   true,
			requestBody:  nil,
			expectedCode: http.StatusBadRequest,
			mockSetup:    func() {},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()

			var body []byte
			var err error
			if tt.name == "Invalid JSON" {
				body = []byte(`invalid json`)
			} else {
				body, err = json.Marshal(tt.requestBody)
				assert.NoError(t, err)
			}

			req := httptest.NewRequest(http.MethodDelete, "/api/user/urls", bytes.NewBuffer(body))
			if tt.withCookie {
				req.AddCookie(&http.Cookie{Name: UserIDCookieName, Value: tt.userID})
			}
			w := httptest.NewRecorder()

			handler.DeleteUserURLsHandler(w, req)

			assert.Equal(t, tt.expectedCode, w.Code)
		})
	}
}
