package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/ekastn/load-stuffing-calculator/internal/api"
	"github.com/ekastn/load-stuffing-calculator/internal/config"
	"github.com/ekastn/load-stuffing-calculator/internal/seeder"
	"github.com/ekastn/load-stuffing-calculator/internal/store"
	"github.com/jackc/pgx/v5/pgxpool"
)

//	@title			Load Stuffing Calculator API
//	@version		1.0
//	@description	API Server for Load Stuffing Calculator application.
//	@termsOfService	http://swagger.io/terms/

//	@contact.name	API Support
//	@contact.url	http://www.swagger.io/support
//	@contact.email	support@swagger.io

//	@license.name	Apache 2.0
//	@license.url	http://www.apache.org/licenses/LICENSE-2.0.html

//	@BasePath	/api/v1

// @securityDefinitions.apikey	BearerAuth
// @in							header
// @name						Authorization
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

	// Seed Founder User
	querier := store.New(dbPool)
	seed := seeder.New(querier, cfg)
	if err := seed.SeedFounder(context.Background()); err != nil {
		log.Printf("Warning: Failed to seed founder user: %v", err)
	}

	app := api.NewApp(cfg, dbPool)
	if err := app.Run(); err != nil {
		log.Fatalf("Server stopped with error: %v", err)
	}
}
