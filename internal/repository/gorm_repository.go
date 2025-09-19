package repository

import (
	"gorm.io/gorm"
)

type gormRepository struct {
	db         *gorm.DB
	unitOfWork UnitOfWork
}

func NewGormRepository(db *gorm.DB) Repository {
	return &gormRepository{
		db:         db,
		unitOfWork: NewUnitOfWork(db),
	}
}

func (r *gormRepository) Scooter() ScooterRepository {
	return &gormScooterRepository{db: r.db}
}

func (r *gormRepository) Trip() TripRepository {
	return &gormTripRepository{db: r.db}
}

func (r *gormRepository) User() UserRepository {
	return &gormUserRepository{db: r.db}
}

func (r *gormRepository) LocationUpdate() LocationUpdateRepository {
	return &gormLocationUpdateRepository{db: r.db}
}

func (r *gormRepository) UnitOfWork() UnitOfWork {
	return r.unitOfWork
}
