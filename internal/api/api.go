package api

import (
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
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
	version string
}

func NewApplication(cfg Config, logger *log.Logger, db *pgxpool.Pool, version string) *Application {
	return &Application{
		config:  cfg,
		logger:  logger,
		db:      db,
		version: version,
	}
}
