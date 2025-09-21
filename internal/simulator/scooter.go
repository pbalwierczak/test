package simulator

import (
	"context"
	"math/rand"
	"time"

	"scootin-aboot/internal/config"
	"scootin-aboot/internal/logger"

	"github.com/google/uuid"
)

type UserTracker interface {
	IsUserActive(userID string) bool
	MarkUserActive(userID string)
	MarkUserInactive(userID string)
	GetAvailableUsers() []string
}

type StatisticsUpdater interface {
	OnTripStarted()
	OnTripEnded()
}

type Scooter struct {
	ID                int
	APIScooterID      string
	Ctx               context.Context
	Publisher         EventPublisher
	Config            *config.Config
	Movement          *Movement
	CurrentTrip       *Trip
	Location          Location
	Status            string
	LastSeen          time.Time
	UserTracker       UserTracker
	StatisticsUpdater StatisticsUpdater
}

type Trip struct {
	ID        string
	UserID    string
	StartTime time.Time
	Direction float64
}

func NewScooterFromAPI(ctx context.Context, publisher EventPublisher, apiScooter APIScooter, cfg *config.Config, userTracker UserTracker, statsUpdater StatisticsUpdater) (*Scooter, error) {
	movement := NewMovement(cfg)

	location := Location{
		Latitude:  apiScooter.Latitude,
		Longitude: apiScooter.Longitude,
	}

	return &Scooter{
		ID:                0,
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

func (s *Scooter) Simulate() {
	logger.Info("Scooter simulation started",
		logger.Int("scooter_id", s.ID),
		logger.String("status", s.Status),
		logger.Float64("lat", s.Location.Latitude),
		logger.Float64("lng", s.Location.Longitude),
	)

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
				if s.shouldStartTrip() {
					s.startRandomTrip()
				}
			} else if s.Status == "occupied" {
				s.updateLocation()

				if s.shouldEndTrip() {
					s.EndCurrentTrip()
				}
			}
		}
	}
}

func (s *Scooter) updateLocation() {
	if s.Status != "occupied" || s.CurrentTrip == nil {
		return
	}

	newLocation := s.Movement.CalculateMovement(s.Location, s.CurrentTrip.Direction, 3*time.Second)

	s.Location = newLocation
	s.LastSeen = time.Now()
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

func (s *Scooter) getScooterID() string {
	return s.APIScooterID
}

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

func (s *Scooter) GetLocation() Location {
	return s.Location
}

func (s *Scooter) GetStatus() string {
	return s.Status
}

func (s *Scooter) IsAvailable() bool {
	return s.Status == "available"
}

func (s *Scooter) GetLastSeen() time.Time {
	return s.LastSeen
}

func (s *Scooter) shouldStartTrip() bool {
	// 60% chance every 3 seconds to start a trip when available
	return rand.Float64() < 0.6
}

func (s *Scooter) shouldEndTrip() bool {
	if s.CurrentTrip == nil {
		return false
	}

	tripDuration := time.Since(s.CurrentTrip.StartTime)
	minDuration := 5 * time.Second
	maxDuration := 15 * time.Second

	if tripDuration > minDuration {
		return rand.Float64() < 0.2 || tripDuration > maxDuration
	}

	return false
}

func (s *Scooter) startRandomTrip() {
	availableUsers := s.UserTracker.GetAvailableUsers()
	if len(availableUsers) == 0 {
		logger.Debug("No available users for trip",
			logger.Int("scooter_id", s.ID),
		)
		return
	}

	userID := availableUsers[rand.Intn(len(availableUsers))]
	tripID := uuid.New().String()

	s.UserTracker.MarkUserActive(userID)

	s.StartTrip(tripID, userID)
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
		s.EndTrip()
		s.UserTracker.MarkUserInactive(userID)
	} else {
		logger.Info("Trip start event published successfully",
			logger.Int("scooter_id", s.ID),
			logger.String("trip_id", s.CurrentTrip.ID),
			logger.String("user_id", userID),
		)

		s.StatisticsUpdater.OnTripStarted()
	}
}

func (s *Scooter) EndCurrentTrip() {
	if s.CurrentTrip == nil {
		return
	}

	tripID := s.CurrentTrip.ID
	userID := s.CurrentTrip.UserID
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

	s.UserTracker.MarkUserInactive(userID)

	s.StatisticsUpdater.OnTripEnded()

	s.EndTrip()
}
