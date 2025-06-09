package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/Gerfey/shortener/internal/app/service"
	"github.com/Gerfey/shortener/internal/app/settings"
	"github.com/Gerfey/shortener/internal/models"
	chi "github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	pgx "github.com/jackc/pgx/v5"
)

// Константы для работы с куками
const (
	UserIDCookieName = "user_id"
)

// URLHandler обрабатывает HTTP-запросы для сервиса сокращения URL
type URLHandler struct {
	shortener  *service.ShortenerService
	url        *service.URLService
	settings   *settings.Settings
	repository models.Repository
}

// NewURLHandler создает новый обработчик URL
func NewURLHandler(shortener *service.ShortenerService, url *service.URLService, s *settings.Settings, r models.Repository) *URLHandler {
	return &URLHandler{
		shortener:  shortener,
		url:        url,
		settings:   s,
		repository: r,
	}
}

// GetUserURLsHandler обрабатывает запросы для получения списка URL пользователя
func (h *URLHandler) GetUserURLsHandler(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie(UserIDCookieName)
	if err != nil || cookie == nil {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	userID := cookie.Value
	urls, err := h.repository.GetUserURLs(r.Context(), userID)
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

// ShortenHandler обрабатывает запросы для сокращения URL
func (h *URLHandler) ShortenHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer func() {
		if closeErr := r.Body.Close(); closeErr != nil {
			fmt.Printf("error closing request body: %v\n", closeErr)
		}
	}()

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

	shortURL, err := h.shortener.ShortenID(r.Context(), originalURL, cookie.Value)
	if err != nil {
		if err == models.ErrURLExists {
			w.Header().Set("Content-Type", "text/plain")
			w.WriteHeader(http.StatusConflict)
			if _, err := w.Write([]byte(h.settings.ShortenerServerAddress() + "/" + shortURL)); err != nil {
				fmt.Printf("error writing response: %v\n", err)
			}
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusCreated)
	if _, err := w.Write([]byte(h.settings.ShortenerServerAddress() + "/" + shortURL)); err != nil {
		fmt.Printf("error writing response: %v\n", err)
	}
}

// RedirectURLHandler обрабатывает запросы для перенаправления по сокращенному URL
func (h *URLHandler) RedirectURLHandler(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if r.Body != nil {
			if err := r.Body.Close(); err != nil {
				fmt.Printf("error closing request body: %v\n", err)
			}
		}
	}()

	id := chi.URLParam(r, "id")
	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	originalURL, found, isDeleted := h.repository.Find(r.Context(), id)
	if !found {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	if isDeleted {
		w.WriteHeader(http.StatusGone)
		return
	}

	w.Header().Set("Location", originalURL)
	w.WriteHeader(http.StatusTemporaryRedirect)
}

// ShortenJSONHandler обрабатывает запросы для сокращения URL в формате JSON
func (h *URLHandler) ShortenJSONHandler(w http.ResponseWriter, r *http.Request) {
	var request struct {
		URL string `json:"url"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	defer func() {
		if closeErr := r.Body.Close(); closeErr != nil {
			fmt.Printf("error closing request body: %v\n", closeErr)
		}
	}()

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

	shortURL, err := h.shortener.ShortenID(r.Context(), request.URL, cookie.Value)
	if err != nil {
		if err == models.ErrURLExists {
			response := struct {
				Result string `json:"result"`
			}{
				Result: h.settings.ShortenerServerAddress() + "/" + shortURL,
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusConflict)
			if encodeErr := json.NewEncoder(w).Encode(response); encodeErr != nil {
				fmt.Printf("error encoding response: %v\n", encodeErr)
			}
			return
		}
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
	if encodeErr := json.NewEncoder(w).Encode(response); encodeErr != nil {
		fmt.Printf("error encoding response: %v\n", encodeErr)
	}
}

// ShortenBatchHandler обрабатывает запросы для пакетного сокращения URL
func (h *URLHandler) ShortenBatchHandler(w http.ResponseWriter, r *http.Request) {
	var request []struct {
		CorrelationID string `json:"correlation_id"`
		OriginalURL   string `json:"original_url"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	defer func() {
		if closeErr := r.Body.Close(); closeErr != nil {
			fmt.Printf("error closing request body: %v\n", closeErr)
		}
	}()

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
		shortURL, err := h.shortener.ShortenID(r.Context(), item.OriginalURL, cookie.Value)
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

	if err := h.repository.SaveBatch(r.Context(), urls, cookie.Value); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if encodeErr := json.NewEncoder(w).Encode(response); encodeErr != nil {
		fmt.Printf("error encoding response: %v\n", encodeErr)
	}
}

// PingHandler проверяет доступность хранилища данных
func (h *URLHandler) PingHandler(w http.ResponseWriter, r *http.Request) {
	err := h.repository.Ping(r.Context())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

// ShortenURLHandler обрабатывает запросы для сокращения URL
func (h *URLHandler) ShortenURLHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	bodySaveURL, _ := io.ReadAll(r.Body)
	defer func() {
		if r.Body != nil {
			if err := r.Body.Close(); err != nil {
				fmt.Printf("error closing request body: %v\n", err)
			}
		}
	}()

	cookie, err := r.Cookie(UserIDCookieName)
	if err != nil {
		http.Error(w, "User ID cookie not found", http.StatusInternalServerError)
		return
	}
	userID := cookie.Value

	shortenID, err := h.shortener.ShortenID(r.Context(), string(bodySaveURL), userID)

	if err != nil {
		if err.Error() == "no rows in result set" {
			shortenerURL, urlErr := h.url.ShortenerURL(shortenID)
			if urlErr != nil {
				return
			}

			w.Header().Set("Content-Type", "text/plain")
			w.WriteHeader(http.StatusConflict)
			_, writeErr := w.Write([]byte(shortenerURL))
			if writeErr != nil {
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

// DeleteUserURLsHandler обрабатывает запросы для удаления URL пользователя
func (h *URLHandler) DeleteUserURLsHandler(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie(UserIDCookieName)
	if err != nil || cookie == nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	defer func() {
		if closeErr := r.Body.Close(); closeErr != nil {
			fmt.Printf("error closing request body: %v\n", closeErr)
		}
	}()

	var shortURLs []string
	if err := json.Unmarshal(body, &shortURLs); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	done := make(chan error, 1)
	go func() {
		done <- h.repository.DeleteUserURLsBatch(r.Context(), shortURLs, cookie.Value)
	}()

	select {
	case err := <-done:
		if err != nil {
			fmt.Printf("Error deleting URLs: %v\n", err)
		}
	case <-r.Context().Done():
		fmt.Printf("Request context cancelled while deleting URLs\n")
	}

	w.WriteHeader(http.StatusAccepted)
}
