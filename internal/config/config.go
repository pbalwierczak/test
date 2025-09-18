package config

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

// Geographic constants for Ottawa and Montreal
const (
	OttawaCenterLat   = 45.4215
	OttawaCenterLng   = -75.6972
	MontrealCenterLat = 45.5017
	MontrealCenterLng = -73.5673
	CityRadiusKm      = 10
)

// Config holds all configuration for our application
type Config struct {
	// API Configuration
	APIKey     string
	ServerPort string
	ServerHost string

	// Database Configuration
	DBHost     string
	DBPort     string
	DBName     string
	DBUser     string
	DBPassword string
	DBSSLMode  string

	// Simulator Configuration
	SimulatorScooters        int
	SimulatorUsers           int
	SimulatorServerURL       string
	SimulatorSpeed           int
	SimulatorTripDurationMin int
	SimulatorTripDurationMax int
	SimulatorRestMin         int
	SimulatorRestMax         int

	// Logging
	LogLevel  string
	LogFormat string
}

// Load loads configuration from environment variables
func Load() (*Config, error) {
	// Load .env file if it exists
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: .env file not found: %v", err)
	}

	config := &Config{
		// API Configuration
		APIKey:     getEnv("API_KEY", "test-api-key-12345"),
		ServerPort: getEnv("SERVER_PORT", "8080"),
		ServerHost: getEnv("SERVER_HOST", "localhost"),

		// Database Configuration
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("DB_PORT", "5432"),
		DBName:     getEnv("DB_NAME", "scootin_aboot"),
		DBUser:     getEnv("DB_USER", "postgres"),
		DBPassword: getEnv("DB_PASSWORD", "password"),
		DBSSLMode:  getEnv("DB_SSLMODE", "disable"),

		// Simulator Configuration
		SimulatorScooters:        getEnvAsInt("SIMULATOR_SCOOTERS", 20),
		SimulatorUsers:           getEnvAsInt("SIMULATOR_USERS", 5),
		SimulatorServerURL:       getEnv("SIMULATOR_SERVER_URL", "http://localhost:8080"),
		SimulatorSpeed:           getEnvAsInt("SIMULATOR_SPEED", 100),
		SimulatorTripDurationMin: getEnvAsInt("SIMULATOR_TRIP_DURATION_MIN", 5),
		SimulatorTripDurationMax: getEnvAsInt("SIMULATOR_TRIP_DURATION_MAX", 10),
		SimulatorRestMin:         getEnvAsInt("SIMULATOR_REST_MIN", 2),
		SimulatorRestMax:         getEnvAsInt("SIMULATOR_REST_MAX", 5),

		// Logging
		LogLevel:  getEnv("LOG_LEVEL", "info"),
		LogFormat: getEnv("LOG_FORMAT", "json"),
	}

	return config, nil
}

// getEnv gets an environment variable with a fallback value
func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

// getEnvAsInt gets an environment variable as integer with a fallback value
func getEnvAsInt(key string, fallback int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return fallback
}

// GetDatabaseDSN returns the database connection string
func (c *Config) GetDatabaseDSN() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		c.DBHost, c.DBPort, c.DBUser, c.DBPassword, c.DBName, c.DBSSLMode)
}
