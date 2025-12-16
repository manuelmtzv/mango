package main

import (
	"context"

	"go.uber.org/zap"

	"github.com/manuelmtzv/mangocatnotes-api/internal/config"
	"github.com/manuelmtzv/mangocatnotes-api/internal/db"
	"github.com/manuelmtzv/mangocatnotes-api/internal/kvstore"
	"github.com/manuelmtzv/mangocatnotes-api/internal/server"
	"github.com/manuelmtzv/mangocatnotes-api/internal/session"
	"github.com/manuelmtzv/mangocatnotes-api/internal/store"
	"github.com/redis/go-redis/v9"
)

func main() {
	logger := zap.Must(zap.NewProduction()).Sugar()
	defer func() {
		_ = logger.Sync()
	}()

	cfg := config.LoadConfig()

	database, err := db.New(cfg.DBAddr, cfg.MaxOpenConns, cfg.MaxIdleConns, cfg.MaxIdleTime)
	if err != nil {
		logger.Fatalw("Failed to connect to database", "error", err)
	}
	defer database.Close()

	redisOpts, err := redis.ParseURL(cfg.RedisAddr)
	if err != nil {
		logger.Fatalw("Failed to parse Redis URL", "error", err)
	}
	redisClient := redis.NewClient(redisOpts)

	if err := redisClient.Ping(context.Background()).Err(); err != nil {
		logger.Fatalw("Failed to connect to Redis", "error", err)
	}

	cache := kvstore.NewRedisStore(redisClient)
	defer redisClient.Close()

	session := session.NewSessionManager(cache)

	s := server.New(cfg, logger, store.NewStorage(database.Pool), session)

	if err := s.Start(); err != nil {
		logger.Fatalw("Server failed", "error", err)
	}
}
