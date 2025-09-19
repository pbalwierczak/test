package repository

import (
	"context"
	"time"

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
	GetByScooterIDOrdered(ctx context.Context, scooterID uuid.UUID) ([]*models.LocationUpdate, error)
	GetLatestByScooterID(ctx context.Context, scooterID uuid.UUID) (*models.LocationUpdate, error)

	// Time-based queries
	GetByDateRange(ctx context.Context, start, end time.Time) ([]*models.LocationUpdate, error)
	GetByScooterIDAndDateRange(ctx context.Context, scooterID uuid.UUID, start, end time.Time) ([]*models.LocationUpdate, error)

	// Geographic queries
	GetInBounds(ctx context.Context, minLat, maxLat, minLng, maxLng float64) ([]*models.LocationUpdate, error)
	GetInRadius(ctx context.Context, latitude, longitude, radiusKm float64) ([]*models.LocationUpdate, error)

	// Statistics
	GetUpdateCount(ctx context.Context) (int64, error)
	GetUpdateCountByScooter(ctx context.Context, scooterID uuid.UUID) (int64, error)
}
