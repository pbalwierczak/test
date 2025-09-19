package repository

import (
	"context"

	"scootin-aboot/internal/models"

	"github.com/google/uuid"
)

type TripRepository interface {
	Create(ctx context.Context, trip *models.Trip) error
	GetByID(ctx context.Context, id uuid.UUID) (*models.Trip, error)
	Update(ctx context.Context, trip *models.Trip) error
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context, limit, offset int) ([]*models.Trip, error)

	UpdateStatus(ctx context.Context, id uuid.UUID, status models.TripStatus) error
	EndTrip(ctx context.Context, id uuid.UUID, endLat, endLng float64) error
	CancelTrip(ctx context.Context, id uuid.UUID) error

	GetActiveByUserID(ctx context.Context, userID uuid.UUID) (*models.Trip, error)
	GetActiveByScooterID(ctx context.Context, scooterID uuid.UUID) (*models.Trip, error)
}
