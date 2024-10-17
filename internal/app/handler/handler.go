package handler

import (
	"io"
	"net/http"

	"github.com/Gerfey/shortener/internal/app/service"
)

type URLHandler struct {
	shortener *service.ShortenerService
	url       *service.URLService
}

func NewURLHandler(shortener *service.ShortenerService, url *service.URLService) *URLHandler {
	return &URLHandler{
		shortener: shortener,
		url:       url,
	}
}

func (e *URLHandler) ShortenURLHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	bodySaveURL, _ := io.ReadAll(r.Body)

	shortenID, err := e.shortener.ShortenID(string(bodySaveURL))
	if err != nil {
		return
	}

	shortenerURL, err := e.url.ShortenerURL(shortenID)
	if err != nil {
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusCreated)

	_, err = w.Write([]byte(shortenerURL))
	if err != nil {
		return
	}
}

func (e *URLHandler) RedirectURLHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	id := r.PathValue("id")
	if len(id) == 0 {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	redirectURL, err := e.shortener.FindURL(id)
	if err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
	}

	w.Header().Set("Location", redirectURL)
	w.WriteHeader(http.StatusTemporaryRedirect)
}
