package app

import (
	"github.com/Gerfey/shortener/internal/app/endpoint"
	"github.com/Gerfey/shortener/internal/app/generator"
	"github.com/Gerfey/shortener/internal/app/store"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"log"
	"net/http"
)

type App struct {
	e *endpoint.Endpoint
	r *chi.Mux
}

func NewApp() (*App, error) {
	application := &App{}

	g := generator.NewGenerator()
	s := store.NewStore()

	application.e = endpoint.NewEndpoint(g, s)

	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	application.r = r

	return application, nil
}

func (a *App) Run() {
	log.Println("Starting server...")

	a.r.Route("/", func(r chi.Router) {
		r.Post("/", a.e.ShortenURLHandler)
		r.Get("/{id}", a.e.RedirectURLHandler)
	})

	err := http.ListenAndServe(":8080", a.r)
	if err != nil {
		log.Fatal(err)
	}
}
