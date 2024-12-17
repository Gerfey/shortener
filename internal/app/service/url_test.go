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
		wantErr bool
	}{
		{
			name:    "valid address",
			address: "http://localhost:8080",
			id:      "12345",
			want:    "http://localhost:8080/12345",
			wantErr: false,
		},
		{
			name:    "empty id",
			address: "http://localhost:8080",
			id:      "",
			want:    "http://localhost:8080/",
			wantErr: false,
		},
		{
			name:    "invalid address",
			address: "://invalid",
			id:      "12345",
			want:    "",
			wantErr: true,
		},
		{
			name:    "address with path",
			address: "http://localhost:8080/api",
			id:      "12345",
			want:    "http://localhost:8080/api/12345",
			wantErr: false,
		},
		{
			name:    "address with query params",
			address: "http://localhost:8080?version=1",
			id:      "12345",
			want:    "http://localhost:8080/12345",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			serverSettings := settings.NewSettings(settings.ServerSettings{
				ServerRunAddress:       "localhost:8080",
				ServerShortenerAddress: tt.address,
			})

			urlService := NewURLService(serverSettings)
			shortenerURL, err := urlService.ShortenerURL(tt.id)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Empty(t, shortenerURL)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, shortenerURL)
			}
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

func TestFormatURL(t *testing.T) {
	tests := []struct {
		name    string
		url     string
		want    string
		wantErr bool
	}{
		{
			name:    "Valid URL",
			url:     "http://example.com",
			want:    "http://example.com",
			wantErr: false,
		},
		{
			name:    "URL with path",
			url:     "http://example.com/path",
			want:    "http://example.com/path",
			wantErr: false,
		},
		{
			name:    "URL with query",
			url:     "http://example.com?param=value",
			want:    "http://example.com?param=value",
			wantErr: false,
		},
		{
			name:    "URL with fragment",
			url:     "http://example.com#section",
			want:    "http://example.com#section",
			wantErr: false,
		},
		{
			name:    "Invalid URL",
			url:     "://invalid",
			want:    "",
			wantErr: true,
		},
		{
			name:    "Empty URL",
			url:     "",
			want:    "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := formatURL(tt.url)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Empty(t, got)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}
