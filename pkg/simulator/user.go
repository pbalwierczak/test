package simulator

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"scootin-aboot/internal/config"
	"scootin-aboot/pkg/utils"

	"go.uber.org/zap"
)

// User simulates a mobile client user
type User struct {
	ID       int
	UserID   string // Store the actual user ID from database
	ctx      context.Context
	client   *APIClient
	config   *config.Config
	movement *Movement
}

// NewUser creates a new user simulator
func NewUser(ctx context.Context, client *APIClient, id int, cfg *config.Config) (*User, error) {
	return &User{
		ID:       id,
		ctx:      ctx,
		client:   client,
		config:   cfg,
		movement: NewMovement(cfg),
	}, nil
}

// NewUserWithID creates a new user simulator with a specific user ID
func NewUserWithID(ctx context.Context, client *APIClient, id int, userID string, cfg *config.Config) (*User, error) {
	return &User{
		ID:       id,
		UserID:   userID,
		ctx:      ctx,
		client:   client,
		config:   cfg,
		movement: NewMovement(cfg),
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
			u.simulateTrip()
		}
	}
}

// simulateTrip simulates a complete trip cycle
func (u *User) simulateTrip() {
	// Find an available scooter
	scooter, err := u.findAvailableScooter()
	if err != nil {
		utils.Error("Failed to find available scooter",
			zap.Int("user_id", u.ID),
			zap.Error(err),
		)
		u.rest()
		return
	}

	// Start trip
	tripID, err := u.startTrip(scooter.ID, scooter.Latitude, scooter.Longitude)
	if err != nil {
		utils.Error("Failed to start trip",
			zap.Int("user_id", u.ID),
			zap.String("scooter_id", scooter.ID),
			zap.Error(err),
		)
		u.rest()
		return
	}

	utils.Info("Trip started",
		zap.Int("user_id", u.ID),
		zap.String("scooter_id", scooter.ID),
		zap.String("trip_id", tripID),
	)

	// Simulate driving and get final location
	finalLocation := u.simulateDriving(scooter.ID, tripID)

	// End trip at the final location
	if err := u.endTrip(scooter.ID, finalLocation.Latitude, finalLocation.Longitude); err != nil {
		utils.Error("Failed to end trip",
			zap.Int("user_id", u.ID),
			zap.String("scooter_id", scooter.ID),
			zap.Error(err),
		)
	}

	utils.Info("Trip ended",
		zap.Int("user_id", u.ID),
		zap.String("scooter_id", scooter.ID),
	)

	// Rest before next trip
	u.rest()
}

// findAvailableScooter finds an available scooter
func (u *User) findAvailableScooter() (*APIScooter, error) {
	scooters, err := u.client.GetAvailableScooters(u.ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get available scooters: %w", err)
	}

	if len(scooters) == 0 {
		return nil, fmt.Errorf("no available scooters")
	}

	// Pick a random scooter
	scooter := scooters[rand.Intn(len(scooters))]
	return &scooter, nil
}

// startTrip starts a trip with the given scooter
func (u *User) startTrip(scooterID string, startLat, startLng float64) (string, error) {
	// Use the actual user ID from the database
	userID := u.UserID
	if userID == "" {
		// Fallback to old logic for backward compatibility
		userID = fmt.Sprintf("550e8400-e29b-41d4-a716-446655440%03d", u.ID)
	}

	response, err := u.client.StartTrip(u.ctx, scooterID, userID, startLat, startLng)
	if err != nil {
		return "", fmt.Errorf("failed to start trip: %w", err)
	}

	return response.TripID, nil
}

// endTrip ends the trip with the given scooter
func (u *User) endTrip(scooterID string, endLat, endLng float64) error {
	// Use the actual user ID from the database
	userID := u.UserID
	if userID == "" {
		// Fallback to old logic for backward compatibility
		userID = fmt.Sprintf("550e8400-e29b-41d4-a716-446655440%03d", u.ID)
	}

	return u.client.EndTrip(u.ctx, scooterID, userID, endLat, endLng)
}

// simulateDriving simulates the driving portion of the trip and returns the final location
func (u *User) simulateDriving(scooterID, tripID string) Location {
	// Calculate trip duration (5-10 seconds)
	duration := time.Duration(u.config.SimulatorTripDurationMin+rand.Intn(u.config.SimulatorTripDurationMax-u.config.SimulatorTripDurationMin+1)) * time.Second

	utils.Info("Driving simulation started",
		zap.Int("user_id", u.ID),
		zap.String("scooter_id", scooterID),
		zap.Duration("duration", duration),
	)

	// Start location (random)
	startLocation := u.movement.GetRandomLocation()
	currentLocation := startLocation

	// Random direction (straight line movement)
	direction := u.movement.GetRandomDirection()

	// Update location every 3 seconds
	ticker := time.NewTicker(3 * time.Second)
	defer ticker.Stop()

	tripTimer := time.NewTimer(duration)
	defer tripTimer.Stop()

	for {
		select {
		case <-u.ctx.Done():
			return currentLocation
		case <-tripTimer.C:
			// Trip duration completed
			return currentLocation
		case <-ticker.C:
			// Update location
			currentLocation = u.movement.CalculateMovement(currentLocation, direction, 3*time.Second)

			// Send location update to server
			if err := u.client.UpdateLocation(u.ctx, scooterID, currentLocation.Latitude, currentLocation.Longitude); err != nil {
				utils.Error("Failed to update location",
					zap.Int("user_id", u.ID),
					zap.String("scooter_id", scooterID),
					zap.Error(err),
				)
			} else {
				utils.Debug("Location updated",
					zap.Int("user_id", u.ID),
					zap.String("scooter_id", scooterID),
					zap.Float64("lat", currentLocation.Latitude),
					zap.Float64("lng", currentLocation.Longitude),
				)
			}
		}
	}
}

// rest simulates the rest period between trips
func (u *User) rest() {
	// Calculate rest duration (2-5 seconds)
	duration := time.Duration(u.config.SimulatorRestMin+rand.Intn(u.config.SimulatorRestMax-u.config.SimulatorRestMin+1)) * time.Second

	utils.Info("User resting",
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
