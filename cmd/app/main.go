package main

import (
	"log"

	"github.com/cheesy008/ffbc-backend/internal/app"
)

func main() {
	if err := app.Run(); err != nil {
		log.Fatalf("run app: %v", err)
	}
}
