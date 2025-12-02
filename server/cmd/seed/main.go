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

	databaseURL := env.GetString("DATABASE_URL", "mongodb://localhost:27017")
	dbName := env.GetString("DB_NAME", "mangocatnotes")

	database, err := db.New(databaseURL, dbName)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer database.Close()

	storage := store.NewStorage(database.DB)
	seeder := seed.NewSeeder(storage)

	ctx := context.Background()
	if err := seeder.Run(ctx); err != nil {
		log.Fatalf("Seeding failed: %v", err)
	}
}
