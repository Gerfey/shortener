package endpoint

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Gerfey/shortener/internal/app/generator"
	"github.com/Gerfey/shortener/internal/app/settings"
	"github.com/Gerfey/shortener/internal/app/store"
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

	for _, tc := range testCase {
		t.Run(tc.method, func(t *testing.T) {
			r := httptest.NewRequest(tc.method, "/", nil)
			w := httptest.NewRecorder()

			c := settings.NewSettings(
				settings.ServerSettings{ServerRunAddress: "", ServerShortenerAddress: ""},
			)
			g := generator.NewGenerator()
			s := store.NewStore()

			e := NewEndpoint(g, s, c)

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

	for _, tc := range testCase {
		t.Run(tc.method, func(t *testing.T) {
			checkKey := "s53dew1"

			s := store.NewStore()

			if tc.setPathValue {
				s.Set(checkKey, tc.expectedURL)
			}

			r := httptest.NewRequest(tc.method, fmt.Sprintf("/%s", checkKey), nil)
			r.SetPathValue("id", checkKey)

			w := httptest.NewRecorder()

			c := settings.NewSettings(
				settings.ServerSettings{ServerRunAddress: "", ServerShortenerAddress: ""},
			)
			g := generator.NewGenerator()

			e := NewEndpoint(g, s, c)

			e.RedirectURLHandler(w, r)

			if tc.expectedURL != "" {
				url := w.Header().Get("Location")
				assert.Equal(t, tc.expectedURL, url, "URL в Header Location не совпадает с ожидаемым")
			}

			assert.Equal(t, tc.expectedCode, w.Code, "Код ответа не совпадает с ожидаемым")
		})
	}
}
