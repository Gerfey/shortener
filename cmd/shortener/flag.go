package main

import (
	"flag"
	"os"
)

var (
	flagRunServerAddress       string
	flagShortenerServerAddress string
)

func parseFlags() {
	flag.StringVar(&flagRunServerAddress, "a", ":8080", "Run server address and port")
	flag.StringVar(&flagShortenerServerAddress, "b", "http://localhost:8080", "Run server address and port")

	flag.Parse()

	if envRunServerAddress := os.Getenv("SERVER_ADDRESS"); envRunServerAddress != "" {
		flagRunServerAddress = envRunServerAddress
	}

	if envShortenerServerAddress := os.Getenv("BASE_URL"); envShortenerServerAddress != "" {
		flagShortenerServerAddress = envShortenerServerAddress
	}
}
