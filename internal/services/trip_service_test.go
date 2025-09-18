package services

import (
	"context"
	"testing"
	"time"

	"scootin-aboot/internal/models"
	"scootin-aboot/internal/repository/mocks"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestTripService_StartTrip(t *testing.T) {
	tests := []struct {
		name          string
		scooterID     uuid.UUID
		userID        uuid.UUID
		lat           float64
		lng           float64
		setupMocks    func(*mocks.MockTripRepository, *mocks.MockScooterRepository, *mocks.MockUserRepository, *mocks.MockLocationUpdateRepository)
		expectedError string
		expectedTrip  *models.Trip
	}{
		{
			name:      "successful trip start",
			scooterID: uuid.New(),
			userID:    uuid.New(),
			lat:       45.4215,
			lng:       -75.6972,
			setupMocks: func(tripRepo *mocks.MockTripRepository, scooterRepo *mocks.MockScooterRepository, userRepo *mocks.MockUserRepository, locationRepo *mocks.MockLocationUpdateRepository) {
				// User exists
				userRepo.On("GetByID", mock.Anything, mock.Anything).Return(&models.User{ID: uuid.New()}, nil)

				// User has no active trip
				tripRepo.On("GetActiveByUserID", mock.Anything, mock.Anything).Return(nil, nil)

				// Scooter exists and is available
				scooter := &models.Scooter{
					ID:               uuid.New(),
					Status:           models.ScooterStatusAvailable,
					CurrentLatitude:  45.4215,
					CurrentLongitude: -75.6972,
				}
				scooterRepo.On("GetByIDForUpdate", mock.Anything, mock.Anything).Return(scooter, nil)

				// Scooter has no active trip
				tripRepo.On("GetActiveByScooterID", mock.Anything, mock.Anything).Return(nil, nil)

				// Trip creation succeeds
				tripRepo.On("Create", mock.Anything, mock.Anything).Return(nil)

				// Scooter status update succeeds
				scooterRepo.On("UpdateStatusWithCheck", mock.Anything, mock.Anything, models.ScooterStatusOccupied, models.ScooterStatusAvailable).Return(nil)

				// Location update succeeds
				scooterRepo.On("UpdateLocation", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)
			},
			expectedError: "",
		},
		{
			name:      "invalid coordinates - latitude too high",
			scooterID: uuid.New(),
			userID:    uuid.New(),
			lat:       91.0,
			lng:       -75.6972,
			setupMocks: func(tripRepo *mocks.MockTripRepository, scooterRepo *mocks.MockScooterRepository, userRepo *mocks.MockUserRepository, locationRepo *mocks.MockLocationUpdateRepository) {
				// No mocks needed - validation fails before any repository calls
			},
			expectedError: "invalid coordinates: invalid latitude: must be between -90 and 90",
		},
		{
			name:      "user not found",
			scooterID: uuid.New(),
			userID:    uuid.New(),
			lat:       45.4215,
			lng:       -75.6972,
			setupMocks: func(tripRepo *mocks.MockTripRepository, scooterRepo *mocks.MockScooterRepository, userRepo *mocks.MockUserRepository, locationRepo *mocks.MockLocationUpdateRepository) {
				userRepo.On("GetByID", mock.Anything, mock.Anything).Return(nil, nil)
			},
			expectedError: "user not found",
		},
		{
			name:      "user already has active trip",
			scooterID: uuid.New(),
			userID:    uuid.New(),
			lat:       45.4215,
			lng:       -75.6972,
			setupMocks: func(tripRepo *mocks.MockTripRepository, scooterRepo *mocks.MockScooterRepository, userRepo *mocks.MockUserRepository, locationRepo *mocks.MockLocationUpdateRepository) {
				userRepo.On("GetByID", mock.Anything, mock.Anything).Return(&models.User{ID: uuid.New()}, nil)
				tripRepo.On("GetActiveByUserID", mock.Anything, mock.Anything).Return(&models.Trip{ID: uuid.New()}, nil)
			},
			expectedError: "user already has an active trip",
		},
		{
			name:      "scooter not found",
			scooterID: uuid.New(),
			userID:    uuid.New(),
			lat:       45.4215,
			lng:       -75.6972,
			setupMocks: func(tripRepo *mocks.MockTripRepository, scooterRepo *mocks.MockScooterRepository, userRepo *mocks.MockUserRepository, locationRepo *mocks.MockLocationUpdateRepository) {
				userRepo.On("GetByID", mock.Anything, mock.Anything).Return(&models.User{ID: uuid.New()}, nil)
				tripRepo.On("GetActiveByUserID", mock.Anything, mock.Anything).Return(nil, nil)
				scooterRepo.On("GetByIDForUpdate", mock.Anything, mock.Anything).Return(nil, nil)
			},
			expectedError: "scooter not found",
		},
		{
			name:      "scooter not available",
			scooterID: uuid.New(),
			userID:    uuid.New(),
			lat:       45.4215,
			lng:       -75.6972,
			setupMocks: func(tripRepo *mocks.MockTripRepository, scooterRepo *mocks.MockScooterRepository, userRepo *mocks.MockUserRepository, locationRepo *mocks.MockLocationUpdateRepository) {
				userRepo.On("GetByID", mock.Anything, mock.Anything).Return(&models.User{ID: uuid.New()}, nil)
				tripRepo.On("GetActiveByUserID", mock.Anything, mock.Anything).Return(nil, nil)

				scooter := &models.Scooter{
					ID:               uuid.New(),
					Status:           models.ScooterStatusOccupied,
					CurrentLatitude:  45.4215,
					CurrentLongitude: -75.6972,
				}
				scooterRepo.On("GetByIDForUpdate", mock.Anything, mock.Anything).Return(scooter, nil)
			},
			expectedError: "scooter is not available",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create mock repositories
			tripRepo := &mocks.MockTripRepository{}
			scooterRepo := &mocks.MockScooterRepository{}
			userRepo := &mocks.MockUserRepository{}
			locationRepo := &mocks.MockLocationUpdateRepository{}

			// Setup mocks
			tt.setupMocks(tripRepo, scooterRepo, userRepo, locationRepo)

			// Create service
			service := NewTripService(tripRepo, scooterRepo, userRepo, locationRepo)

			// Call method
			trip, err := service.StartTrip(context.Background(), tt.scooterID, tt.userID, tt.lat, tt.lng)

			// Assertions
			if tt.expectedError != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError)
				assert.Nil(t, trip)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, trip)
				assert.Equal(t, tt.scooterID, trip.ScooterID)
				assert.Equal(t, tt.userID, trip.UserID)
				assert.Equal(t, tt.lat, trip.StartLatitude)
				assert.Equal(t, tt.lng, trip.StartLongitude)
				assert.Equal(t, models.TripStatusActive, trip.Status)
			}

			// Verify all expectations were met
			tripRepo.AssertExpectations(t)
			scooterRepo.AssertExpectations(t)
			userRepo.AssertExpectations(t)
			locationRepo.AssertExpectations(t)
		})
	}
}

