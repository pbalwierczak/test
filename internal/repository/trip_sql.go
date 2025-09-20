package repository

import (
	"context"
	"database/sql"
	"time"

	"scootin-aboot/internal/models"

	"github.com/google/uuid"
)

type sqlTripRepository struct {
	db SQLExecutor
}

func (r *sqlTripRepository) Create(ctx context.Context, trip *models.Trip) error {
	trip.SetID()
	if err := trip.ValidateAndSetTimestamps(); err != nil {
		return err
	}

	query := `
		INSERT INTO trips (id, scooter_id, user_id, start_time, end_time, start_latitude, start_longitude, end_latitude, end_longitude, status, created_at, updated_at, deleted_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)`

	_, err := r.db.ExecContext(ctx, query,
		trip.ID,
		trip.ScooterID,
		trip.UserID,
		trip.StartTime,
		trip.EndTime,
		trip.StartLatitude,
		trip.StartLongitude,
		trip.EndLatitude,
		trip.EndLongitude,
		trip.Status,
		trip.CreatedAt,
		trip.UpdatedAt,
		trip.DeletedAt,
	)
	return err
}

func (r *sqlTripRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Trip, error) {
	query := `
		SELECT id, scooter_id, user_id, start_time, end_time, start_latitude, start_longitude, end_latitude, end_longitude, status, created_at, updated_at, deleted_at
		FROM trips
		WHERE id = $1 AND deleted_at IS NULL`

	trip := &models.Trip{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&trip.ID,
		&trip.ScooterID,
		&trip.UserID,
		&trip.StartTime,
		&trip.EndTime,
		&trip.StartLatitude,
		&trip.StartLongitude,
		&trip.EndLatitude,
		&trip.EndLongitude,
		&trip.Status,
		&trip.CreatedAt,
		&trip.UpdatedAt,
		&trip.DeletedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return trip, nil
}

func (r *sqlTripRepository) Update(ctx context.Context, trip *models.Trip) error {
	if err := trip.ValidateAndSetTimestamps(); err != nil {
		return err
	}

	query := `
		UPDATE trips
		SET scooter_id = $2, user_id = $3, start_time = $4, end_time = $5, start_latitude = $6, start_longitude = $7, end_latitude = $8, end_longitude = $9, status = $10, updated_at = $11, deleted_at = $12
		WHERE id = $1 AND deleted_at IS NULL`

	result, err := r.db.ExecContext(ctx, query,
		trip.ID,
		trip.ScooterID,
		trip.UserID,
		trip.StartTime,
		trip.EndTime,
		trip.StartLatitude,
		trip.StartLongitude,
		trip.EndLatitude,
		trip.EndLongitude,
		trip.Status,
		trip.UpdatedAt,
		trip.DeletedAt,
	)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrTripNotFound
	}

	return nil
}

func (r *sqlTripRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `
		UPDATE trips
		SET deleted_at = NOW()
		WHERE id = $1 AND deleted_at IS NULL`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrTripNotFound
	}

	return nil
}

func (r *sqlTripRepository) List(ctx context.Context, limit, offset int) ([]*models.Trip, error) {
	query := `
		SELECT id, scooter_id, user_id, start_time, end_time, start_latitude, start_longitude, end_latitude, end_longitude, status, created_at, updated_at, deleted_at
		FROM trips
		WHERE deleted_at IS NULL
		ORDER BY created_at DESC`

	if limit > 0 {
		query += " LIMIT $1"
		if offset > 0 {
			query += " OFFSET $2"
		}
	}

	var rows *sql.Rows
	var err error

	if limit > 0 && offset > 0 {
		rows, err = r.db.QueryContext(ctx, query, limit, offset)
	} else if limit > 0 {
		rows, err = r.db.QueryContext(ctx, query, limit)
	} else {
		rows, err = r.db.QueryContext(ctx, query)
	}

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var trips []*models.Trip
	for rows.Next() {
		trip := &models.Trip{}
		err := rows.Scan(
			&trip.ID,
			&trip.ScooterID,
			&trip.UserID,
			&trip.StartTime,
			&trip.EndTime,
			&trip.StartLatitude,
			&trip.StartLongitude,
			&trip.EndLatitude,
			&trip.EndLongitude,
			&trip.Status,
			&trip.CreatedAt,
			&trip.UpdatedAt,
			&trip.DeletedAt,
		)
		if err != nil {
			return nil, err
		}
		trips = append(trips, trip)
	}

	return trips, rows.Err()
}

