package endpoint

import (
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/Gerfey/shortener/internal/app/generator"
	"github.com/Gerfey/shortener/internal/app/settings"
	"github.com/Gerfey/shortener/internal/app/store"
)

type Endpoint struct {
	generator generator.Generated
	storage   store.Stored
	settings  *settings.Settings
}

func NewEndpoint(generator generator.Generated, storage store.Stored, settings *settings.Settings) *Endpoint {
	return &Endpoint{
		generator: generator,
		storage:   storage,
		settings:  settings,
	}
}

func (e *Endpoint) ShortenURLHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	id := e.generator.RandomString(8)

	bodyBytes, _ := io.ReadAll(r.Body)

	e.storage.Set(id, string(bodyBytes))

	urlFormat, err := formatURL(e.settings.ShortenerServerAddress())
	if err != nil {
		panic("shortener server address " + urlFormat)
	}

	redirectURL := fmt.Sprintf("%v/%v", urlFormat, id)

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusCreated)
	_, err = w.Write([]byte(redirectURL))
	if err != nil {
		return
	}
}

func (e *Endpoint) RedirectURLHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	id := r.PathValue("id")
	if len(id) == 0 {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	redirectURL, exists := e.storage.Get(id)
	if !exists {
		http.Error(w, "Invalid request", http.StatusBadRequest)
	}

	w.Header().Set("Location", redirectURL)
	w.WriteHeader(http.StatusTemporaryRedirect)
}

func formatURL(URL string) (string, error) {
	urlParsed, err := url.Parse(URL)

	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%v", urlParsed.String()), err
}