func TestTripService_EndTrip(t *testing.T) {
	tests := []struct {
		name          string
		scooterID     uuid.UUID
		lat           float64
		lng           float64
		setupMocks    func(*mocks.MockTripRepository, *mocks.MockScooterRepository, *mocks.MockUserRepository, *mocks.MockLocationUpdateRepository)
		expectedError string
	}{
		{
			name:      "successful trip end",
			scooterID: uuid.New(),
			lat:       45.4216,
			lng:       -75.6973,
			setupMocks: func(tripRepo *mocks.MockTripRepository, scooterRepo *mocks.MockScooterRepository, userRepo *mocks.MockUserRepository, locationRepo *mocks.MockLocationUpdateRepository) {
				// Active trip exists
				trip := &models.Trip{
					ID:             uuid.New(),
					ScooterID:      uuid.New(),
					UserID:         uuid.New(),
					StartTime:      time.Now().Add(-10 * time.Minute),
					Status:         models.TripStatusActive,
					StartLatitude:  45.4215,
					StartLongitude: -75.6972,
				}
				tripRepo.On("GetActiveByScooterID", mock.Anything, mock.Anything).Return(trip, nil)

				// Trip end succeeds
				tripRepo.On("EndTrip", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)

				// Scooter status update succeeds
				scooterRepo.On("UpdateStatusWithCheck", mock.Anything, mock.Anything, models.ScooterStatusAvailable, models.ScooterStatusOccupied).Return(nil)

				// Location update succeeds
				scooterRepo.On("UpdateLocation", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)
			},
			expectedError: "",
		},
		{
			name:      "no active trip found",
			scooterID: uuid.New(),
			lat:       45.4216,
			lng:       -75.6973,
			setupMocks: func(tripRepo *mocks.MockTripRepository, scooterRepo *mocks.MockScooterRepository, userRepo *mocks.MockUserRepository, locationRepo *mocks.MockLocationUpdateRepository) {
				tripRepo.On("GetActiveByScooterID", mock.Anything, mock.Anything).Return(nil, nil)
			},
			expectedError: "no active trip found for scooter",
		},
		{
			name:      "invalid coordinates",
			scooterID: uuid.New(),
			lat:       91.0,
			lng:       -75.6973,
			setupMocks: func(tripRepo *mocks.MockTripRepository, scooterRepo *mocks.MockScooterRepository, userRepo *mocks.MockUserRepository, locationRepo *mocks.MockLocationUpdateRepository) {
				// No mocks needed - validation fails before any repository calls
			},
			expectedError: "invalid coordinates: invalid latitude: must be between -90 and 90",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create mock repositories
			tripRepo := &mocks.MockTripRepository{}
			scooterRepo := &mocks.MockScooterRepository{}
			userRepo := &mocks.MockUserRepository{}
			locationRepo := &mocks.MockLocationUpdateRepository{}

			// Setup mocks
			tt.setupMocks(tripRepo, scooterRepo, userRepo, locationRepo)

			// Create service
			service := NewTripService(tripRepo, scooterRepo, userRepo, locationRepo)

			// Call method
			trip, err := service.EndTrip(context.Background(), tt.scooterID, tt.lat, tt.lng)

			// Assertions
			if tt.expectedError != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError)
				assert.Nil(t, trip)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, trip)
				assert.Equal(t, models.TripStatusCompleted, trip.Status)
				assert.NotNil(t, trip.EndTime)
				assert.NotNil(t, trip.EndLatitude)
				assert.NotNil(t, trip.EndLongitude)
			}

			// Verify all expectations were met
			tripRepo.AssertExpectations(t)
			scooterRepo.AssertExpectations(t)
			userRepo.AssertExpectations(t)
			locationRepo.AssertExpectations(t)
		})
	}
}

