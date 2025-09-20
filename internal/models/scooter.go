package models

import (
	"errors"
	"time"

	"scootin-aboot/pkg/validation"

	"github.com/google/uuid"
)

type ScooterStatus string

const (
	ScooterStatusAvailable ScooterStatus = "available"
	ScooterStatusOccupied  ScooterStatus = "occupied"
)

type Scooter struct {
	ID               uuid.UUID     `json:"id" db:"id"`
	Status           ScooterStatus `json:"status" db:"status"`
	CurrentLatitude  float64       `json:"current_latitude" db:"current_latitude"`
	CurrentLongitude float64       `json:"current_longitude" db:"current_longitude"`
	CreatedAt        time.Time     `json:"created_at" db:"created_at"`
	UpdatedAt        time.Time     `json:"updated_at" db:"updated_at"`
	LastSeen         time.Time     `json:"last_seen" db:"last_seen"`
	DeletedAt        *time.Time    `json:"deleted_at,omitempty" db:"deleted_at"`

	// Relationships
	Trips []Trip `json:"trips,omitempty"`
}

func (Scooter) TableName() string {
	return "scooters"
}

// SetTimestamps sets the created_at, updated_at, and last_seen timestamps
func (s *Scooter) SetTimestamps() {
	now := time.Now()
	if s.CreatedAt.IsZero() {
		s.CreatedAt = now
	}
	s.UpdatedAt = now
	if s.LastSeen.IsZero() {
		s.LastSeen = now
	}
}

// SetID sets the ID if not already set
func (s *Scooter) SetID() {
	if s.ID == uuid.Nil {
		s.ID = uuid.New()
	}
}

// ValidateAndSetTimestamps validates coordinates and sets timestamps
func (s *Scooter) ValidateAndSetTimestamps() error {
	if err := validation.ValidateCoordinates(s.CurrentLatitude, s.CurrentLongitude); err != nil {
		return err
	}
	s.SetTimestamps()
	return nil
}

func (s *Scooter) IsAvailable() bool {
	return s.Status == ScooterStatusAvailable
}

func (s *Scooter) IsOccupied() bool {
	return s.Status == ScooterStatusOccupied
}

func (s *Scooter) ValidateCoordinates() error {
	return validation.ValidateCoordinates(s.CurrentLatitude, s.CurrentLongitude)
}

func (s *Scooter) UpdateLocation(latitude, longitude float64) error {
	s.CurrentLatitude = latitude
	s.CurrentLongitude = longitude
	s.LastSeen = time.Now()

	return validation.ValidateCoordinates(latitude, longitude)
}

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

func (s *Scooter) GetLatitude() float64 {
	return s.CurrentLatitude
}

func (s *Scooter) GetLongitude() float64 {
	return s.CurrentLongitude
}
