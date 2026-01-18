package config

import (
	"os"
)

type Config struct {
	DatabaseURL       string
	ServerPort        string
	PackingServiceURL string
}

func Load() *Config {
	return &Config{
		DatabaseURL:       getEnv("DATABASE_URL", "postgres://localhost:5432/loadstuff"),
		ServerPort:        getEnv("SERVER_PORT", "8080"),
		PackingServiceURL: getEnv("PACKING_SERVICE_URL", "http://localhost:5000"),
	}
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
