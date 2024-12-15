package handler

import (
	"encoding/json"
	"fmt"
	"github.com/Gerfey/shortener/internal/app/service"
	"github.com/Gerfey/shortener/internal/app/settings"
	"github.com/Gerfey/shortener/internal/models"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/jackc/pgx"
	"io"
	"net/http"
)

const (
	UserIDCookieName = "user_id"
)

type URLHandler struct {
	shortener  *service.ShortenerService
	url        *service.URLService
	settings   *settings.Settings
	repository models.Repository
}

func NewURLHandler(shortener *service.ShortenerService, url *service.URLService, s *settings.Settings, r models.Repository) *URLHandler {
	return &URLHandler{
		shortener:  shortener,
		url:        url,
		settings:   s,
		repository: r,
	}
}

func (h *URLHandler) GetUserURLsHandler(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie(UserIDCookieName)
	if err != nil || cookie == nil {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	userID := cookie.Value
	urls, err := h.repository.GetUserURLs(userID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if len(urls) == 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	baseURL := h.settings.ShortenerServerAddress()
	for i := range urls {
		urls[i].ShortURL = baseURL + "/" + urls[i].ShortURL
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(urls); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (h *URLHandler) ShortenHandler(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	if len(body) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	originalURL := string(body)
	if !h.url.IsValidURL(originalURL) {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	cookie, err := r.Cookie(UserIDCookieName)
	if err != nil || cookie == nil {
		userID := uuid.New().String()
		cookie = &http.Cookie{
			Name:     UserIDCookieName,
			Value:    userID,
			Path:     "/",
			HttpOnly: true,
		}
		http.SetCookie(w, cookie)
	}

	shortURL, err := h.shortener.ShortenID(originalURL, cookie.Value)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusCreated)
	if _, err := w.Write([]byte(h.settings.ShortenerServerAddress() + "/" + shortURL)); err != nil {
		fmt.Printf("error writing response: %v\n", err)
	}
}

func (h *URLHandler) RedirectURLHandler(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	originalURL, found := h.repository.Find(id)
	if !found {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.Header().Set("Location", originalURL)
	w.WriteHeader(http.StatusTemporaryRedirect)
}

func (h *URLHandler) ShortenJSONHandler(w http.ResponseWriter, r *http.Request) {
	var request struct {
		URL string `json:"url"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	if request.URL == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	cookie, err := r.Cookie(UserIDCookieName)
	if err != nil || cookie == nil {
		userID := uuid.New().String()
		cookie = &http.Cookie{
			Name:     UserIDCookieName,
			Value:    userID,
			Path:     "/",
			HttpOnly: true,
		}
		http.SetCookie(w, cookie)
	}

	shortURL, err := h.shortener.ShortenID(request.URL, cookie.Value)
	if err != nil {
		if err == pgx.ErrNoRows {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	response := struct {
		Result string `json:"result"`
	}{
		Result: h.settings.ShortenerServerAddress() + "/" + shortURL,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		fmt.Printf("error encoding response: %v\n", err)
	}
}

func (h *URLHandler) ShortenBatchHandler(w http.ResponseWriter, r *http.Request) {
	var request []struct {
		CorrelationID string `json:"correlation_id"`
		OriginalURL   string `json:"original_url"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	if len(request) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	for _, item := range request {
		if !h.url.IsValidURL(item.OriginalURL) {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}

	cookie, err := r.Cookie(UserIDCookieName)
	if err != nil || cookie == nil {
		userID := uuid.New().String()
		cookie = &http.Cookie{
			Name:     UserIDCookieName,
			Value:    userID,
			Path:     "/",
			HttpOnly: true,
		}
		http.SetCookie(w, cookie)
	}

	urls := make(map[string]string)
	response := make([]struct {
		CorrelationID string `json:"correlation_id"`
		ShortURL      string `json:"short_url"`
	}, len(request))

	for i, item := range request {
		shortURL, err := h.shortener.ShortenID(item.OriginalURL, cookie.Value)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		urls[shortURL] = item.OriginalURL
		response[i] = struct {
			CorrelationID string `json:"correlation_id"`
			ShortURL      string `json:"short_url"`
		}{
			CorrelationID: item.CorrelationID,
			ShortURL:      h.settings.ShortenerServerAddress() + "/" + shortURL,
		}
	}

	if err := h.repository.SaveBatch(urls, cookie.Value); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		fmt.Printf("error encoding response: %v\n", err)
	}
}

func (h *URLHandler) PingHandler(w http.ResponseWriter, r *http.Request) {
	err := h.repository.Ping()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (h *URLHandler) ShortenURLHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	bodySaveURL, _ := io.ReadAll(r.Body)

	cookie, err := r.Cookie(UserIDCookieName)
	if err != nil {
		http.Error(w, "User ID cookie not found", http.StatusInternalServerError)
		return
	}
	userID := cookie.Value

	shortenID, err := h.shortener.ShortenID(string(bodySaveURL), userID)

	if err != nil {
		if err.Error() == "no rows in result set" {
			shortenerURL, err := h.url.ShortenerURL(shortenID)
			if err != nil {
				return
			}

			w.Header().Set("Content-Type", "text/plain")
			w.WriteHeader(http.StatusConflict)
			_, err = w.Write([]byte(shortenerURL))
			if err != nil {
				return
			}
			return
		}
		http.Error(w, "Failed to shorten URL", http.StatusInternalServerError)
		return
	}

	shortenerURL, err := h.url.ShortenerURL(shortenID)
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
