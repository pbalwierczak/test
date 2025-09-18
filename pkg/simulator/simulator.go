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
	utils.Info("Stopping simulation...")
	s.cancel()
	s.wg.Wait()
	utils.Info("Simulation stopped")
}

// initializeScooters creates and initializes scooter instances
func (s *Simulator) initializeScooters() error {
	s.scooters = make([]*Scooter, s.config.SimulatorScooters)

	for i := 0; i < s.config.SimulatorScooters; i++ {
		scooter, err := NewScooter(s.ctx, s.client, i+1, s.config)
		if err != nil {
			return fmt.Errorf("failed to create scooter %d: %w", i+1, err)
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

// initializeUsers creates and initializes user instances
func (s *Simulator) initializeUsers() error {
	s.users = make([]*User, s.config.SimulatorUsers)

	for i := 0; i < s.config.SimulatorUsers; i++ {
		user, err := NewUser(s.ctx, s.client, i+1, s.config)
		if err != nil {
			return fmt.Errorf("failed to create user %d: %w", i+1, err)
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
