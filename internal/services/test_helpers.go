package services

import (
	"context"
	"time"

	"scootin-aboot/internal/models"
	"scootin-aboot/internal/repository/mocks"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

// TestFixtures contains common test data
var TestFixtures = struct {
	ValidCoordinates struct {
		Latitude  float64
		Longitude float64
	}
	InvalidCoordinates struct {
		HighLat float64
		HighLng float64
		LowLat  float64
		LowLng  float64
	}
	ValidTime time.Time
}{
	ValidCoordinates: struct {
		Latitude  float64
		Longitude float64
	}{
		Latitude:  45.4215,
		Longitude: -75.6972,
	},
	InvalidCoordinates: struct {
		HighLat float64
		HighLng float64
		LowLat  float64
		LowLng  float64
	}{
		HighLat: 91.0,
		HighLng: 181.0,
		LowLat:  -91.0,
		LowLng:  -181.0,
	},
	ValidTime: time.Now(),
}

// TestScooterBuilder provides a fluent interface for building test scooters
type TestScooterBuilder struct {
	scooter *models.Scooter
}

// NewTestScooterBuilder creates a new scooter builder with default values
func NewTestScooterBuilder() *TestScooterBuilder {
	return &TestScooterBuilder{
		scooter: &models.Scooter{
			ID:               uuid.New(),
			Status:           models.ScooterStatusAvailable,
			CurrentLatitude:  TestFixtures.ValidCoordinates.Latitude,
			CurrentLongitude: TestFixtures.ValidCoordinates.Longitude,
			LastSeen:         TestFixtures.ValidTime,
			CreatedAt:        TestFixtures.ValidTime,
			UpdatedAt:        TestFixtures.ValidTime,
		},
	}
}

func (b *TestScooterBuilder) WithID(id uuid.UUID) *TestScooterBuilder {
	b.scooter.ID = id
	return b
}

func (b *TestScooterBuilder) WithStatus(status models.ScooterStatus) *TestScooterBuilder {
	b.scooter.Status = status
	return b
}

func (b *TestScooterBuilder) WithLocation(lat, lng float64) *TestScooterBuilder {
	b.scooter.CurrentLatitude = lat
	b.scooter.CurrentLongitude = lng
	return b
}

func (b *TestScooterBuilder) WithLastSeen(lastSeen time.Time) *TestScooterBuilder {
	b.scooter.LastSeen = lastSeen
	return b
}

func (b *TestScooterBuilder) Build() *models.Scooter {
	return b.scooter
}

// TestTripBuilder provides a fluent interface for building test trips
type TestTripBuilder struct {
	trip *models.Trip
}

// NewTestTripBuilder creates a new trip builder with default values
func NewTestTripBuilder() *TestTripBuilder {
	return &TestTripBuilder{
		trip: &models.Trip{
			ID:             uuid.New(),
			ScooterID:      uuid.New(),
			UserID:         uuid.New(),
			Status:         models.TripStatusActive,
			StartTime:      TestFixtures.ValidTime,
			StartLatitude:  TestFixtures.ValidCoordinates.Latitude,
			StartLongitude: TestFixtures.ValidCoordinates.Longitude,
		},
	}
}

func (b *TestTripBuilder) WithID(id uuid.UUID) *TestTripBuilder {
	b.trip.ID = id
	return b
}

func (b *TestTripBuilder) WithScooterID(scooterID uuid.UUID) *TestTripBuilder {
	b.trip.ScooterID = scooterID
	return b
}

func (b *TestTripBuilder) WithUserID(userID uuid.UUID) *TestTripBuilder {
	b.trip.UserID = userID
	return b
}

func (b *TestTripBuilder) WithStatus(status models.TripStatus) *TestTripBuilder {
	b.trip.Status = status
	return b
}

func (b *TestTripBuilder) WithStartLocation(lat, lng float64) *TestTripBuilder {
	b.trip.StartLatitude = lat
	b.trip.StartLongitude = lng
	return b
}

func (b *TestTripBuilder) WithEndLocation(lat, lng float64) *TestTripBuilder {
	b.trip.EndLatitude = &lat
	b.trip.EndLongitude = &lng
	return b
}

func (b *TestTripBuilder) WithStartTime(startTime time.Time) *TestTripBuilder {
	b.trip.StartTime = startTime
	return b
}

func (b *TestTripBuilder) WithEndTime(endTime time.Time) *TestTripBuilder {
	b.trip.EndTime = &endTime
	return b
}

func (b *TestTripBuilder) Build() *models.Trip {
	return b.trip
}

// TestUserBuilder provides a fluent interface for building test users
type TestUserBuilder struct {
	user *models.User
}

// NewTestUserBuilder creates a new user builder with default values
func NewTestUserBuilder() *TestUserBuilder {
	return &TestUserBuilder{
		user: &models.User{
			ID:        uuid.New(),
			CreatedAt: TestFixtures.ValidTime,
			UpdatedAt: TestFixtures.ValidTime,
		},
	}
}

func (b *TestUserBuilder) WithID(id uuid.UUID) *TestUserBuilder {
	b.user.ID = id
	return b
}

func (b *TestUserBuilder) Build() *models.User {
	return b.user
}

// MockSetup provides utilities for setting up mocks
type MockSetup struct{}

// SetupScooterServiceMocks creates all necessary mocks for scooter service tests
func (m *MockSetup) SetupScooterServiceMocks() (*mocks.MockScooterRepository, *mocks.MockTripRepository, *mocks.MockLocationUpdateRepository, *mocks.MockUnitOfWork) {
	return &mocks.MockScooterRepository{},
		&mocks.MockTripRepository{},
		&mocks.MockLocationUpdateRepository{},
		&mocks.MockUnitOfWork{}
}

// SetupTripServiceMocks creates all necessary mocks for trip service tests
func (m *MockSetup) SetupTripServiceMocks() (*mocks.MockTripRepository, *mocks.MockScooterRepository, *mocks.MockUserRepository, *mocks.MockLocationUpdateRepository, *mocks.MockUnitOfWork) {
	return &mocks.MockTripRepository{},
		&mocks.MockScooterRepository{},
		&mocks.MockUserRepository{},
		&mocks.MockLocationUpdateRepository{},
		&mocks.MockUnitOfWork{}
}

// SetupBasicUnitOfWork sets up a basic unit of work with transaction
func (m *MockSetup) SetupBasicUnitOfWork(unitOfWork *mocks.MockUnitOfWork, tripRepo *mocks.MockTripRepository, scooterRepo *mocks.MockScooterRepository, userRepo *mocks.MockUserRepository, locationRepo *mocks.MockLocationUpdateRepository) *mocks.MockUnitOfWorkTx {
	mockTx := &mocks.MockUnitOfWorkTx{}
	unitOfWork.On("Begin", mock.Anything).Return(mockTx, nil)
	mockTx.On("UserRepository").Return(userRepo)
	mockTx.On("TripRepository").Return(tripRepo)
	mockTx.On("ScooterRepository").Return(scooterRepo)
	mockTx.On("LocationUpdateRepository").Return(locationRepo)
	mockTx.On("Commit").Return(nil)
	mockTx.On("Rollback").Return(nil)
	return mockTx
}

// SetupScooterServiceUnitOfWork sets up unit of work specifically for scooter service
func (m *MockSetup) SetupScooterServiceUnitOfWork(unitOfWork *mocks.MockUnitOfWork, scooterRepo *mocks.MockScooterRepository, locationRepo *mocks.MockLocationUpdateRepository) *mocks.MockUnitOfWorkTx {
	mockTx := &mocks.MockUnitOfWorkTx{}
	unitOfWork.On("Begin", mock.Anything).Return(mockTx, nil)
	mockTx.On("ScooterRepository").Return(scooterRepo)
	mockTx.On("LocationUpdateRepository").Return(locationRepo)
	mockTx.On("Commit").Return(nil)
	mockTx.On("Rollback").Return(nil)
	return mockTx
}

// CreateTestScooterService creates a scooter service with mocks
func (m *MockSetup) CreateTestScooterService() (ScooterService, *mocks.MockScooterRepository, *mocks.MockTripRepository, *mocks.MockLocationUpdateRepository, *mocks.MockUnitOfWork) {
	scooterRepo, tripRepo, locationRepo, unitOfWork := m.SetupScooterServiceMocks()
	service := NewScooterService(scooterRepo, tripRepo, locationRepo, unitOfWork)
	return service, scooterRepo, tripRepo, locationRepo, unitOfWork
}

// CreateTestTripService creates a trip service with mocks
func (m *MockSetup) CreateTestTripService() (TripService, *mocks.MockTripRepository, *mocks.MockScooterRepository, *mocks.MockUserRepository, *mocks.MockLocationUpdateRepository, *mocks.MockUnitOfWork) {
	tripRepo, scooterRepo, userRepo, locationRepo, unitOfWork := m.SetupTripServiceMocks()
	service := NewTripService(tripRepo, scooterRepo, userRepo, locationRepo, unitOfWork)
	return service, tripRepo, scooterRepo, userRepo, locationRepo, unitOfWork
}

// TestContext provides a test context
func TestContext() context.Context {
	return context.Background()
}
