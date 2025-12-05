package main

import (
	"context"
	"log"

	"github.com/manuelmtzv/mangocatnotes-api/internal/config"
	"github.com/manuelmtzv/mangocatnotes-api/internal/db"
	"github.com/manuelmtzv/mangocatnotes-api/internal/kvstore"
	"github.com/manuelmtzv/mangocatnotes-api/internal/server"
	"github.com/manuelmtzv/mangocatnotes-api/internal/session"
	"github.com/manuelmtzv/mangocatnotes-api/internal/store"
	"github.com/redis/go-redis/v9"
)

func main() {
	cfg := config.LoadConfig()

	database, err := db.New(cfg.DBAddr, cfg.MaxOpenConns, cfg.MaxIdleConns, cfg.MaxIdleTime)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer database.Close()

	redisClient := redis.NewClient(&redis.Options{
		Addr:     cfg.RedisAddr,
		Password: cfg.RedisPassword,
		DB:       0,
	})

	if err := redisClient.Ping(context.Background()).Err(); err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}

	cache := kvstore.NewRedisStore(redisClient)
	defer redisClient.Close()

	session := session.NewSessionManager(cache)

	s := server.New(cfg, store.NewStorage(database.Pool), session)

	if err := s.Start(); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
