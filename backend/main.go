package main

import (
	"log"
	"transfer-tracker/config"
	"transfer-tracker/database"
	"transfer-tracker/routes"
)

func main() {
	cfg := config.Load()
	database.Init(cfg)

	r := routes.Setup()

	log.Printf("Server starting on port %s...", cfg.ServerPort)
	if err := r.Run(":" + cfg.ServerPort); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
