package db

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/android-project-46group/core-api/config"
	"github.com/android-project-46group/core-api/repository"
	_ "github.com/lib/pq"
)

type database struct {
	db *sql.DB
}

func New(ctx context.Context, c config.Config) (repository.Database, error) {
	source := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		c.DBHost, c.DBPort, c.DBUser, c.DBPassword, c.DBName,
	)
	sqlDB, err := sql.Open(c.DBDriver, source)
	if err != nil {
		return nil, fmt.Errorf("failed to open sql: %w", err)
	}

	db := &database{
		db: sqlDB,
	}

	return db, nil
}
