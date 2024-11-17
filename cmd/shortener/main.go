package main

import (
	"github.com/Gerfey/shortener/internal/app/database"
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
			DefaultFilePath:        flags.FlagDefaultFilePath,
			DefaultDatabaseDSN:     flags.FlagDefaultDatabaseDSN,
		},
	)

	db, err := database.NewDatabase(flags.FlagDefaultDatabaseDSN)
	if err != nil {
		log.Fatalf("Ошибка инициализации клиента БД: %v", err)
	}

	application, err := app.NewShortenerApp(configApplication, db)
	if err != nil {
		log.Fatal(err)
	}

	application.Run()
}
