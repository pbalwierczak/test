package repository

import (
	"gorm.io/gorm"
)

// NewRepository creates a new repository instance
func NewRepository(db *gorm.DB) Repository {
	return NewGormRepository(db)
}
