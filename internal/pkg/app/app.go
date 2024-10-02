package app

import (
	"github.com/Gerfey/shortener/internal/app/endpoint"
	"github.com/Gerfey/shortener/internal/app/generator"
	"github.com/Gerfey/shortener/internal/app/middleware"
	"log"
	"net/http"
)

type App struct {
	e   *endpoint.Endpoint
	mux *http.ServeMux
}

func NewApp() (*App, error) {
	a := &App{}

	g := generator.NewGenerator()

	a.e = endpoint.NewEndpoint(g)

	a.mux = http.NewServeMux()

	a.mux.HandleFunc("/", a.e.ShortenUrlHandler)
	a.mux.HandleFunc("/{id}", a.e.RedirectUrlHandler)

	return a, nil
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
