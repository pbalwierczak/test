package simulator

import (
	"context"
	"fmt"
	"time"

	"scootin-aboot/internal/config"
	"scootin-aboot/pkg/utils"

	"go.uber.org/zap"
)

// Scooter simulates a scooter's behavior
type Scooter struct {
	ID          int
	Ctx         context.Context
	Client      *APIClient
	Config      *config.Config
	Movement    *Movement
	CurrentTrip *Trip
	Location    Location
	Status      string
	LastSeen    time.Time
}

// Trip represents an active trip
type Trip struct {
	ID        string
	UserID    string
	StartTime time.Time
	Direction float64
}

// NewScooter creates a new scooter simulator
func NewScooter(ctx context.Context, client *APIClient, id int, cfg *config.Config) (*Scooter, error) {
	movement := NewMovement(cfg)

	// Start with random location
	location := movement.GetRandomLocation()

	return &Scooter{
		ID:       id,
		Ctx:      ctx,
		Client:   client,
		Config:   cfg,
		Movement: movement,
		Location: location,
		Status:   "available",
		LastSeen: time.Now(),
	}, nil
}

// Simulate runs the scooter simulation loop
func (s *Scooter) Simulate() {
	utils.Info("Scooter simulation started",
		zap.Int("scooter_id", s.ID),
		zap.String("status", s.Status),
		zap.Float64("lat", s.Location.Latitude),
		zap.Float64("lng", s.Location.Longitude),
	)

	// Update location every 3 seconds
	ticker := time.NewTicker(3 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-s.Ctx.Done():
			utils.Info("Scooter simulation stopped",
				zap.Int("scooter_id", s.ID),
			)
			return
		case <-ticker.C:
			s.updateLocation()
		}
	}
}

// updateLocation updates the scooter's location
func (s *Scooter) updateLocation() {
	// Only update location if scooter is occupied (has an active trip)
	if s.Status != "occupied" || s.CurrentTrip == nil {
		return
	}

	// Calculate new location based on movement
	newLocation := s.Movement.CalculateMovement(s.Location, s.CurrentTrip.Direction, 3*time.Second)

	// Update location
	s.Location = newLocation
	s.LastSeen = time.Now()

	// Send location update to server
	if err := s.Client.UpdateLocation(s.Ctx, s.getScooterID(), s.Location.Latitude, s.Location.Longitude); err != nil {
		utils.Error("Failed to update scooter location",
			zap.Int("scooter_id", s.ID),
			zap.Error(err),
		)
	} else {
		utils.Debug("Scooter location updated",
			zap.Int("scooter_id", s.ID),
			zap.String("trip_id", s.CurrentTrip.ID),
			zap.Float64("lat", s.Location.Latitude),
			zap.Float64("lng", s.Location.Longitude),
		)
	}
}

// getScooterID returns the scooter ID as a string
func (s *Scooter) getScooterID() string {
	// Use seeded scooter IDs from the database
	// Ottawa scooters: 650e8400-e29b-41d4-a716-446655440001 to 010
	// Montreal scooters: 750e8400-e29b-41d4-a716-446655440001 to 010
	if s.ID <= 10 {
		// Ottawa scooters
		return fmt.Sprintf("650e8400-e29b-41d4-a716-446655440%03d", s.ID)
	} else {
		// Montreal scooters
		return fmt.Sprintf("750e8400-e29b-41d4-a716-446655440%03d", s.ID-10)
	}
}

// StartTrip starts a trip for this scooter
func (s *Scooter) StartTrip(tripID, userID string) {
	s.CurrentTrip = &Trip{
		ID:        tripID,
		UserID:    userID,
		StartTime: time.Now(),
		Direction: s.Movement.GetRandomDirection(),
	}
	s.Status = "occupied"

	utils.Info("Scooter trip started",
		zap.Int("scooter_id", s.ID),
		zap.String("trip_id", tripID),
		zap.String("user_id", userID),
		zap.Float64("direction", s.CurrentTrip.Direction),
	)
}

// EndTrip ends the current trip for this scooter
func (s *Scooter) EndTrip() {
	if s.CurrentTrip != nil {
		duration := time.Since(s.CurrentTrip.StartTime)

		utils.Info("Scooter trip ended",
			zap.Int("scooter_id", s.ID),
			zap.String("trip_id", s.CurrentTrip.ID),
			zap.String("user_id", s.CurrentTrip.UserID),
			zap.Duration("duration", duration),
		)
	}

	s.CurrentTrip = nil
	s.Status = "available"
}

// GetLocation returns the current location of the scooter
func (s *Scooter) GetLocation() Location {
	return s.Location
}

// GetStatus returns the current status of the scooter
func (s *Scooter) GetStatus() string {
	return s.Status
}

// IsAvailable returns true if the scooter is available for a new trip
func (s *Scooter) IsAvailable() bool {
	return s.Status == "available"
}

// GetLastSeen returns the last time the scooter was seen
func (s *Scooter) GetLastSeen() time.Time {
	return s.LastSeen
}
