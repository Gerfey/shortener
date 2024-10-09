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
