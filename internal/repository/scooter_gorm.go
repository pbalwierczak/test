package repository

import (
	"context"
	"fmt"
	"time"

	"scootin-aboot/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// gormScooterRepository implements ScooterRepository using GORM
type gormScooterRepository struct {
	db *gorm.DB
}

// Create creates a new scooter
func (r *gormScooterRepository) Create(ctx context.Context, scooter *models.Scooter) error {
	return r.db.WithContext(ctx).Create(scooter).Error
}

// GetByID retrieves a scooter by ID
func (r *gormScooterRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Scooter, error) {
	var scooter models.Scooter
	err := r.db.WithContext(ctx).First(&scooter, "id = ?", id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &scooter, nil
}

// Update updates a scooter
func (r *gormScooterRepository) Update(ctx context.Context, scooter *models.Scooter) error {
	return r.db.WithContext(ctx).Save(scooter).Error
}

// Delete deletes a scooter by ID
func (r *gormScooterRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&models.Scooter{}, "id = ?", id).Error
}

// List retrieves scooters with pagination
func (r *gormScooterRepository) List(ctx context.Context, limit, offset int) ([]*models.Scooter, error) {
	var scooters []*models.Scooter
	query := r.db.WithContext(ctx)

	if limit > 0 {
		query = query.Limit(limit)
	}
	if offset > 0 {
		query = query.Offset(offset)
	}

	err := query.Find(&scooters).Error
	return scooters, err
}

// UpdateStatus updates a scooter's status
func (r *gormScooterRepository) UpdateStatus(ctx context.Context, id uuid.UUID, status models.ScooterStatus) error {
	return r.db.WithContext(ctx).Model(&models.Scooter{}).
		Where("id = ?", id).
		Update("status", status).Error
}

// UpdateLocation updates a scooter's location
func (r *gormScooterRepository) UpdateLocation(ctx context.Context, id uuid.UUID, latitude, longitude float64) error {
	return r.db.WithContext(ctx).Model(&models.Scooter{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"current_latitude":  latitude,
			"current_longitude": longitude,
			"last_seen":         time.Now(),
		}).Error
}

// GetByStatus retrieves scooters by status
func (r *gormScooterRepository) GetByStatus(ctx context.Context, status models.ScooterStatus) ([]*models.Scooter, error) {
	var scooters []*models.Scooter
	err := r.db.WithContext(ctx).Where("status = ?", status).Find(&scooters).Error
	return scooters, err
}

// GetAvailable retrieves available scooters
func (r *gormScooterRepository) GetAvailable(ctx context.Context) ([]*models.Scooter, error) {
	return r.GetByStatus(ctx, models.ScooterStatusAvailable)
}

// GetOccupied retrieves occupied scooters
func (r *gormScooterRepository) GetOccupied(ctx context.Context) ([]*models.Scooter, error) {
	return r.GetByStatus(ctx, models.ScooterStatusOccupied)
}

// GetInBounds retrieves scooters within geographic bounds
func (r *gormScooterRepository) GetInBounds(ctx context.Context, minLat, maxLat, minLng, maxLng float64) ([]*models.Scooter, error) {
	var scooters []*models.Scooter
	err := r.db.WithContext(ctx).
		Where("current_latitude BETWEEN ? AND ? AND current_longitude BETWEEN ? AND ?",
			minLat, maxLat, minLng, maxLng).
		Find(&scooters).Error
	return scooters, err
}

// GetClosest retrieves the closest scooters to a given location
func (r *gormScooterRepository) GetClosest(ctx context.Context, latitude, longitude float64, limit int) ([]*models.Scooter, error) {
	var scooters []*models.Scooter

	// Calculate distance using Haversine formula in SQL
	// This is a simplified version - for production, consider using PostGIS
	distanceQuery := fmt.Sprintf(`
		(6371 * acos(
			cos(radians(%f)) * 
			cos(radians(current_latitude)) * 
			cos(radians(current_longitude) - radians(%f)) + 
			sin(radians(%f)) * 
			sin(radians(current_latitude))
		)) AS distance`,
		latitude, longitude, latitude)

	query := r.db.WithContext(ctx).
		Select("*, " + distanceQuery).
		Order("distance ASC")

	if limit > 0 {
		query = query.Limit(limit)
	}

	err := query.Find(&scooters).Error
	return scooters, err
}

// GetInRadius retrieves scooters within a radius of a given location
func (r *gormScooterRepository) GetInRadius(ctx context.Context, latitude, longitude, radiusKm float64) ([]*models.Scooter, error) {
	var scooters []*models.Scooter

	// Calculate distance using Haversine formula
	distanceQuery := fmt.Sprintf(`
		(6371 * acos(
			cos(radians(%f)) * 
			cos(radians(current_latitude)) * 
			cos(radians(current_longitude) - radians(%f)) + 
			sin(radians(%f)) * 
			sin(radians(current_latitude))
		)) AS distance`,
		latitude, longitude, latitude)

	err := r.db.WithContext(ctx).
		Select("*, "+distanceQuery).
		Having("distance <= ?", radiusKm).
		Order("distance ASC").
		Find(&scooters).Error
	return scooters, err
}

// GetAvailableInBounds retrieves available scooters within geographic bounds
func (r *gormScooterRepository) GetAvailableInBounds(ctx context.Context, minLat, maxLat, minLng, maxLng float64) ([]*models.Scooter, error) {
	var scooters []*models.Scooter
	err := r.db.WithContext(ctx).
		Where("status = ? AND current_latitude BETWEEN ? AND ? AND current_longitude BETWEEN ? AND ?",
			models.ScooterStatusAvailable, minLat, maxLat, minLng, maxLng).
		Find(&scooters).Error
	return scooters, err
}

// GetByStatusInBounds retrieves scooters by status within geographic bounds
func (r *gormScooterRepository) GetByStatusInBounds(ctx context.Context, status models.ScooterStatus, minLat, maxLat, minLng, maxLng float64) ([]*models.Scooter, error) {
	var scooters []*models.Scooter
	err := r.db.WithContext(ctx).
		Where("status = ? AND current_latitude BETWEEN ? AND ? AND current_longitude BETWEEN ? AND ?",
			status, minLat, maxLat, minLng, maxLng).
		Find(&scooters).Error
	return scooters, err
}
