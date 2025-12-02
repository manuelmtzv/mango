package main

import (
	"github.com/manuelmtzv/mangocatnotes-api/internal/env"
)

type Config struct {
	Port            string
	DatabaseURL     string
	DBName          string
	JWTSecret       string
	DefaultTagLimit int64
}

func LoadConfig() Config {
	env.Load()

	return Config{
		Port:            env.GetString("PORT", "3000"),
		DatabaseURL:     env.GetString("DATABASE_URL", "mongodb://localhost:27017"),
		DBName:          env.GetString("DB_NAME", "mangocatnotes"),
		JWTSecret:       env.GetString("JWT_SECRET", "supersecret"),
		DefaultTagLimit: env.GetInt64("DEFAULT_TAG_LIMIT", 100),
	}
}
