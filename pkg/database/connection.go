package database

import (
	"context"
	"fmt"
	"log"
	"time"

	"scootin-aboot/internal/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

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

	sqlDB.SetMaxOpenConns(10)
	sqlDB.SetMaxIdleConns(5)
	sqlDB.SetConnMaxLifetime(0)
	sqlDB.SetConnMaxIdleTime(0)

	log.Println("Successfully connected to database")
	return db, nil
}

func AutoMigrate(db *gorm.DB) error {
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

func StartHealthCheck(db *gorm.DB, interval time.Duration) func() {
	stopChan := make(chan struct{})

	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				if err := checkConnection(db); err != nil {
					log.Printf("Database health check failed: %v", err)
				}
			case <-stopChan:
				log.Println("Database health check stopped")
				return
			}
		}
	}()

	return func() {
		close(stopChan)
	}
}

func checkConnection(db *gorm.DB) error {
	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("failed to get underlying sql.DB: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := sqlDB.PingContext(ctx); err != nil {
		return fmt.Errorf("database ping failed: %w", err)
	}

	return nil
}
