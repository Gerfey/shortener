package main

import (
	"fmt"
	"os"

	"github.com/Gerfey/shortener/internal/app/settings"
	"github.com/Gerfey/shortener/internal/app/strategy"
	"github.com/Gerfey/shortener/internal/models"
	"github.com/Gerfey/shortener/internal/pkg/app"
	"github.com/sirupsen/logrus"
)

var (
	testMode     = false
	testDoneCh   = make(chan struct{})
	buildVersion string
	buildDate    string
	buildCommit  string
)

func main() {
	version := buildVersion
	if version == "" {
		version = "N/A"
	}

	date := buildDate
	if date == "" {
		date = "N/A"
	}

	commit := buildCommit
	if commit == "" {
		commit = "N/A"
	}

	fmt.Printf("Build version: %s\n", version)
	fmt.Printf("Build date: %s\n", date)
	fmt.Printf("Build commit: %s\n", commit)

	flags := parseFlags(os.Args[1:])

	appSettings := settings.NewSettings(
		settings.ServerSettings{
			ServerRunAddress:       flags.FlagServerRunAddress,
			ServerShortenerAddress: flags.FlagServerShortenerAddress,
			DefaultFilePath:        flags.FlagDefaultFilePath,
			DefaultDatabaseDSN:     flags.FlagDefaultDatabaseDSN,
			EnableHTTPS:            flags.FlagEnableHTTPS,
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

	appDone := make(chan struct{})

	go func() {
		application.Run()
		appDone <- struct{}{}
	}()

	select {
	case <-appDone:
		logrus.Info("Приложение завершило работу")
	case <-testDoneCh:
		if testMode {
			logrus.Info("Завершение в тестовом режиме")
			return
		}
	}
}
