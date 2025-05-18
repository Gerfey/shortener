package service

import (
	"context"
	"testing"

	"github.com/Gerfey/shortener/internal/app/repository"
)

func BenchmarkShortenID(b *testing.B) {
	ctx := context.Background()
	repo := repository.NewMemoryRepository()
	service := NewShortenerService(repo)
	userID := "test-user-id"
	originalURL := "https://example.com/very/long/url/path/to/some/resource?param1=value1&param2=value2"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = service.ShortenID(ctx, originalURL, userID)
	}
}

func BenchmarkFindOriginalURL(b *testing.B) {
	ctx := context.Background()
	repo := repository.NewMemoryRepository()
	service := NewShortenerService(repo)
	userID := "test-user-id"
	originalURL := "https://example.com/very/long/url/path/to/some/resource?param1=value1&param2=value2"

	shortURL, _ := service.ShortenID(ctx, originalURL, userID)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _, _ = repo.Find(ctx, shortURL)
	}
}
