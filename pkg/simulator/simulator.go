package simulator

import (
	"context"
	"fmt"
	"sync"
	"time"

	"scootin-aboot/internal/config"
	"scootin-aboot/pkg/kafka"
	"scootin-aboot/pkg/logger"
)

// Simulator orchestrates the entire simulation
type Simulator struct {
	config        *config.Config
	client        *APIClient
	publisher     EventPublisher
	users         []*User
	scooters      []*Scooter
	ctx           context.Context
	cancel        context.CancelFunc
	wg            sync.WaitGroup
	stats         *Statistics
	activeUsers   map[string]bool // Track which users are currently in trips
	activeUsersMu sync.RWMutex    // Mutex for activeUsers map
}

// Statistics tracks simulation metrics
type Statistics struct {
	mu                sync.RWMutex
	ActiveTrips       int
	CompletedTrips    int
	TotalUsers        int
	TotalScooters     int
	AvailableScooters int
	OccupiedScooters  int
	StartTime         time.Time
}

// NewSimulator creates a new simulator instance
func NewSimulator(cfg *config.Config) (*Simulator, error) {
	ctx, cancel := context.WithCancel(context.Background())

	// Create API client
	client := NewAPIClient(cfg.SimulatorServerURL, cfg.APIKey)

	// Create appropriate publisher based on configuration
	var publisher EventPublisher

	if cfg.SimulatorMode == "kafka" {
		// Create Kafka producer
		kafkaProducer, err := kafka.NewKafkaProducer(&cfg.KafkaConfig)
		if err != nil {
			cancel() // Clean up context on error
			return nil, fmt.Errorf("failed to create Kafka producer: %w", err)
		}
		publisher = NewKafkaEventPublisher(kafkaProducer)
		logger.Info("Using Kafka event publisher")
	} else {
		// Use REST publisher
		publisher = NewRESTEventPublisher(client)
		logger.Info("Using REST event publisher")
	}

	return &Simulator{
		config:      cfg,
		client:      client,
		publisher:   publisher,
		ctx:         ctx,
		cancel:      cancel,
		activeUsers: make(map[string]bool),
		stats: &Statistics{
			StartTime: time.Now(),
		},
	}, nil
}

// Start begins the simulation
func (s *Simulator) Start() error {
	logger.Info("Starting simulation",
		logger.Int("scooters", s.config.SimulatorScooters),
		logger.Int("users", s.config.SimulatorUsers),
		logger.String("server_url", s.config.SimulatorServerURL),
	)

	// Initialize scooters
	if err := s.initializeScooters(); err != nil {
		return fmt.Errorf("failed to initialize scooters: %w", err)
	}

	// Initialize users
	if err := s.initializeUsers(); err != nil {
		return fmt.Errorf("failed to initialize users: %w", err)
	}

	// Start scooter simulations
	s.startScooterSimulations()

	// Start user simulations
	s.startUserSimulations()

	// Start statistics reporting
	s.startStatisticsReporting()

	logger.Info("Simulation started successfully")
	return nil
}

// Stop gracefully stops the simulation
func (s *Simulator) Stop() {
	logger.Info("Stopping simulation gracefully...")

	// Check for active trips before shutdown
	activeTrips := s.getActiveTripsCount()
	if activeTrips > 0 {
		logger.Info("Active trips detected, ending them gracefully", logger.Int("active_trips", activeTrips))

		// End all active trips before shutting down
		s.endAllActiveTrips()

		// Wait for trips to complete with timeout
		timeout := 10 * time.Second
		start := time.Now()

		for {
			remainingTrips := s.getActiveTripsCount()
			if remainingTrips == 0 {
				logger.Info("All trips ended successfully")
				break
			}

			if time.Since(start) > timeout {
				logger.Info("Timeout reached, forcing shutdown",
					logger.Int("remaining_trips", remainingTrips),
					logger.Duration("timeout", timeout))
				break
			}

			// Wait a bit before checking again
			time.Sleep(500 * time.Millisecond)
		}
	}

	logger.Info("Cancelling context to signal shutdown to all goroutines...")
	s.cancel()

	logger.Info("Waiting for all goroutines to complete...")
	s.wg.Wait()

	// Close the publisher
	if err := s.publisher.Close(); err != nil {
		logger.Error("Error closing publisher", logger.ErrorField(err))
	}

	logger.Info("Simulation stopped gracefully - all trips completed")
}

