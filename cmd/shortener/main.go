package main

import (
	"github.com/Gerfey/shortener/internal/pkg/app"
	"log"
)

func main() {
	application, err := app.NewApp()
	if err != nil {
		log.Fatal(err)
	}

	application.Run()
}
