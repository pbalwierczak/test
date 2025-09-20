package repository

import (
	"context"
	"database/sql"
)

// SQLExecutor is a common interface for both *sql.DB and *sql.Tx
type SQLExecutor interface {
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row
}

// Ensure *sql.DB implements SQLExecutor
var _ SQLExecutor = (*sql.DB)(nil)

// Ensure *sql.Tx implements SQLExecutor
var _ SQLExecutor = (*sql.Tx)(nil)
