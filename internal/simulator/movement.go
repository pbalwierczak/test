package simulator

import (
	"math"
	"math/rand"
	"time"

	"scootin-aboot/internal/config"
)

// Movement handles realistic movement calculations for scooters
type Movement struct {
	config *config.Config
}

// NewMovement creates a new movement calculator
func NewMovement(cfg *config.Config) *Movement {
	return &Movement{
		config: cfg,
	}
}

// Location represents a GPS coordinate
type Location struct {
	Latitude  float64
	Longitude float64
}

// City represents a city with its center coordinates
type City struct {
	Name      string
	CenterLat float64
	CenterLng float64
	RadiusKm  float64
}

// GetCities returns the available cities for simulation
func (m *Movement) GetCities() []City {
	return []City{
		{
			Name:      "Ottawa",
			CenterLat: config.OttawaCenterLat,
			CenterLng: config.OttawaCenterLng,
			RadiusKm:  config.CityRadiusKm,
		},
		{
			Name:      "Montreal",
			CenterLat: config.MontrealCenterLat,
			CenterLng: config.MontrealCenterLng,
			RadiusKm:  config.CityRadiusKm,
		},
	}
}

// GetRandomLocationInCity returns a random location within a city's boundaries
func (m *Movement) GetRandomLocationInCity(city City) Location {
	// Generate random angle and distance
	angle := rand.Float64() * 2 * math.Pi
	distance := rand.Float64() * city.RadiusKm

	// Convert distance from km to degrees
	// 1 degree of latitude ≈ 111 km
	// 1 degree of longitude ≈ 111 km * cos(latitude)
	latOffset := distance / 111.0
	lngOffset := distance / (111.0 * math.Cos(city.CenterLat*math.Pi/180))

	// Calculate new coordinates
	newLat := city.CenterLat + latOffset*math.Cos(angle)
	newLng := city.CenterLng + lngOffset*math.Sin(angle)

	return Location{
		Latitude:  newLat,
		Longitude: newLng,
	}
}

// GetRandomLocation returns a random location in any available city
func (m *Movement) GetRandomLocation() Location {
	cities := m.GetCities()
	city := cities[rand.Intn(len(cities))]
	return m.GetRandomLocationInCity(city)
}

// CalculateMovement calculates the new location after moving for a given duration
func (m *Movement) CalculateMovement(start Location, direction float64, duration time.Duration) Location {
	// Convert speed from km/h to m/s
	speedMs := float64(m.config.SimulatorSpeed) * 1000.0 / 3600.0

	// Calculate distance traveled in meters
	distanceM := speedMs * duration.Seconds()

	// Convert distance from meters to degrees
	// 1 degree of latitude ≈ 111,000 meters
	// 1 degree of longitude ≈ 111,000 meters * cos(latitude)
	latOffset := distanceM / 111000.0
	lngOffset := distanceM / (111000.0 * math.Cos(start.Latitude*math.Pi/180))

	// Calculate new coordinates
	newLat := start.Latitude + latOffset*math.Cos(direction*math.Pi/180)
	newLng := start.Longitude + lngOffset*math.Sin(direction*math.Pi/180)

	return Location{
		Latitude:  newLat,
		Longitude: newLng,
	}
}

// GetRandomDirection returns a random direction in degrees (0-360)
func (m *Movement) GetRandomDirection() float64 {
	return rand.Float64() * 360.0
}

// CalculateDistance calculates the distance between two locations in meters
func (m *Movement) CalculateDistance(loc1, loc2 Location) float64 {
	const earthRadius = 6371000 // Earth's radius in meters

	lat1Rad := loc1.Latitude * math.Pi / 180
	lat2Rad := loc2.Latitude * math.Pi / 180
	deltaLat := (loc2.Latitude - loc1.Latitude) * math.Pi / 180
	deltaLng := (loc2.Longitude - loc1.Longitude) * math.Pi / 180

	a := math.Sin(deltaLat/2)*math.Sin(deltaLat/2) +
		math.Cos(lat1Rad)*math.Cos(lat2Rad)*
			math.Sin(deltaLng/2)*math.Sin(deltaLng/2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	return earthRadius * c
}

// IsWithinCityBounds checks if a location is within a city's boundaries
func (m *Movement) IsWithinCityBounds(location Location, city City) bool {
	distance := m.CalculateDistance(location, Location{
		Latitude:  city.CenterLat,
		Longitude: city.CenterLng,
	})

	return distance <= city.RadiusKm*1000 // Convert km to meters
}

// GetClosestCity returns the closest city to a given location
func (m *Movement) GetClosestCity(location Location) City {
	cities := m.GetCities()
	closest := cities[0]
	minDistance := m.CalculateDistance(location, Location{
		Latitude:  closest.CenterLat,
		Longitude: closest.CenterLng,
	})

	for _, city := range cities[1:] {
		distance := m.CalculateDistance(location, Location{
			Latitude:  city.CenterLat,
			Longitude: city.CenterLng,
		})

		if distance < minDistance {
			minDistance = distance
			closest = city
		}
	}

	return closest
}
