package simulator

import (
	"context"
	"fmt"
	"sync"
	"time"

	"scootin-aboot/internal/config"
	"scootin-aboot/pkg/utils"

	"go.uber.org/zap"
)

// Simulator orchestrates the entire simulation
type Simulator struct {
	config   *config.Config
	client   *APIClient
	users    []*User
	scooters []*Scooter
	ctx      context.Context
	cancel   context.CancelFunc
	wg       sync.WaitGroup
	stats    *Statistics
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
func NewSimulator(cfg *config.Config) *Simulator {
	ctx, cancel := context.WithCancel(context.Background())

	return &Simulator{
		config: cfg,
		client: NewAPIClient(cfg.SimulatorServerURL, cfg.APIKey),
		ctx:    ctx,
		cancel: cancel,
		stats: &Statistics{
			StartTime: time.Now(),
		},
	}
}

// Start begins the simulation
func (s *Simulator) Start() error {
	utils.Info("Starting simulation",
		zap.Int("scooters", s.config.SimulatorScooters),
		zap.Int("users", s.config.SimulatorUsers),
		zap.String("server_url", s.config.SimulatorServerURL),
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

	utils.Info("Simulation started successfully")
	return nil
}

// Stop gracefully stops the simulation
func (s *Simulator) Stop() {
	utils.Info("Stopping simulation gracefully...")

	// Check for active trips before shutdown
	activeTrips := s.getActiveTripsCount()
	if activeTrips > 0 {
		utils.Info("Active trips detected, waiting for completion", zap.Int("active_trips", activeTrips))
	}

	utils.Info("Cancelling context to signal shutdown to all goroutines...")
	s.cancel()

	utils.Info("Waiting for all goroutines to complete...")
	s.wg.Wait()

	utils.Info("Simulation stopped gracefully - all trips completed")
}

// getActiveTripsCount returns the number of users currently in trips
func (s *Simulator) getActiveTripsCount() int {
	count := 0
	for _, user := range s.users {
		if user.getInTrip() {
			count++
		}
	}
	return count
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
		utils.Info("Limited scooters to available count",
			zap.Int("requested", s.config.SimulatorScooters),
			zap.Int("available", len(apiScooters)),
			zap.Int("using", maxScooters))
	}

	s.scooters = make([]*Scooter, maxScooters)

	for i := 0; i < maxScooters; i++ {
		apiScooter := apiScooters[i]
		scooter, err := NewScooterFromAPI(s.ctx, s.client, apiScooter, s.config)
		if err != nil {
			return fmt.Errorf("failed to create scooter %s: %w", apiScooter.ID, err)
		}
		s.scooters[i] = scooter
	}

	s.stats.mu.Lock()
	s.stats.TotalScooters = len(s.scooters)
	s.stats.AvailableScooters = len(s.scooters)
	s.stats.mu.Unlock()

	utils.Info("Initialized scooters", zap.Int("count", len(s.scooters)))
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
		utils.Info("Limited users to available seeded count",
			zap.Int("requested", s.config.SimulatorUsers),
			zap.Int("available", len(seededUserIDs)),
			zap.Int("using", maxUsers))
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

	utils.Info("Initialized users", zap.Int("count", len(s.users)))
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

	utils.Info("Simulation Statistics",
		zap.Duration("uptime", uptime),
		zap.Int("active_trips", activeTrips),
		zap.Int("completed_trips", completedTrips),
		zap.Int("available_scooters", availableScooters),
		zap.Int("occupied_scooters", occupiedScooters),
		zap.Int("total_users", totalUsers),
		zap.Int("total_scooters", totalScooters),
	)
}

// UpdateStats updates simulation statistics
func (s *Simulator) UpdateStats(update func(*Statistics)) {
	s.stats.mu.Lock()
	update(s.stats)
	s.stats.mu.Unlock()
}
