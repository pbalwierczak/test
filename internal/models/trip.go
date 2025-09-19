package models

import (
	"errors"
	"time"

	"scootin-aboot/pkg/validation"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TripStatus string

const (
	TripStatusActive    TripStatus = "active"
	TripStatusCompleted TripStatus = "completed"
	TripStatusCancelled TripStatus = "cancelled"
)

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

	Scooter Scooter `json:"scooter,omitempty" gorm:"foreignKey:ScooterID"`
	User    User    `json:"user,omitempty" gorm:"foreignKey:UserID"`
}

func (Trip) TableName() string {
	return "trips"
}

func (t *Trip) BeforeCreate(tx *gorm.DB) error {
	if t.ID == uuid.Nil {
		t.ID = uuid.New()
	}

	if err := t.ValidateStartCoordinates(); err != nil {
		return err
	}

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

func (t *Trip) BeforeUpdate(tx *gorm.DB) error {
	if err := t.ValidateStartCoordinates(); err != nil {
		return err
	}

	if t.EndTime != nil && t.EndLatitude != nil && t.EndLongitude != nil {
		if err := t.ValidateEndCoordinates(); err != nil {
			return err
		}
	}

	t.UpdatedAt = time.Now()

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
