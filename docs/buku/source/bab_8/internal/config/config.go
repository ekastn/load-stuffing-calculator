package config

import (
	"os"
	"strings"
)

type Config struct {
	DatabaseURL        string
	ServerPort         string
	PackingServiceURL  string
	CORSAllowedOrigins []string
	JWTSecret          string
}

func Load() *Config {
	return &Config{
		DatabaseURL:        getEnv("DATABASE_URL", "postgres://localhost:5432/loadstuff"),
		ServerPort:         getEnv("SERVER_PORT", "8080"),
		PackingServiceURL:  getEnv("PACKING_SERVICE_URL", "http://localhost:5000"),
		CORSAllowedOrigins: strings.Split(getEnv("CORS_ALLOWED_ORIGINS", "http://localhost:3000"), ","),
		JWTSecret:          getEnv("JWT_SECRET", "verysecret"),
	}
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