func TestTripService_UpdateLocation(t *testing.T) {
	tests := []struct {
		name          string
		scooterID     uuid.UUID
		lat           float64
		lng           float64
		setupMocks    func(*mocks.MockTripRepository, *mocks.MockScooterRepository, *mocks.MockUserRepository, *mocks.MockLocationUpdateRepository)
		expectedError string
	}{
		{
			name:      "successful location update",
			scooterID: uuid.New(),
			lat:       45.4216,
			lng:       -75.6973,
			setupMocks: func(tripRepo *mocks.MockTripRepository, scooterRepo *mocks.MockScooterRepository, userRepo *mocks.MockUserRepository, locationRepo *mocks.MockLocationUpdateRepository) {
				// Active trip exists
				trip := &models.Trip{
					ID:        uuid.New(),
					ScooterID: uuid.New(),
					UserID:    uuid.New(),
					Status:    models.TripStatusActive,
				}
				tripRepo.On("GetActiveByScooterID", mock.Anything, mock.Anything).Return(trip, nil)

				// Location update creation succeeds
				locationRepo.On("Create", mock.Anything, mock.Anything).Return(nil)

				// Scooter location update succeeds
				scooterRepo.On("UpdateLocation", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)
			},
			expectedError: "",
		},
		{
			name:      "no active trip found",
			scooterID: uuid.New(),
			lat:       45.4216,
			lng:       -75.6973,
			setupMocks: func(tripRepo *mocks.MockTripRepository, scooterRepo *mocks.MockScooterRepository, userRepo *mocks.MockUserRepository, locationRepo *mocks.MockLocationUpdateRepository) {
				tripRepo.On("GetActiveByScooterID", mock.Anything, mock.Anything).Return(nil, nil)
			},
			expectedError: "no active trip found for scooter",
		},
		{
			name:      "invalid coordinates",
			scooterID: uuid.New(),
			lat:       91.0,
			lng:       -75.6973,
			setupMocks: func(tripRepo *mocks.MockTripRepository, scooterRepo *mocks.MockScooterRepository, userRepo *mocks.MockUserRepository, locationRepo *mocks.MockLocationUpdateRepository) {
				// No mocks needed - validation fails before any repository calls
			},
			expectedError: "invalid coordinates: invalid latitude: must be between -90 and 90",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create mock repositories
			tripRepo := &mocks.MockTripRepository{}
			scooterRepo := &mocks.MockScooterRepository{}
			userRepo := &mocks.MockUserRepository{}
			locationRepo := &mocks.MockLocationUpdateRepository{}

			// Setup mocks
			tt.setupMocks(tripRepo, scooterRepo, userRepo, locationRepo)

			// Create service
			service := NewTripService(tripRepo, scooterRepo, userRepo, locationRepo)

			// Call method
			err := service.UpdateLocation(context.Background(), tt.scooterID, tt.lat, tt.lng)

			// Assertions
			if tt.expectedError != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError)
			} else {
				assert.NoError(t, err)
			}

			// Verify all expectations were met
			tripRepo.AssertExpectations(t)
			scooterRepo.AssertExpectations(t)
			userRepo.AssertExpectations(t)
			locationRepo.AssertExpectations(t)
		})
	}
}

