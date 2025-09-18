package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestScooterModel(t *testing.T) {
	t.Run("CreateScooter", func(t *testing.T) {
		scooter := &Scooter{
			CurrentLatitude:  45.4215, // Ottawa coordinates
			CurrentLongitude: -75.6972,
			Status:           ScooterStatusAvailable,
		}

		// Test validation
		err := scooter.ValidateCoordinates()
		assert.NoError(t, err)

		// Test status methods
		assert.True(t, scooter.IsAvailable())
		assert.False(t, scooter.IsOccupied())

		// Test status change
		err = scooter.SetStatus(ScooterStatusOccupied)
		assert.NoError(t, err)
		assert.True(t, scooter.IsOccupied())
		assert.False(t, scooter.IsAvailable())

		// Test location update
		err = scooter.UpdateLocation(45.4216, -75.6973)
		assert.NoError(t, err)
		assert.Equal(t, 45.4216, scooter.CurrentLatitude)
		assert.Equal(t, -75.6973, scooter.CurrentLongitude)
	})

	t.Run("InvalidCoordinates", func(t *testing.T) {
		scooter := &Scooter{
			CurrentLatitude:  91.0, // Invalid latitude
			CurrentLongitude: -75.6972,
		}

		err := scooter.ValidateCoordinates()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid latitude")

		scooter.CurrentLatitude = 45.4215
		scooter.CurrentLongitude = 181.0 // Invalid longitude

		err = scooter.ValidateCoordinates()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid longitude")
	})

	t.Run("SetStatus", func(t *testing.T) {
		scooter := &Scooter{
			CurrentLatitude:  45.4215,
			CurrentLongitude: -75.6972,
			Status:           ScooterStatusAvailable,
		}

		// Test valid status change
		err := scooter.SetStatus(ScooterStatusOccupied)
		assert.NoError(t, err)
		assert.Equal(t, ScooterStatusOccupied, scooter.Status)

		// Test invalid status
		err = scooter.SetStatus("invalid_status")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid scooter status")
	})

	t.Run("UpdateLocation", func(t *testing.T) {
		scooter := &Scooter{
			CurrentLatitude:  45.4215,
			CurrentLongitude: -75.6972,
			Status:           ScooterStatusAvailable,
		}

		// Test valid location update
		err := scooter.UpdateLocation(45.4216, -75.6973)
		assert.NoError(t, err)
		assert.Equal(t, 45.4216, scooter.CurrentLatitude)
		assert.Equal(t, -75.6973, scooter.CurrentLongitude)

		// Test invalid location update
		err = scooter.UpdateLocation(91.0, -75.6973)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid latitude")
	})
}
