package repository

import (
	"context"
	"database/sql"

	"scootin-aboot/internal/models"

	"github.com/google/uuid"
)

type sqlScooterRepository struct {
	db SQLExecutor
}

func (r *sqlScooterRepository) Create(ctx context.Context, scooter *models.Scooter) error {
	scooter.SetID()
	if err := scooter.ValidateAndSetTimestamps(); err != nil {
		return err
	}

	query := `
		INSERT INTO scooters (id, status, current_latitude, current_longitude, created_at, updated_at, last_seen, deleted_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`

	_, err := r.db.ExecContext(ctx, query,
		scooter.ID,
		scooter.Status,
		scooter.CurrentLatitude,
		scooter.CurrentLongitude,
		scooter.CreatedAt,
		scooter.UpdatedAt,
		scooter.LastSeen,
		scooter.DeletedAt,
	)
	return err
}

func (r *sqlScooterRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Scooter, error) {
	query := `
		SELECT id, status, current_latitude, current_longitude, created_at, updated_at, last_seen, deleted_at
		FROM scooters
		WHERE id = $1 AND deleted_at IS NULL`

	scooter := &models.Scooter{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&scooter.ID,
		&scooter.Status,
		&scooter.CurrentLatitude,
		&scooter.CurrentLongitude,
		&scooter.CreatedAt,
		&scooter.UpdatedAt,
		&scooter.LastSeen,
		&scooter.DeletedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrScooterNotFound
		}
		return nil, err
	}

	return scooter, nil
}

func (r *sqlScooterRepository) Update(ctx context.Context, scooter *models.Scooter) error {
	if err := scooter.ValidateAndSetTimestamps(); err != nil {
		return err
	}

	query := `
		UPDATE scooters
		SET status = $2, current_latitude = $3, current_longitude = $4, updated_at = $5, last_seen = $6, deleted_at = $7
		WHERE id = $1 AND deleted_at IS NULL`

	result, err := r.db.ExecContext(ctx, query,
		scooter.ID,
		scooter.Status,
		scooter.CurrentLatitude,
		scooter.CurrentLongitude,
		scooter.UpdatedAt,
		scooter.LastSeen,
		scooter.DeletedAt,
	)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrScooterNotFound
	}

	return nil
}

func (r *sqlScooterRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `
		UPDATE scooters
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
		return ErrScooterNotFound
	}

	return nil
}

func (r *sqlScooterRepository) List(ctx context.Context, limit, offset int) ([]*models.Scooter, error) {
	query := `
		SELECT id, status, current_latitude, current_longitude, created_at, updated_at, last_seen, deleted_at
		FROM scooters
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

	var scooters []*models.Scooter
	for rows.Next() {
		scooter := &models.Scooter{}
		err := rows.Scan(
			&scooter.ID,
			&scooter.Status,
			&scooter.CurrentLatitude,
			&scooter.CurrentLongitude,
			&scooter.CreatedAt,
			&scooter.UpdatedAt,
			&scooter.LastSeen,
			&scooter.DeletedAt,
		)
		if err != nil {
			return nil, err
		}
		scooters = append(scooters, scooter)
	}

	return scooters, rows.Err()
}

