package app

import (
	"log"
	"net/http"

	"github.com/Gerfey/shortener/internal/app/endpoint"
	"github.com/Gerfey/shortener/internal/app/generator"
	"github.com/Gerfey/shortener/internal/app/settings"
	"github.com/Gerfey/shortener/internal/app/store"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type ShortenerApp struct {
	settings *settings.Settings
	endpoint *endpoint.Endpoint
	router   *chi.Mux
}

func NewShortenerApp(c *settings.Settings) (*ShortenerApp, error) {
	application := &ShortenerApp{}

	application.settings = c

	g := generator.NewGenerator()
	s := store.NewStore()

	application.endpoint = endpoint.NewEndpoint(g, s, c)

	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	application.router = r

	return application, nil
}

func (a *ShortenerApp) Run() {
	log.Println("Starting server...")

	a.router.Route("/", func(r chi.Router) {
		r.Post("/", a.endpoint.ShortenURLHandler)
		r.Get("/{id}", a.endpoint.RedirectURLHandler)
	})

	err := http.ListenAndServe(a.settings.ServerAddress(), a.router)
	if err != nil {
		log.Fatal(err)
	}
}
