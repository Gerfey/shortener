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

func (e *Endpoint) ShortenUrlHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	id := e.g.RandomString(8)

	bodyBytes, _ := io.ReadAll(r.Body)

	e.s.Set(id, string(bodyBytes))

	shortUrl := getCorrectUrl(r) + id

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(shortUrl))
}

func (e *Endpoint) RedirectUrlHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	id := r.PathValue("id")
	if len(id) == 0 {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	redirectUrl, exists := e.s.Get(id)
	if exists == false {
		http.Error(w, "Invalid request", http.StatusBadRequest)
	}

	w.Header().Set("Location", redirectUrl)
	w.WriteHeader(http.StatusTemporaryRedirect)
}

func getCorrectUrl(r *http.Request) string {
	scheme := "http"
	if r.TLS != nil {
		scheme = "https"
	}

	return fmt.Sprintf("%v://%v%v", scheme, r.Host, r.RequestURI)
}
