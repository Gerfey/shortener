package handler

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/Gerfey/shortener/internal/app/repository"
	"github.com/Gerfey/shortener/internal/app/service"
	"github.com/Gerfey/shortener/internal/app/settings"
)

func BenchmarkShortenHandler(b *testing.B) {
	ctx := context.Background()

	serverSettings := settings.ServerSettings{
		ServerRunAddress:       "localhost:8080",
		ServerShortenerAddress: "http://localhost:8080",
		DefaultFilePath:        "/tmp/short-url-db.json",
		ShutdownTimeout:        5 * time.Second,
	}
	s := settings.NewSettings(serverSettings)

	repo := repository.NewMemoryRepository()
	shortenerService := service.NewShortenerService(repo)
	urlService := service.NewURLService(s)
	handler := NewURLHandler(shortenerService, urlService, s, repo)

	originalURL := "https://example.com/very/long/url/path/to/some/resource?param1=value1&param2=value2"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBufferString(originalURL))
		req = req.WithContext(ctx)
		w := httptest.NewRecorder()
		handler.ShortenHandler(w, req)
	}
}

func BenchmarkRedirectURLHandler(b *testing.B) {
	ctx := context.Background()

	serverSettings := settings.ServerSettings{
		ServerRunAddress:       "localhost:8080",
		ServerShortenerAddress: "http://localhost:8080",
		DefaultFilePath:        "/tmp/short-url-db.json",
		ShutdownTimeout:        5 * time.Second,
	}
	s := settings.NewSettings(serverSettings)

	repo := repository.NewMemoryRepository()
	shortenerService := service.NewShortenerService(repo)
	urlService := service.NewURLService(s)
	handler := NewURLHandler(shortenerService, urlService, s, repo)

	originalURL := "https://example.com/very/long/url/path/to/some/resource?param1=value1&param2=value2"
	userID := "test-user-id"

	shortURL, _ := shortenerService.ShortenID(ctx, originalURL, userID)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		req := httptest.NewRequest(http.MethodGet, "/"+shortURL, nil)
		req = req.WithContext(ctx)
		w := httptest.NewRecorder()

		type ctxKey string
		const idKey ctxKey = "id"

		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			handler.RedirectURLHandler(w, r.WithContext(
				context.WithValue(r.Context(), idKey, map[string]string{"id": shortURL}),
			))
		}).ServeHTTP(w, req)
	}
}
