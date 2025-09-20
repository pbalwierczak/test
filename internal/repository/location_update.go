package repository

import (
	"context"

	"scootin-aboot/internal/models"

	"github.com/google/uuid"
)

type LocationUpdateRepository interface {
	Create(ctx context.Context, update *models.LocationUpdate) error
	GetByID(ctx context.Context, id uuid.UUID) (*models.LocationUpdate, error)
	Update(ctx context.Context, update *models.LocationUpdate) error
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context, limit, offset int) ([]*models.LocationUpdate, error)

	GetByScooterID(ctx context.Context, scooterID uuid.UUID) ([]*models.LocationUpdate, error)
}
