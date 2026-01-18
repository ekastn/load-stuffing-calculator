package main

import (
	"context"
	"log"

	"load-stuffing-calculator/internal/api"
	"load-stuffing-calculator/internal/config"

	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
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
