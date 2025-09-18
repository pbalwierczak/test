package database

import (
	"fmt"
	"log"
	"time"

	"scootin-aboot/internal/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// ConnectDatabase establishes a connection to the database using the provided DSN
func ConnectDatabase(dsn string) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Test the connection
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get underlying sql.DB: %w", err)
	}

	if err := sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	// Configure connection pool settings
	sqlDB.SetMaxOpenConns(25)
	sqlDB.SetMaxIdleConns(25)
	sqlDB.SetConnMaxLifetime(time.Hour)

	log.Println("Successfully connected to database")
	return db, nil
}

// AutoMigrate runs GORM auto-migration for all models
func AutoMigrate(db *gorm.DB) error {
	// Auto-migrate all models
	err := db.AutoMigrate(
		&models.User{},
		&models.Scooter{},
		&models.Trip{},
		&models.LocationUpdate{},
	)
	if err != nil {
		return fmt.Errorf("failed to auto-migrate models: %w", err)
	}

	log.Println("Database auto-migration completed successfully")
	return nil
}
