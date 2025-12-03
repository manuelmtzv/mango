package config

import (
	"github.com/manuelmtzv/mangocatnotes-api/internal/env"
)

type Config struct {
	Port            string
	DBAddr          string
	MaxOpenConns    int
	MaxIdleConns    int
	MaxIdleTime     string
	JWTSecret       string
	DefaultTagLimit int64
}

func LoadConfig() *Config {
	env.Load()

	return &Config{
		Port:            env.GetString("PORT", "3000"),
		DBAddr:          env.GetRequired("DB_ADDR"),
		MaxOpenConns:    env.GetInt("DB_MAX_OPEN_CONNS", 30),
		MaxIdleConns:    env.GetInt("DB_MAX_IDLE_CONNS", 30),
		MaxIdleTime:     env.GetString("DB_MAX_IDLE_TIME", "15m"),
		JWTSecret:       env.GetRequired("JWT_SECRET"),
		DefaultTagLimit: env.GetInt64("DEFAULT_TAG_LIMIT", 100),
	}
}
