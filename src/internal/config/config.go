package config

import (
	"os"
	"time"
)

type Config struct {
	HTTPPort        string
	DatabaseURL     string
	ShutdownTimeout time.Duration
}

func Load() Config {
	return Config{
		HTTPPort:        env("HTTP_PORT", "8080"),
		DatabaseURL:     env("DATABASE_URL", "postgres://storeapi:storeapi@localhost:5432/storeapi?sslmode=disable"),
		ShutdownTimeout: durationEnv("SHUTDOWN_TIMEOUT", 10*time.Second),
	}
}

func env(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

func durationEnv(key string, fallback time.Duration) time.Duration {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	d, err := time.ParseDuration(value)
	if err != nil {
		return fallback
	}
	return d
}
