package repository

import (
	"database/sql"
)

func NewRepository(db *sql.DB) Repository {
	return &sqlRepository{
		db:         db,
		unitOfWork: NewSQLUnitOfWork(db),
	}
}
