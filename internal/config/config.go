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
	}
}
