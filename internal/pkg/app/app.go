package app

import (
	"net/http"

	middleware2 "github.com/Gerfey/shortener/internal/app/middleware"
	log "github.com/sirupsen/logrus"

	"github.com/Gerfey/shortener/internal/app/handler"
	"github.com/Gerfey/shortener/internal/app/repository/memory"
	"github.com/Gerfey/shortener/internal/app/service"
	"github.com/Gerfey/shortener/internal/app/settings"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type ShortenerApp struct {
	settings *settings.Settings
	handler  *handler.URLHandler
	router   *chi.Mux
}

func NewShortenerApp(s *settings.Settings) (*ShortenerApp, error) {

	log.SetFormatter(&log.TextFormatter{
		FullTimestamp: true,
	})
	log.SetLevel(log.InfoLevel)

	application := &ShortenerApp{}

	application.settings = s

	repository := memory.NewURLMemoryRepository()
	shortenerService := service.NewShortenerService(repository)
	URLService := service.NewURLService(s)

	application.handler = handler.NewURLHandler(shortenerService, URLService)

	r := chi.NewRouter()

	r.Use(middleware2.LoggingMiddleware)
	r.Use(middleware.Recoverer)

	application.router = r

	return application, nil
}

func (a *ShortenerApp) Run() {
	log.Println("Starting server...")

	a.router.Route("/", func(r chi.Router) {
		r.Post("/", a.handler.ShortenURLHandler)
		r.Get("/{id}", a.handler.RedirectURLHandler)
	})

	err := http.ListenAndServe(a.settings.ServerAddress(), a.router)
	if err != nil {
		log.Fatal(err)
	}
}
