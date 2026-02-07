package main

import (
	"context"
	"log"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/unknownmemory/img-processing/internal/rabbitmq"
)

func main() {
	err := godotenv.Load()
	failOnError(err, "Error loading .env file")

	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)
	logger.Printf("Starting worker")

	db, err := pgxpool.New(context.Background(), os.Getenv("DB_DSN"))
	if err != nil {
		log.Fatalf("Unable to create connection pool: %v\n", err)
	}

	defer db.Close()

	worker := rabbitmq.NewWorker(os.Getenv("RABBIT_MQ"), logger, db)
	worker.Listen()
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", err, msg)
	}
}
