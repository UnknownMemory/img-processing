package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

const version = "0.1.0"

type config struct {
	port int
	mode string
	db   struct {
		dsn string
	}
}

type application struct {
	config config
	logger *log.Logger
	db     *pgxpool.Pool
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	var cfg config

	cfg.port, _ = strconv.Atoi(os.Getenv("PORT"))
	cfg.mode = os.Getenv("MODE")
	cfg.db.dsn = os.Getenv("DB_DSN")

	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)

	db, err := pgxpool.New(context.Background(), cfg.db.dsn)
	if err != nil {
		log.Fatalf("Unable to create connection pool: %v\n", err)
		os.Exit(1)
	}

	defer db.Close()

	app := &application{
		config: cfg,
		logger: logger,
		db:     db,
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/api/v1/healthcheck", app.healthcheckHandler)

	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.port),
		Handler:      mux,
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	logger.Printf("starting %s server on %s", cfg.mode, server.Addr)
	err = server.ListenAndServe()
	logger.Fatal(err)
}
