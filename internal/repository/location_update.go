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

	// Trip-specific queries
	GetByTripID(ctx context.Context, tripID uuid.UUID) ([]*models.LocationUpdate, error)
	GetByTripIDOrdered(ctx context.Context, tripID uuid.UUID) ([]*models.LocationUpdate, error)
	GetLatestByTripID(ctx context.Context, tripID uuid.UUID) (*models.LocationUpdate, error)

	// Time-based queries
	GetByDateRange(ctx context.Context, start, end time.Time) ([]*models.LocationUpdate, error)
	GetByTripIDAndDateRange(ctx context.Context, tripID uuid.UUID, start, end time.Time) ([]*models.LocationUpdate, error)

	// Geographic queries
	GetInBounds(ctx context.Context, minLat, maxLat, minLng, maxLng float64) ([]*models.LocationUpdate, error)
	GetInRadius(ctx context.Context, latitude, longitude, radiusKm float64) ([]*models.LocationUpdate, error)

	// Statistics
	GetUpdateCount(ctx context.Context) (int64, error)
	GetUpdateCountByTrip(ctx context.Context, tripID uuid.UUID) (int64, error)
}
