package repository

import (
	"errors"
	"fmt"
	"math"
	"sort"

	"scootin-aboot/pkg/validation"
)

// HaversineDistance calculates the distance between two points on Earth using the Haversine formula
// Returns distance in kilometers
func HaversineDistance(lat1, lon1, lat2, lon2 float64) float64 {
	const earthRadius = 6371 // Earth's radius in kilometers

	// Convert degrees to radians
	lat1Rad := lat1 * math.Pi / 180
	lon1Rad := lon1 * math.Pi / 180
	lat2Rad := lat2 * math.Pi / 180
	lon2Rad := lon2 * math.Pi / 180

	// Haversine formula
	dlat := lat2Rad - lat1Rad
	dlon := lon2Rad - lon1Rad

	a := math.Sin(dlat/2)*math.Sin(dlat/2) +
		math.Cos(lat1Rad)*math.Cos(lat2Rad)*
			math.Sin(dlon/2)*math.Sin(dlon/2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	return earthRadius * c
}

type BoundingBox struct {
	MinLat float64
	MaxLat float64
	MinLng float64
	MaxLng float64
}

func NewBoundingBox(centerLat, centerLng, radiusKm float64) BoundingBox {
	// Approximate conversion: 1 degree latitude ≈ 111 km
	// 1 degree longitude ≈ 111 km * cos(latitude)
	latDelta := radiusKm / 111.0
	lngDelta := radiusKm / (111.0 * math.Cos(centerLat*math.Pi/180))

	return BoundingBox{
		MinLat: centerLat - latDelta,
		MaxLat: centerLat + latDelta,
		MinLng: centerLng - lngDelta,
		MaxLng: centerLng + lngDelta,
	}
}

func ValidateGeographicBounds(minLat, maxLat, minLng, maxLng float64) error {
	if minLat == 0 && maxLat == 0 && minLng == 0 && maxLng == 0 {
		return nil
	}

	if err := validation.ValidateCoordinates(minLat, minLng); err != nil {
		return fmt.Errorf("invalid min bounds: %w", err)
	}
	if err := validation.ValidateCoordinates(maxLat, maxLng); err != nil {
		return fmt.Errorf("invalid max bounds: %w", err)
	}

	if minLat >= maxLat {
		return errors.New("min_lat must be less than max_lat")
	}
	if minLng >= maxLng {
		return errors.New("min_lng must be less than max_lng")
	}

	latDiff := maxLat - minLat
	lngDiff := maxLng - minLng
	if latDiff > 10 || lngDiff > 10 {
		return errors.New("geographic bounds are too large (max 10 degrees)")
	}

	return nil
}

type LocationProvider interface {
	GetLatitude() float64
	GetLongitude() float64
}

type ItemWithDistance[T LocationProvider] struct {
	Item     T
	Distance float64
}

func FilterAndSortByDistance[T LocationProvider](items []T, centerLat, centerLng, radius float64, limit int) []T {
	if len(items) == 0 {
		return items
	}

	var itemsWithDistance []ItemWithDistance[T]
	for _, item := range items {
		distance := HaversineDistance(centerLat, centerLng, item.GetLatitude(), item.GetLongitude())
		if radius <= 0 || distance <= radius {
			itemsWithDistance = append(itemsWithDistance, ItemWithDistance[T]{
				Item:     item,
				Distance: distance,
			})
		}
	}

	sort.Slice(itemsWithDistance, func(i, j int) bool {
		return itemsWithDistance[i].Distance < itemsWithDistance[j].Distance
	})

	if limit > 0 && len(itemsWithDistance) > limit {
		itemsWithDistance = itemsWithDistance[:limit]
	}

	result := make([]T, len(itemsWithDistance))
	for i, itemWithDistance := range itemsWithDistance {
		result[i] = itemWithDistance.Item
	}

	return result
}
