package main

import (
	"log"

	"github.com/VetKA-org/vetka/internal/app"
	"github.com/VetKA-org/vetka/internal/config"
)

func main() {
	// Configuration
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalf("Config error: %s", err)
	}

	// Run
	app.Run(cfg)
}
