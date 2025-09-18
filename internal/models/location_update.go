package models

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// LocationUpdate represents a location update during a trip
type LocationUpdate struct {
	ID        uuid.UUID      `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	TripID    uuid.UUID      `json:"trip_id" gorm:"type:uuid;not null"`
	Latitude  float64        `json:"latitude" gorm:"type:decimal(10,8);not null"`
	Longitude float64        `json:"longitude" gorm:"type:decimal(11,8);not null"`
	Timestamp time.Time      `json:"timestamp" gorm:"not null"`
	CreatedAt time.Time      `json:"created_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"index"`

	// Relationships
	Trip Trip `json:"trip,omitempty" gorm:"foreignKey:TripID"`
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
	if err := lu.ValidateCoordinates(); err != nil {
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
	// Validate latitude (-90 to 90)
	if lu.Latitude < -90 || lu.Latitude > 90 {
		return errors.New("invalid latitude: must be between -90 and 90")
	}

	// Validate longitude (-180 to 180)
	if lu.Longitude < -180 || lu.Longitude > 180 {
		return errors.New("invalid longitude: must be between -180 and 180")
	}

	return nil
}

// CreateLocationUpdate creates a new location update
func CreateLocationUpdate(tripID uuid.UUID, latitude, longitude float64, timestamp time.Time) (*LocationUpdate, error) {
	lu := &LocationUpdate{
		TripID:    tripID,
		Latitude:  latitude,
		Longitude: longitude,
		Timestamp: timestamp,
	}

	if err := lu.ValidateCoordinates(); err != nil {
		return nil, err
	}

	return lu, nil
}
