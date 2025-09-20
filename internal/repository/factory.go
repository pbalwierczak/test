package repository

import (
	"database/sql"
)

// NewRepository creates a new repository using native SQL
func NewRepository(db *sql.DB) Repository {
	return &sqlRepository{
		db:         db,
		unitOfWork: NewSQLUnitOfWork(db),
	}
}
