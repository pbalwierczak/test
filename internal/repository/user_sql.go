package repository

import (
	"context"
	"database/sql"

	"scootin-aboot/internal/models"

	"github.com/google/uuid"
)

type sqlUserRepository struct {
	db SQLExecutor
}

func (r *sqlUserRepository) Create(ctx context.Context, user *models.User) error {
	user.SetID()
	user.SetTimestamps()

	query := `
		INSERT INTO users (id, created_at, updated_at, deleted_at)
		VALUES ($1, $2, $3, $4)`

	_, err := r.db.ExecContext(ctx, query, user.ID, user.CreatedAt, user.UpdatedAt, user.DeletedAt)
	return err
}

func (r *sqlUserRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.User, error) {
	query := `
		SELECT id, created_at, updated_at, deleted_at
		FROM users
		WHERE id = $1 AND deleted_at IS NULL`

	user := &models.User{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&user.ID,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.DeletedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return user, nil
}

func (r *sqlUserRepository) Update(ctx context.Context, user *models.User) error {
	user.SetTimestamps()

	query := `
		UPDATE users
		SET updated_at = $2, deleted_at = $3
		WHERE id = $1 AND deleted_at IS NULL`

	result, err := r.db.ExecContext(ctx, query, user.ID, user.UpdatedAt, user.DeletedAt)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrUserNotFound
	}

	return nil
}

func (r *sqlUserRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `
		UPDATE users
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
		return ErrUserNotFound
	}

	return nil
}

func (r *sqlUserRepository) List(ctx context.Context, limit, offset int) ([]*models.User, error) {
	query := `
		SELECT id, created_at, updated_at, deleted_at
		FROM users
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

	var users []*models.User
	for rows.Next() {
		user := &models.User{}
		err := rows.Scan(
			&user.ID,
			&user.CreatedAt,
			&user.UpdatedAt,
			&user.DeletedAt,
		)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, rows.Err()
}
