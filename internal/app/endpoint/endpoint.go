package endpoint

import (
	"fmt"
	"github.com/Gerfey/shortener/internal/app/generator"
	"github.com/Gerfey/shortener/internal/app/store"
	"io"
	"net/http"
)

type Endpoint struct {
	g generator.Generator
	s *store.Store
}

func NewEndpoint(g generator.Generator, s *store.Store) *Endpoint {
	return &Endpoint{
		g: g,
		s: s,
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

	shortURL := getCorrectURL(r) + id

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(shortURL))
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

func getCorrectURL(r *http.Request) string {
	scheme := "http"
	if r.TLS != nil {
		scheme = "https"
	}

	return fmt.Sprintf("%v://%v%v", scheme, r.Host, r.RequestURI)
}
