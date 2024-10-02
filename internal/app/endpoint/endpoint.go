package endpoint

import (
	"fmt"
	"github.com/Gerfey/shortener/internal/app/generator"
	"net/http"
)

type Endpoint struct {
	g generator.Generator
}

func NewEndpoint(g generator.Generator) *Endpoint {
	return &Endpoint{
		g: g,
	}
}

func (e *Endpoint) ShortenUrlHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	id := e.g.RandomString(8)

	shortURL := fmt.Sprintf("http://%s/%s", r.Host, id)

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(shortURL))
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

	w.Header().Set("Location", "https://practicum.yandex.ru/")
	w.WriteHeader(http.StatusTemporaryRedirect)
}
