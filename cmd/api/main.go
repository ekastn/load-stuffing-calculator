package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/ekastn/load-stuffing-calculator/internal/api"
	"github.com/ekastn/load-stuffing-calculator/internal/config"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	cfg := config.Load()

	dbPool, err := pgxpool.New(context.Background(), cfg.DatabaseURL)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer dbPool.Close()

	err = dbPool.Ping(context.Background())
	if err != nil {
		fmt.Fprintf(os.Stderr, "Database ping failed: %v\n", err)
		os.Exit(1)
	}
	log.Println("Database connected successfully!")

	app := api.NewApp(cfg, dbPool)
	if err := app.Run(); err != nil {
		log.Fatalf("Server stopped with error: %v", err)
	}
}
