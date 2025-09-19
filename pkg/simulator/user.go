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
			// Only search for scooters
			u.searchForScooters()
			u.rest()
		}
	}
}

// searchForScooters searches for available scooters
func (u *User) searchForScooters() {
	// Find available scooters
	scooters, err := u.findAvailableScooters()
	if err != nil {
		utils.Error("Failed to find available scooters",
			zap.Int("user_id", u.ID),
			zap.Error(err),
		)
		return
	}

	if len(scooters) == 0 {
		utils.Debug("No available scooters found",
			zap.Int("user_id", u.ID),
		)
		return
	}

	// Log that we found scooters (simulating user browsing)
	utils.Info("Found available scooters",
		zap.Int("user_id", u.ID),
		zap.Int("count", len(scooters)),
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
