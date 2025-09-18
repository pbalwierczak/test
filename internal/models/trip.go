package models

import (
	"errors"
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

// BeforeCreate hook to set the ID if not already set and validate data
func (t *Trip) BeforeCreate(tx *gorm.DB) error {
	if t.ID == uuid.Nil {
		t.ID = uuid.New()
	}

	// Validate coordinates
	if err := t.ValidateStartCoordinates(); err != nil {
		return err
	}

	// Set timestamps
	now := time.Now()
	if t.CreatedAt.IsZero() {
		t.CreatedAt = now
	}
	if t.UpdatedAt.IsZero() {
		t.UpdatedAt = now
	}
	if t.StartTime.IsZero() {
		t.StartTime = now
	}

	return nil
}

// BeforeUpdate hook to validate data before update
func (t *Trip) BeforeUpdate(tx *gorm.DB) error {
	// Validate coordinates
	if err := t.ValidateStartCoordinates(); err != nil {
		return err
	}

	// If ending the trip, validate end coordinates
	if t.EndTime != nil && t.EndLatitude != nil && t.EndLongitude != nil {
		if err := t.ValidateEndCoordinates(); err != nil {
			return err
		}
	}

	// Update timestamp
	t.UpdatedAt = time.Now()

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

// ValidateStartCoordinates validates the trip's start coordinates
func (t *Trip) ValidateStartCoordinates() error {
	// Validate latitude (-90 to 90)
	if t.StartLatitude < -90 || t.StartLatitude > 90 {
		return errors.New("invalid start latitude: must be between -90 and 90")
	}

	// Validate longitude (-180 to 180)
	if t.StartLongitude < -180 || t.StartLongitude > 180 {
		return errors.New("invalid start longitude: must be between -180 and 180")
	}

	return nil
}

// ValidateEndCoordinates validates the trip's end coordinates
func (t *Trip) ValidateEndCoordinates() error {
	if t.EndLatitude == nil || t.EndLongitude == nil {
		return errors.New("end coordinates cannot be nil")
	}

	// Validate latitude (-90 to 90)
	if *t.EndLatitude < -90 || *t.EndLatitude > 90 {
		return errors.New("invalid end latitude: must be between -90 and 90")
	}

	// Validate longitude (-180 to 180)
	if *t.EndLongitude < -180 || *t.EndLongitude > 180 {
		return errors.New("invalid end longitude: must be between -180 and 180")
	}

	return nil
}

// StartTrip initializes a new trip
func (t *Trip) StartTrip(scooterID, userID uuid.UUID, latitude, longitude float64) error {
	t.ScooterID = scooterID
	t.UserID = userID
	t.StartLatitude = latitude
	t.StartLongitude = longitude
	t.StartTime = time.Now()
	t.Status = TripStatusActive

	return t.ValidateStartCoordinates()
}

// EndTrip completes a trip
func (t *Trip) EndTrip(latitude, longitude float64) error {
	if !t.IsActive() {
		return errors.New("cannot end a trip that is not active")
	}

	now := time.Now()
	t.EndTime = &now
	t.EndLatitude = &latitude
	t.EndLongitude = &longitude
	t.Status = TripStatusCompleted

	return t.ValidateEndCoordinates()
}

// CancelTrip cancels an active trip
func (t *Trip) CancelTrip() error {
	if !t.IsActive() {
		return errors.New("cannot cancel a trip that is not active")
	}

	t.Status = TripStatusCancelled
	t.UpdatedAt = time.Now()

	return nil
}
