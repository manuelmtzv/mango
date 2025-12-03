package main

import (
	"log"

	"github.com/manuelmtzv/mangocatnotes-api/internal/config"
	"github.com/manuelmtzv/mangocatnotes-api/internal/db"
	"github.com/manuelmtzv/mangocatnotes-api/internal/server"
	"github.com/manuelmtzv/mangocatnotes-api/internal/store"
)

func main() {
	cfg := config.LoadConfig()

	database, err := db.New(cfg.DBAddr, cfg.MaxOpenConns, cfg.MaxIdleConns, cfg.MaxIdleTime)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer database.Close()

	s := server.New(cfg, store.NewStorage(database.Pool))

	if err := s.Start(); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
