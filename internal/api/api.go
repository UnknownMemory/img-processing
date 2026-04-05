package api

import (
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/unknownmemory/img-processing/internal/aws"
	"github.com/unknownmemory/img-processing/internal/rabbitmq"
)

type Config struct {
	Port int
	Mode string
	DB   struct {
		DSN string
	}
}

type Application struct {
	config  Config
	logger  *log.Logger
	db      *pgxpool.Pool
	s3      *aws.S3Client
	rmq     *rabbitmq.RabbitMQ
	version string
}

func NewApplication(cfg Config, logger *log.Logger, db *pgxpool.Pool, rmq *rabbitmq.RabbitMQ, version string) *Application {
	return &Application{
		config:  cfg,
		logger:  logger,
		db:      db,
		s3:      aws.NewS3Client(),
		rmq:     rmq,
		version: version,
	}
}
