package models

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestLocationUpdate_TableName(t *testing.T) {
	locationUpdate := LocationUpdate{}
	assert.Equal(t, "location_updates", locationUpdate.TableName())
}

func TestLocationUpdate_BeforeCreate(t *testing.T) {
	tests := []struct {
		name           string
		locationUpdate *LocationUpdate
		expectError    bool
		errorMsg       string
	}{
		{
			name: "successful creation with valid data",
			locationUpdate: &LocationUpdate{
				Latitude:  45.4215,
				Longitude: -75.6972,
			},
			expectError: false,
		},
		{
			name: "successful creation with existing ID",
			locationUpdate: &LocationUpdate{
				ID:        uuid.New(),
				Latitude:  45.4215,
				Longitude: -75.6972,
			},
			expectError: false,
		},
		{
			name: "successful creation with existing timestamps",
			locationUpdate: &LocationUpdate{
				Latitude:  45.4215,
				Longitude: -75.6972,
				Timestamp: time.Now().Add(-1 * time.Hour),
				CreatedAt: time.Now().Add(-30 * time.Minute),
			},
			expectError: false,
		},
		{
			name: "invalid coordinates - latitude too high",
			locationUpdate: &LocationUpdate{
				Latitude:  91.0,
				Longitude: -75.6972,
			},
			expectError: true,
			errorMsg:    "invalid latitude",
		},
		{
			name: "invalid coordinates - latitude too low",
			locationUpdate: &LocationUpdate{
				Latitude:  -91.0,
				Longitude: -75.6972,
			},
			expectError: true,
			errorMsg:    "invalid latitude",
		},
		{
			name: "invalid coordinates - longitude too high",
			locationUpdate: &LocationUpdate{
				Latitude:  45.4215,
				Longitude: 181.0,
			},
			expectError: true,
			errorMsg:    "invalid longitude",
		},
		{
			name: "invalid coordinates - longitude too low",
			locationUpdate: &LocationUpdate{
				Latitude:  45.4215,
				Longitude: -181.0,
			},
			expectError: true,
			errorMsg:    "invalid longitude",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a mock GORM DB
			var tx *gorm.DB

			err := tt.locationUpdate.BeforeCreate(tx)

			if tt.expectError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errorMsg)
			} else {
				assert.NoError(t, err)

				// Check that ID was set if it was nil
				if tt.locationUpdate.ID == uuid.Nil {
					assert.NotEqual(t, uuid.Nil, tt.locationUpdate.ID)
				}

				// Check that timestamps were set if they were zero
				if tt.locationUpdate.Timestamp.IsZero() {
					assert.False(t, tt.locationUpdate.Timestamp.IsZero())
				}
				if tt.locationUpdate.CreatedAt.IsZero() {
					assert.False(t, tt.locationUpdate.CreatedAt.IsZero())
				}
			}
		})
	}
}
