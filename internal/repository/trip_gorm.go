package repository

import (
	"context"
	"time"

	"scootin-aboot/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// gormTripRepository implements TripRepository using GORM
type gormTripRepository struct {
	db *gorm.DB
}

// Create creates a new trip
func (r *gormTripRepository) Create(ctx context.Context, trip *models.Trip) error {
	return r.db.WithContext(ctx).Create(trip).Error
}

// GetByID retrieves a trip by ID
func (r *gormTripRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Trip, error) {
	var trip models.Trip
	err := r.db.WithContext(ctx).First(&trip, "id = ?", id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &trip, nil
}

// Update updates a trip
func (r *gormTripRepository) Update(ctx context.Context, trip *models.Trip) error {
	return r.db.WithContext(ctx).Save(trip).Error
}

// Delete deletes a trip by ID
func (r *gormTripRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&models.Trip{}, "id = ?", id).Error
}

// List retrieves trips with pagination
func (r *gormTripRepository) List(ctx context.Context, limit, offset int) ([]*models.Trip, error) {
	var trips []*models.Trip
	query := r.db.WithContext(ctx)

	if limit > 0 {
		query = query.Limit(limit)
	}
	if offset > 0 {
		query = query.Offset(offset)
	}

	err := query.Find(&trips).Error
	return trips, err
}

// UpdateStatus updates a trip's status
func (r *gormTripRepository) UpdateStatus(ctx context.Context, id uuid.UUID, status models.TripStatus) error {
	return r.db.WithContext(ctx).Model(&models.Trip{}).
		Where("id = ?", id).
		Update("status", status).Error
}

// EndTrip ends a trip with end coordinates
func (r *gormTripRepository) EndTrip(ctx context.Context, id uuid.UUID, endLat, endLng float64) error {
	now := time.Now()
	return r.db.WithContext(ctx).Model(&models.Trip{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"end_time":      now,
			"end_latitude":  endLat,
			"end_longitude": endLng,
			"status":        models.TripStatusCompleted,
		}).Error
}

// CancelTrip cancels a trip
func (r *gormTripRepository) CancelTrip(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Model(&models.Trip{}).
		Where("id = ?", id).
		Update("status", models.TripStatusCancelled).Error
}

// GetByStatus retrieves trips by status
func (r *gormTripRepository) GetByStatus(ctx context.Context, status models.TripStatus) ([]*models.Trip, error) {
	var trips []*models.Trip
	err := r.db.WithContext(ctx).Where("status = ?", status).Find(&trips).Error
	return trips, err
}

// GetActive retrieves active trips
func (r *gormTripRepository) GetActive(ctx context.Context) ([]*models.Trip, error) {
	return r.GetByStatus(ctx, models.TripStatusActive)
}

// GetCompleted retrieves completed trips
func (r *gormTripRepository) GetCompleted(ctx context.Context) ([]*models.Trip, error) {
	return r.GetByStatus(ctx, models.TripStatusCompleted)
}

// GetCancelled retrieves cancelled trips
func (r *gormTripRepository) GetCancelled(ctx context.Context) ([]*models.Trip, error) {
	return r.GetByStatus(ctx, models.TripStatusCancelled)
}

// GetByUserID retrieves trips by user ID
func (r *gormTripRepository) GetByUserID(ctx context.Context, userID uuid.UUID) ([]*models.Trip, error) {
	var trips []*models.Trip
	err := r.db.WithContext(ctx).Where("user_id = ?", userID).Find(&trips).Error
	return trips, err
}

// GetByScooterID retrieves trips by scooter ID
func (r *gormTripRepository) GetByScooterID(ctx context.Context, scooterID uuid.UUID) ([]*models.Trip, error) {
	var trips []*models.Trip
	err := r.db.WithContext(ctx).Where("scooter_id = ?", scooterID).Find(&trips).Error
	return trips, err
}

// GetActiveByUserID retrieves the active trip for a user
func (r *gormTripRepository) GetActiveByUserID(ctx context.Context, userID uuid.UUID) (*models.Trip, error) {
	var trip models.Trip
	err := r.db.WithContext(ctx).
		Where("user_id = ? AND status = ?", userID, models.TripStatusActive).
		First(&trip).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &trip, nil
}

// GetActiveByScooterID retrieves the active trip for a scooter
func (r *gormTripRepository) GetActiveByScooterID(ctx context.Context, scooterID uuid.UUID) (*models.Trip, error) {
	var trip models.Trip
	err := r.db.WithContext(ctx).
		Where("scooter_id = ? AND status = ?", scooterID, models.TripStatusActive).
		First(&trip).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &trip, nil
}

// GetByDateRange retrieves trips within a date range
func (r *gormTripRepository) GetByDateRange(ctx context.Context, start, end time.Time) ([]*models.Trip, error) {
	var trips []*models.Trip
	err := r.db.WithContext(ctx).
		Where("start_time BETWEEN ? AND ?", start, end).
		Find(&trips).Error
	return trips, err
}

// GetByUserIDAndDateRange retrieves trips by user within a date range
func (r *gormTripRepository) GetByUserIDAndDateRange(ctx context.Context, userID uuid.UUID, start, end time.Time) ([]*models.Trip, error) {
	var trips []*models.Trip
	err := r.db.WithContext(ctx).
		Where("user_id = ? AND start_time BETWEEN ? AND ?", userID, start, end).
		Find(&trips).Error
	return trips, err
}

// GetByScooterIDAndDateRange retrieves trips by scooter within a date range
func (r *gormTripRepository) GetByScooterIDAndDateRange(ctx context.Context, scooterID uuid.UUID, start, end time.Time) ([]*models.Trip, error) {
	var trips []*models.Trip
	err := r.db.WithContext(ctx).
		Where("scooter_id = ? AND start_time BETWEEN ? AND ?", scooterID, start, end).
		Find(&trips).Error
	return trips, err
}

// GetTripCount returns the total number of trips
func (r *gormTripRepository) GetTripCount(ctx context.Context) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&models.Trip{}).Count(&count).Error
	return count, err
}

// GetTripCountByUser returns the number of trips for a user
func (r *gormTripRepository) GetTripCountByUser(ctx context.Context, userID uuid.UUID) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&models.Trip{}).
		Where("user_id = ?", userID).
		Count(&count).Error
	return count, err
}

// GetTripCountByScooter returns the number of trips for a scooter
func (r *gormTripRepository) GetTripCountByScooter(ctx context.Context, scooterID uuid.UUID) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&models.Trip{}).
		Where("scooter_id = ?", scooterID).
		Count(&count).Error
	return count, err
}
