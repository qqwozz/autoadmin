package main

import (
	"autoadmin/internal/bot"
	"log"
	"os"
)

func main() {
	apiURL := os.Getenv("API_URL")
	if apiURL == "" {
		apiURL = "http://localhost:8080"
	}

	b, err := bot.New(apiURL)
	if err != nil {
		log.Fatal(err)
	}

	b.Start()
}
