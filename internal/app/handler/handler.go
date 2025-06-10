package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/Gerfey/shortener/internal/app/service"
	"github.com/Gerfey/shortener/internal/app/settings"
	"github.com/Gerfey/shortener/internal/app/usecase"
	"github.com/Gerfey/shortener/internal/models"
	chi "github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

// Константы для работы с куками
const (
	UserIDCookieName = "user_id"
)

// URLHandler обрабатывает HTTP-запросы для сервиса сокращения URL
type URLHandler struct {
	shortenUseCase   *usecase.ShortenUseCase
	userURLsUseCase  *usecase.UserURLsUseCase
	statsUseCase     *usecase.StatsUseCase
	redirectUseCase  *usecase.RedirectUseCase
	pingUseCase      *usecase.PingUseCase
}

// NewURLHandler создает новый обработчик URL
func NewURLHandler(shortener *service.ShortenerService, s *settings.Settings, r models.Repository) *URLHandler {
	return &URLHandler{
		shortenUseCase:   usecase.NewShortenUseCase(shortener, s),
		userURLsUseCase:  usecase.NewUserURLsUseCase(r, s),
		statsUseCase:     usecase.NewStatsUseCase(r, s),
		redirectUseCase:  usecase.NewRedirectUseCase(r),
		pingUseCase:      usecase.NewPingUseCase(r),
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
	
	urls, err := h.userURLsUseCase.GetUserURLs(r.Context(), userID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if len(urls) == 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

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
	if originalURL == "" {
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

	result, err := h.shortenUseCase.ShortenURL(r.Context(), originalURL, cookie.Value)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	if result.AlreadyExists {
		w.WriteHeader(http.StatusConflict)
	} else {
		w.WriteHeader(http.StatusCreated)
	}

	if _, err := w.Write([]byte(result.FullShortURL)); err != nil {
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

	result, err := h.redirectUseCase.GetOriginalURL(r.Context(), id)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	if result.IsDeleted {
		w.WriteHeader(http.StatusGone)
		return
	}

	w.Header().Set("Location", result.OriginalURL)
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

	result, err := h.shortenUseCase.ShortenURL(r.Context(), request.URL, cookie.Value)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	response := struct {
		Result string `json:"result"`
	}{
		Result: result.FullShortURL,
	}

	w.Header().Set("Content-Type", "application/json")
	if result.AlreadyExists {
		w.WriteHeader(http.StatusConflict)
	} else {
		w.WriteHeader(http.StatusCreated)
	}
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Ошибка кодирования JSON", http.StatusInternalServerError)
		return
	}
}

// ShortenBatchHandler обрабатывает запросы для пакетного сокращения URL
func (h *URLHandler) ShortenBatchHandler(w http.ResponseWriter, r *http.Request) {
	var requestItems []models.BatchRequestItem

	if err := json.NewDecoder(r.Body).Decode(&requestItems); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	defer func() {
		if closeErr := r.Body.Close(); closeErr != nil {
			fmt.Printf("error closing request body: %v\n", closeErr)
		}
	}()

	if len(requestItems) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	for _, item := range requestItems {
		if item.OriginalURL == "" {
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

	results, err := h.shortenUseCase.ShortenBatch(r.Context(), requestItems, cookie.Value)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if encodeErr := json.NewEncoder(w).Encode(results); encodeErr != nil {
		fmt.Printf("error encoding response: %v\n", encodeErr)
	}
}

// PingHandler проверяет доступность хранилища данных
func (h *URLHandler) PingHandler(w http.ResponseWriter, r *http.Request) {
	err := h.pingUseCase.Ping(r.Context())
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

	result, err := h.shortenUseCase.ShortenURL(r.Context(), string(bodySaveURL), userID)
	if err != nil {
		http.Error(w, "Failed to shorten URL", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	if result.AlreadyExists {
		w.WriteHeader(http.StatusConflict)
	} else {
		w.WriteHeader(http.StatusCreated)
	}

	_, err = w.Write([]byte(result.FullShortURL))
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
		done <- h.userURLsUseCase.DeleteUserURLs(r.Context(), cookie.Value, shortURLs)
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
