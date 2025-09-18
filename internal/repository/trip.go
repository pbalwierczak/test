package repository

import (
	"context"
	"time"

	"scootin-aboot/internal/models"

	"github.com/google/uuid"
)

// TripRepository defines the interface for trip data operations
type TripRepository interface {
	// Basic CRUD operations
	Create(ctx context.Context, trip *models.Trip) error
	GetByID(ctx context.Context, id uuid.UUID) (*models.Trip, error)
	Update(ctx context.Context, trip *models.Trip) error
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context, limit, offset int) ([]*models.Trip, error)

	// Status operations
	UpdateStatus(ctx context.Context, id uuid.UUID, status models.TripStatus) error
	EndTrip(ctx context.Context, id uuid.UUID, endLat, endLng float64) error
	CancelTrip(ctx context.Context, id uuid.UUID) error

	// Query operations
	GetByStatus(ctx context.Context, status models.TripStatus) ([]*models.Trip, error)
	GetActive(ctx context.Context) ([]*models.Trip, error)
	GetCompleted(ctx context.Context) ([]*models.Trip, error)
	GetCancelled(ctx context.Context) ([]*models.Trip, error)

	// User and scooter specific queries
	GetByUserID(ctx context.Context, userID uuid.UUID) ([]*models.Trip, error)
	GetByScooterID(ctx context.Context, scooterID uuid.UUID) ([]*models.Trip, error)
	GetActiveByUserID(ctx context.Context, userID uuid.UUID) (*models.Trip, error)
	GetActiveByScooterID(ctx context.Context, scooterID uuid.UUID) (*models.Trip, error)

	// Time-based queries
	GetByDateRange(ctx context.Context, start, end time.Time) ([]*models.Trip, error)
	GetByUserIDAndDateRange(ctx context.Context, userID uuid.UUID, start, end time.Time) ([]*models.Trip, error)
	GetByScooterIDAndDateRange(ctx context.Context, scooterID uuid.UUID, start, end time.Time) ([]*models.Trip, error)

	// Statistics
	GetTripCount(ctx context.Context) (int64, error)
	GetTripCountByUser(ctx context.Context, userID uuid.UUID) (int64, error)
	GetTripCountByScooter(ctx context.Context, scooterID uuid.UUID) (int64, error)
}
