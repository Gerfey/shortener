package endpoint

import (
	"fmt"
	"github.com/Gerfey/shortener/internal/app/config"
	"github.com/Gerfey/shortener/internal/app/generator"
	"github.com/Gerfey/shortener/internal/app/store"
	"io"
	"net/http"
	"net/url"
)

type Endpoint struct {
	g generator.Generator
	s *store.Store
	c *config.Config
}

func NewEndpoint(g generator.Generator, s *store.Store, c *config.Config) *Endpoint {
	return &Endpoint{
		g: g,
		s: s,
		c: c,
	}
}

func (e *Endpoint) ShortenURLHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	id := e.g.RandomString(8)

	bodyBytes, _ := io.ReadAll(r.Body)

	e.s.Set(id, string(bodyBytes))

	urlFormat, err := formatURL(e.c.GetShortenerServerAddress())
	if err != nil {
		panic("shortener server address " + urlFormat)
	}

	redirectURL := fmt.Sprintf("%v/%v", urlFormat, id)

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(redirectURL))
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

	redirectURL, exists := e.s.Get(id)
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

	return urlParsed.String(), err
}
