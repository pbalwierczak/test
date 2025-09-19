package repository

import (
	"context"
	"errors"
	"time"

	"scootin-aboot/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type gormScooterRepository struct {
	db *gorm.DB
}

func (r *gormScooterRepository) Create(ctx context.Context, scooter *models.Scooter) error {
	return r.db.WithContext(ctx).Create(scooter).Error
}

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

func (r *gormScooterRepository) Update(ctx context.Context, scooter *models.Scooter) error {
	return r.db.WithContext(ctx).Save(scooter).Error
}

func (r *gormScooterRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&models.Scooter{}, "id = ?", id).Error
}

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

func (r *gormScooterRepository) UpdateStatus(ctx context.Context, id uuid.UUID, status models.ScooterStatus) error {
	return r.db.WithContext(ctx).Model(&models.Scooter{}).
		Where("id = ?", id).
		Update("status", status).Error
}

func (r *gormScooterRepository) UpdateLocation(ctx context.Context, id uuid.UUID, latitude, longitude float64) error {
	return r.db.WithContext(ctx).Model(&models.Scooter{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"current_latitude":  latitude,
			"current_longitude": longitude,
			"last_seen":         time.Now(),
		}).Error
}

func (r *gormScooterRepository) GetByStatus(ctx context.Context, status models.ScooterStatus) ([]*models.Scooter, error) {
	var scooters []*models.Scooter
	err := r.db.WithContext(ctx).Where("status = ?", status).Find(&scooters).Error
	return scooters, err
}

func (r *gormScooterRepository) GetInBounds(ctx context.Context, minLat, maxLat, minLng, maxLng float64) ([]*models.Scooter, error) {
	var scooters []*models.Scooter
	err := r.db.WithContext(ctx).
		Where("current_latitude BETWEEN ? AND ? AND current_longitude BETWEEN ? AND ?",
			minLat, maxLat, minLng, maxLng).
		Find(&scooters).Error
	return scooters, err
}

func (r *gormScooterRepository) GetClosest(ctx context.Context, latitude, longitude float64, limit int) ([]*models.Scooter, error) {
	return r.GetClosestWithRadius(ctx, latitude, longitude, 0, "", limit)
}

func (r *gormScooterRepository) GetClosestWithRadius(ctx context.Context, latitude, longitude, radius float64, status string, limit int) ([]*models.Scooter, error) {
	var scooters []*models.Scooter

	var query *gorm.DB
	if radius > 0 {
		bbox := NewBoundingBox(latitude, longitude, radius)

		query = r.db.WithContext(ctx).Model(&models.Scooter{}).
			Where("current_latitude BETWEEN ? AND ?", bbox.MinLat, bbox.MaxLat).
			Where("current_longitude BETWEEN ? AND ?", bbox.MinLng, bbox.MaxLng)
	} else {
		query = r.db.WithContext(ctx).Model(&models.Scooter{})
	}

	if status != "" {
		query = query.Where("status = ?", status)
	}

	err := query.Find(&scooters).Error
	if err != nil {
		return nil, err
	}

	filteredScooters := FilterAndSortByDistance(scooters, latitude, longitude, radius, limit)

	return filteredScooters, nil
}

func (r *gormScooterRepository) GetByStatusInBounds(ctx context.Context, status models.ScooterStatus, minLat, maxLat, minLng, maxLng float64) ([]*models.Scooter, error) {
	var scooters []*models.Scooter
	err := r.db.WithContext(ctx).
		Where("status = ? AND current_latitude BETWEEN ? AND ? AND current_longitude BETWEEN ? AND ?",
			status, minLat, maxLat, minLng, maxLng).
		Find(&scooters).Error
	return scooters, err
}

func (r *gormScooterRepository) GetByIDForUpdate(ctx context.Context, id uuid.UUID) (*models.Scooter, error) {
	var scooter models.Scooter
	err := r.db.WithContext(ctx).Set("gorm:query_option", "FOR UPDATE").
		First(&scooter, "id = ?", id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &scooter, nil
}

func (r *gormScooterRepository) UpdateStatusWithCheck(ctx context.Context, id uuid.UUID, newStatus models.ScooterStatus, expectedStatus models.ScooterStatus) error {
	result := r.db.WithContext(ctx).Model(&models.Scooter{}).
		Where("id = ? AND status = ?", id, expectedStatus).
		Update("status", newStatus)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("scooter status update failed: status mismatch or scooter not found")
	}

	return nil
}
