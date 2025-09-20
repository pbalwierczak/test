package repository

import (
	"context"
	"database/sql"
)

// sqlUnitOfWork implements UnitOfWork using native SQL transactions
type sqlUnitOfWork struct {
	db *sql.DB
}

// NewSQLUnitOfWork creates a new UnitOfWork instance using native SQL
func NewSQLUnitOfWork(db *sql.DB) UnitOfWork {
	return &sqlUnitOfWork{db: db}
}

// Begin starts a new transaction
func (u *sqlUnitOfWork) Begin(ctx context.Context) (UnitOfWorkTx, error) {
	tx, err := u.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}

	return &sqlUnitOfWorkTx{
		tx: tx,
	}, nil
}

// sqlUnitOfWorkTx implements UnitOfWorkTx using native SQL transaction
type sqlUnitOfWorkTx struct {
	tx *sql.Tx
}

// ScooterRepository returns a scooter repository that uses this transaction
func (u *sqlUnitOfWorkTx) ScooterRepository() ScooterRepository {
	return &sqlScooterRepository{db: u.tx}
}

// TripRepository returns a trip repository that uses this transaction
func (u *sqlUnitOfWorkTx) TripRepository() TripRepository {
	return &sqlTripRepository{db: u.tx}
}

// UserRepository returns a user repository that uses this transaction
func (u *sqlUnitOfWorkTx) UserRepository() UserRepository {
	return &sqlUserRepository{db: u.tx}
}

// LocationUpdateRepository returns a location update repository that uses this transaction
func (u *sqlUnitOfWorkTx) LocationUpdateRepository() LocationUpdateRepository {
	return &sqlLocationUpdateRepository{db: u.tx}
}

// Commit commits the transaction
func (u *sqlUnitOfWorkTx) Commit() error {
	return u.tx.Commit()
}

// Rollback rolls back the transaction
func (u *sqlUnitOfWorkTx) Rollback() error {
	return u.tx.Rollback()
}
