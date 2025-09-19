package repository

import (
	"context"

	"gorm.io/gorm"
)

// UnitOfWork defines the interface for managing database transactions
type UnitOfWork interface {
	Begin(ctx context.Context) (UnitOfWorkTx, error)
}

// UnitOfWorkTx represents a database transaction
type UnitOfWorkTx interface {
	// Repository accessors
	ScooterRepository() ScooterRepository
	TripRepository() TripRepository
	UserRepository() UserRepository
	LocationUpdateRepository() LocationUpdateRepository

	// Transaction control
	Commit() error
	Rollback() error
}

// gormUnitOfWork implements UnitOfWork using GORM
type gormUnitOfWork struct {
	db *gorm.DB
}

// NewUnitOfWork creates a new UnitOfWork instance
func NewUnitOfWork(db *gorm.DB) UnitOfWork {
	return &gormUnitOfWork{db: db}
}

// Begin starts a new transaction
func (u *gormUnitOfWork) Begin(ctx context.Context) (UnitOfWorkTx, error) {
	tx := u.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return nil, tx.Error
	}

	return &gormUnitOfWorkTx{
		tx: tx,
	}, nil
}

// gormUnitOfWorkTx implements UnitOfWorkTx using GORM transaction
type gormUnitOfWorkTx struct {
	tx *gorm.DB
}

// ScooterRepository returns a scooter repository that uses this transaction
func (u *gormUnitOfWorkTx) ScooterRepository() ScooterRepository {
	return &gormScooterRepository{db: u.tx}
}

// TripRepository returns a trip repository that uses this transaction
func (u *gormUnitOfWorkTx) TripRepository() TripRepository {
	return &gormTripRepository{db: u.tx}
}

// UserRepository returns a user repository that uses this transaction
func (u *gormUnitOfWorkTx) UserRepository() UserRepository {
	return &gormUserRepository{db: u.tx}
}

// LocationUpdateRepository returns a location update repository that uses this transaction
func (u *gormUnitOfWorkTx) LocationUpdateRepository() LocationUpdateRepository {
	return &gormLocationUpdateRepository{db: u.tx}
}

// Commit commits the transaction
func (u *gormUnitOfWorkTx) Commit() error {
	return u.tx.Commit().Error
}

// Rollback rolls back the transaction
func (u *gormUnitOfWorkTx) Rollback() error {
	return u.tx.Rollback().Error
}
