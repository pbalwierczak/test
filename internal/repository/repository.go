package repository

import "errors"

// Repository defines the interface for all repository operations
type Repository interface {
	Scooter() ScooterRepository
	Trip() TripRepository
	User() UserRepository
	LocationUpdate() LocationUpdateRepository
	UnitOfWork() UnitOfWork
}

// Common repository errors
var (
	ErrScooterNotFound        = errors.New("scooter not found")
	ErrTripNotFound           = errors.New("trip not found")
	ErrUserNotFound           = errors.New("user not found")
	ErrLocationUpdateNotFound = errors.New("location update not found")
)
