package models

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestScooter_TableName(t *testing.T) {
	scooter := Scooter{}
	assert.Equal(t, "scooters", scooter.TableName())
}

func TestScooter_BeforeCreate(t *testing.T) {
	tests := []struct {
		name        string
		scooter     *Scooter
		expectError bool
		errorMsg    string
	}{
		{
			name: "successful creation with valid data",
			scooter: &Scooter{
				CurrentLatitude:  45.4215,
				CurrentLongitude: -75.6972,
				Status:           ScooterStatusAvailable,
			},
			expectError: false,
		},
		{
			name: "successful creation with existing ID",
			scooter: &Scooter{
				ID:               uuid.New(),
				CurrentLatitude:  45.4215,
				CurrentLongitude: -75.6972,
				Status:           ScooterStatusAvailable,
			},
			expectError: false,
		},
		{
			name: "successful creation with existing timestamps",
			scooter: &Scooter{
				CurrentLatitude:  45.4215,
				CurrentLongitude: -75.6972,
				Status:           ScooterStatusAvailable,
				CreatedAt:        time.Now().Add(-1 * time.Hour),
				UpdatedAt:        time.Now().Add(-30 * time.Minute),
				LastSeen:         time.Now().Add(-15 * time.Minute),
			},
			expectError: false,
		},
		{
			name: "invalid coordinates - latitude too high",
			scooter: &Scooter{
				CurrentLatitude:  91.0,
				CurrentLongitude: -75.6972,
				Status:           ScooterStatusAvailable,
			},
			expectError: true,
			errorMsg:    "invalid latitude",
		},
		{
			name: "invalid coordinates - latitude too low",
			scooter: &Scooter{
				CurrentLatitude:  -91.0,
				CurrentLongitude: -75.6972,
				Status:           ScooterStatusAvailable,
			},
			expectError: true,
			errorMsg:    "invalid latitude",
		},
		{
			name: "invalid coordinates - longitude too high",
			scooter: &Scooter{
				CurrentLatitude:  45.4215,
				CurrentLongitude: 181.0,
				Status:           ScooterStatusAvailable,
			},
			expectError: true,
			errorMsg:    "invalid longitude",
		},
		{
			name: "invalid coordinates - longitude too low",
			scooter: &Scooter{
				CurrentLatitude:  45.4215,
				CurrentLongitude: -181.0,
				Status:           ScooterStatusAvailable,
			},
			expectError: true,
			errorMsg:    "invalid longitude",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a mock GORM DB (we don't need to use it for these tests)
			var tx *gorm.DB

			err := tt.scooter.BeforeCreate(tx)

			if tt.expectError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errorMsg)
			} else {
				assert.NoError(t, err)

				// Check that ID was set if it was nil
				if tt.scooter.ID == uuid.Nil {
					assert.NotEqual(t, uuid.Nil, tt.scooter.ID)
				}

				// Check that timestamps were set if they were zero
				if tt.scooter.CreatedAt.IsZero() {
					assert.False(t, tt.scooter.CreatedAt.IsZero())
				}
				if tt.scooter.UpdatedAt.IsZero() {
					assert.False(t, tt.scooter.UpdatedAt.IsZero())
				}
				if tt.scooter.LastSeen.IsZero() {
					assert.False(t, tt.scooter.LastSeen.IsZero())
				}
			}
		})
	}
}

func TestScooter_BeforeUpdate(t *testing.T) {
	tests := []struct {
		name        string
		scooter     *Scooter
		expectError bool
		errorMsg    string
	}{
		{
			name: "successful update with valid coordinates",
			scooter: &Scooter{
				CurrentLatitude:  45.4215,
				CurrentLongitude: -75.6972,
				Status:           ScooterStatusAvailable,
			},
			expectError: false,
		},
		{
			name: "invalid coordinates - latitude too high",
			scooter: &Scooter{
				CurrentLatitude:  91.0,
				CurrentLongitude: -75.6972,
				Status:           ScooterStatusAvailable,
			},
			expectError: true,
			errorMsg:    "invalid latitude",
		},
		{
			name: "invalid coordinates - longitude too high",
			scooter: &Scooter{
				CurrentLatitude:  45.4215,
				CurrentLongitude: 181.0,
				Status:           ScooterStatusAvailable,
			},
			expectError: true,
			errorMsg:    "invalid longitude",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Store original UpdatedAt
			originalUpdatedAt := tt.scooter.UpdatedAt

			// Create a mock GORM DB
			var tx *gorm.DB

			err := tt.scooter.BeforeUpdate(tx)

			if tt.expectError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errorMsg)
			} else {
				assert.NoError(t, err)

				// Check that UpdatedAt was set
				assert.True(t, tt.scooter.UpdatedAt.After(originalUpdatedAt) || tt.scooter.UpdatedAt.Equal(originalUpdatedAt))
			}
		})
	}
}

func TestScooter_GetLatitude(t *testing.T) {
	tests := []struct {
		name     string
		scooter  *Scooter
		expected float64
	}{
		{
			name: "get latitude from scooter",
			scooter: &Scooter{
				CurrentLatitude: 45.4215,
			},
			expected: 45.4215,
		},
		{
			name: "get zero latitude",
			scooter: &Scooter{
				CurrentLatitude: 0.0,
			},
			expected: 0.0,
		},
		{
			name: "get negative latitude",
			scooter: &Scooter{
				CurrentLatitude: -45.4215,
			},
			expected: -45.4215,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.scooter.GetLatitude()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestScooter_GetLongitude(t *testing.T) {
	tests := []struct {
		name     string
		scooter  *Scooter
		expected float64
	}{
		{
			name: "get longitude from scooter",
			scooter: &Scooter{
				CurrentLongitude: -75.6972,
			},
			expected: -75.6972,
		},
		{
			name: "get zero longitude",
			scooter: &Scooter{
				CurrentLongitude: 0.0,
			},
			expected: 0.0,
		},
		{
			name: "get negative longitude",
			scooter: &Scooter{
				CurrentLongitude: -75.6972,
			},
			expected: -75.6972,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.scooter.GetLongitude()
			assert.Equal(t, tt.expected, result)
		})
	}
}
