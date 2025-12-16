package config

import (
	"github.com/manuelmtzv/mangocatnotes-api/internal/env"
)

type Config struct {
	Port                 string
	DBAddr               string
	MaxOpenConns         int
	MaxIdleConns         int
	MaxIdleTime          string
	RedisAddr            string
	RedisPassword        string
	DefaultTagLimit      int64
	BaseURL              string
	SessionDurationHours int
	IsProd               bool
	AllowedOrigins       []string
}

func LoadConfig() *Config {
	env.Load()

	return &Config{
		Port:                 env.GetString("PORT", "3000"),
		DBAddr:               env.GetRequired("DB_ADDR"),
		MaxOpenConns:         env.GetInt("DB_MAX_OPEN_CONNS", 30),
		MaxIdleConns:         env.GetInt("DB_MAX_IDLE_CONNS", 30),
		MaxIdleTime:          env.GetString("DB_MAX_IDLE_TIME", "15m"),
		RedisAddr:            env.GetString("REDIS_ADDR", "localhost:6379"),
		RedisPassword:        env.GetString("REDIS_PASSWORD", ""),
		DefaultTagLimit:      env.GetInt64("DEFAULT_TAG_LIMIT", 100),
		BaseURL:              env.GetString("BASE_URL", "http://localhost:8080"),
		SessionDurationHours: env.GetInt("SESSION_DURATION_HOURS", 168),
		IsProd:               env.GetBool("IS_PROD", false),
		AllowedOrigins:       env.GetSlice("ALLOWED_ORIGINS", []string{"http://localhost:3000"}),
	}
}