func (r *sqlTripRepository) UpdateStatus(ctx context.Context, id uuid.UUID, status models.TripStatus) error {
	query := `
		UPDATE trips
		SET status = $2, updated_at = NOW()
		WHERE id = $1 AND deleted_at IS NULL`

	result, err := r.db.ExecContext(ctx, query, id, status)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrTripNotFound
	}

	return nil
}

func (r *sqlTripRepository) EndTrip(ctx context.Context, id uuid.UUID, endLat, endLng float64) error {
	now := time.Now()
	query := `
		UPDATE trips
		SET end_time = $2, end_latitude = $3, end_longitude = $4, status = $5, updated_at = NOW()
		WHERE id = $1 AND deleted_at IS NULL`

	result, err := r.db.ExecContext(ctx, query, id, now, endLat, endLng, models.TripStatusCompleted)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrTripNotFound
	}

	return nil
}

func (r *sqlTripRepository) CancelTrip(ctx context.Context, id uuid.UUID) error {
	query := `
		UPDATE trips
		SET status = $2, updated_at = NOW()
		WHERE id = $1 AND deleted_at IS NULL`

	result, err := r.db.ExecContext(ctx, query, id, models.TripStatusCancelled)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrTripNotFound
	}

	return nil
}

func (r *sqlTripRepository) GetActiveByUserID(ctx context.Context, userID uuid.UUID) (*models.Trip, error) {
	query := `
		SELECT id, scooter_id, user_id, start_time, end_time, start_latitude, start_longitude, end_latitude, end_longitude, status, created_at, updated_at, deleted_at
		FROM trips
		WHERE user_id = $1 AND status = $2 AND deleted_at IS NULL
		ORDER BY created_at DESC
		LIMIT 1`

	trip := &models.Trip{}
	err := r.db.QueryRowContext(ctx, query, userID, models.TripStatusActive).Scan(
		&trip.ID,
		&trip.ScooterID,
		&trip.UserID,
		&trip.StartTime,
		&trip.EndTime,
		&trip.StartLatitude,
		&trip.StartLongitude,
		&trip.EndLatitude,
		&trip.EndLongitude,
		&trip.Status,
		&trip.CreatedAt,
		&trip.UpdatedAt,
		&trip.DeletedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return trip, nil
}

func (r *sqlTripRepository) GetActiveByScooterID(ctx context.Context, scooterID uuid.UUID) (*models.Trip, error) {
	query := `
		SELECT id, scooter_id, user_id, start_time, end_time, start_latitude, start_longitude, end_latitude, end_longitude, status, created_at, updated_at, deleted_at
		FROM trips
		WHERE scooter_id = $1 AND status = $2 AND deleted_at IS NULL
		ORDER BY created_at DESC
		LIMIT 1`

	trip := &models.Trip{}
	err := r.db.QueryRowContext(ctx, query, scooterID, models.TripStatusActive).Scan(
		&trip.ID,
		&trip.ScooterID,
		&trip.UserID,
		&trip.StartTime,
		&trip.EndTime,
		&trip.StartLatitude,
		&trip.StartLongitude,
		&trip.EndLatitude,
		&trip.EndLongitude,
		&trip.Status,
		&trip.CreatedAt,
		&trip.UpdatedAt,
		&trip.DeletedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return trip, nil
}
