package models

import (
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

// BeforeCreate hook to set the ID if not already set
func (lu *LocationUpdate) BeforeCreate(tx *gorm.DB) error {
	if lu.ID == uuid.Nil {
		lu.ID = uuid.New()
	}
	return nil
}
