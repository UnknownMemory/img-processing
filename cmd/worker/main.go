package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	failOnError(err, "Error loading .env file")

	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)
	logger.Printf("Starting worker")

	worker := NewWorker(os.Getenv("RABBIT_MQ"), logger)
	worker.Connect()
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", err, msg)
	}
}
