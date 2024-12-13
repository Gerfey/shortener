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

func parseFlags(args []string) Flags {
	var flagServerRunAddress, flagServerShortenerAddress, flagDefaultFilePath, flagDefaultDatabaseDSN string

	fs := flag.NewFlagSet("shortener", flag.ExitOnError)

	fs.StringVar(&flagServerRunAddress, "a", ":8080", "Run server address and port")
	fs.StringVar(&flagServerShortenerAddress, "b", "http://localhost:8080", "Run server address and port")
	fs.StringVar(&flagDefaultFilePath, "f", "", "Path to the file where URLs are stored")
	fs.StringVar(&flagDefaultDatabaseDSN, "d", "", "Database connection DSN")

	err := fs.Parse(args)
	if err != nil {
		return Flags{}
	}

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
