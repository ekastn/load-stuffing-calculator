package config

import (
	"log"

	"github.com/ekastn/load-stuffing-calculator/internal/env"
	"github.com/joho/godotenv"
)

type Config struct {
	Addr        string
	DatabaseURL string
	JWTSecret   string

	PackingServiceURL string

	FounderUsername string
	FounderEmail    string
	FounderPassword string
}

func Load() Config {
	if env.GetString("SRV_ENV", "dev") == "dev" {
		err := godotenv.Load()
		if err != nil {
			log.Println("Error loading .env file in development mode, continuing without it.")
		}
	}
	return Config{
		Addr:        env.GetString("SRV_ADDR", ":8080"),
		DatabaseURL: env.GetString("DATABASE_URL", ""),
		JWTSecret:   env.GetString("JWT_SECRET", "secret"),

		PackingServiceURL: env.GetString("PACKING_SERVICE_URL", "http://localhost:5051"),

		// Founder bootstrap (backwards compatible with ADMIN_*).
		FounderUsername: env.GetString("FOUNDER_USERNAME", env.GetString("ADMIN_USERNAME", "admin")),
		FounderEmail:    env.GetString("FOUNDER_EMAIL", env.GetString("ADMIN_EMAIL", "admin@example.com")),
		FounderPassword: env.GetString("FOUNDER_PASSWORD", env.GetString("ADMIN_PASSWORD", "admin123")),
	}
}
