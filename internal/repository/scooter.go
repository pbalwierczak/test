package repository

import (
	"context"

	"scootin-aboot/internal/models"

	"github.com/google/uuid"
)

// ScooterRepository defines the interface for scooter data operations
type ScooterRepository interface {
	// Basic CRUD operations
	Create(ctx context.Context, scooter *models.Scooter) error
	GetByID(ctx context.Context, id uuid.UUID) (*models.Scooter, error)
	Update(ctx context.Context, scooter *models.Scooter) error
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context, limit, offset int) ([]*models.Scooter, error)

	// Status operations
	UpdateStatus(ctx context.Context, id uuid.UUID, status models.ScooterStatus) error
	UpdateLocation(ctx context.Context, id uuid.UUID, latitude, longitude float64) error

	// Query operations
	GetByStatus(ctx context.Context, status models.ScooterStatus) ([]*models.Scooter, error)
	GetAvailable(ctx context.Context) ([]*models.Scooter, error)
	GetOccupied(ctx context.Context) ([]*models.Scooter, error)

	// Geographic operations
	GetInBounds(ctx context.Context, minLat, maxLat, minLng, maxLng float64) ([]*models.Scooter, error)
	GetClosest(ctx context.Context, latitude, longitude float64, limit int) ([]*models.Scooter, error)
	GetClosestWithRadius(ctx context.Context, latitude, longitude, radius float64, status string, limit int) ([]*models.Scooter, error)
	GetInRadius(ctx context.Context, latitude, longitude, radiusKm float64) ([]*models.Scooter, error)

	// Combined filters
	GetAvailableInBounds(ctx context.Context, minLat, maxLat, minLng, maxLng float64) ([]*models.Scooter, error)
	GetByStatusInBounds(ctx context.Context, status models.ScooterStatus, minLat, maxLat, minLng, maxLng float64) ([]*models.Scooter, error)
}
