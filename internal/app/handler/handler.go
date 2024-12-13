package handler

import (
	"encoding/json"
	"fmt"
	"github.com/Gerfey/shortener/internal/app/repository"
	"github.com/Gerfey/shortener/internal/app/service"
	"github.com/Gerfey/shortener/internal/app/settings"
	"github.com/Gerfey/shortener/internal/models"
	"github.com/jackc/pgx"
	"io"
	"net/http"
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

func (e *URLHandler) ShortenBatchHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var batchRequest []models.BatchRequestItem
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&batchRequest); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if len(batchRequest) == 0 {
		http.Error(w, "Batch is empty", http.StatusBadRequest)
		return
	}

	urlsToSave := make(map[string]string)
	batchResponse := make([]models.BatchResponseItem, 0, len(batchRequest))

	for _, item := range batchRequest {
		shortURL, err := e.shortener.ShortenID(item.OriginalURL)
		if err != nil {
			http.Error(w, "Failed to shorten URL", http.StatusInternalServerError)
			return
		}

		urlsToSave[shortURL] = item.OriginalURL
		batchResponse = append(batchResponse, models.BatchResponseItem{
			CorrelationID: item.CorrelationID,
			ShortURL:      fmt.Sprintf("%s/%s", e.settings.ShortenerServerAddress(), shortURL),
		})
	}

	err := e.repository.SaveBatch(urlsToSave)
	if err != nil {
		http.Error(w, "Failed to save URLs", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	encoder := json.NewEncoder(w)
	if err := encoder.Encode(batchResponse); err != nil {
		http.Error(w, "Failed to write response", http.StatusInternalServerError)
		return
	}
}

func (e *URLHandler) PingHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	config, err := pgx.ParseConnectionString(e.settings.Server.DefaultDatabaseDSN)
	if err != nil {
		http.Error(w, "PostgresRepository connection failed", http.StatusInternalServerError)
		return
	}

	conn, err := pgx.Connect(config)
	if err != nil {
		http.Error(w, "PostgresRepository connection failed", http.StatusInternalServerError)
		return
	}

	_, err = repository.NewPostgresRepository(conn)
	if err != nil {
		http.Error(w, "PostgresRepository connection failed", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (e *URLHandler) ShortenURLHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	bodySaveURL, _ := io.ReadAll(r.Body)

	shortenID, err := e.shortener.ShortenID(string(bodySaveURL))

	if err != nil {
		if err.Error() == "no rows in result set" {
			shortenerURL, err := e.url.ShortenerURL(shortenID)
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

func (e *URLHandler) ShortenJSONHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req models.ShortenRequest
	dec := json.NewDecoder(r.Body)
	if err := dec.Decode(&req); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if req.URL == "" {
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	shortenID, err := e.shortener.ShortenID(req.URL)

	if err != nil {
		if err.Error() == "no rows in result set" {
			shortenerURL, err := e.url.ShortenerURL(shortenID)
			if err != nil {
				return
			}

			resp := models.ShortenResponse{Result: shortenerURL}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusConflict)
			_ = json.NewEncoder(w).Encode(resp)
			return
		}

		http.Error(w, "Failed to shorten URL", http.StatusInternalServerError)
		return
	}

	shortenerURL, err := e.url.ShortenerURL(shortenID)
	if err != nil {
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	resp := models.ShortenResponse{
		Result: shortenerURL,
	}

	enc := json.NewEncoder(w)
	if err := enc.Encode(resp); err != nil {
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
