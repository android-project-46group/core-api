package db

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/android-project-46group/core-api/config"
	"github.com/android-project-46group/core-api/repository"
	"github.com/android-project-46group/core-api/util/logger"
	_ "github.com/lib/pq" // postgres driver
)

type database struct {
	db     *sql.DB
	logger logger.Logger
}

func New(ctx context.Context, cfg config.Config, logger logger.Logger) (repository.Database, error) {
	source := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName,
	)

	sqlDB, err := sql.Open(cfg.DBDriver, source)
	if err != nil {
		return nil, fmt.Errorf("failed to open sql: %w", err)
	}

	db := &database{
		db:     sqlDB,
		logger: logger,
	}

	return db, nil
}
