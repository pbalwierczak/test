package models

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTripModel(t *testing.T) {
	t.Run("CreateTrip", func(t *testing.T) {
		scooterID := uuid.New()
		userID := uuid.New()

		trip := &Trip{}
		err := trip.StartTrip(scooterID, userID, 45.4215, -75.6972)
		require.NoError(t, err)

		assert.Equal(t, scooterID, trip.ScooterID)
		assert.Equal(t, userID, trip.UserID)
		assert.Equal(t, 45.4215, trip.StartLatitude)
		assert.Equal(t, -75.6972, trip.StartLongitude)
		assert.Equal(t, TripStatusActive, trip.Status)
		assert.True(t, trip.IsActive())

		// Test ending trip
		err = trip.EndTrip(45.4216, -75.6973)
		require.NoError(t, err)

		assert.Equal(t, TripStatusCompleted, trip.Status)
		assert.True(t, trip.IsCompleted())
		assert.False(t, trip.IsActive())
		assert.NotNil(t, trip.EndTime)
		assert.Equal(t, 45.4216, *trip.EndLatitude)
		assert.Equal(t, -75.6973, *trip.EndLongitude)

		// Test duration
		duration := trip.Duration()
		assert.NotNil(t, duration)
		assert.True(t, *duration > 0)
	})

	t.Run("CancelTrip", func(t *testing.T) {
		scooterID := uuid.New()
		userID := uuid.New()

		trip := &Trip{}
		err := trip.StartTrip(scooterID, userID, 45.4215, -75.6972)
		require.NoError(t, err)

		err = trip.CancelTrip()
		require.NoError(t, err)

		assert.Equal(t, TripStatusCancelled, trip.Status)
		assert.True(t, trip.IsCancelled())
		assert.False(t, trip.IsActive())
	})

	t.Run("InvalidStartCoordinates", func(t *testing.T) {
		trip := &Trip{
			StartLatitude:  91.0, // Invalid latitude
			StartLongitude: -75.6972,
		}

		err := trip.ValidateStartCoordinates()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid latitude")

		trip.StartLatitude = 45.4215
		trip.StartLongitude = 181.0 // Invalid longitude

		err = trip.ValidateStartCoordinates()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid longitude")
	})

	t.Run("InvalidEndCoordinates", func(t *testing.T) {
		latitude := 91.0
		longitude := -75.6972

		trip := &Trip{
			EndLatitude:  &latitude,
			EndLongitude: &longitude,
		}

		err := trip.ValidateEndCoordinates()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid latitude")

		latitude = 45.4215
		longitude = 181.0
		trip.EndLatitude = &latitude
		trip.EndLongitude = &longitude

		err = trip.ValidateEndCoordinates()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid longitude")
	})

	t.Run("EndTripWithoutEndCoordinates", func(t *testing.T) {
		trip := &Trip{
			EndLatitude:  nil,
			EndLongitude: nil,
		}

		err := trip.ValidateEndCoordinates()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "end coordinates cannot be nil")
	})

	t.Run("EndInactiveTrip", func(t *testing.T) {
		trip := &Trip{
			Status: TripStatusCompleted,
		}

		err := trip.EndTrip(45.4216, -75.6973)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "cannot end a trip that is not active")
	})

	t.Run("CancelInactiveTrip", func(t *testing.T) {
		trip := &Trip{
			Status: TripStatusCompleted,
		}

		err := trip.CancelTrip()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "cannot cancel a trip that is not active")
	})

	t.Run("DurationCalculation", func(t *testing.T) {
		// Test trip without end time
		trip := &Trip{
			StartTime: time.Now(),
		}

		duration := trip.Duration()
		assert.Nil(t, duration)

		// Test trip with end time
		startTime := time.Now()
		endTime := startTime.Add(5 * time.Minute)
		trip = &Trip{
			StartTime: startTime,
			EndTime:   &endTime,
		}

		duration = trip.Duration()
		assert.NotNil(t, duration)
		assert.Equal(t, 5*time.Minute, *duration)
	})

	t.Run("StatusChecks", func(t *testing.T) {
		// Test active trip
		trip := &Trip{Status: TripStatusActive}
		assert.True(t, trip.IsActive())
		assert.False(t, trip.IsCompleted())
		assert.False(t, trip.IsCancelled())

		// Test completed trip
		trip = &Trip{Status: TripStatusCompleted}
		assert.False(t, trip.IsActive())
		assert.True(t, trip.IsCompleted())
		assert.False(t, trip.IsCancelled())

		// Test cancelled trip
		trip = &Trip{Status: TripStatusCancelled}
		assert.False(t, trip.IsActive())
		assert.False(t, trip.IsCompleted())
		assert.True(t, trip.IsCancelled())
	})
}
