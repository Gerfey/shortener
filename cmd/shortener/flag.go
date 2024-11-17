package main

import (
	"flag"
	"os"
)

type Flags struct {
	FlagServerRunAddress       string
	FlagServerShortenerAddress string
	FlagDefaultFilePath        string
	FlagDefaultDatabaseDSN     string
}

func parseFlags() Flags {
	var flagServerRunAddress, flagServerShortenerAddress, flagDefaultFilePath, flagDefaultDatabaseDSN string

	flag.StringVar(&flagServerRunAddress, "a", ":8080", "Run server address and port")
	flag.StringVar(&flagServerShortenerAddress, "b", "http://localhost:8080", "Run server address and port")
	flag.StringVar(&flagDefaultFilePath, "f", "url_store.json", "Path to the file where URLs are stored")
	flag.StringVar(&flagDefaultDatabaseDSN, "d", "host=localhost port=5432 user=shortener password=shortener dbname=shortener sslmode=disable", "Database connection DSN")

	flag.Parse()

	if envServerRunAddress := os.Getenv("SERVER_ADDRESS"); envServerRunAddress != "" {
		flagServerRunAddress = envServerRunAddress
	}

	if envServerShortenerAddress := os.Getenv("BASE_URL"); envServerShortenerAddress != "" {
		flagServerShortenerAddress = envServerShortenerAddress
	}

	if envDefaultFilePath := os.Getenv("FILE_STORAGE_PATH"); envDefaultFilePath != "" {
		flagDefaultFilePath = envDefaultFilePath
	}

	if envDefaultDatabaseDSN := os.Getenv("DATABASE_DSN"); envDefaultDatabaseDSN != "" {
		flagDefaultDatabaseDSN = envDefaultDatabaseDSN
	}

	return Flags{flagServerRunAddress, flagServerShortenerAddress, flagDefaultFilePath, flagDefaultDatabaseDSN}
}
