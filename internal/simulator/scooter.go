package simulator

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"scootin-aboot/internal/config"
	"scootin-aboot/internal/logger"

	"github.com/google/uuid"
)

// UserTracker interface for tracking active users
type UserTracker interface {
	IsUserActive(userID string) bool
	MarkUserActive(userID string)
	MarkUserInactive(userID string)
	GetAvailableUsers() []string
}

// StatisticsUpdater interface for updating simulation statistics
type StatisticsUpdater interface {
	OnTripStarted()
	OnTripEnded()
}

// Scooter simulates a scooter's behavior
type Scooter struct {
	ID                int
	APIScooterID      string // Store the actual API scooter ID
	Ctx               context.Context
	Publisher         EventPublisher
	Config            *config.Config
	Movement          *Movement
	CurrentTrip       *Trip
	Location          Location
	Status            string
	LastSeen          time.Time
	UserTracker       UserTracker       // Interface for tracking active users
	StatisticsUpdater StatisticsUpdater // Interface for updating statistics
}

// Trip represents an active trip
type Trip struct {
	ID        string
	UserID    string
	StartTime time.Time
	Direction float64
}

// NewScooter creates a new scooter simulator
func NewScooter(ctx context.Context, publisher EventPublisher, id int, cfg *config.Config, userTracker UserTracker, statsUpdater StatisticsUpdater) (*Scooter, error) {
	movement := NewMovement(cfg)

	// Start with random location
	location := movement.GetRandomLocation()

	return &Scooter{
		ID:                id,
		Ctx:               ctx,
		Publisher:         publisher,
		Config:            cfg,
		Movement:          movement,
		Location:          location,
		Status:            "available",
		LastSeen:          time.Now(),
		UserTracker:       userTracker,
		StatisticsUpdater: statsUpdater,
	}, nil
}

// NewScooterFromAPI creates a new scooter simulator from API data
func NewScooterFromAPI(ctx context.Context, publisher EventPublisher, apiScooter APIScooter, cfg *config.Config, userTracker UserTracker, statsUpdater StatisticsUpdater) (*Scooter, error) {
	movement := NewMovement(cfg)

	// Use the scooter's current location from the API
	location := Location{
		Latitude:  apiScooter.Latitude,
		Longitude: apiScooter.Longitude,
	}

	return &Scooter{
		ID:                0, // Internal ID for simulation
		APIScooterID:      apiScooter.ID,
		Ctx:               ctx,
		Publisher:         publisher,
		Config:            cfg,
		Movement:          movement,
		Location:          location,
		Status:            apiScooter.Status,
		LastSeen:          time.Now(),
		UserTracker:       userTracker,
		StatisticsUpdater: statsUpdater,
	}, nil
}

