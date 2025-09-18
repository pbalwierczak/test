package repository

import (
	"math"
	"testing"
	"time"

	"scootin-aboot/internal/models"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestHaversineDistance(t *testing.T) {
	tests := []struct {
		name      string
		lat1      float64
		lon1      float64
		lat2      float64
		lon2      float64
		expected  float64
		tolerance float64
	}{
		{
			name:      "Same point",
			lat1:      45.4215,
			lon1:      -75.6972,
			lat2:      45.4215,
			lon2:      -75.6972,
			expected:  0.0,
			tolerance: 0.001,
		},
		{
			name:      "Ottawa to Montreal",
			lat1:      45.4215, // Ottawa
			lon1:      -75.6972,
			lat2:      45.5017, // Montreal
			lon2:      -73.5673,
			expected:  166.0, // Approximately 166 km
			tolerance: 10.0,
		},
		{
			name:      "Short distance",
			lat1:      45.4215,
			lon1:      -75.6972,
			lat2:      45.4225,
			lon2:      -75.6982,
			expected:  0.1, // Approximately 100 meters
			tolerance: 0.05,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			distance := HaversineDistance(tt.lat1, tt.lon1, tt.lat2, tt.lon2)
			assert.InDelta(t, tt.expected, distance, tt.tolerance)
		})
	}
}

