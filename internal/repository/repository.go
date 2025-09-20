package repository

import (
	"context"
	"errors"
)

// Repository defines the interface for all repository operations
type Repository interface {
	Scooter() ScooterRepository
	Trip() TripRepository
	User() UserRepository
	LocationUpdate() LocationUpdateRepository
	UnitOfWork() UnitOfWork
}

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

// Common repository errors
var (
	ErrScooterNotFound        = errors.New("scooter not found")
	ErrTripNotFound           = errors.New("trip not found")
	ErrUserNotFound           = errors.New("user not found")
	ErrLocationUpdateNotFound = errors.New("location update not found")
)