// getActiveTripsCount returns the number of scooters currently in trips
func (s *Simulator) getActiveTripsCount() int {
	count := 0
	for _, scooter := range s.scooters {
		if scooter.Status == "occupied" {
			count++
		}
	}
	return count
}

// endAllActiveTrips ends all currently active trips
func (s *Simulator) endAllActiveTrips() {
	logger.Info("Ending all active trips...")

	for _, scooter := range s.scooters {
		if scooter.Status == "occupied" && scooter.CurrentTrip != nil {
			logger.Info("Ending trip for scooter",
				logger.Int("scooter_id", scooter.ID),
				logger.String("trip_id", scooter.CurrentTrip.ID),
				logger.String("user_id", scooter.CurrentTrip.UserID),
			)

			// End the trip by calling the scooter's EndCurrentTrip method
			scooter.EndCurrentTrip()
		}
	}

	logger.Info("All active trips ended")
}

// initializeScooters fetches existing scooters from the API and creates scooter instances
func (s *Simulator) initializeScooters() error {
	// Fetch all scooters from the API
	apiScooters, err := s.client.GetAllScooters(s.ctx)
	if err != nil {
		return fmt.Errorf("failed to fetch scooters from API: %w", err)
	}

	if len(apiScooters) == 0 {
		return fmt.Errorf("no scooters found in database")
	}

	// Limit to configured number of scooters
	maxScooters := s.config.SimulatorScooters
	if len(apiScooters) < maxScooters {
		maxScooters = len(apiScooters)
		logger.Info("Limited scooters to available count",
			logger.Int("requested", s.config.SimulatorScooters),
			logger.Int("available", len(apiScooters)),
			logger.Int("using", maxScooters))
	}

	s.scooters = make([]*Scooter, maxScooters)

	for i := 0; i < maxScooters; i++ {
		apiScooter := apiScooters[i]
		scooter, err := NewScooterFromAPI(s.ctx, s.publisher, apiScooter, s.config, s, s)
		if err != nil {
			return fmt.Errorf("failed to create scooter %s: %w", apiScooter.ID, err)
		}
		s.scooters[i] = scooter
	}

	s.stats.mu.Lock()
	s.stats.TotalScooters = len(s.scooters)
	s.stats.AvailableScooters = len(s.scooters)
	s.stats.mu.Unlock()

	logger.Info("Initialized scooters", logger.Int("count", len(s.scooters)))
	return nil
}

// initializeUsers creates and initializes user instances using seeded user IDs
func (s *Simulator) initializeUsers() error {
	// Seeded user IDs from seeds/users.sql
	seededUserIDs := []string{
		"550e8400-e29b-41d4-a716-446655440001",
		"550e8400-e29b-41d4-a716-446655440002",
		"550e8400-e29b-41d4-a716-446655440003",
		"550e8400-e29b-41d4-a716-446655440004",
		"550e8400-e29b-41d4-a716-446655440005",
		"550e8400-e29b-41d4-a716-446655440006",
		"550e8400-e29b-41d4-a716-446655440007",
		"550e8400-e29b-41d4-a716-446655440008",
		"550e8400-e29b-41d4-a716-446655440009",
		"550e8400-e29b-41d4-a716-446655440010",
	}

	// Limit to configured number of users
	maxUsers := s.config.SimulatorUsers
	if len(seededUserIDs) < maxUsers {
		maxUsers = len(seededUserIDs)
		logger.Info("Limited users to available seeded count",
			logger.Int("requested", s.config.SimulatorUsers),
			logger.Int("available", len(seededUserIDs)),
			logger.Int("using", maxUsers))
	}

	s.users = make([]*User, maxUsers)

	for i := 0; i < maxUsers; i++ {
		user, err := NewUserWithID(s.ctx, s.client, i+1, seededUserIDs[i], s.config)
		if err != nil {
			return fmt.Errorf("failed to create user %s: %w", seededUserIDs[i], err)
		}
		s.users[i] = user
	}

	s.stats.mu.Lock()
	s.stats.TotalUsers = len(s.users)
	s.stats.mu.Unlock()

	logger.Info("Initialized users", logger.Int("count", len(s.users)))
	return nil
}

// startScooterSimulations starts all scooter simulation goroutines
func (s *Simulator) startScooterSimulations() {
	for _, scooter := range s.scooters {
		s.wg.Add(1)
		go func(scooter *Scooter) {
			defer s.wg.Done()
			scooter.Simulate()
		}(scooter)
	}
}