func TestIsPointInBounds(t *testing.T) {
	tests := []struct {
		name     string
		lat      float64
		lng      float64
		minLat   float64
		maxLat   float64
		minLng   float64
		maxLng   float64
		expected bool
	}{
		{
			name:     "Point inside bounds",
			lat:      45.4215,
			lng:      -75.6972,
			minLat:   45.0,
			maxLat:   46.0,
			minLng:   -76.0,
			maxLng:   -75.0,
			expected: true,
		},
		{
			name:     "Point outside bounds - latitude too high",
			lat:      46.5,
			lng:      -75.6972,
			minLat:   45.0,
			maxLat:   46.0,
			minLng:   -76.0,
			maxLng:   -75.0,
			expected: false,
		},
		{
			name:     "Point outside bounds - latitude too low",
			lat:      44.5,
			lng:      -75.6972,
			minLat:   45.0,
			maxLat:   46.0,
			minLng:   -76.0,
			maxLng:   -75.0,
			expected: false,
		},
		{
			name:     "Point outside bounds - longitude too high",
			lat:      45.4215,
			lng:      -74.5,
			minLat:   45.0,
			maxLat:   46.0,
			minLng:   -76.0,
			maxLng:   -75.0,
			expected: false,
		},
		{
			name:     "Point outside bounds - longitude too low",
			lat:      45.4215,
			lng:      -76.5,
			minLat:   45.0,
			maxLat:   46.0,
			minLng:   -76.0,
			maxLng:   -75.0,
			expected: false,
		},
		{
			name:     "Point on boundary",
			lat:      45.0,
			lng:      -75.0,
			minLat:   45.0,
			maxLat:   46.0,
			minLng:   -76.0,
			maxLng:   -75.0,
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsPointInBounds(tt.lat, tt.lng, tt.minLat, tt.maxLat, tt.minLng, tt.maxLng)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestIsPointInRadius(t *testing.T) {
	tests := []struct {
		name      string
		centerLat float64
		centerLng float64
		pointLat  float64
		pointLng  float64
		radiusKm  float64
		expected  bool
	}{
		{
			name:      "Point within radius",
			centerLat: 45.4215,
			centerLng: -75.6972,
			pointLat:  45.4225,
			pointLng:  -75.6982,
			radiusKm:  1.0,
			expected:  true,
		},
		{
			name:      "Point outside radius",
			centerLat: 45.4215,
			centerLng: -75.6972,
			pointLat:  45.4315,
			pointLng:  -75.7072,
			radiusKm:  0.1,
			expected:  false,
		},
		{
			name:      "Point exactly on radius boundary",
			centerLat: 45.4215,
			centerLng: -75.6972,
			pointLat:  45.4215,
			pointLng:  -75.6972,
			radiusKm:  0.0,
			expected:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsPointInRadius(tt.centerLat, tt.centerLng, tt.pointLat, tt.pointLng, tt.radiusKm)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestBoundingBox(t *testing.T) {
	t.Run("NewBoundingBox", func(t *testing.T) {
		centerLat := 45.4215
		centerLng := -75.6972
		radiusKm := 1.0

		bb := NewBoundingBox(centerLat, centerLng, radiusKm)

		// Check that the center point is within the bounding box
		assert.True(t, bb.Contains(centerLat, centerLng))

		// Check that the bounding box is approximately correct
		// The radius should be roughly 1/111 degrees (1 km ≈ 1/111 degrees)
		expectedLatDelta := radiusKm / 111.0
		expectedLngDelta := radiusKm / (111.0 * math.Cos(centerLat*math.Pi/180))

		assert.InDelta(t, centerLat-expectedLatDelta, bb.MinLat, 0.01)
		assert.InDelta(t, centerLat+expectedLatDelta, bb.MaxLat, 0.01)
		assert.InDelta(t, centerLng-expectedLngDelta, bb.MinLng, 0.01)
		assert.InDelta(t, centerLng+expectedLngDelta, bb.MaxLng, 0.01)
	})

	t.Run("Contains", func(t *testing.T) {
		bb := BoundingBox{
			MinLat: 45.0,
			MaxLat: 46.0,
			MinLng: -76.0,
			MaxLng: -75.0,
		}

		// Point inside
		assert.True(t, bb.Contains(45.5, -75.5))

		// Point outside
		assert.False(t, bb.Contains(44.5, -75.5))
		assert.False(t, bb.Contains(46.5, -75.5))
		assert.False(t, bb.Contains(45.5, -76.5))
		assert.False(t, bb.Contains(45.5, -74.5))

		// Point on boundary
		assert.True(t, bb.Contains(45.0, -75.0))
		assert.True(t, bb.Contains(46.0, -76.0))
	})

	t.Run("Expand", func(t *testing.T) {
		bb := BoundingBox{
			MinLat: 45.0,
			MaxLat: 46.0,
			MinLng: -76.0,
			MaxLng: -75.0,
		}

		expanded := bb.Expand(1.0) // Expand by 1 km

		// Original box should be contained within expanded box
		assert.True(t, expanded.Contains(45.5, -75.5))

		// Expanded box should be larger
		assert.Less(t, expanded.MinLat, bb.MinLat)
		assert.Greater(t, expanded.MaxLat, bb.MaxLat)
		assert.Less(t, expanded.MinLng, bb.MinLng)
		assert.Greater(t, expanded.MaxLng, bb.MaxLng)
	})
}

// createTestScooters creates test scooters for FilterAndSortByDistance tests
func createTestScooters() []*models.Scooter {
	now := time.Now()
	return []*models.Scooter{
		{
			ID:               uuid.New(),
			Status:           models.ScooterStatusAvailable,
			CurrentLatitude:  45.4215, // Ottawa center
			CurrentLongitude: -75.6972,
			CreatedAt:        now,
			UpdatedAt:        now,
			LastSeen:         now,
		},
		{
			ID:               uuid.New(),
			Status:           models.ScooterStatusAvailable,
			CurrentLatitude:  45.4225, // ~100m from center
			CurrentLongitude: -75.6982,
			CreatedAt:        now,
			UpdatedAt:        now,
			LastSeen:         now,
		},
		{
			ID:               uuid.New(),
			Status:           models.ScooterStatusOccupied,
			CurrentLatitude:  45.4315, // ~1km from center
			CurrentLongitude: -75.7072,
			CreatedAt:        now,
			UpdatedAt:        now,
			LastSeen:         now,
		},
		{
			ID:               uuid.New(),
			Status:           models.ScooterStatusAvailable,
			CurrentLatitude:  45.5017, // Montreal (~166km from Ottawa)
			CurrentLongitude: -73.5673,
			CreatedAt:        now,
			UpdatedAt:        now,
			LastSeen:         now,
		},
		{
			ID:               uuid.New(),
			Status:           models.ScooterStatusAvailable,
			CurrentLatitude:  45.4200, // ~200m from center
			CurrentLongitude: -75.6960,
			CreatedAt:        now,
			UpdatedAt:        now,
			LastSeen:         now,
		},
	}
}

func TestFilterAndSortByDistance(t *testing.T) {
	tests := []struct {
		name          string
		centerLat     float64
		centerLng     float64
		radius        float64
		limit         int
		expectedCount int
		expectedFirst *models.Scooter
	}{
		{
			name:          "Filter scooters within 1km radius",
			centerLat:     45.4215,
			centerLng:     -75.6972,
			radius:        1.0,
			limit:         0,
			expectedCount: 3, // 3 scooters within 1km
			expectedFirst: &models.Scooter{
				CurrentLatitude:  45.4215, // Should be first (exact match)
				CurrentLongitude: -75.6972,
			},
		},
		{
			name:          "Filter scooters within 0.1km radius (very small)",
			centerLat:     45.4215,
			centerLng:     -75.6972,
			radius:        0.1,
			limit:         0,
			expectedCount: 1, // Only the exact match
			expectedFirst: &models.Scooter{
				CurrentLatitude:  45.4215, // Should be first (exact match)
				CurrentLongitude: -75.6972,
			},
		},
		{
			name:          "Filter scooters with large radius (includes all)",
			centerLat:     45.4215,
			centerLng:     -75.6972,
			radius:        200.0,
			limit:         0,
			expectedCount: 5, // All scooters within 200km
			expectedFirst: &models.Scooter{
				CurrentLatitude:  45.4215, // Should be first (exact match)
				CurrentLongitude: -75.6972,
			},
		},
		{
			name:          "Filter scooters with limit",
			centerLat:     45.4215,
			centerLng:     -75.6972,
			radius:        1.0,
			limit:         2,
			expectedCount: 2, // Limited to 2 results
			expectedFirst: &models.Scooter{
				CurrentLatitude:  45.4215, // Should be first (exact match)
				CurrentLongitude: -75.6972,
			},
		},
		{
			name:          "No radius filter (get all)",
			centerLat:     45.4215,
			centerLng:     -75.6972,
			radius:        0,
			limit:         0,
			expectedCount: 5, // All scooters
			expectedFirst: &models.Scooter{
				CurrentLatitude:  45.4215, // Should be first (exact match)
				CurrentLongitude: -75.6972,
			},
		},
		{
			name:          "Negative radius (treated as 0)",
			centerLat:     45.4215,
			centerLng:     -75.6972,
			radius:        -1.0,
			limit:         0,
			expectedCount: 5, // All scooters (negative radius treated as 0)
			expectedFirst: &models.Scooter{
				CurrentLatitude:  45.4215, // Should be first (exact match)
				CurrentLongitude: -75.6972,
			},
		},
		{
			name:          "Empty input slice",
			centerLat:     45.4215,
			centerLng:     -75.6972,
			radius:        1.0,
			limit:         0,
			expectedCount: 0,
			expectedFirst: nil,
		},
		{
			name:          "Limit larger than available items",
			centerLat:     45.4215,
			centerLng:     -75.6972,
			radius:        1.0,
			limit:         10,
			expectedCount: 3, // Only 3 items available within radius
			expectedFirst: &models.Scooter{
				CurrentLatitude:  45.4215, // Should be first (exact match)
				CurrentLongitude: -75.6972,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var scooters []*models.Scooter
			if tt.name != "Empty input slice" {
				scooters = createTestScooters()
			}
			filteredScooters := FilterAndSortByDistance(scooters, tt.centerLat, tt.centerLng, tt.radius, tt.limit)

			// Assertions
			assert.Len(t, filteredScooters, tt.expectedCount)

			if tt.expectedFirst != nil && len(filteredScooters) > 0 {
				assert.Equal(t, tt.expectedFirst.CurrentLatitude, filteredScooters[0].CurrentLatitude)
				assert.Equal(t, tt.expectedFirst.CurrentLongitude, filteredScooters[0].CurrentLongitude)
			}

			// Verify scooters are sorted by distance
			if len(filteredScooters) > 1 {
				for i := 1; i < len(filteredScooters); i++ {
					dist1 := HaversineDistance(tt.centerLat, tt.centerLng, filteredScooters[i-1].CurrentLatitude, filteredScooters[i-1].CurrentLongitude)
					dist2 := HaversineDistance(tt.centerLat, tt.centerLng, filteredScooters[i].CurrentLatitude, filteredScooters[i].CurrentLongitude)
					assert.LessOrEqual(t, dist1, dist2, "Scooters should be sorted by distance")
				}
			}

			// Verify all returned scooters are within radius
			for _, scooter := range filteredScooters {
				distance := HaversineDistance(tt.centerLat, tt.centerLng, scooter.CurrentLatitude, scooter.CurrentLongitude)
				if tt.radius > 0 {
					assert.LessOrEqual(t, distance, tt.radius, "Scooter should be within radius")
				}
			}
		})
	}
}

func TestFilterAndSortByDistanceEdgeCases(t *testing.T) {
	t.Run("Nil input slice", func(t *testing.T) {
		var scooters []*models.Scooter
		result := FilterAndSortByDistance(scooters, 45.4215, -75.6972, 1.0, 0)
		assert.Empty(t, result)
	})

	t.Run("Zero limit", func(t *testing.T) {
		scooters := createTestScooters()
		result := FilterAndSortByDistance(scooters, 45.4215, -75.6972, 1.0, 0)
		assert.Len(t, result, 3) // Should return all items within radius
	})

	t.Run("Very small radius", func(t *testing.T) {
		scooters := createTestScooters()
		result := FilterAndSortByDistance(scooters, 45.4215, -75.6972, 0.001, 0)
		assert.Len(t, result, 1) // Only exact match
	})

	t.Run("Very large radius", func(t *testing.T) {
		scooters := createTestScooters()
		result := FilterAndSortByDistance(scooters, 45.4215, -75.6972, 1000.0, 0)
		assert.Len(t, result, 5) // All scooters
	})
}
