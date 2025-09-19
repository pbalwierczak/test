package repository

import (
	"context"

	"scootin-aboot/internal/models"

	"github.com/google/uuid"
)

// LocationUpdateRepository defines the interface for location update data operations
type LocationUpdateRepository interface {
	// Basic CRUD operations
	Create(ctx context.Context, update *models.LocationUpdate) error
	GetByID(ctx context.Context, id uuid.UUID) (*models.LocationUpdate, error)
	Update(ctx context.Context, update *models.LocationUpdate) error
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context, limit, offset int) ([]*models.LocationUpdate, error)

	// Scooter-specific queries
	GetByScooterID(ctx context.Context, scooterID uuid.UUID) ([]*models.LocationUpdate, error)
}
