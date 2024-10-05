package main

import (
	"github.com/Gerfey/shortener/internal/app/config"
	"github.com/Gerfey/shortener/internal/pkg/app"
	"log"
)

func main() {
	parseFlags()

	configApplication := config.NewConfig(
		config.ServerConfig{
			RunServerAddress:       flagRunServerAddress,
			ShortenerServerAddress: flagShortenerServerAddress,
		},
	)

	application, err := app.NewApp(configApplication)
	if err != nil {
		log.Fatal(err)
	}

	application.Run()
}
