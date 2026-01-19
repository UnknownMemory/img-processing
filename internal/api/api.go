package api

import (
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/unknownmemory/img-processing/internal/aws"
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
	version string
}

func NewApplication(cfg Config, logger *log.Logger, db *pgxpool.Pool, version string) *Application {
	return &Application{
		config:  cfg,
		logger:  logger,
		db:      db,
		s3:      aws.NewS3Client(),
		version: version,
	}
}
