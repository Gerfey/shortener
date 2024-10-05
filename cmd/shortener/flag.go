package main

import "flag"

var (
	flagRunServerAddress       string
	flagShortenerServerAddress string
)

func parseFlags() {
	flag.StringVar(&flagRunServerAddress, "a", ":8080", "Run server address and port")
	flag.StringVar(&flagShortenerServerAddress, "b", "localhost:8080", "Run server address and port")

	flag.Parse()
}
