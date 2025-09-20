package models

import (
	"time"

	"scootin-aboot/pkg/validation"

	"github.com/google/uuid"
)

type LocationUpdate struct {
	ID        uuid.UUID  `json:"id" db:"id"`
	ScooterID uuid.UUID  `json:"scooter_id" db:"scooter_id"`
	Latitude  float64    `json:"latitude" db:"latitude"`
	Longitude float64    `json:"longitude" db:"longitude"`
	Timestamp time.Time  `json:"timestamp" db:"timestamp"`
	CreatedAt time.Time  `json:"created_at" db:"created_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty" db:"deleted_at"`

	Scooter Scooter `json:"scooter,omitempty"`
}

func (LocationUpdate) TableName() string {
	return "location_updates"
}

// SetTimestamps sets the created_at timestamp
func (lu *LocationUpdate) SetTimestamps() {
	if lu.CreatedAt.IsZero() {
		lu.CreatedAt = time.Now()
	}
}

// SetID sets the ID if not already set
func (lu *LocationUpdate) SetID() {
	if lu.ID == uuid.Nil {
		lu.ID = uuid.New()
	}
}

// ValidateAndSetTimestamps validates coordinates and sets timestamps
func (lu *LocationUpdate) ValidateAndSetTimestamps() error {
	if err := validation.ValidateCoordinates(lu.Latitude, lu.Longitude); err != nil {
		return err
	}

	if lu.Timestamp.IsZero() {
		lu.Timestamp = time.Now()
	}

	lu.SetTimestamps()
	return nil
}

func (lu *LocationUpdate) ValidateCoordinates() error {
	return validation.ValidateCoordinates(lu.Latitude, lu.Longitude)
}

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
