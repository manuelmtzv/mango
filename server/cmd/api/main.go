package main

import (
	"log"

	"github.com/manuelmtzv/mangocatnotes-api/internal/db"
	"github.com/manuelmtzv/mangocatnotes-api/internal/store"
)

func main() {
	cfg := LoadConfig()

	database, err := db.New(cfg.DatabaseURL, cfg.DBName)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer database.Close()

	app := &application{
		config: cfg,
		store:  store.NewStorage(database.DB),
	}

	if err := app.serve(); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
