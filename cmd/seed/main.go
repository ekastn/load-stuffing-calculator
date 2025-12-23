package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/ekastn/load-stuffing-calculator/internal/config"
	"github.com/jackc/pgx/v5/pgxpool"
)

type stringSliceFlag []string

func (s *stringSliceFlag) String() string {
	return fmt.Sprint([]string(*s))
}

func (s *stringSliceFlag) Set(value string) error {
	if value == "" {
		return errors.New("file path cannot be empty")
	}
	*s = append(*s, value)
	return nil
}

func main() {
	var files stringSliceFlag
	var dryRun bool

	flag.Var(&files, "file", "SQL file to execute (repeatable)")
	flag.BoolVar(&dryRun, "dry-run", false, "Print execution order without running")
	flag.Parse()

	if len(files) == 0 {
		fmt.Fprintln(os.Stderr, "error: at least one --file is required")
		flag.Usage()
		os.Exit(2)
	}

	log.SetFlags(0)
	start := time.Now()

	if dryRun {
		for i, file := range files {
			log.Printf("dry-run [%d/%d]: %s", i+1, len(files), file)
		}
		return
	}

	cfg := config.Load()
	if cfg.DatabaseURL == "" {
		fmt.Fprintln(os.Stderr, "error: DATABASE_URL is required")
		os.Exit(2)
	}

	ctx := context.Background()
	log.Printf("connecting to database")
	pool, err := pgxpool.New(ctx, cfg.DatabaseURL)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: failed to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer pool.Close()

	if err := pool.Ping(ctx); err != nil {
		fmt.Fprintf(os.Stderr, "error: database ping failed: %v\n", err)
		os.Exit(1)
	}
	log.Printf("database connected")

	log.Printf("begin transaction")
	tx, err := pool.Begin(ctx)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: failed to begin transaction: %v\n", err)
		os.Exit(1)
	}

	committed := false
	defer func() {
		if committed {
			return
		}
		log.Printf("rollback transaction")
		_ = tx.Rollback(ctx)
	}()

	for i, file := range files {
		log.Printf("executing [%d/%d]: %s", i+1, len(files), file)
		b, err := os.ReadFile(file)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: failed to read %s: %v\n", file, err)
			os.Exit(1)
		}
		if len(b) == 0 {
			fmt.Fprintf(os.Stderr, "error: file is empty: %s\n", file)
			os.Exit(1)
		}
		if _, err := tx.Exec(ctx, string(b)); err != nil {
			fmt.Fprintf(os.Stderr, "error: failed executing %s: %v\n", file, err)
			os.Exit(1)
		}
		log.Printf("done [%d/%d]: %s", i+1, len(files), file)
	}

	log.Printf("commit transaction")
	if err := tx.Commit(ctx); err != nil {
		fmt.Fprintf(os.Stderr, "error: commit failed: %v\n", err)
		os.Exit(1)
	}
	committed = true

	log.Printf("ok (%s)", time.Since(start).Round(time.Millisecond))
}