// Simulate runs the scooter simulation loop
func (s *Scooter) Simulate() {
	logger.Info("Scooter simulation started",
		logger.Int("scooter_id", s.ID),
		logger.String("status", s.Status),
		logger.Float64("lat", s.Location.Latitude),
		logger.Float64("lng", s.Location.Longitude),
	)

	// Update location every 3 seconds
	ticker := time.NewTicker(3 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-s.Ctx.Done():
			logger.Info("Scooter simulation stopped",
				logger.Int("scooter_id", s.ID),
			)
			return
		case <-ticker.C:
			if s.Status == "available" {
				// Randomly decide to start a trip
				if s.shouldStartTrip() {
					s.startRandomTrip()
				}
			} else if s.Status == "occupied" {
				// Update location during trip
				s.updateLocation()

				// Randomly decide to end trip
				if s.shouldEndTrip() {
					s.EndCurrentTrip()
				}
			}
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

	// Send location update via publisher
	heading := s.CurrentTrip.Direction
	speed := 15.0 // Mock speed for now

	logger.Debug("Publishing location update event",
		logger.Int("scooter_id", s.ID),
		logger.String("trip_id", s.CurrentTrip.ID),
		logger.Float64("lat", s.Location.Latitude),
		logger.Float64("lng", s.Location.Longitude),
		logger.Float64("heading", heading),
		logger.Float64("speed", speed),
	)

	if err := s.Publisher.PublishLocationUpdated(s.Ctx, s.getScooterID(), s.CurrentTrip.ID, s.Location.Latitude, s.Location.Longitude, heading, speed); err != nil {
		logger.Error("Failed to publish location update event",
			logger.Int("scooter_id", s.ID),
			logger.String("trip_id", s.CurrentTrip.ID),
			logger.ErrorField(err),
		)
	} else {
		logger.Debug("Location update event published successfully",
			logger.Int("scooter_id", s.ID),
			logger.String("trip_id", s.CurrentTrip.ID),
			logger.Float64("lat", s.Location.Latitude),
			logger.Float64("lng", s.Location.Longitude),
		)
	}
}

// getScooterID returns the scooter ID as a string
func (s *Scooter) getScooterID() string {
	// Use the actual API scooter ID from the database
	if s.APIScooterID != "" {
		return s.APIScooterID
	}

	// Fallback to old logic for backward compatibility
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

	logger.Info("Scooter trip state changed to started",
		logger.Int("scooter_id", s.ID),
		logger.String("trip_id", tripID),
		logger.String("user_id", userID),
		logger.Float64("direction", s.CurrentTrip.Direction),
	)
}

// EndTrip ends the current trip for this scooter
func (s *Scooter) EndTrip() {
	if s.CurrentTrip != nil {
		duration := time.Since(s.CurrentTrip.StartTime)

		logger.Info("Scooter trip state changed to ended",
			logger.Int("scooter_id", s.ID),
			logger.String("trip_id", s.CurrentTrip.ID),
			logger.String("user_id", s.CurrentTrip.UserID),
			logger.Duration("duration", duration),
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

// shouldStartTrip determines if the scooter should start a trip
func (s *Scooter) shouldStartTrip() bool {
	// 60% chance every 3 seconds to start a trip when available
	return rand.Float64() < 0.6
}

// shouldEndTrip determines if the scooter should end the current trip
func (s *Scooter) shouldEndTrip() bool {
	if s.CurrentTrip == nil {
		return false
	}

	// End trip after 5-15 seconds (simulating trip duration)
	tripDuration := time.Since(s.CurrentTrip.StartTime)
	minDuration := 5 * time.Second
	maxDuration := 15 * time.Second

	// 20% chance every 3 seconds after minimum duration
	if tripDuration > minDuration {
		return rand.Float64() < 0.2 || tripDuration > maxDuration
	}

	return false
}

// startRandomTrip starts a trip with a random available user
func (s *Scooter) startRandomTrip() {
	// Get available users (not currently in trips)
	availableUsers := s.UserTracker.GetAvailableUsers()
	if len(availableUsers) == 0 {
		logger.Debug("No available users for trip",
			logger.Int("scooter_id", s.ID),
		)
		return
	}

	// Pick a random available user
	userID := availableUsers[rand.Intn(len(availableUsers))]
	tripID := uuid.New().String()

	// Mark user as active before starting trip
	s.UserTracker.MarkUserActive(userID)

	// Start trip locally
	s.StartTrip(tripID, userID)

	// Send trip start event via publisher
	logger.Info("Publishing trip start event",
		logger.Int("scooter_id", s.ID),
		logger.String("trip_id", tripID),
		logger.String("user_id", userID),
		logger.Float64("lat", s.Location.Latitude),
		logger.Float64("lng", s.Location.Longitude),
	)

	if err := s.Publisher.PublishTripStarted(s.Ctx, tripID, s.getScooterID(), userID, s.Location.Latitude, s.Location.Longitude); err != nil {
		logger.Error("Failed to publish trip start event",
			logger.Int("scooter_id", s.ID),
			logger.String("trip_id", tripID),
			logger.String("user_id", userID),
			logger.ErrorField(err),
		)
		// Revert local state and user tracking if publisher call failed
		s.EndTrip()
		s.UserTracker.MarkUserInactive(userID)
	} else {
		logger.Info("Trip start event published successfully",
			logger.Int("scooter_id", s.ID),
			logger.String("trip_id", s.CurrentTrip.ID),
			logger.String("user_id", userID),
		)

		// Update statistics - trip started
		s.StatisticsUpdater.OnTripStarted()
	}
}

// EndCurrentTrip ends the current trip
func (s *Scooter) EndCurrentTrip() {
	if s.CurrentTrip == nil {
		return
	}

	tripID := s.CurrentTrip.ID
	userID := s.CurrentTrip.UserID

	// Send trip end event via publisher
	logger.Info("Publishing trip end event",
		logger.Int("scooter_id", s.ID),
		logger.String("trip_id", tripID),
		logger.String("user_id", userID),
		logger.Float64("lat", s.Location.Latitude),
		logger.Float64("lng", s.Location.Longitude),
	)

	if err := s.Publisher.PublishTripEnded(s.Ctx, tripID, s.getScooterID(), userID, s.Location.Latitude, s.Location.Longitude, s.CurrentTrip.StartTime); err != nil {
		logger.Error("Failed to publish trip end event",
			logger.Int("scooter_id", s.ID),
			logger.String("trip_id", tripID),
			logger.String("user_id", userID),
			logger.ErrorField(err),
		)
	} else {
		logger.Info("Trip end event published successfully",
			logger.Int("scooter_id", s.ID),
			logger.String("trip_id", tripID),
			logger.String("user_id", userID),
		)
	}

	// Mark user as inactive
	s.UserTracker.MarkUserInactive(userID)

	// Update statistics - trip ended
	s.StatisticsUpdater.OnTripEnded()

	// End trip locally
	s.EndTrip()
}
