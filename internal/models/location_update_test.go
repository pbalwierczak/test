package models

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLocationUpdateModel(t *testing.T) {
	t.Run("CreateLocationUpdate", func(t *testing.T) {
		tripID := uuid.New()
		latitude := 45.4215
		longitude := -75.6972
		timestamp := time.Now()

		lu, err := CreateLocationUpdate(tripID, latitude, longitude, timestamp)
		require.NoError(t, err)

		assert.Equal(t, tripID, lu.TripID)
		assert.Equal(t, latitude, lu.Latitude)
		assert.Equal(t, longitude, lu.Longitude)
		assert.Equal(t, timestamp, lu.Timestamp)

		// Test validation
		err = lu.ValidateCoordinates()
		assert.NoError(t, err)
	})

	t.Run("CreateLocationUpdateWithInvalidCoordinates", func(t *testing.T) {
		tripID := uuid.New()
		latitude := 91.0 // Invalid latitude
		longitude := -75.6972
		timestamp := time.Now()

		lu, err := CreateLocationUpdate(tripID, latitude, longitude, timestamp)
		assert.Error(t, err)
		assert.Nil(t, lu)
		assert.Contains(t, err.Error(), "invalid latitude")
	})

	t.Run("InvalidCoordinates", func(t *testing.T) {
		lu := &LocationUpdate{
			Latitude:  91.0, // Invalid latitude
			Longitude: -75.6972,
		}

		err := lu.ValidateCoordinates()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid latitude")

		lu.Latitude = 45.4215
		lu.Longitude = 181.0 // Invalid longitude

		err = lu.ValidateCoordinates()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid longitude")
	})

	t.Run("ValidCoordinates", func(t *testing.T) {
		testCases := []struct {
			name      string
			latitude  float64
			longitude float64
		}{
			{
				name:      "Ottawa coordinates",
				latitude:  45.4215,
				longitude: -75.6972,
			},
			{
				name:      "Montreal coordinates",
				latitude:  45.5017,
				longitude: -73.5673,
			},
			{
				name:      "Equator coordinates",
				latitude:  0.0,
				longitude: 0.0,
			},
			{
				name:      "North pole coordinates",
				latitude:  90.0,
				longitude: 0.0,
			},
			{
				name:      "South pole coordinates",
				latitude:  -90.0,
				longitude: 0.0,
			},
			{
				name:      "International date line",
				latitude:  0.0,
				longitude: 180.0,
			},
			{
				name:      "Negative longitude",
				latitude:  0.0,
				longitude: -180.0,
			},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				lu := &LocationUpdate{
					Latitude:  tc.latitude,
					Longitude: tc.longitude,
				}

				err := lu.ValidateCoordinates()
				assert.NoError(t, err)
			})
		}
	})

	t.Run("BoundaryCoordinates", func(t *testing.T) {
		// Test exact boundary values
		lu := &LocationUpdate{
			Latitude:  90.0,  // Maximum latitude
			Longitude: 180.0, // Maximum longitude
		}

		err := lu.ValidateCoordinates()
		assert.NoError(t, err)

		lu = &LocationUpdate{
			Latitude:  -90.0,  // Minimum latitude
			Longitude: -180.0, // Minimum longitude
		}

		err = lu.ValidateCoordinates()
		assert.NoError(t, err)
	})

	t.Run("OutOfBoundsCoordinates", func(t *testing.T) {
		testCases := []struct {
			name      string
			latitude  float64
			longitude float64
			expected  string
		}{
			{
				name:      "Latitude too high",
				latitude:  90.1,
				longitude: 0.0,
				expected:  "invalid latitude",
			},
			{
				name:      "Latitude too low",
				latitude:  -90.1,
				longitude: 0.0,
				expected:  "invalid latitude",
			},
			{
				name:      "Longitude too high",
				latitude:  0.0,
				longitude: 180.1,
				expected:  "invalid longitude",
			},
			{
				name:      "Longitude too low",
				latitude:  0.0,
				longitude: -180.1,
				expected:  "invalid longitude",
			},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				lu := &LocationUpdate{
					Latitude:  tc.latitude,
					Longitude: tc.longitude,
				}

				err := lu.ValidateCoordinates()
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tc.expected)
			})
		}
	})
}
