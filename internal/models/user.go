package models

import (
	"time"

	"github.com/google/uuid"
)

// User represents a user in the system
type User struct {
	ID        uuid.UUID  `json:"id" db:"id"`
	CreatedAt time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt time.Time  `json:"updated_at" db:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty" db:"deleted_at"`

	// Relationships
	Trips []Trip `json:"trips,omitempty"`
}

// TableName returns the table name for the User model
func (User) TableName() string {
	return "users"
}

// SetTimestamps sets the created_at and updated_at timestamps
func (u *User) SetTimestamps() {
	now := time.Now()
	if u.CreatedAt.IsZero() {
		u.CreatedAt = now
	}
	u.UpdatedAt = now
}

// SetID sets the ID if not already set
func (u *User) SetID() {
	if u.ID == uuid.Nil {
		u.ID = uuid.New()
	}
}

// CreateUser creates a new user
func CreateUser() *User {
	return &User{
		ID: uuid.New(),
	}
}
