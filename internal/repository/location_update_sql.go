package repository

import (
	"context"
	"database/sql"

	"scootin-aboot/internal/models"

	"github.com/google/uuid"
)

type sqlLocationUpdateRepository struct {
	db SQLExecutor
}

func (r *sqlLocationUpdateRepository) Create(ctx context.Context, update *models.LocationUpdate) error {
	update.SetID()
	if err := update.ValidateAndSetTimestamps(); err != nil {
		return err
	}

	query := `
		INSERT INTO location_updates (id, scooter_id, latitude, longitude, timestamp, created_at, deleted_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)`

	_, err := r.db.ExecContext(ctx, query,
		update.ID,
		update.ScooterID,
		update.Latitude,
		update.Longitude,
		update.Timestamp,
		update.CreatedAt,
		update.DeletedAt,
	)
	return err
}

func (r *sqlLocationUpdateRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.LocationUpdate, error) {
	query := `
		SELECT id, scooter_id, latitude, longitude, timestamp, created_at, deleted_at
		FROM location_updates
		WHERE id = $1 AND deleted_at IS NULL`

	update := &models.LocationUpdate{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&update.ID,
		&update.ScooterID,
		&update.Latitude,
		&update.Longitude,
		&update.Timestamp,
		&update.CreatedAt,
		&update.DeletedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return update, nil
}

func (r *sqlLocationUpdateRepository) Update(ctx context.Context, update *models.LocationUpdate) error {
	if err := update.ValidateAndSetTimestamps(); err != nil {
		return err
	}

	query := `
		UPDATE location_updates
		SET scooter_id = $2, latitude = $3, longitude = $4, timestamp = $5, deleted_at = $6
		WHERE id = $1 AND deleted_at IS NULL`

	result, err := r.db.ExecContext(ctx, query,
		update.ID,
		update.ScooterID,
		update.Latitude,
		update.Longitude,
		update.Timestamp,
		update.DeletedAt,
	)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrLocationUpdateNotFound
	}

	return nil
}

func (r *sqlLocationUpdateRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `
		UPDATE location_updates
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
		return ErrLocationUpdateNotFound
	}

	return nil
}

func (r *sqlLocationUpdateRepository) List(ctx context.Context, limit, offset int) ([]*models.LocationUpdate, error) {
	query := `
		SELECT id, scooter_id, latitude, longitude, timestamp, created_at, deleted_at
		FROM location_updates
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

	var updates []*models.LocationUpdate
	for rows.Next() {
		update := &models.LocationUpdate{}
		err := rows.Scan(
			&update.ID,
			&update.ScooterID,
			&update.Latitude,
			&update.Longitude,
			&update.Timestamp,
			&update.CreatedAt,
			&update.DeletedAt,
		)
		if err != nil {
			return nil, err
		}
		updates = append(updates, update)
	}

	return updates, rows.Err()
}

func (r *sqlLocationUpdateRepository) GetByScooterID(ctx context.Context, scooterID uuid.UUID) ([]*models.LocationUpdate, error) {
	query := `
		SELECT id, scooter_id, latitude, longitude, timestamp, created_at, deleted_at
		FROM location_updates
		WHERE scooter_id = $1 AND deleted_at IS NULL
		ORDER BY timestamp DESC`

	rows, err := r.db.QueryContext(ctx, query, scooterID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var updates []*models.LocationUpdate
	for rows.Next() {
		update := &models.LocationUpdate{}
		err := rows.Scan(
			&update.ID,
			&update.ScooterID,
			&update.Latitude,
			&update.Longitude,
			&update.Timestamp,
			&update.CreatedAt,
			&update.DeletedAt,
		)
		if err != nil {
			return nil, err
		}
		updates = append(updates, update)
	}

	return updates, rows.Err()
}