// startUserSimulations starts all user simulation goroutines
func (s *Simulator) startUserSimulations() {
	for _, user := range s.users {
		s.wg.Add(1)
		go func(user *User) {
			defer s.wg.Done()
			user.Simulate()
		}(user)
	}
}

// startStatisticsReporting starts the statistics reporting goroutine
func (s *Simulator) startStatisticsReporting() {
	s.wg.Add(1)
	go func() {
		defer s.wg.Done()
		ticker := time.NewTicker(10 * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-s.ctx.Done():
				return
			case <-ticker.C:
				s.reportStatistics()
			}
		}
	}()
}

// reportStatistics logs current simulation statistics
func (s *Simulator) reportStatistics() {
	s.stats.mu.RLock()
	activeTrips := s.stats.ActiveTrips
	completedTrips := s.stats.CompletedTrips
	availableScooters := s.stats.AvailableScooters
	occupiedScooters := s.stats.OccupiedScooters
	totalUsers := s.stats.TotalUsers
	totalScooters := s.stats.TotalScooters
	startTime := s.stats.StartTime
	s.stats.mu.RUnlock()

	uptime := time.Since(startTime)

	logger.Info("Simulation Statistics",
		logger.Duration("uptime", uptime),
		logger.Int("active_trips", activeTrips),
		logger.Int("completed_trips", completedTrips),
		logger.Int("available_scooters", availableScooters),
		logger.Int("occupied_scooters", occupiedScooters),
		logger.Int("total_users", totalUsers),
		logger.Int("total_scooters", totalScooters),
	)
}

// UpdateStats updates simulation statistics
func (s *Simulator) UpdateStats(update func(*Statistics)) {
	s.stats.mu.Lock()
	update(s.stats)
	s.stats.mu.Unlock()
}

// IsUserActive checks if a user is currently in a trip
func (s *Simulator) IsUserActive(userID string) bool {
	s.activeUsersMu.RLock()
	defer s.activeUsersMu.RUnlock()
	return s.activeUsers[userID]
}

// MarkUserActive marks a user as active (in a trip)
func (s *Simulator) MarkUserActive(userID string) {
	s.activeUsersMu.Lock()
	defer s.activeUsersMu.Unlock()
	s.activeUsers[userID] = true
}

// MarkUserInactive marks a user as inactive (not in a trip)
func (s *Simulator) MarkUserInactive(userID string) {
	s.activeUsersMu.Lock()
	defer s.activeUsersMu.Unlock()
	delete(s.activeUsers, userID)
}

// GetAvailableUsers returns a list of user IDs that are not currently in trips
func (s *Simulator) GetAvailableUsers() []string {
	// Seeded user IDs from seeds/users.sql
	allUserIDs := []string{
		"550e8400-e29b-41d4-a716-446655440001",
		"550e8400-e29b-41d4-a716-446655440002",
		"550e8400-e29b-41d4-a716-446655440003",
		"550e8400-e29b-41d4-a716-446655440004",
		"550e8400-e29b-41d4-a716-446655440005",
		"550e8400-e29b-41d4-a716-446655440006",
		"550e8400-e29b-41d4-a716-446655440007",
		"550e8400-e29b-41d4-a716-446655440008",
		"550e8400-e29b-41d4-a716-446655440009",
		"550e8400-e29b-41d4-a716-446655440010",
	}

	s.activeUsersMu.RLock()
	defer s.activeUsersMu.RUnlock()

	var availableUsers []string
	for _, userID := range allUserIDs {
		if !s.activeUsers[userID] {
			availableUsers = append(availableUsers, userID)
		}
	}

	return availableUsers
}

// OnTripStarted is called when a scooter starts a trip
func (s *Simulator) OnTripStarted() {
	s.stats.mu.Lock()
	s.stats.ActiveTrips++
	s.stats.AvailableScooters--
	s.stats.OccupiedScooters++
	s.stats.mu.Unlock()
}

// OnTripEnded is called when a scooter ends a trip
func (s *Simulator) OnTripEnded() {
	s.stats.mu.Lock()
	s.stats.ActiveTrips--
	s.stats.CompletedTrips++
	s.stats.AvailableScooters++
	s.stats.OccupiedScooters--
	s.stats.mu.Unlock()
}
