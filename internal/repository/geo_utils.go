package repository

import "math"

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

// IsPointInBounds checks if a point is within the given rectangular bounds
func IsPointInBounds(lat, lng, minLat, maxLat, minLng, maxLng float64) bool {
	return lat >= minLat && lat <= maxLat && lng >= minLng && lng <= maxLng
}

// IsPointInRadius checks if a point is within the given radius of a center point
func IsPointInRadius(centerLat, centerLng, pointLat, pointLng, radiusKm float64) bool {
	distance := HaversineDistance(centerLat, centerLng, pointLat, pointLng)
	return distance <= radiusKm
}

// BoundingBox represents a rectangular geographic area
type BoundingBox struct {
	MinLat float64
	MaxLat float64
	MinLng float64
	MaxLng float64
}

// NewBoundingBox creates a new bounding box from center point and radius
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

// Contains checks if a point is within the bounding box
func (bb BoundingBox) Contains(lat, lng float64) bool {
	return IsPointInBounds(lat, lng, bb.MinLat, bb.MaxLat, bb.MinLng, bb.MaxLng)
}

// Expand expands the bounding box by the given radius
func (bb BoundingBox) Expand(radiusKm float64) BoundingBox {
	latDelta := radiusKm / 111.0
	lngDelta := radiusKm / 111.0 // Use average for simplicity

	return BoundingBox{
		MinLat: bb.MinLat - latDelta,
		MaxLat: bb.MaxLat + latDelta,
		MinLng: bb.MinLng - lngDelta,
		MaxLng: bb.MaxLng + lngDelta,
	}
}
