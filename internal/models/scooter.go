package models

import (
	"errors"
	"time"

	"scootin-aboot/pkg/validation"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// ScooterStatus represents the possible statuses of a scooter
type ScooterStatus string

const (
	ScooterStatusAvailable ScooterStatus = "available"
	ScooterStatusOccupied  ScooterStatus = "occupied"
)

// Scooter represents a scooter in the system
type Scooter struct {
	ID               uuid.UUID      `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	Status           ScooterStatus  `json:"status" gorm:"type:varchar(20);not null;default:'available';check:status IN ('available','occupied')"`
	CurrentLatitude  float64        `json:"current_latitude" gorm:"type:decimal(10,8);not null"`
	CurrentLongitude float64        `json:"current_longitude" gorm:"type:decimal(11,8);not null"`
	CreatedAt        time.Time      `json:"created_at"`
	UpdatedAt        time.Time      `json:"updated_at"`
	LastSeen         time.Time      `json:"last_seen" gorm:"default:CURRENT_TIMESTAMP"`
	DeletedAt        gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"index"`

	// Relationships
	Trips []Trip `json:"trips,omitempty" gorm:"foreignKey:ScooterID"`
}

// TableName returns the table name for the Scooter model
func (Scooter) TableName() string {
	return "scooters"
}

// BeforeCreate hook to set the ID if not already set and validate data
func (s *Scooter) BeforeCreate(tx *gorm.DB) error {
	if s.ID == uuid.Nil {
		s.ID = uuid.New()
	}

	// Validate coordinates
	if err := validation.ValidateCoordinates(s.CurrentLatitude, s.CurrentLongitude); err != nil {
		return err
	}

	// Set timestamps
	now := time.Now()
	if s.CreatedAt.IsZero() {
		s.CreatedAt = now
	}
	if s.UpdatedAt.IsZero() {
		s.UpdatedAt = now
	}
	if s.LastSeen.IsZero() {
		s.LastSeen = now
	}

	return nil
}

// BeforeUpdate hook to validate data before update
func (s *Scooter) BeforeUpdate(tx *gorm.DB) error {
	// Validate coordinates
	if err := validation.ValidateCoordinates(s.CurrentLatitude, s.CurrentLongitude); err != nil {
		return err
	}

	// Update timestamp
	s.UpdatedAt = time.Now()

	return nil
}

// IsAvailable checks if the scooter is available for use
func (s *Scooter) IsAvailable() bool {
	return s.Status == ScooterStatusAvailable
}

// IsOccupied checks if the scooter is currently occupied
func (s *Scooter) IsOccupied() bool {
	return s.Status == ScooterStatusOccupied
}

// ValidateCoordinates validates the scooter's coordinates
func (s *Scooter) ValidateCoordinates() error {
	return validation.ValidateCoordinates(s.CurrentLatitude, s.CurrentLongitude)
}

// UpdateLocation updates the scooter's location and last seen timestamp
func (s *Scooter) UpdateLocation(latitude, longitude float64) error {
	s.CurrentLatitude = latitude
	s.CurrentLongitude = longitude
	s.LastSeen = time.Now()

	return validation.ValidateCoordinates(latitude, longitude)
}

// SetStatus updates the scooter's status
func (s *Scooter) SetStatus(status ScooterStatus) error {
	switch status {
	case ScooterStatusAvailable, ScooterStatusOccupied:
		s.Status = status
		s.UpdatedAt = time.Now()
		return nil
	default:
		return errors.New("invalid scooter status")
	}
}

// GetLatitude returns the scooter's current latitude
func (s *Scooter) GetLatitude() float64 {
	return s.CurrentLatitude
}

// GetLongitude returns the scooter's current longitude
func (s *Scooter) GetLongitude() float64 {
	return s.CurrentLongitude
}
