package handler

import (
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Gerfey/shortener/internal/app/service"
	"github.com/Gerfey/shortener/internal/app/settings"
	"github.com/Gerfey/shortener/internal/mock"
	"go.uber.org/mock/gomock"
)

// ErrorReadCloser имитирует ошибку при закрытии тела запроса
type ErrorReadCloser struct {
	io.Reader
}

func (e *ErrorReadCloser) Close() error {
	return errors.New("close error")
}

func TestBodyCloseError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock.NewMockRepository(ctrl)

	mockRepo.EXPECT().FindShortURL(gomock.Any(), gomock.Any()).Return("", errors.New("not found")).AnyTimes()
	mockRepo.EXPECT().Save(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return("abc123", nil).AnyTimes()
	mockRepo.EXPECT().SaveBatch(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
	mockRepo.EXPECT().DeleteUserURLsBatch(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).AnyTimes()

	shortener := service.NewShortenerService(mockRepo)
	appSettings := settings.NewSettings(settings.ServerSettings{
		ServerRunAddress:       "localhost:8080",
		ServerShortenerAddress: "http://localhost:8080",
	})
	urlService := service.NewURLService(appSettings)
	handler := NewURLHandler(shortener, urlService, appSettings, mockRepo)

	t.Run("ShortenJSONHandler Close Error", func(t *testing.T) {
		body := strings.NewReader(`{"url": "https://example.com"}`)
		req := httptest.NewRequest(http.MethodPost, "/api/shorten", body)
		req.Body = &ErrorReadCloser{body}
		req.AddCookie(&http.Cookie{Name: UserIDCookieName, Value: "user123"})
		w := httptest.NewRecorder()

		handler.ShortenJSONHandler(w, req)
	})

	t.Run("ShortenBatchHandler Close Error", func(t *testing.T) {
		body := strings.NewReader(`[{"correlation_id": "1", "original_url": "https://example.com"}]`)
		req := httptest.NewRequest(http.MethodPost, "/api/shorten/batch", body)
		req.Body = &ErrorReadCloser{body}
		req.AddCookie(&http.Cookie{Name: UserIDCookieName, Value: "user123"})
		w := httptest.NewRecorder()

		handler.ShortenBatchHandler(w, req)
	})

	t.Run("DeleteUserURLsHandler Close Error", func(t *testing.T) {
		body := strings.NewReader(`["abc123", "def456"]`)
		req := httptest.NewRequest(http.MethodDelete, "/api/user/urls", body)
		req.Body = &ErrorReadCloser{body}
		req.AddCookie(&http.Cookie{Name: UserIDCookieName, Value: "user123"})
		w := httptest.NewRecorder()

		handler.DeleteUserURLsHandler(w, req)
	})
}
