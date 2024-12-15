package service

import (
	"testing"

	"github.com/Gerfey/shortener/internal/app/settings"
	"github.com/stretchr/testify/assert"
)

func TestURLServiceShortenerURL(t *testing.T) {
	tests := []struct {
		name    string
		address string
		id      string
		want    string
	}{
		{
			name:    "valid address",
			address: "http://localhost:8080",
			id:      "12345",
			want:    "http://localhost:8080/12345",
		},
		{
			name:    "empty id",
			address: "http://localhost:8080",
			id:      "",
			want:    "http://localhost:8080/",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			serverSettings := settings.NewSettings(settings.ServerSettings{
				ServerRunAddress:       "localhost:8080",
				ServerShortenerAddress: tt.address,
			})

			urlService := NewURLService(serverSettings)
			shortenerURL, _ := urlService.ShortenerURL(tt.id)

			assert.Equal(t, shortenerURL, tt.want)
		})
	}
}

func TestURLService_IsValidURL(t *testing.T) {
	tests := []struct {
		name string
		url  string
		want bool
	}{
		{
			name: "Valid HTTP URL",
			url:  "http://example.com",
			want: true,
		},
		{
			name: "Valid HTTPS URL",
			url:  "https://example.com",
			want: true,
		},
		{
			name: "Valid URL with path",
			url:  "https://example.com/path",
			want: true,
		},
		{
			name: "Valid URL with query",
			url:  "https://example.com/path?query=value",
			want: true,
		},
		{
			name: "Empty URL",
			url:  "",
			want: false,
		},
		{
			name: "Invalid URL - no scheme",
			url:  "example.com",
			want: false,
		},
		{
			name: "Invalid URL - no host",
			url:  "http://",
			want: false,
		},
		{
			name: "Invalid URL - malformed",
			url:  "http:///invalid",
			want: false,
		},
		{
			name: "Invalid URL - wrong scheme",
			url:  "ftp://example.com",
			want: true,
		},
	}

	serverSettings := settings.NewSettings(settings.ServerSettings{})
	urlService := NewURLService(serverSettings)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := urlService.IsValidURL(tt.url)
			assert.Equal(t, tt.want, got)
		})
	}
}
