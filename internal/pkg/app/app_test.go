package app

import (
	"bytes"
	"encoding/json"
	"github.com/Gerfey/shortener/internal/app/handler"
	"github.com/Gerfey/shortener/internal/app/repository"
	"github.com/Gerfey/shortener/internal/app/service"
	"github.com/Gerfey/shortener/internal/app/settings"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestShortenerApp_Run(t *testing.T) {
	configApplication := settings.NewSettings(
		settings.ServerSettings{
			ServerShortenerAddress: "http://localhost:8080",
		},
	)

	mockRepo := repository.NewMemoryRepository()
	mockShortenerService := service.NewShortenerService(mockRepo)
	mockURLService := service.NewURLService(configApplication)

	urlHandler := handler.NewURLHandler(mockShortenerService, mockURLService, configApplication)

	router := chi.NewRouter()

	router.Post("/api/shorten", urlHandler.ShortenJSONHandler)

	server := httptest.NewServer(router)
	defer server.Close()

	reqBody := map[string]string{"url": "https://example.com"}
	reqBytes, _ := json.Marshal(reqBody)

	resp, err := http.Post(server.URL+"/api/shorten", "application/json", bytes.NewReader(reqBytes))
	assert.NoError(t, err)

	defer resp.Body.Close()

	assert.Equal(t, http.StatusCreated, resp.StatusCode)

	var shortenResp map[string]string
	err = json.NewDecoder(resp.Body).Decode(&shortenResp)
	assert.NoError(t, err)

	assert.Contains(t, shortenResp["result"], "http://localhost:8080")
}
