package repository

import (
	"database/sql"
)

type sqlRepository struct {
	db         *sql.DB
	unitOfWork UnitOfWork
}

func (r *sqlRepository) Scooter() ScooterRepository {
	return &sqlScooterRepository{db: r.db}
}

func (r *sqlRepository) Trip() TripRepository {
	return &sqlTripRepository{db: r.db}
}

func (r *sqlRepository) User() UserRepository {
	return &sqlUserRepository{db: r.db}
}

func (r *sqlRepository) LocationUpdate() LocationUpdateRepository {
	return &sqlLocationUpdateRepository{db: r.db}
}

func (r *sqlRepository) UnitOfWork() UnitOfWork {
	return r.unitOfWork
}
