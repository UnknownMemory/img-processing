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
	"github.com/unknownmemory/img-processing/internal/api"
)

const version = "0.1.0"

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	var cfg api.Config

	cfg.Port, _ = strconv.Atoi(os.Getenv("PORT"))
	cfg.Mode = os.Getenv("MODE")
	cfg.DB.DSN = os.Getenv("DB_DSN")

	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)

	db, err := pgxpool.New(context.Background(), cfg.DB.DSN)
	if err != nil {
		log.Fatalf("Unable to create connection pool: %v\n", err)
	}

	defer db.Close()

	app := api.NewApplication(cfg, logger, db, version)

	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.Port),
		Handler:      app.Routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	logger.Printf("starting %s server on %s", cfg.Mode, server.Addr)
	err = server.ListenAndServe()
	logger.Fatal(err)
}