func TestValidateCoordinates(t *testing.T) {
	tests := []struct {
		name        string
		lat         float64
		lng         float64
		expectError bool
		errorMsg    string
	}{
		{
			name:        "valid coordinates",
			lat:         45.4215,
			lng:         -75.6972,
			expectError: false,
		},
		{
			name:        "latitude too high",
			lat:         91.0,
			lng:         -75.6972,
			expectError: true,
			errorMsg:    "invalid latitude: must be between -90 and 90",
		},
		{
			name:        "latitude too low",
			lat:         -91.0,
			lng:         -75.6972,
			expectError: true,
			errorMsg:    "invalid latitude: must be between -90 and 90",
		},
		{
			name:        "longitude too high",
			lat:         45.4215,
			lng:         181.0,
			expectError: true,
			errorMsg:    "invalid longitude: must be between -180 and 180",
		},
		{
			name:        "longitude too low",
			lat:         45.4215,
			lng:         -181.0,
			expectError: true,
			errorMsg:    "invalid longitude: must be between -180 and 180",
		},
		{
			name:        "edge case - latitude 90",
			lat:         90.0,
			lng:         -75.6972,
			expectError: false,
		},
		{
			name:        "edge case - latitude -90",
			lat:         -90.0,
			lng:         -75.6972,
			expectError: false,
		},
		{
			name:        "edge case - longitude 180",
			lat:         45.4215,
			lng:         180.0,
			expectError: false,
		},
		{
			name:        "edge case - longitude -180",
			lat:         45.4215,
			lng:         -180.0,
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateCoordinates(tt.lat, tt.lng)

			if tt.expectError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errorMsg)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
