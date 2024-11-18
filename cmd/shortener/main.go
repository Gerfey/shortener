package main

import (
	"github.com/Gerfey/shortener/internal/app/database"
	"log"

	"github.com/Gerfey/shortener/internal/app/settings"
	"github.com/Gerfey/shortener/internal/pkg/app"
)

func main() {
	flags := parseFlags()

	err := run(flags)
	if err != nil {
		log.Fatalf("Ошибка: %v", err)
	}
}

func run(flags Flags) error {
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
		return err
	}

	application, err := app.NewShortenerApp(configApplication, db)
	if err != nil {
		return err
	}

	application.Run()
	return nil
}
