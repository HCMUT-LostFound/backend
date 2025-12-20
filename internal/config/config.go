package config

import (
	"os"
)

type Config struct {
	DBUrl string
}

func Load() *Config {
	return &Config{
		DBUrl: getEnv(
			"DATABASE_URL",
			"postgres://postgres:postgres@localhost:5432/lostfound_app?sslmode=disable",
		),
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
