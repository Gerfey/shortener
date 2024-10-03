package app

import (
	"github.com/Gerfey/shortener/internal/app/endpoint"
	"github.com/Gerfey/shortener/internal/app/generator"
	"github.com/Gerfey/shortener/internal/app/middleware"
	"github.com/Gerfey/shortener/internal/app/store"
	"log"
	"net/http"
)

type App struct {
	e   *endpoint.Endpoint
	mux *http.ServeMux
}

func NewApp() (*App, error) {
	application := &App{}

	g := generator.NewGenerator()

	s := store.NewStore()

	application.e = endpoint.NewEndpoint(g, s)

	application.mux = http.NewServeMux()

	application.mux.HandleFunc("/", application.e.ShortenUrlHandler)
	application.mux.HandleFunc("/{id}", application.e.RedirectUrlHandler)

	return application, nil
}

func (a *App) Run() {
	log.Println("Starting server...")

	handler := middleware.Logging(a.mux)
	handler = middleware.PanicRecovery(handler)

	err := http.ListenAndServe(":8080", handler)
	if err != nil {
		log.Fatal(err)
	}
}
