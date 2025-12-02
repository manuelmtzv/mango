package env

import (
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

func Load() error {
	return godotenv.Load(".env")
}

func GetString(key, fallback string) string {
	val, ok := os.LookupEnv(key)
	if !ok {
		return fallback
	}

	return val
}

func GetInt(key string, fallback int) int {
	val, ok := os.LookupEnv(key)
	if !ok {
		return fallback
	}

	number, err := strconv.Atoi(val)
	if err != nil {
		return fallback
	}

	return number
}

func GetInt64(key string, fallback int64) int64 {
	val, ok := os.LookupEnv(key)
	if !ok {
		return fallback
	}

	number, err := strconv.ParseInt(val, 10, 64)
	if err != nil {
		return fallback
	}

	return number
}

func GetBool(key string, fallback bool) bool {
	val, ok := os.LookupEnv(key)
	if !ok {
		return fallback
	}

	result, err := strconv.ParseBool(val)
	if err != nil {
		return fallback
	}

	return result
}

func GetDuration(key string, fallback time.Duration) time.Duration {
	val, ok := os.LookupEnv(key)
	if !ok {
		return fallback
	}

	duration, err := time.ParseDuration(val)
	if err != nil {
		return fallback
	}

	return duration
}

func GetSlice(key string, fallback []string) []string {
	val, ok := os.LookupEnv(key)
	if !ok {
		return fallback
	}

	return strings.Split(val, ",")
}
