package handler

import (
	"fmt"
	"github.com/Gerfey/shortener/internal/app/repository"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Gerfey/shortener/internal/app/service"
	"github.com/Gerfey/shortener/internal/app/settings"
	"github.com/stretchr/testify/assert"
)

func TestShortenURLHandler(t *testing.T) {
	testCase := []struct {
		method       string
		expectedCode int
	}{
		{method: http.MethodGet, expectedCode: http.StatusMethodNotAllowed},
		{method: http.MethodPut, expectedCode: http.StatusMethodNotAllowed},
		{method: http.MethodDelete, expectedCode: http.StatusMethodNotAllowed},
		{method: http.MethodPost, expectedCode: http.StatusCreated},
	}

	path := "test.json"

	for _, tc := range testCase {
		t.Run(tc.method, func(t *testing.T) {
			r := httptest.NewRequest(tc.method, "/", nil)
			w := httptest.NewRecorder()

			s := settings.NewSettings(
				settings.ServerSettings{ServerRunAddress: "", ServerShortenerAddress: "", DefaultFilePath: path},
			)

			fileStorageService := service.NewFileStorage(path)
			memoryRepository := repository.NewURLMemoryRepository()
			shortenerService := service.NewShortenerService(memoryRepository, fileStorageService)
			URLService := service.NewURLService(s)

			e := NewURLHandler(shortenerService, URLService)

			e.ShortenURLHandler(w, r)

			assert.Equal(t, tc.expectedCode, w.Code, "Код ответа не совпадает с ожидаемым")
		})
	}
}

func TestRedirectURLHandler(t *testing.T) {
	testCase := []struct {
		method       string
		expectedCode int
		setPathValue bool
		expectedURL  string
	}{
		{method: http.MethodGet, expectedCode: http.StatusTemporaryRedirect, setPathValue: true, expectedURL: "https://example.com"},
		{method: http.MethodGet, expectedCode: http.StatusBadRequest, setPathValue: false, expectedURL: ""},
		{method: http.MethodPut, expectedCode: http.StatusMethodNotAllowed, setPathValue: false, expectedURL: ""},
		{method: http.MethodDelete, expectedCode: http.StatusMethodNotAllowed, setPathValue: false, expectedURL: ""},
		{method: http.MethodPost, expectedCode: http.StatusMethodNotAllowed, setPathValue: false, expectedURL: ""},
	}

	path := "test.json"

	for _, tc := range testCase {
		t.Run(tc.method, func(t *testing.T) {
			checkKey := "s53dew1"

			fileStorageService := service.NewFileStorage(path)
			memoryRepository := repository.NewURLMemoryRepository()

			if tc.setPathValue {
				_ = memoryRepository.Save(checkKey, tc.expectedURL)
			}

			r := httptest.NewRequest(tc.method, fmt.Sprintf("/%s", checkKey), nil)
			r.SetPathValue("id", checkKey)

			w := httptest.NewRecorder()

			s := settings.NewSettings(
				settings.ServerSettings{ServerRunAddress: "", ServerShortenerAddress: ""},
			)

			shortenerService := service.NewShortenerService(memoryRepository, fileStorageService)
			URLService := service.NewURLService(s)

			e := NewURLHandler(shortenerService, URLService)

			e.RedirectURLHandler(w, r)

			if tc.expectedURL != "" {
				url := w.Header().Get("Location")
				assert.Equal(t, tc.expectedURL, url, "URL в Header Location не совпадает с ожидаемым")
			}

			assert.Equal(t, tc.expectedCode, w.Code, "Код ответа не совпадает с ожидаемым")
		})
	}
}

func TestShortenJsonHandler(t *testing.T) {
	testCase := []struct {
		method       string
		body         string
		expectedCode int
	}{
		{method: http.MethodGet, expectedCode: http.StatusMethodNotAllowed},
		{method: http.MethodPost, body: `{"url": "https://practicum.yandex.ru"}`, expectedCode: http.StatusCreated},
	}

	path := "test.json"

	for _, tc := range testCase {
		t.Run(tc.method, func(t *testing.T) {
			r := httptest.NewRequest(tc.method, "/", strings.NewReader(tc.body))
			w := httptest.NewRecorder()

			s := settings.NewSettings(
				settings.ServerSettings{ServerRunAddress: "", ServerShortenerAddress: ""},
			)

			fileStorageService := service.NewFileStorage(path)
			memoryRepository := repository.NewURLMemoryRepository()
			shortenerService := service.NewShortenerService(memoryRepository, fileStorageService)
			URLService := service.NewURLService(s)

			e := NewURLHandler(shortenerService, URLService)

			e.ShortenJSONHandler(w, r)

			assert.Equal(t, tc.expectedCode, w.Code, "Код ответа не совпадает с ожидаемым")
		})
	}
}
