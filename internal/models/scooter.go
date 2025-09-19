package models

import (
	"errors"
	"time"

	"scootin-aboot/pkg/validation"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ScooterStatus string

const (
	ScooterStatusAvailable ScooterStatus = "available"
	ScooterStatusOccupied  ScooterStatus = "occupied"
)

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

func (Scooter) TableName() string {
	return "scooters"
}

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

func (s *Scooter) BeforeUpdate(tx *gorm.DB) error {
	if err := validation.ValidateCoordinates(s.CurrentLatitude, s.CurrentLongitude); err != nil {
		return err
	}

	s.UpdatedAt = time.Now()

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
