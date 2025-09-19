package models

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestTrip_TableName(t *testing.T) {
	trip := Trip{}
	assert.Equal(t, "trips", trip.TableName())
}

func TestTrip_BeforeCreate(t *testing.T) {
	tests := []struct {
		name        string
		trip        *Trip
		expectError bool
		errorMsg    string
	}{
		{
			name: "successful creation with valid data",
			trip: &Trip{
				StartLatitude:  45.4215,
				StartLongitude: -75.6972,
				Status:         TripStatusActive,
			},
			expectError: false,
		},
		{
			name: "successful creation with existing ID",
			trip: &Trip{
				ID:             uuid.New(),
				StartLatitude:  45.4215,
				StartLongitude: -75.6972,
				Status:         TripStatusActive,
			},
			expectError: false,
		},
		{
			name: "successful creation with existing timestamps",
			trip: &Trip{
				StartLatitude:  45.4215,
				StartLongitude: -75.6972,
				Status:         TripStatusActive,
				CreatedAt:      time.Now().Add(-1 * time.Hour),
				UpdatedAt:      time.Now().Add(-30 * time.Minute),
				StartTime:      time.Now().Add(-15 * time.Minute),
			},
			expectError: false,
		},
		{
			name: "invalid start coordinates - latitude too high",
			trip: &Trip{
				StartLatitude:  91.0,
				StartLongitude: -75.6972,
				Status:         TripStatusActive,
			},
			expectError: true,
			errorMsg:    "invalid latitude",
		},
		{
			name: "invalid start coordinates - longitude too high",
			trip: &Trip{
				StartLatitude:  45.4215,
				StartLongitude: 181.0,
				Status:         TripStatusActive,
			},
			expectError: true,
			errorMsg:    "invalid longitude",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a mock GORM DB
			var tx *gorm.DB

			err := tt.trip.BeforeCreate(tx)

			if tt.expectError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errorMsg)
			} else {
				assert.NoError(t, err)

				// Check that ID was set if it was nil
				if tt.trip.ID == uuid.Nil {
					assert.NotEqual(t, uuid.Nil, tt.trip.ID)
				}

				// Check that timestamps were set if they were zero
				if tt.trip.CreatedAt.IsZero() {
					assert.False(t, tt.trip.CreatedAt.IsZero())
				}
				if tt.trip.UpdatedAt.IsZero() {
					assert.False(t, tt.trip.UpdatedAt.IsZero())
				}
				if tt.trip.StartTime.IsZero() {
					assert.False(t, tt.trip.StartTime.IsZero())
				}
			}
		})
	}
}

func TestTrip_BeforeUpdate(t *testing.T) {
	tests := []struct {
		name        string
		trip        *Trip
		expectError bool
		errorMsg    string
	}{
		{
			name: "successful update with valid start coordinates only",
			trip: &Trip{
				StartLatitude:  45.4215,
				StartLongitude: -75.6972,
				Status:         TripStatusActive,
			},
			expectError: false,
		},
		{
			name: "successful update with valid start and end coordinates",
			trip: &Trip{
				StartLatitude:  45.4215,
				StartLongitude: -75.6972,
				EndLatitude:    &[]float64{45.4216}[0],
				EndLongitude:   &[]float64{-75.6973}[0],
				EndTime:        &[]time.Time{time.Now()}[0],
				Status:         TripStatusCompleted,
			},
			expectError: false,
		},
		{
			name: "invalid start coordinates - latitude too high",
			trip: &Trip{
				StartLatitude:  91.0,
				StartLongitude: -75.6972,
				Status:         TripStatusActive,
			},
			expectError: true,
			errorMsg:    "invalid latitude",
		},
		{
			name: "invalid end coordinates - longitude too high",
			trip: &Trip{
				StartLatitude:  45.4215,
				StartLongitude: -75.6972,
				EndLatitude:    &[]float64{45.4216}[0],
				EndLongitude:   &[]float64{181.0}[0],
				EndTime:        &[]time.Time{time.Now()}[0],
				Status:         TripStatusCompleted,
			},
			expectError: true,
			errorMsg:    "invalid longitude",
		},
		{
			name: "partial end coordinates - should not validate end coordinates",
			trip: &Trip{
				StartLatitude:  45.4215,
				StartLongitude: -75.6972,
				EndLatitude:    &[]float64{45.4216}[0],
				// Missing EndLongitude and EndTime
				Status: TripStatusActive,
			},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a mock GORM DB
			var tx *gorm.DB

			err := tt.trip.BeforeUpdate(tx)

			if tt.expectError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errorMsg)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
