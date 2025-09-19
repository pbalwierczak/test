package repository

import (
	"context"

	"scootin-aboot/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type gormLocationUpdateRepository struct {
	db *gorm.DB
}

func (r *gormLocationUpdateRepository) Create(ctx context.Context, update *models.LocationUpdate) error {
	return r.db.WithContext(ctx).Create(update).Error
}

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

func (r *gormLocationUpdateRepository) Update(ctx context.Context, update *models.LocationUpdate) error {
	return r.db.WithContext(ctx).Save(update).Error
}

func (r *gormLocationUpdateRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&models.LocationUpdate{}, "id = ?", id).Error
}

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

func (r *gormLocationUpdateRepository) GetByScooterID(ctx context.Context, scooterID uuid.UUID) ([]*models.LocationUpdate, error) {
	var updates []*models.LocationUpdate
	err := r.db.WithContext(ctx).Where("scooter_id = ?", scooterID).Find(&updates).Error
	return updates, err
}
