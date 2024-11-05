package main

import (
	"flag"
	"os"
)

type Flags struct {
	FlagServerRunAddress       string
	FlagServerShortenerAddress string
	FlagDefaultFilePath        string
}

func parseFlags() Flags {
	var flagServerRunAddress, flagServerShortenerAddress, flagDefaultFilePath string

	flag.StringVar(&flagServerRunAddress, "a", ":8080", "Run server address and port")
	flag.StringVar(&flagServerShortenerAddress, "b", "http://localhost:8080", "Run server address and port")
	flag.StringVar(&flagDefaultFilePath, "f", "url_store.json", "Path to the file where URLs are stored")

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

	return Flags{flagServerRunAddress, flagServerShortenerAddress, flagDefaultFilePath}
}
