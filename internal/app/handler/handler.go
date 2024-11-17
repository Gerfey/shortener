package handler

import (
	"encoding/json"
	"github.com/Gerfey/shortener/internal/app/database"
	"github.com/Gerfey/shortener/internal/app/service"
	"github.com/Gerfey/shortener/internal/models"
	"io"
	"net/http"
)

type URLHandler struct {
	shortener *service.ShortenerService
	url       *service.URLService
	database  *database.Database
}

func NewURLHandler(shortener *service.ShortenerService, url *service.URLService, db *database.Database) *URLHandler {
	return &URLHandler{
		shortener: shortener,
		url:       url,
		database:  db,
	}
}

func (e *URLHandler) Ping(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	connect, err := e.database.Connect()
	if err != nil {
		http.Error(w, "Database connection failed", http.StatusInternalServerError)
		return
	}
	defer connect.Close()

	err = connect.Ping(r.Context())
	if err != nil {
		http.Error(w, "Database PING failed", http.StatusInternalServerError)
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
