package main

import (
	"log"

	"github.com/Gerfey/shortener/internal/app/settings"
	"github.com/Gerfey/shortener/internal/pkg/app"
)

func main() {
	flags := parseFlags()

	configApplication := settings.NewSettings(
		settings.ServerSettings{
			ServerRunAddress:       flags.FlagServerRunAddress,
			ServerShortenerAddress: flags.FlagServerShortenerAddress,
		},
	)

	application, err := app.NewShortenerApp(configApplication)
	if err != nil {
		log.Fatal(err)
	}

	application.Run()
}
