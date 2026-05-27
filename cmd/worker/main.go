package main

import (
	"context"
	"log"
	"os"
	"strconv"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/unknownmemory/img-processing/internal/aws"
	db "github.com/unknownmemory/img-processing/internal/database"
	process "github.com/unknownmemory/img-processing/internal/proc"
	"github.com/unknownmemory/img-processing/internal/rabbitmq"
)

func main() {
	err := godotenv.Load()
	failOnError(err, "Error loading .env file")

	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)
	logger.Printf("Starting worker")

	dbPool, err := pgxpool.New(context.Background(), os.Getenv("DB_DSN"))
	if err != nil {
		log.Fatalf("Unable to create connection pool: %v\n", err)
	}
	defer dbPool.Close()

	storage := aws.NewS3Client()
	proc := process.NewProcessor()
	database := db.New(dbPool)

	worker, err := rabbitmq.NewWorker(os.Getenv("RABBIT_MQ"), logger, storage, proc, database)
	if err != nil {
		log.Fatalf("Unable to connect to RabbitMQ: %v\n", err)
	}
	defer worker.Close()

	workerPool, err := strconv.Atoi(os.Getenv("WORKER_POOL"))
	if err != nil {
		log.Fatalf("Unable to parse WORKER_POOL: %v\n", err)
	}

	worker.Listen(workerPool)
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", err, msg)
	}
}
