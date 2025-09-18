package models

import (
	"time"

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

// BeforeCreate hook to set the ID if not already set
func (s *Scooter) BeforeCreate(tx *gorm.DB) error {
	if s.ID == uuid.Nil {
		s.ID = uuid.New()
	}
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
