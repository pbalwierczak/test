package repository

import (
	"context"
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
	return r.GetClosestWithRadius(ctx, latitude, longitude, 0, "", limit)
}

// GetClosestWithRadius retrieves the closest scooters to a given location within a radius
func (r *gormScooterRepository) GetClosestWithRadius(ctx context.Context, latitude, longitude, radius float64, status string, limit int) ([]*models.Scooter, error) {
	var scooters []*models.Scooter

	// If no radius specified or negative radius, get all scooters
	var query *gorm.DB
	if radius > 0 {
		// Create bounding box for efficient database query using indexes
		bbox := NewBoundingBox(latitude, longitude, radius)

		query = r.db.WithContext(ctx).Model(&models.Scooter{}).
			Where("current_latitude BETWEEN ? AND ?", bbox.MinLat, bbox.MaxLat).
			Where("current_longitude BETWEEN ? AND ?", bbox.MinLng, bbox.MaxLng)
	} else {
		// No radius filter, get all scooters
		query = r.db.WithContext(ctx).Model(&models.Scooter{})
	}

	// Add status filter if specified
	if status != "" {
		query = query.Where("status = ?", status)
	}

	// Get all scooters in the bounding box
	err := query.Find(&scooters).Error
	if err != nil {
		return nil, err
	}

	// Filter by radius, sort by distance, and apply limit using the factored function
	filteredScooters := FilterAndSortByDistance(scooters, latitude, longitude, radius, limit)

	return filteredScooters, nil
}

// GetInRadius retrieves scooters within a radius of a given location
func (r *gormScooterRepository) GetInRadius(ctx context.Context, latitude, longitude, radiusKm float64) ([]*models.Scooter, error) {
	return r.GetClosestWithRadius(ctx, latitude, longitude, radiusKm, "", 0)
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
