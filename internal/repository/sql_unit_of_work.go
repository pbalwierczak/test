package repository

import (
	"context"
	"database/sql"
)

type sqlUnitOfWork struct {
	db *sql.DB
}

func NewSQLUnitOfWork(db *sql.DB) UnitOfWork {
	return &sqlUnitOfWork{db: db}
}

func (u *sqlUnitOfWork) Begin(ctx context.Context) (UnitOfWorkTx, error) {
	tx, err := u.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}

	return &sqlUnitOfWorkTx{
		tx: tx,
	}, nil
}

type sqlUnitOfWorkTx struct {
	tx *sql.Tx
}

func (u *sqlUnitOfWorkTx) ScooterRepository() ScooterRepository {
	return &sqlScooterRepository{db: u.tx}
}

func (u *sqlUnitOfWorkTx) TripRepository() TripRepository {
	return &sqlTripRepository{db: u.tx}
}

func (u *sqlUnitOfWorkTx) UserRepository() UserRepository {
	return &sqlUserRepository{db: u.tx}
}

func (u *sqlUnitOfWorkTx) LocationUpdateRepository() LocationUpdateRepository {
	return &sqlLocationUpdateRepository{db: u.tx}
}

func (u *sqlUnitOfWorkTx) Commit() error {
	return u.tx.Commit()
}

func (u *sqlUnitOfWorkTx) Rollback() error {
	return u.tx.Rollback()
}
