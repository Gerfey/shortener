package main

import (
	"flag"
	"os"
)

type Flags struct {
	FlagServerRunAddress       string
	FlagServerShortenerAddress string
}

func parseFlags() Flags {
	var flagServerRunAddress, flagServerShortenerAddress string

	flag.StringVar(&flagServerRunAddress, "a", ":8080", "Run server address and port")
	flag.StringVar(&flagServerShortenerAddress, "b", "http://localhost:8080", "Run server address and port")

	flag.Parse()

	if envServerRunAddress := os.Getenv("SERVER_ADDRESS"); envServerRunAddress != "" {
		flagServerRunAddress = envServerRunAddress
	}

	if envServerShortenerAddress := os.Getenv("BASE_URL"); envServerShortenerAddress != "" {
		flagServerShortenerAddress = envServerShortenerAddress
	}

	return Flags{flagServerRunAddress, flagServerShortenerAddress}
}
