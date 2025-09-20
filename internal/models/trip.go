package models

import (
	"errors"
	"time"

	"scootin-aboot/internal/validation"

	"github.com/google/uuid"
)

type TripStatus string

const (
	TripStatusActive    TripStatus = "active"
	TripStatusCompleted TripStatus = "completed"
	TripStatusCancelled TripStatus = "cancelled"
)

type Trip struct {
	ID             uuid.UUID  `json:"id" db:"id"`
	ScooterID      uuid.UUID  `json:"scooter_id" db:"scooter_id"`
	UserID         uuid.UUID  `json:"user_id" db:"user_id"`
	StartTime      time.Time  `json:"start_time" db:"start_time"`
	EndTime        *time.Time `json:"end_time,omitempty" db:"end_time"`
	StartLatitude  float64    `json:"start_latitude" db:"start_latitude"`
	StartLongitude float64    `json:"start_longitude" db:"start_longitude"`
	EndLatitude    *float64   `json:"end_latitude,omitempty" db:"end_latitude"`
	EndLongitude   *float64   `json:"end_longitude,omitempty" db:"end_longitude"`
	Status         TripStatus `json:"status" db:"status"`
	CreatedAt      time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at" db:"updated_at"`
	DeletedAt      *time.Time `json:"deleted_at,omitempty" db:"deleted_at"`

	Scooter Scooter `json:"scooter,omitempty"`
	User    User    `json:"user,omitempty"`
}

func (Trip) TableName() string {
	return "trips"
}

// SetTimestamps sets the created_at and updated_at timestamps
func (t *Trip) SetTimestamps() {
	now := time.Now()
	if t.CreatedAt.IsZero() {
		t.CreatedAt = now
	}
	t.UpdatedAt = now
}

// SetID sets the ID if not already set
func (t *Trip) SetID() {
	if t.ID == uuid.Nil {
		t.ID = uuid.New()
	}
}

// ValidateAndSetTimestamps validates coordinates and sets timestamps
func (t *Trip) ValidateAndSetTimestamps() error {
	if err := t.ValidateStartCoordinates(); err != nil {
		return err
	}

	if t.EndTime != nil && t.EndLatitude != nil && t.EndLongitude != nil {
		if err := t.ValidateEndCoordinates(); err != nil {
			return err
		}
	}

	t.SetTimestamps()
	return nil
}

func (t *Trip) IsActive() bool {
	return t.Status == TripStatusActive
}

func (t *Trip) IsCompleted() bool {
	return t.Status == TripStatusCompleted
}

func (t *Trip) IsCancelled() bool {
	return t.Status == TripStatusCancelled
}

func (t *Trip) Duration() *time.Duration {
	if t.EndTime == nil {
		return nil
	}
	duration := t.EndTime.Sub(t.StartTime)
	return &duration
}

func (t *Trip) ValidateStartCoordinates() error {
	if err := validation.ValidateCoordinates(t.StartLatitude, t.StartLongitude); err != nil {
		return err
	}
	return nil
}

func (t *Trip) ValidateEndCoordinates() error {
	if t.EndLatitude == nil || t.EndLongitude == nil {
		return errors.New("end coordinates cannot be nil")
	}

	if err := validation.ValidateCoordinates(*t.EndLatitude, *t.EndLongitude); err != nil {
		return err
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
