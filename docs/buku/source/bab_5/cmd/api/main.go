package main

import (
	"context"
	"log"

	"load-stuffing-calculator/internal/api"
	"load-stuffing-calculator/internal/config"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

func main() {
	// Load .env file jika ada (opsional untuk development)
	// Di production, environment variables di-set langsung oleh orchestrator
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	cfg := config.Load()

	pool, err := pgxpool.New(context.Background(), cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v", err)
	}
	defer pool.Close()

	app := api.NewApp(cfg, pool)
	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}
