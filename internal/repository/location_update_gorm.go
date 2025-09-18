package repository

import (
	"context"
	"fmt"
	"time"

	"scootin-aboot/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// gormLocationUpdateRepository implements LocationUpdateRepository using GORM
type gormLocationUpdateRepository struct {
	db *gorm.DB
}

// Create creates a new location update
func (r *gormLocationUpdateRepository) Create(ctx context.Context, update *models.LocationUpdate) error {
	return r.db.WithContext(ctx).Create(update).Error
}

// GetByID retrieves a location update by ID
func (r *gormLocationUpdateRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.LocationUpdate, error) {
	var update models.LocationUpdate
	err := r.db.WithContext(ctx).First(&update, "id = ?", id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &update, nil
}

// Update updates a location update
func (r *gormLocationUpdateRepository) Update(ctx context.Context, update *models.LocationUpdate) error {
	return r.db.WithContext(ctx).Save(update).Error
}

// Delete deletes a location update by ID
func (r *gormLocationUpdateRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&models.LocationUpdate{}, "id = ?", id).Error
}

// List retrieves location updates with pagination
func (r *gormLocationUpdateRepository) List(ctx context.Context, limit, offset int) ([]*models.LocationUpdate, error) {
	var updates []*models.LocationUpdate
	query := r.db.WithContext(ctx)

	if limit > 0 {
		query = query.Limit(limit)
	}
	if offset > 0 {
		query = query.Offset(offset)
	}

	err := query.Find(&updates).Error
	return updates, err
}

// GetByTripID retrieves location updates for a trip
func (r *gormLocationUpdateRepository) GetByTripID(ctx context.Context, tripID uuid.UUID) ([]*models.LocationUpdate, error) {
	var updates []*models.LocationUpdate
	err := r.db.WithContext(ctx).Where("trip_id = ?", tripID).Find(&updates).Error
	return updates, err
}

// GetByTripIDOrdered retrieves location updates for a trip ordered by timestamp
func (r *gormLocationUpdateRepository) GetByTripIDOrdered(ctx context.Context, tripID uuid.UUID) ([]*models.LocationUpdate, error) {
	var updates []*models.LocationUpdate
	err := r.db.WithContext(ctx).
		Where("trip_id = ?", tripID).
		Order("timestamp ASC").
		Find(&updates).Error
	return updates, err
}

// GetLatestByTripID retrieves the latest location update for a trip
func (r *gormLocationUpdateRepository) GetLatestByTripID(ctx context.Context, tripID uuid.UUID) (*models.LocationUpdate, error) {
	var update models.LocationUpdate
	err := r.db.WithContext(ctx).
		Where("trip_id = ?", tripID).
		Order("timestamp DESC").
		First(&update).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &update, nil
}

// GetByDateRange retrieves location updates within a date range
func (r *gormLocationUpdateRepository) GetByDateRange(ctx context.Context, start, end time.Time) ([]*models.LocationUpdate, error) {
	var updates []*models.LocationUpdate
	err := r.db.WithContext(ctx).
		Where("timestamp BETWEEN ? AND ?", start, end).
		Find(&updates).Error
	return updates, err
}

// GetByTripIDAndDateRange retrieves location updates for a trip within a date range
func (r *gormLocationUpdateRepository) GetByTripIDAndDateRange(ctx context.Context, tripID uuid.UUID, start, end time.Time) ([]*models.LocationUpdate, error) {
	var updates []*models.LocationUpdate
	err := r.db.WithContext(ctx).
		Where("trip_id = ? AND timestamp BETWEEN ? AND ?", tripID, start, end).
		Find(&updates).Error
	return updates, err
}

// GetInBounds retrieves location updates within geographic bounds
func (r *gormLocationUpdateRepository) GetInBounds(ctx context.Context, minLat, maxLat, minLng, maxLng float64) ([]*models.LocationUpdate, error) {
	var updates []*models.LocationUpdate
	err := r.db.WithContext(ctx).
		Where("latitude BETWEEN ? AND ? AND longitude BETWEEN ? AND ?",
			minLat, maxLat, minLng, maxLng).
		Find(&updates).Error
	return updates, err
}

// GetInRadius retrieves location updates within a radius of a given location
func (r *gormLocationUpdateRepository) GetInRadius(ctx context.Context, latitude, longitude, radiusKm float64) ([]*models.LocationUpdate, error) {
	var updates []*models.LocationUpdate

	// Calculate distance using Haversine formula
	distanceQuery := fmt.Sprintf(`
		(6371 * acos(
			cos(radians(%f)) * 
			cos(radians(latitude)) * 
			cos(radians(longitude) - radians(%f)) + 
			sin(radians(%f)) * 
			sin(radians(latitude))
		)) AS distance`,
		latitude, longitude, latitude)

	err := r.db.WithContext(ctx).
		Select("*, "+distanceQuery).
		Having("distance <= ?", radiusKm).
		Order("distance ASC").
		Find(&updates).Error
	return updates, err
}

// GetUpdateCount returns the total number of location updates
func (r *gormLocationUpdateRepository) GetUpdateCount(ctx context.Context) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&models.LocationUpdate{}).Count(&count).Error
	return count, err
}

// GetUpdateCountByTrip returns the number of location updates for a trip
func (r *gormLocationUpdateRepository) GetUpdateCountByTrip(ctx context.Context, tripID uuid.UUID) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&models.LocationUpdate{}).
		Where("trip_id = ?", tripID).
		Count(&count).Error
	return count, err
}
