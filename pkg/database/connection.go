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

	// Configure connection pool settings for persistent connections
	sqlDB.SetMaxOpenConns(10)
	sqlDB.SetMaxIdleConns(5)
	// Set to 0 to disable connection lifetime limit (connections stay open indefinitely)
	sqlDB.SetConnMaxLifetime(0)
	// Set to 0 to disable idle connection timeout (idle connections stay open indefinitely)
	sqlDB.SetConnMaxIdleTime(0)

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

// StartHealthCheck starts a periodic health check for the database connection
// Returns a stop function that can be called to gracefully stop the health check
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

	// Return stop function
	return func() {
		close(stopChan)
	}
}

// checkConnection performs a simple ping to verify the database connection is alive
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

// ReconnectDatabase attempts to reconnect to the database if the connection is lost
func ReconnectDatabase(dsn string, maxRetries int, retryDelay time.Duration) (*gorm.DB, error) {
	var lastErr error

	for i := 0; i < maxRetries; i++ {
		db, err := ConnectDatabase(dsn)
		if err == nil {
			log.Printf("Successfully reconnected to database after %d attempts", i+1)
			return db, nil
		}

		lastErr = err
		log.Printf("Reconnection attempt %d/%d failed: %v", i+1, maxRetries, err)

		if i < maxRetries-1 {
			time.Sleep(retryDelay)
		}
	}

	return nil, fmt.Errorf("failed to reconnect after %d attempts: %w", maxRetries, lastErr)
}
