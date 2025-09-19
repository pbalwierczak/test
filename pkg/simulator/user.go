package simulator

import (
	"context"
	"fmt"
	"math"
	"math/rand"
	"time"

	"scootin-aboot/internal/config"
	"scootin-aboot/pkg/utils"
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
	utils.Info("User simulation started", utils.Int("user_id", u.ID))

	for {
		select {
		case <-u.ctx.Done():
			utils.Info("User simulation stopped", utils.Int("user_id", u.ID))
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
		utils.Error("Failed to find available scooters",
			utils.Int("user_id", u.ID),
			utils.String("search_type", searchType),
			utils.ErrorField(err),
		)
		return
	}

	if len(scooters) == 0 {
		utils.Debug("No available scooters found",
			utils.String("test_user_id", u.UserID),
			utils.Int("user_id", u.ID),
			utils.String("search_type", searchType),
			utils.Float64("lat", u.Location.Latitude),
			utils.Float64("lng", u.Location.Longitude),
			utils.Int("radius", u.SearchRadius),
		)
		return
	}

	utils.Info("Found available scooters",
		utils.Int("test_user_id", u.ID),
		utils.String("test_user_id", u.UserID),
		utils.String("search_type", searchType),
		utils.Int("count", len(scooters)),
		utils.Float64("lat", u.Location.Latitude),
		utils.Float64("lng", u.Location.Longitude),
		utils.Int("radius", u.SearchRadius),
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

	utils.Debug("User moved to new location",
		utils.Int("user_id", u.ID),
		utils.Float64("old_lat", oldLocation.Latitude),
		utils.Float64("old_lng", oldLocation.Longitude),
		utils.Float64("new_lat", u.Location.Latitude),
		utils.Float64("new_lng", u.Location.Longitude),
	)
}

func (u *User) updateSearchRadius() {
	// Random radius between 500m and 3000m
	radii := []int{500, 750, 1000, 1500, 2000, 2500, 3000}
	u.SearchRadius = radii[rand.Intn(len(radii))]

	utils.Debug("User updated search radius",
		utils.Int("user_id", u.ID),
		utils.Int("new_radius", u.SearchRadius),
	)
}

func (u *User) rest() {
	// Calculate rest duration (2-5 seconds)
	duration := time.Duration(u.config.SimulatorRestMin+rand.Intn(u.config.SimulatorRestMax-u.config.SimulatorRestMin+1)) * time.Second

	utils.Debug("User resting between searches",
		utils.Int("user_id", u.ID),
		utils.Duration("duration", duration),
	)

	select {
	case <-u.ctx.Done():
		return
	case <-time.After(duration):
		// Rest period completed
	}
}
