package repository

import (
	"context"
	"errors"
)

type Repository interface {
	Scooter() ScooterRepository
	Trip() TripRepository
	User() UserRepository
	LocationUpdate() LocationUpdateRepository
	UnitOfWork() UnitOfWork
}

type UnitOfWork interface {
	Begin(ctx context.Context) (UnitOfWorkTx, error)
}

type UnitOfWorkTx interface {
	ScooterRepository() ScooterRepository
	TripRepository() TripRepository
	UserRepository() UserRepository
	LocationUpdateRepository() LocationUpdateRepository

	Commit() error
	Rollback() error
}

var (
	ErrScooterNotFound        = errors.New("scooter not found")
	ErrTripNotFound           = errors.New("trip not found")
	ErrUserNotFound           = errors.New("user not found")
	ErrLocationUpdateNotFound = errors.New("location update not found")
)
