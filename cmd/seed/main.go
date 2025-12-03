package main

import (
	"context"
	"log"

	"github.com/manuelmtzv/mangocatnotes-api/internal/db"
	"github.com/manuelmtzv/mangocatnotes-api/internal/db/seed"
	"github.com/manuelmtzv/mangocatnotes-api/internal/env"
	"github.com/manuelmtzv/mangocatnotes-api/internal/store"
)

func main() {
	env.Load()

	addr := env.GetRequired("DB_ADDR")
	maxOpenConns := env.GetInt("DB_MAX_OPEN_CONNS", 30)
	maxIdleConns := env.GetInt("DB_MAX_IDLE_CONNS", 30)
	maxIdleTime := env.GetString("DB_MAX_IDLE_TIME", "15m")

	database, err := db.New(addr, maxOpenConns, maxIdleConns, maxIdleTime)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer database.Close()

	storage := store.NewStorage(database.Pool)
	seeder := seed.NewSeeder(storage)

	ctx := context.Background()
	if err := seeder.Run(ctx); err != nil {
		log.Fatalf("Seeding failed: %v", err)
	}
}