func (r *sqlScooterRepository) UpdateStatus(ctx context.Context, id uuid.UUID, status models.ScooterStatus) error {
	query := `
		UPDATE scooters
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
		return ErrScooterNotFound
	}

	return nil
}

func (r *sqlScooterRepository) UpdateLocation(ctx context.Context, id uuid.UUID, latitude, longitude float64) error {
	query := `
		UPDATE scooters
		SET current_latitude = $2, current_longitude = $3, last_seen = NOW(), updated_at = NOW()
		WHERE id = $1 AND deleted_at IS NULL`

	result, err := r.db.ExecContext(ctx, query, id, latitude, longitude)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrScooterNotFound
	}

	return nil
}

func (r *sqlScooterRepository) GetByStatus(ctx context.Context, status models.ScooterStatus) ([]*models.Scooter, error) {
	query := `
		SELECT id, status, current_latitude, current_longitude, created_at, updated_at, last_seen, deleted_at
		FROM scooters
		WHERE status = $1 AND deleted_at IS NULL
		ORDER BY created_at DESC`

	rows, err := r.db.QueryContext(ctx, query, status)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var scooters []*models.Scooter
	for rows.Next() {
		scooter := &models.Scooter{}
		err := rows.Scan(
			&scooter.ID,
			&scooter.Status,
			&scooter.CurrentLatitude,
			&scooter.CurrentLongitude,
			&scooter.CreatedAt,
			&scooter.UpdatedAt,
			&scooter.LastSeen,
			&scooter.DeletedAt,
		)
		if err != nil {
			return nil, err
		}
		scooters = append(scooters, scooter)
	}

	return scooters, rows.Err()
}

func (r *sqlScooterRepository) GetInBounds(ctx context.Context, minLat, maxLat, minLng, maxLng float64) ([]*models.Scooter, error) {
	query := `
		SELECT id, status, current_latitude, current_longitude, created_at, updated_at, last_seen, deleted_at
		FROM scooters
		WHERE current_latitude BETWEEN $1 AND $2
		AND current_longitude BETWEEN $3 AND $4
		AND deleted_at IS NULL
		ORDER BY created_at DESC`

	rows, err := r.db.QueryContext(ctx, query, minLat, maxLat, minLng, maxLng)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var scooters []*models.Scooter
	for rows.Next() {
		scooter := &models.Scooter{}
		err := rows.Scan(
			&scooter.ID,
			&scooter.Status,
			&scooter.CurrentLatitude,
			&scooter.CurrentLongitude,
			&scooter.CreatedAt,
			&scooter.UpdatedAt,
			&scooter.LastSeen,
			&scooter.DeletedAt,
		)
		if err != nil {
			return nil, err
		}
		scooters = append(scooters, scooter)
	}

	return scooters, rows.Err()
}

func (r *sqlScooterRepository) GetClosest(ctx context.Context, latitude, longitude float64, limit int) ([]*models.Scooter, error) {
	// Using Haversine formula for distance calculation
	query := `
		SELECT id, status, current_latitude, current_longitude, created_at, updated_at, last_seen, deleted_at,
		(6371 * acos(cos(radians($1)) * cos(radians(current_latitude)) * 
		cos(radians(current_longitude) - radians($2)) + 
		sin(radians($1)) * sin(radians(current_latitude)))) AS distance
		FROM scooters
		WHERE deleted_at IS NULL
		ORDER BY distance
		LIMIT $3`

	rows, err := r.db.QueryContext(ctx, query, latitude, longitude, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var scooters []*models.Scooter
	for rows.Next() {
		scooter := &models.Scooter{}
		var distance float64
		err := rows.Scan(
			&scooter.ID,
			&scooter.Status,
			&scooter.CurrentLatitude,
			&scooter.CurrentLongitude,
			&scooter.CreatedAt,
			&scooter.UpdatedAt,
			&scooter.LastSeen,
			&scooter.DeletedAt,
			&distance,
		)
		if err != nil {
			return nil, err
		}
		scooters = append(scooters, scooter)
	}

	return scooters, rows.Err()
}

func (r *sqlScooterRepository) GetClosestWithRadius(ctx context.Context, latitude, longitude, radius float64, status string, limit int) ([]*models.Scooter, error) {
	// Using Haversine formula for distance calculation with radius filtering
	query := `
		SELECT id, status, current_latitude, current_longitude, created_at, updated_at, last_seen, deleted_at,
		(6371 * acos(cos(radians($1)) * cos(radians(current_latitude)) * 
		cos(radians(current_longitude) - radians($2)) + 
		sin(radians($1)) * sin(radians(current_latitude)))) AS distance
		FROM scooters
		WHERE deleted_at IS NULL
		AND status = $4
		AND (6371 * acos(cos(radians($1)) * cos(radians(current_latitude)) * 
		cos(radians(current_longitude) - radians($2)) + 
		sin(radians($1)) * sin(radians(current_latitude)))) <= $3
		ORDER BY distance
		LIMIT $5`

	rows, err := r.db.QueryContext(ctx, query, latitude, longitude, radius, status, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var scooters []*models.Scooter
	for rows.Next() {
		scooter := &models.Scooter{}
		var distance float64
		err := rows.Scan(
			&scooter.ID,
			&scooter.Status,
			&scooter.CurrentLatitude,
			&scooter.CurrentLongitude,
			&scooter.CreatedAt,
			&scooter.UpdatedAt,
			&scooter.LastSeen,
			&scooter.DeletedAt,
			&distance,
		)
		if err != nil {
			return nil, err
		}
		scooters = append(scooters, scooter)
	}

	return scooters, rows.Err()
}

func (r *sqlScooterRepository) GetByStatusInBounds(ctx context.Context, status models.ScooterStatus, minLat, maxLat, minLng, maxLng float64) ([]*models.Scooter, error) {
	query := `
		SELECT id, status, current_latitude, current_longitude, created_at, updated_at, last_seen, deleted_at
		FROM scooters
		WHERE status = $1
		AND current_latitude BETWEEN $2 AND $3
		AND current_longitude BETWEEN $4 AND $5
		AND deleted_at IS NULL
		ORDER BY created_at DESC`

	rows, err := r.db.QueryContext(ctx, query, status, minLat, maxLat, minLng, maxLng)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var scooters []*models.Scooter
	for rows.Next() {
		scooter := &models.Scooter{}
		err := rows.Scan(
			&scooter.ID,
			&scooter.Status,
			&scooter.CurrentLatitude,
			&scooter.CurrentLongitude,
			&scooter.CreatedAt,
			&scooter.UpdatedAt,
			&scooter.LastSeen,
			&scooter.DeletedAt,
		)
		if err != nil {
			return nil, err
		}
		scooters = append(scooters, scooter)
	}

	return scooters, rows.Err()
}

func (r *sqlScooterRepository) GetByIDForUpdate(ctx context.Context, id uuid.UUID) (*models.Scooter, error) {
	query := `
		SELECT id, status, current_latitude, current_longitude, created_at, updated_at, last_seen, deleted_at
		FROM scooters
		WHERE id = $1 AND deleted_at IS NULL
		FOR UPDATE`

	scooter := &models.Scooter{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&scooter.ID,
		&scooter.Status,
		&scooter.CurrentLatitude,
		&scooter.CurrentLongitude,
		&scooter.CreatedAt,
		&scooter.UpdatedAt,
		&scooter.LastSeen,
		&scooter.DeletedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrScooterNotFound
		}
		return nil, err
	}

	return scooter, nil
}

func (r *sqlScooterRepository) UpdateStatusWithCheck(ctx context.Context, id uuid.UUID, newStatus models.ScooterStatus, expectedStatus models.ScooterStatus) error {
	query := `
		UPDATE scooters
		SET status = $2, updated_at = NOW()
		WHERE id = $1 AND status = $3 AND deleted_at IS NULL`

	result, err := r.db.ExecContext(ctx, query, id, newStatus, expectedStatus)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrScooterNotFound
	}

	return nil
}
