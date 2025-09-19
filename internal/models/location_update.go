package models

import (
	"time"

	"scootin-aboot/pkg/validation"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// LocationUpdate represents a location update for a scooter
type LocationUpdate struct {
	ID        uuid.UUID      `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	ScooterID uuid.UUID      `json:"scooter_id" gorm:"type:uuid;not null"`
	Latitude  float64        `json:"latitude" gorm:"type:decimal(10,8);not null"`
	Longitude float64        `json:"longitude" gorm:"type:decimal(11,8);not null"`
	Timestamp time.Time      `json:"timestamp" gorm:"not null"`
	CreatedAt time.Time      `json:"created_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"index"`

	// Relationships
	Scooter Scooter `json:"scooter,omitempty" gorm:"foreignKey:ScooterID"`
}

// TableName returns the table name for the LocationUpdate model
func (LocationUpdate) TableName() string {
	return "location_updates"
}

// BeforeCreate hook to set the ID if not already set and validate data
func (lu *LocationUpdate) BeforeCreate(tx *gorm.DB) error {
	if lu.ID == uuid.Nil {
		lu.ID = uuid.New()
	}

	// Validate coordinates
	if err := validation.ValidateCoordinates(lu.Latitude, lu.Longitude); err != nil {
		return err
	}

	// Set timestamp if not provided
	if lu.Timestamp.IsZero() {
		lu.Timestamp = time.Now()
	}

	// Set created timestamp
	if lu.CreatedAt.IsZero() {
		lu.CreatedAt = time.Now()
	}

	return nil
}

// ValidateCoordinates validates the location update coordinates
func (lu *LocationUpdate) ValidateCoordinates() error {
	return validation.ValidateCoordinates(lu.Latitude, lu.Longitude)
}

// CreateLocationUpdate creates a new location update
func CreateLocationUpdate(scooterID uuid.UUID, latitude, longitude float64, timestamp time.Time) (*LocationUpdate, error) {
	lu := &LocationUpdate{
		ScooterID: scooterID,
		Latitude:  latitude,
		Longitude: longitude,
		Timestamp: timestamp,
	}

	if err := validation.ValidateCoordinates(latitude, longitude); err != nil {
		return nil, err
	}

	return lu, nil
}
