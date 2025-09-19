package models

import (
	"time"

	"scootin-aboot/pkg/validation"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type LocationUpdate struct {
	ID        uuid.UUID      `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	ScooterID uuid.UUID      `json:"scooter_id" gorm:"type:uuid;not null"`
	Latitude  float64        `json:"latitude" gorm:"type:decimal(10,8);not null"`
	Longitude float64        `json:"longitude" gorm:"type:decimal(11,8);not null"`
	Timestamp time.Time      `json:"timestamp" gorm:"not null"`
	CreatedAt time.Time      `json:"created_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"index"`

	Scooter Scooter `json:"scooter,omitempty" gorm:"foreignKey:ScooterID"`
}

func (LocationUpdate) TableName() string {
	return "location_updates"
}

func (lu *LocationUpdate) BeforeCreate(tx *gorm.DB) error {
	if lu.ID == uuid.Nil {
		lu.ID = uuid.New()
	}

	if err := validation.ValidateCoordinates(lu.Latitude, lu.Longitude); err != nil {
		return err
	}

	if lu.Timestamp.IsZero() {
		lu.Timestamp = time.Now()
	}

	if lu.CreatedAt.IsZero() {
		lu.CreatedAt = time.Now()
	}

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
