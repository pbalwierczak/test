package services

import (
	"scootin-aboot/internal/models"
	"time"

	"github.com/google/uuid"
)

// TestData contains common test data structures
var TestData = struct {
	// Valid UUIDs for testing
	ValidScooterID uuid.UUID
	ValidUserID    uuid.UUID
	ValidTripID    uuid.UUID

	// Valid coordinates
	ValidLatitude  float64
	ValidLongitude float64

	// Invalid coordinates
	InvalidLatitudeHigh  float64
	InvalidLatitudeLow   float64
	InvalidLongitudeHigh float64
	InvalidLongitudeLow  float64

	// Valid time
	ValidTime time.Time

	// Test limits and offsets
	ValidLimit     int
	ValidOffset    int
	MaxLimit       int
	ExcessiveLimit int

	// Test radius
	ValidRadius     float64
	ExcessiveRadius float64
}{
	ValidScooterID: uuid.MustParse("550e8400-e29b-41d4-a716-446655440000"),
	ValidUserID:    uuid.MustParse("550e8400-e29b-41d4-a716-446655440001"),
	ValidTripID:    uuid.MustParse("550e8400-e29b-41d4-a716-446655440002"),

	ValidLatitude:  45.4215,
	ValidLongitude: -75.6972,

	InvalidLatitudeHigh:  91.0,
	InvalidLatitudeLow:   -91.0,
	InvalidLongitudeHigh: 181.0,
	InvalidLongitudeLow:  -181.0,

	ValidTime: time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC),

	ValidLimit:     10,
	ValidOffset:    0,
	MaxLimit:       100,
	ExcessiveLimit: 101,

	ValidRadius:     1000.0,
	ExcessiveRadius: 60000.0,
}

// GetValidScooterQueryParams returns valid scooter query parameters
func GetValidScooterQueryParams() ScooterQueryParams {
	return ScooterQueryParams{
		Limit:  TestData.ValidLimit,
		Offset: TestData.ValidOffset,
	}
}

// GetValidScooterQueryParamsWithStatus returns valid scooter query parameters with status
func GetValidScooterQueryParamsWithStatus(status string) ScooterQueryParams {
	params := GetValidScooterQueryParams()
	params.Status = status
	return params
}

// GetValidScooterQueryParamsWithBounds returns valid scooter query parameters with geographic bounds
func GetValidScooterQueryParamsWithBounds() ScooterQueryParams {
	params := GetValidScooterQueryParams()
	params.MinLat = TestData.ValidLatitude - 0.1
	params.MaxLat = TestData.ValidLatitude + 0.1
	params.MinLng = TestData.ValidLongitude - 0.1
	params.MaxLng = TestData.ValidLongitude + 0.1
	return params
}

// GetValidClosestScootersQueryParams returns valid closest scooters query parameters
func GetValidClosestScootersQueryParams() ClosestScootersQueryParams {
	return ClosestScootersQueryParams{
		Latitude:  TestData.ValidLatitude,
		Longitude: TestData.ValidLongitude,
		Radius:    TestData.ValidRadius,
		Limit:     TestData.ValidLimit,
	}
}

// GetValidClosestScootersQueryParamsWithStatus returns valid closest scooters query parameters with status
func GetValidClosestScootersQueryParamsWithStatus(status string) ClosestScootersQueryParams {
	params := GetValidClosestScootersQueryParams()
	params.Status = status
	return params
}

// GetInvalidScooterQueryParams returns various invalid scooter query parameters for testing
func GetInvalidScooterQueryParams() []struct {
	Name   string
	Params ScooterQueryParams
} {
	return []struct {
		Name   string
		Params ScooterQueryParams
	}{
		{
			Name: "invalid status",
			Params: ScooterQueryParams{
				Status: "invalid",
				Limit:  TestData.ValidLimit,
				Offset: TestData.ValidOffset,
			},
		},
		{
			Name: "negative limit",
			Params: ScooterQueryParams{
				Limit:  -1,
				Offset: TestData.ValidOffset,
			},
		},
		{
			Name: "negative offset",
			Params: ScooterQueryParams{
				Limit:  TestData.ValidLimit,
				Offset: -1,
			},
		},
		{
			Name: "excessive limit",
			Params: ScooterQueryParams{
				Limit:  TestData.ExcessiveLimit,
				Offset: TestData.ValidOffset,
			},
		},
	}
}

