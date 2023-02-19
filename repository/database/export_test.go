package db

import (
	"database/sql"

	"github.com/android-project-46group/core-api/repository"
)

func NewMockDatabase(db *sql.DB) repository.Database {
	//nolint
	return &database{
		db: db,
	}
}
