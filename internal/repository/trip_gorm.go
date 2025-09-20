package repository

import (
	"context"
	"time"

	"scootin-aboot/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type gormTripRepository struct {
	db *gorm.DB
}

func (r *gormTripRepository) Create(ctx context.Context, trip *models.Trip) error {
	return r.db.WithContext(ctx).Create(trip).Error
}

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

func (r *gormTripRepository) Update(ctx context.Context, trip *models.Trip) error {
	return r.db.WithContext(ctx).Save(trip).Error
}

func (r *gormTripRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&models.Trip{}, "id = ?", id).Error
}

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

func (r *gormTripRepository) UpdateStatus(ctx context.Context, id uuid.UUID, status models.TripStatus) error {
	return r.db.WithContext(ctx).Model(&models.Trip{}).
		Where("id = ?", id).
		Update("status", status).Error
}

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

func (r *gormTripRepository) CancelTrip(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Model(&models.Trip{}).
		Where("id = ?", id).
		Update("status", models.TripStatusCancelled).Error
}

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
