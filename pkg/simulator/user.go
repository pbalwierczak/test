package simulator

import (
	"context"
	"fmt"
	"math"
	"math/rand"
	"time"

	"scootin-aboot/internal/config"
	"scootin-aboot/pkg/logger"
)

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

func (u *User) Simulate() {
	logger.Info("User simulation started", logger.Int("user_id", u.ID))

	for {
		select {
		case <-u.ctx.Done():
			logger.Info("User simulation stopped", logger.Int("user_id", u.ID))
			return
		default:
			u.searchForScooters()
			u.rest()
		}
	}
}

func (u *User) searchForScooters() {
	u.SearchCount++

	if u.SearchCount%5 == 0 {
		u.moveToNewLocation()
	}

	if u.SearchCount%3 == 0 {
		u.updateSearchRadius()
	}

	var scooters []APIScooter
	var err error
	var searchType string

	switch u.SearchCount % 4 {
	case 0:
		scooters, err = u.findClosestScooters()
		searchType = "closest"
	case 1:
		scooters, err = u.findScootersInBounds()
		searchType = "bounds"
	case 2:
		scooters, err = u.findClosestScootersWithRadius(u.SearchRadius * 2)
		searchType = "expanded_radius"
	default:
		scooters, err = u.findAvailableScooters()
		searchType = "all_available"
	}

	if err != nil {
		logger.Error("Failed to find available scooters",
			logger.Int("user_id", u.ID),
			logger.String("search_type", searchType),
			logger.ErrorField(err),
		)
		return
	}

	if len(scooters) == 0 {
		logger.Debug("No available scooters found",
			logger.String("test_user_id", u.UserID),
			logger.Int("user_id", u.ID),
			logger.String("search_type", searchType),
			logger.Float64("lat", u.Location.Latitude),
			logger.Float64("lng", u.Location.Longitude),
			logger.Int("radius", u.SearchRadius),
		)
		return
	}

	logger.Info("Found available scooters",
		logger.Int("test_user_id", u.ID),
		logger.String("test_user_id", u.UserID),
		logger.String("search_type", searchType),
		logger.Int("count", len(scooters)),
		logger.Float64("lat", u.Location.Latitude),
		logger.Float64("lng", u.Location.Longitude),
		logger.Int("radius", u.SearchRadius),
	)
}

func (u *User) findAvailableScooters() ([]APIScooter, error) {
	scooters, err := u.client.GetAvailableScooters(u.ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get available scooters: %w", err)
	}

	return scooters, nil
}

func (u *User) findClosestScooters() ([]APIScooter, error) {
	scooters, err := u.client.GetClosestScooters(u.ctx, u.Location.Latitude, u.Location.Longitude, u.SearchRadius, 10)
	if err != nil {
		return nil, fmt.Errorf("failed to get closest scooters: %w", err)
	}

	return scooters, nil
}

func (u *User) findClosestScootersWithRadius(radius int) ([]APIScooter, error) {
	scooters, err := u.client.GetClosestScooters(u.ctx, u.Location.Latitude, u.Location.Longitude, radius, 10)
	if err != nil {
		return nil, fmt.Errorf("failed to get closest scooters: %w", err)
	}

	return scooters, nil
}

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

func (u *User) moveToNewLocation() {
	oldLocation := u.Location
	u.Location = u.movement.GetRandomLocation()

	logger.Debug("User moved to new location",
		logger.Int("user_id", u.ID),
		logger.Float64("old_lat", oldLocation.Latitude),
		logger.Float64("old_lng", oldLocation.Longitude),
		logger.Float64("new_lat", u.Location.Latitude),
		logger.Float64("new_lng", u.Location.Longitude),
	)
}

func (u *User) updateSearchRadius() {
	// Random radius between 500m and 3000m
	radii := []int{500, 750, 1000, 1500, 2000, 2500, 3000}
	u.SearchRadius = radii[rand.Intn(len(radii))]

	logger.Debug("User updated search radius",
		logger.Int("user_id", u.ID),
		logger.Int("new_radius", u.SearchRadius),
	)
}

func (u *User) rest() {
	// Calculate rest duration (2-5 seconds)
	duration := time.Duration(u.config.SimulatorRestMin+rand.Intn(u.config.SimulatorRestMax-u.config.SimulatorRestMin+1)) * time.Second

	logger.Debug("User resting between searches",
		logger.Int("user_id", u.ID),
		logger.Duration("duration", duration),
	)

	select {
	case <-u.ctx.Done():
		return
	case <-time.After(duration):
		// Rest period completed
	}
}
