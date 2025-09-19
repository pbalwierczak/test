package config

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	APIKey     string
	ServerPort string
	ServerHost string

	DBHost     string
	DBPort     string
	DBName     string
	DBUser     string
	DBPassword string
	DBSSLMode  string

	SimulatorScooters        int
	SimulatorUsers           int
	SimulatorServerURL       string
	SimulatorSpeed           int
	SimulatorTripDurationMin int
	SimulatorTripDurationMax int
	SimulatorRestMin         int
	SimulatorRestMax         int

	KafkaConfig KafkaConfig

	LogLevel  string
	LogFormat string
}

type KafkaConfig struct {
	Brokers          []string
	ClientID         string
	SecurityProtocol string
	Topics           KafkaTopics
}

type KafkaTopics struct {
	TripStarted     string
	TripEnded       string
	LocationUpdated string
}

// City configuration constants
const (
	OttawaCenterLat   = 45.4215
	OttawaCenterLng   = -75.6972
	MontrealCenterLat = 45.5017
	MontrealCenterLng = -73.5673
	CityRadiusKm      = 15.0
)

func Load() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: .env file not found: %v", err)
	}

	config := &Config{
		APIKey:     getEnv("API_KEY", "test-api-key-12345"),
		ServerPort: getEnv("SERVER_PORT", "8080"),
		ServerHost: getEnv("SERVER_HOST", "localhost"),

		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("DB_PORT", "5432"),
		DBName:     getEnv("DB_NAME", "scootin_aboot"),
		DBUser:     getEnv("DB_USER", "postgres"),
		DBPassword: getEnv("DB_PASSWORD", "password"),
		DBSSLMode:  getEnv("DB_SSLMODE", "disable"),

		SimulatorScooters:        getEnvAsInt("SIMULATOR_SCOOTERS", 20),
		SimulatorUsers:           getEnvAsInt("SIMULATOR_USERS", 5),
		SimulatorServerURL:       getEnv("SIMULATOR_SERVER_URL", "http://localhost:8080"),
		SimulatorSpeed:           getEnvAsInt("SIMULATOR_SPEED", 100),
		SimulatorTripDurationMin: getEnvAsInt("SIMULATOR_TRIP_DURATION_MIN", 5),
		SimulatorTripDurationMax: getEnvAsInt("SIMULATOR_TRIP_DURATION_MAX", 10),
		SimulatorRestMin:         getEnvAsInt("SIMULATOR_REST_MIN", 2),
		SimulatorRestMax:         getEnvAsInt("SIMULATOR_REST_MAX", 5),

		KafkaConfig: KafkaConfig{
			Brokers:          getEnvAsStringSlice("KAFKA_BROKERS", []string{"localhost:9092"}),
			ClientID:         getEnv("KAFKA_CLIENT_ID", "scooter-simulator"),
			SecurityProtocol: getEnv("KAFKA_SECURITY_PROTOCOL", "PLAINTEXT"),
			Topics: KafkaTopics{
				TripStarted:     getEnv("KAFKA_TOPIC_TRIP_STARTED", "scooter.trip.started"),
				TripEnded:       getEnv("KAFKA_TOPIC_TRIP_ENDED", "scooter.trip.ended"),
				LocationUpdated: getEnv("KAFKA_TOPIC_LOCATION_UPDATED", "scooter.location.updated"),
			},
		},

		LogLevel:  getEnv("LOG_LEVEL", "info"),
		LogFormat: getEnv("LOG_FORMAT", "json"),
	}

	return config, nil
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

func getEnvAsInt(key string, fallback int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return fallback
}

func getEnvAsStringSlice(key string, fallback []string) []string {
	if value := os.Getenv(key); value != "" {
		// Split by comma and trim whitespace
		parts := strings.Split(value, ",")
		result := make([]string, 0, len(parts))
		for _, part := range parts {
			if trimmed := strings.TrimSpace(part); trimmed != "" {
				result = append(result, trimmed)
			}
		}
		if len(result) > 0 {
			return result
		}
	}
	return fallback
}

func (c *Config) GetDatabaseDSN() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		c.DBHost, c.DBPort, c.DBUser, c.DBPassword, c.DBName, c.DBSSLMode)
}
