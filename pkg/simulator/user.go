package simulator

import (
	"context"
	"fmt"
	"math"
	"math/rand"
	"time"

	"scootin-aboot/internal/config"
	"scootin-aboot/pkg/utils"

	"go.uber.org/zap"
)

// User simulates a mobile client user
type User struct {
	ID           int
	UserID       string // Store the actual user ID from database
	ctx          context.Context
	client       *APIClient
	config       *config.Config
	movement     *Movement
	Location     Location // Current user location
	SearchRadius int      // Current search radius in meters
	SearchCount  int      // Number of searches performed
}

// NewUserWithID creates a new user simulator with a specific user ID
func NewUserWithID(ctx context.Context, client *APIClient, id int, userID string, cfg *config.Config) (*User, error) {
	movement := NewMovement(cfg)
	return &User{
		ID:           id,
		UserID:       userID,
		ctx:          ctx,
		client:       client,
		config:       cfg,
		movement:     movement,
		Location:     movement.GetRandomLocation(),
		SearchRadius: 1000, // Default 1km radius
		SearchCount:  0,
	}, nil
}

// Simulate runs the user simulation loop
func (u *User) Simulate() {
	utils.Info("User simulation started", zap.Int("user_id", u.ID))

	for {
		select {
		case <-u.ctx.Done():
			utils.Info("User simulation stopped", zap.Int("user_id", u.ID))
			return
		default:
			// Only search for scooters
			u.searchForScooters()
			u.rest()
		}
	}
}

// searchForScooters searches for available scooters using different strategies
func (u *User) searchForScooters() {
	u.SearchCount++

	// Move user location occasionally (simulate walking around)
	if u.SearchCount%5 == 0 {
		u.moveToNewLocation()
	}

	// Vary search radius occasionally
	if u.SearchCount%3 == 0 {
		u.updateSearchRadius()
	}

	// Choose search strategy based on search count
	var scooters []APIScooter
	var err error
	var searchType string

	switch u.SearchCount % 4 {
	case 0:
		// Search for closest scooters with current radius
		scooters, err = u.findClosestScooters()
		searchType = "closest"
	case 1:
		// Search for scooters in geographic bounds
		scooters, err = u.findScootersInBounds()
		searchType = "bounds"
	case 2:
		// Search with a different radius
		scooters, err = u.findClosestScootersWithRadius(u.SearchRadius * 2)
		searchType = "expanded_radius"
	default:
		// Search all available scooters
		scooters, err = u.findAvailableScooters()
		searchType = "all_available"
	}

	if err != nil {
		utils.Error("Failed to find available scooters",
			zap.Int("user_id", u.ID),
			zap.String("search_type", searchType),
			zap.Error(err),
		)
		return
	}

	if len(scooters) == 0 {
		utils.Debug("No available scooters found",
			zap.String("test_user_id", u.UserID),
			zap.Int("user_id", u.ID),
			zap.String("search_type", searchType),
			zap.Float64("lat", u.Location.Latitude),
			zap.Float64("lng", u.Location.Longitude),
			zap.Int("radius", u.SearchRadius),
		)
		return
	}

	// Log that we found scooters (simulating user browsing)
	utils.Info("Found available scooters",
		zap.Int("test_user_id", u.ID),
		zap.String("test_user_id", u.UserID),
		zap.String("search_type", searchType),
		zap.Int("count", len(scooters)),
		zap.Float64("lat", u.Location.Latitude),
		zap.Float64("lng", u.Location.Longitude),
		zap.Int("radius", u.SearchRadius),
	)
}

// findAvailableScooters finds available scooters
func (u *User) findAvailableScooters() ([]APIScooter, error) {
	scooters, err := u.client.GetAvailableScooters(u.ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get available scooters: %w", err)
	}

	return scooters, nil
}

// findClosestScooters finds the closest scooters to the user's current location
func (u *User) findClosestScooters() ([]APIScooter, error) {
	scooters, err := u.client.GetClosestScooters(u.ctx, u.Location.Latitude, u.Location.Longitude, u.SearchRadius, 10)
	if err != nil {
		return nil, fmt.Errorf("failed to get closest scooters: %w", err)
	}

	return scooters, nil
}

// findClosestScootersWithRadius finds the closest scooters with a specific radius
func (u *User) findClosestScootersWithRadius(radius int) ([]APIScooter, error) {
	scooters, err := u.client.GetClosestScooters(u.ctx, u.Location.Latitude, u.Location.Longitude, radius, 10)
	if err != nil {
		return nil, fmt.Errorf("failed to get closest scooters: %w", err)
	}

	return scooters, nil
}

// findScootersInBounds finds scooters within geographic bounds around the user
func (u *User) findScootersInBounds() ([]APIScooter, error) {
	// Create a bounding box around the user's location
	// Convert radius from meters to approximate degrees
	latOffset := float64(u.SearchRadius) / 111000.0 // 1 degree latitude ≈ 111km
	lngOffset := float64(u.SearchRadius) / (111000.0 * math.Cos(u.Location.Latitude*math.Pi/180))

	minLat := u.Location.Latitude - latOffset
	maxLat := u.Location.Latitude + latOffset
	minLng := u.Location.Longitude - lngOffset
	maxLng := u.Location.Longitude + lngOffset

	scooters, err := u.client.GetScootersInBounds(u.ctx, minLat, maxLat, minLng, maxLng, 20)
	if err != nil {
		return nil, fmt.Errorf("failed to get scooters in bounds: %w", err)
	}

	return scooters, nil
}

// moveToNewLocation moves the user to a new random location
func (u *User) moveToNewLocation() {
	oldLocation := u.Location
	u.Location = u.movement.GetRandomLocation()

	utils.Debug("User moved to new location",
		zap.Int("user_id", u.ID),
		zap.Float64("old_lat", oldLocation.Latitude),
		zap.Float64("old_lng", oldLocation.Longitude),
		zap.Float64("new_lat", u.Location.Latitude),
		zap.Float64("new_lng", u.Location.Longitude),
	)
}

// updateSearchRadius updates the search radius to a random value
func (u *User) updateSearchRadius() {
	// Random radius between 500m and 3000m
	radii := []int{500, 750, 1000, 1500, 2000, 2500, 3000}
	u.SearchRadius = radii[rand.Intn(len(radii))]

	utils.Debug("User updated search radius",
		zap.Int("user_id", u.ID),
		zap.Int("new_radius", u.SearchRadius),
	)
}

// rest simulates the rest period between searches
func (u *User) rest() {
	// Calculate rest duration (2-5 seconds)
	duration := time.Duration(u.config.SimulatorRestMin+rand.Intn(u.config.SimulatorRestMax-u.config.SimulatorRestMin+1)) * time.Second

	utils.Debug("User resting between searches",
		zap.Int("user_id", u.ID),
		zap.Duration("duration", duration),
	)

	select {
	case <-u.ctx.Done():
		return
	case <-time.After(duration):
		// Rest period completed
	}
}
