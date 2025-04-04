package main

import (
	"github.com/Gerfey/shortener/internal/app/settings"
	"github.com/Gerfey/shortener/internal/app/strategy"
	"github.com/Gerfey/shortener/internal/models"
	"github.com/Gerfey/shortener/internal/pkg/app"
	"github.com/sirupsen/logrus"
	"os"
	"os/signal"
	"syscall"
)

var (
	testMode   = false
	testDoneCh = make(chan struct{})
)

func main() {
	flags := parseFlags(os.Args[1:])

	appSettings := settings.NewSettings(
		settings.ServerSettings{
			ServerRunAddress:       flags.FlagServerRunAddress,
			ServerShortenerAddress: flags.FlagServerShortenerAddress,
			DefaultFilePath:        flags.FlagDefaultFilePath,
			DefaultDatabaseDSN:     flags.FlagDefaultDatabaseDSN,
		})

	var storageStrategy models.StorageStrategy
	if appSettings.Server.DefaultDatabaseDSN != "" {
		storageStrategy = strategy.NewPostgresStrategy(appSettings.Server.DefaultDatabaseDSN)
	} else if flags.FlagDefaultFilePath != "" {
		storageStrategy = strategy.NewFileStrategy(appSettings.Server.DefaultFilePath)
	} else {
		storageStrategy = strategy.NewMemoryStrategy()
	}

	application, err := app.NewShortenerApp(appSettings, storageStrategy)
	if err != nil {
		logrus.Fatal(err)
	}

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	application.Run()

	select {
	case <-sigChan:
		logrus.Info("Получен сигнал завершения")
	case <-testDoneCh:
		if testMode {
			logrus.Info("Завершение в тестовом режиме")
			return
		}
	}
}