// GetInvalidClosestScootersQueryParams returns various invalid closest scooters query parameters for testing
func GetInvalidClosestScootersQueryParams() []struct {
	Name   string
	Params ClosestScootersQueryParams
} {
	return []struct {
		Name   string
		Params ClosestScootersQueryParams
	}{
		{
			Name: "invalid latitude high",
			Params: ClosestScootersQueryParams{
				Latitude:  TestData.InvalidLatitudeHigh,
				Longitude: TestData.ValidLongitude,
				Radius:    TestData.ValidRadius,
				Limit:     TestData.ValidLimit,
			},
		},
		{
			Name: "invalid latitude low",
			Params: ClosestScootersQueryParams{
				Latitude:  TestData.InvalidLatitudeLow,
				Longitude: TestData.ValidLongitude,
				Radius:    TestData.ValidRadius,
				Limit:     TestData.ValidLimit,
			},
		},
		{
			Name: "invalid longitude high",
			Params: ClosestScootersQueryParams{
				Latitude:  TestData.ValidLatitude,
				Longitude: TestData.InvalidLongitudeHigh,
				Radius:    TestData.ValidRadius,
				Limit:     TestData.ValidLimit,
			},
		},
		{
			Name: "invalid longitude low",
			Params: ClosestScootersQueryParams{
				Latitude:  TestData.ValidLatitude,
				Longitude: TestData.InvalidLongitudeLow,
				Radius:    TestData.ValidRadius,
				Limit:     TestData.ValidLimit,
			},
		},
		{
			Name: "negative radius",
			Params: ClosestScootersQueryParams{
				Latitude:  TestData.ValidLatitude,
				Longitude: TestData.ValidLongitude,
				Radius:    -100,
				Limit:     TestData.ValidLimit,
			},
		},
		{
			Name: "excessive radius",
			Params: ClosestScootersQueryParams{
				Latitude:  TestData.ValidLatitude,
				Longitude: TestData.ValidLongitude,
				Radius:    TestData.ExcessiveRadius,
				Limit:     TestData.ValidLimit,
			},
		},
		{
			Name: "negative limit",
			Params: ClosestScootersQueryParams{
				Latitude:  TestData.ValidLatitude,
				Longitude: TestData.ValidLongitude,
				Radius:    TestData.ValidRadius,
				Limit:     -1,
			},
		},
		{
			Name: "excessive limit",
			Params: ClosestScootersQueryParams{
				Latitude:  TestData.ValidLatitude,
				Longitude: TestData.ValidLongitude,
				Radius:    TestData.ValidRadius,
				Limit:     51,
			},
		},
		{
			Name: "invalid status",
			Params: ClosestScootersQueryParams{
				Latitude:  TestData.ValidLatitude,
				Longitude: TestData.ValidLongitude,
				Radius:    TestData.ValidRadius,
				Limit:     TestData.ValidLimit,
				Status:    "invalid",
			},
		},
	}
}

// GetTestScooters returns a slice of test scooters
func GetTestScooters(count int) []*models.Scooter {
	scooters := make([]*models.Scooter, count)
	for i := 0; i < count; i++ {
		scooters[i] = NewTestScooterBuilder().
			WithID(uuid.New()).
			WithLocation(TestData.ValidLatitude+float64(i)*0.001, TestData.ValidLongitude+float64(i)*0.001).
			Build()
	}
	return scooters
}

// GetTestScootersWithStatus returns a slice of test scooters with specific status
func GetTestScootersWithStatus(count int, status models.ScooterStatus) []*models.Scooter {
	scooters := make([]*models.Scooter, count)
	for i := 0; i < count; i++ {
		scooters[i] = NewTestScooterBuilder().
			WithID(uuid.New()).
			WithStatus(status).
			WithLocation(TestData.ValidLatitude+float64(i)*0.001, TestData.ValidLongitude+float64(i)*0.001).
			Build()
	}
	return scooters
}
