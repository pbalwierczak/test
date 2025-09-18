package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// TripStatus represents the possible statuses of a trip
type TripStatus string

const (
	TripStatusActive    TripStatus = "active"
	TripStatusCompleted TripStatus = "completed"
	TripStatusCancelled TripStatus = "cancelled"
)

// Trip represents a trip in the system
type Trip struct {
	ID             uuid.UUID      `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	ScooterID      uuid.UUID      `json:"scooter_id" gorm:"type:uuid;not null"`
	UserID         uuid.UUID      `json:"user_id" gorm:"type:uuid;not null"`
	StartTime      time.Time      `json:"start_time" gorm:"not null"`
	EndTime        *time.Time     `json:"end_time,omitempty"`
	StartLatitude  float64        `json:"start_latitude" gorm:"type:decimal(10,8);not null"`
	StartLongitude float64        `json:"start_longitude" gorm:"type:decimal(11,8);not null"`
	EndLatitude    *float64       `json:"end_latitude,omitempty" gorm:"type:decimal(10,8)"`
	EndLongitude   *float64       `json:"end_longitude,omitempty" gorm:"type:decimal(11,8)"`
	Status         TripStatus     `json:"status" gorm:"type:varchar(20);not null;default:'active';check:status IN ('active','completed','cancelled')"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	DeletedAt      gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"index"`

	// Relationships
	Scooter         Scooter          `json:"scooter,omitempty" gorm:"foreignKey:ScooterID"`
	User            User             `json:"user,omitempty" gorm:"foreignKey:UserID"`
	LocationUpdates []LocationUpdate `json:"location_updates,omitempty" gorm:"foreignKey:TripID"`
}

// TableName returns the table name for the Trip model
func (Trip) TableName() string {
	return "trips"
}

// BeforeCreate hook to set the ID if not already set
func (t *Trip) BeforeCreate(tx *gorm.DB) error {
	if t.ID == uuid.Nil {
		t.ID = uuid.New()
	}
	return nil
}

// IsActive checks if the trip is currently active
func (t *Trip) IsActive() bool {
	return t.Status == TripStatusActive
}

// IsCompleted checks if the trip has been completed
func (t *Trip) IsCompleted() bool {
	return t.Status == TripStatusCompleted
}

// IsCancelled checks if the trip has been cancelled
func (t *Trip) IsCancelled() bool {
	return t.Status == TripStatusCancelled
}

// Duration returns the duration of the trip if it has ended
func (t *Trip) Duration() *time.Duration {
	if t.EndTime == nil {
		return nil
	}
	duration := t.EndTime.Sub(t.StartTime)
	return &duration
}
