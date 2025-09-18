package repository

import (
	"gorm.io/gorm"
)

// gormRepository implements the Repository interface using GORM
type gormRepository struct {
	db *gorm.DB
}

// NewGormRepository creates a new GORM-based repository
func NewGormRepository(db *gorm.DB) Repository {
	return &gormRepository{db: db}
}

// Scooter returns the scooter repository
func (r *gormRepository) Scooter() ScooterRepository {
	return &gormScooterRepository{db: r.db}
}

// Trip returns the trip repository
func (r *gormRepository) Trip() TripRepository {
	return &gormTripRepository{db: r.db}
}

// User returns the user repository
func (r *gormRepository) User() UserRepository {
	return &gormUserRepository{db: r.db}
}

// LocationUpdate returns the location update repository
func (r *gormRepository) LocationUpdate() LocationUpdateRepository {
	return &gormLocationUpdateRepository{db: r.db}
}
