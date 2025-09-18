package database

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

// SeedData loads seed data from SQL files
func SeedData(db *sql.DB, seedsPath string) error {
	seedFiles := []string{
		"users.sql",
		"scooters.sql",
		"sample_trips.sql",
	}

	for _, filename := range seedFiles {
		filePath := filepath.Join(seedsPath, filename)

		// Check if file exists
		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			log.Printf("Warning: Seed file not found: %s", filePath)
			continue
		}

		// Read SQL file
		sqlContent, err := ioutil.ReadFile(filePath)
		if err != nil {
			return fmt.Errorf("failed to read seed file %s: %w", filename, err)
		}

		// Execute SQL
		if _, err := db.Exec(string(sqlContent)); err != nil {
			return fmt.Errorf("failed to execute seed file %s: %w", filename, err)
		}

		log.Printf("Successfully loaded seed data from %s", filename)
	}

	log.Println("All seed data loaded successfully")
	return nil
}

// GetSeedsPath returns the absolute path to the seeds directory
func GetSeedsPath() (string, error) {
	wd, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("failed to get working directory: %w", err)
	}

	seedsPath := filepath.Join(wd, "seeds")
	if _, err := os.Stat(seedsPath); os.IsNotExist(err) {
		return "", fmt.Errorf("seeds directory not found: %s", seedsPath)
	}

	return seedsPath, nil
}
