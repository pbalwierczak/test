package simulator

import (
	"context"
	"fmt"
	"sync"
	"time"

	"scootin-aboot/internal/config"
	"scootin-aboot/internal/events"
	"scootin-aboot/internal/logger"
)

// SeededUserIDs contains the hardcoded user IDs from seeds/users.sql
// These should be kept in sync with the database seed data
var SeededUserIDs = []string{
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

func NewSimulator(cfg *config.Config) (*Simulator, error) {
	ctx, cancel := context.WithCancel(context.Background())

	client := NewAPIClient(cfg.SimulatorServerURL, cfg.APIKey)

	kafkaProducer, err := events.NewKafkaProducer(&cfg.KafkaConfig)
	if err != nil {
		cancel()
		return nil, fmt.Errorf("failed to create Kafka producer: %w", err)
	}
	publisher := NewKafkaEventPublisher(kafkaProducer)
	logger.Info("Using Kafka event publisher")

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

func (s *Simulator) Start() error {
	logger.Info("Starting simulation",
		logger.Int("scooters", s.config.SimulatorScooters),
		logger.Int("users", s.config.SimulatorUsers),
		logger.String("server_url", s.config.SimulatorServerURL),
	)

	if err := s.initializeScooters(); err != nil {
		return fmt.Errorf("failed to initialize scooters: %w", err)
	}

	if err := s.initializeUsers(); err != nil {
		return fmt.Errorf("failed to initialize users: %w", err)
	}

	s.startScooterSimulations()
	s.startUserSimulations()
	s.startStatisticsReporting()

	logger.Info("Simulation started successfully")
	return nil
}

func (s *Simulator) Stop() {
	logger.Info("Stopping simulation gracefully...")

	activeTrips := s.getActiveTripsCount()
	if activeTrips > 0 {
		logger.Info("Active trips detected, ending them gracefully", logger.Int("active_trips", activeTrips))

		s.endAllActiveTrips()

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

			time.Sleep(500 * time.Millisecond)
		}
	}

	logger.Info("Cancelling context to signal shutdown to all goroutines...")
	s.cancel()

	logger.Info("Waiting for all goroutines to complete...")
	s.wg.Wait()

	if err := s.publisher.Close(); err != nil {
		logger.Error("Error closing publisher", logger.ErrorField(err))
	}

	logger.Info("Simulation stopped gracefully - all trips completed")
}

func (s *Simulator) getActiveTripsCount() int {
	count := 0
	for _, scooter := range s.scooters {
		if scooter.Status == "occupied" {
			count++
		}
	}
	return count
}

func (s *Simulator) endAllActiveTrips() {
	logger.Info("Ending all active trips...")

	for _, scooter := range s.scooters {
		if scooter.Status == "occupied" && scooter.CurrentTrip != nil {
			logger.Info("Ending trip for scooter",
				logger.Int("scooter_id", scooter.ID),
				logger.String("trip_id", scooter.CurrentTrip.ID),
				logger.String("user_id", scooter.CurrentTrip.UserID),
			)

			scooter.EndCurrentTrip()
		}
	}

	logger.Info("All active trips ended")
}

func (s *Simulator) initializeScooters() error {
	apiScooters, err := s.client.GetAllScooters(s.ctx)
	if err != nil {
		return fmt.Errorf("failed to fetch scooters from API: %w", err)
	}

	if len(apiScooters) == 0 {
		return fmt.Errorf("no scooters found in database")
	}

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

	s.withStatsLock(func(stats *Statistics) {
		stats.TotalScooters = len(s.scooters)
		stats.AvailableScooters = len(s.scooters)
	})

	logger.Info("Initialized scooters", logger.Int("count", len(s.scooters)))
	return nil
}

func (s *Simulator) initializeUsers() error {
	maxUsers := s.config.SimulatorUsers
	if len(SeededUserIDs) < maxUsers {
		maxUsers = len(SeededUserIDs)
		logger.Info("Limited users to available seeded count",
			logger.Int("requested", s.config.SimulatorUsers),
			logger.Int("available", len(SeededUserIDs)),
			logger.Int("using", maxUsers))
	}

	s.users = make([]*User, maxUsers)

	for i := 0; i < maxUsers; i++ {
		user, err := NewUserWithID(s.ctx, s.client, i+1, SeededUserIDs[i], s.config)
		if err != nil {
			return fmt.Errorf("failed to create user %s: %w", SeededUserIDs[i], err)
		}
		s.users[i] = user
	}

	s.withStatsLock(func(stats *Statistics) {
		stats.TotalUsers = len(s.users)
	})

	logger.Info("Initialized users", logger.Int("count", len(s.users)))
	return nil
}

func (s *Simulator) startScooterSimulations() {
	for _, scooter := range s.scooters {
		s.wg.Add(1)
		go func(scooter *Scooter) {
			defer s.wg.Done()
			scooter.Simulate()
		}(scooter)
	}
}

func (s *Simulator) startUserSimulations() {
	for _, user := range s.users {
		s.wg.Add(1)
		go func(user *User) {
			defer s.wg.Done()
			user.Simulate()
		}(user)
	}
}

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

func (s *Simulator) reportStatistics() {
	var activeTrips, completedTrips, availableScooters, occupiedScooters, totalUsers, totalScooters int
	var startTime time.Time

	s.withStatsReadLock(func(stats *Statistics) {
		activeTrips = stats.ActiveTrips
		completedTrips = stats.CompletedTrips
		availableScooters = stats.AvailableScooters
		occupiedScooters = stats.OccupiedScooters
		totalUsers = stats.TotalUsers
		totalScooters = stats.TotalScooters
		startTime = stats.StartTime
	})

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

func (s *Simulator) UpdateStats(update func(*Statistics)) {
	s.stats.mu.Lock()
	update(s.stats)
	s.stats.mu.Unlock()
}

// withStatsLock executes a function with stats mutex locked for writing
func (s *Simulator) withStatsLock(fn func(*Statistics)) {
	s.stats.mu.Lock()
	fn(s.stats)
	s.stats.mu.Unlock()
}

// withStatsReadLock executes a function with stats mutex locked for reading
func (s *Simulator) withStatsReadLock(fn func(*Statistics)) {
	s.stats.mu.RLock()
	fn(s.stats)
	s.stats.mu.RUnlock()
}

// withActiveUsersReadLock executes a function with activeUsers mutex locked for reading
func (s *Simulator) withActiveUsersReadLock(fn func(map[string]bool)) {
	s.activeUsersMu.RLock()
	fn(s.activeUsers)
	s.activeUsersMu.RUnlock()
}

// withActiveUsersWriteLock executes a function with activeUsers mutex locked for writing
func (s *Simulator) withActiveUsersWriteLock(fn func(map[string]bool)) {
	s.activeUsersMu.Lock()
	fn(s.activeUsers)
	s.activeUsersMu.Unlock()
}

func (s *Simulator) IsUserActive(userID string) bool {
	var isActive bool
	s.withActiveUsersReadLock(func(activeUsers map[string]bool) {
		isActive = activeUsers[userID]
	})
	return isActive
}

func (s *Simulator) MarkUserActive(userID string) {
	s.withActiveUsersWriteLock(func(activeUsers map[string]bool) {
		activeUsers[userID] = true
	})
}

func (s *Simulator) MarkUserInactive(userID string) {
	s.withActiveUsersWriteLock(func(activeUsers map[string]bool) {
		delete(activeUsers, userID)
	})
}

func (s *Simulator) GetAvailableUsers() []string {
	return s.filterAvailableUsers(SeededUserIDs)
}

// filterAvailableUsers returns user IDs that are not currently active
func (s *Simulator) filterAvailableUsers(userIDs []string) []string {
	var availableUsers []string
	s.withActiveUsersReadLock(func(activeUsers map[string]bool) {
		for _, userID := range userIDs {
			if !activeUsers[userID] {
				availableUsers = append(availableUsers, userID)
			}
		}
	})
	return availableUsers
}

func (s *Simulator) OnTripStarted() {
	s.withStatsLock(func(stats *Statistics) {
		stats.ActiveTrips++
		stats.AvailableScooters--
		stats.OccupiedScooters++
	})
}

func (s *Simulator) OnTripEnded() {
	s.withStatsLock(func(stats *Statistics) {
		stats.ActiveTrips--
		stats.CompletedTrips++
		stats.AvailableScooters++
		stats.OccupiedScooters--
	})
}
