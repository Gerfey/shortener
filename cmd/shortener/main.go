package main

import (
	"github.com/Gerfey/shortener/internal/app/strategy"
	"github.com/Gerfey/shortener/internal/models"
	"log"
	"os"

	"github.com/Gerfey/shortener/internal/app/settings"
	"github.com/Gerfey/shortener/internal/pkg/app"
)

func main() {
	flags := parseFlags(os.Args[1:])

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

	var storageStrategy models.StorageStrategy

	if flags.FlagDefaultDatabaseDSN != "" {
		storageStrategy = strategy.NewPostgresStrategy(flags.FlagDefaultDatabaseDSN)
	} else if flags.FlagDefaultFilePath != "" {
		storageStrategy = strategy.NewFileStrategy(flags.FlagDefaultFilePath)
	} else {
		storageStrategy = strategy.NewMemoryStrategy()
	}

	repo, err := storageStrategy.Initialize()
	if err != nil {
		return err
	}
	defer storageStrategy.Close()

	application, err := app.NewShortenerApp(configApplication, repo)
	if err != nil {
		return err
	}

	application.Run()

	return nil
}
