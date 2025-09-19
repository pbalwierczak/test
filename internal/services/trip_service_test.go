package services

import (
	"context"
	"errors"
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

func TestTripService_CancelTrip(t *testing.T) {
	tests := []struct {
		name          string
		scooterID     uuid.UUID
		setupMocks    func(*mocks.MockTripRepository, *mocks.MockScooterRepository, *mocks.MockUserRepository, *mocks.MockLocationUpdateRepository)
		expectedError string
		expectedTrip  *models.Trip
	}{
		{
			name:      "successful trip cancellation",
			scooterID: uuid.New(),
			setupMocks: func(tripRepo *mocks.MockTripRepository, scooterRepo *mocks.MockScooterRepository, userRepo *mocks.MockUserRepository, locationRepo *mocks.MockLocationUpdateRepository) {
				// Active trip exists
				trip := &models.Trip{
					ID:        uuid.New(),
					ScooterID: uuid.New(),
					UserID:    uuid.New(),
					Status:    models.TripStatusActive,
				}
				tripRepo.On("GetActiveByScooterID", mock.Anything, mock.Anything).Return(trip, nil)

				// Trip cancellation succeeds
				tripRepo.On("CancelTrip", mock.Anything, mock.Anything).Return(nil)

				// Scooter status update succeeds
				scooterRepo.On("UpdateStatusWithCheck", mock.Anything, mock.Anything, models.ScooterStatusAvailable, models.ScooterStatusOccupied).Return(nil)
			},
			expectedError: "",
		},
		{
			name:      "no active trip found",
			scooterID: uuid.New(),
			setupMocks: func(tripRepo *mocks.MockTripRepository, scooterRepo *mocks.MockScooterRepository, userRepo *mocks.MockUserRepository, locationRepo *mocks.MockLocationUpdateRepository) {
				tripRepo.On("GetActiveByScooterID", mock.Anything, mock.Anything).Return(nil, nil)
			},
			expectedError: "no active trip found for scooter",
		},
		{
			name:      "repository error when getting active trip",
			scooterID: uuid.New(),
			setupMocks: func(tripRepo *mocks.MockTripRepository, scooterRepo *mocks.MockScooterRepository, userRepo *mocks.MockUserRepository, locationRepo *mocks.MockLocationUpdateRepository) {
				tripRepo.On("GetActiveByScooterID", mock.Anything, mock.Anything).Return(nil, errors.New("database error"))
			},
			expectedError: "failed to get active trip: database error",
		},
		{
			name:      "repository error when cancelling trip",
			scooterID: uuid.New(),
			setupMocks: func(tripRepo *mocks.MockTripRepository, scooterRepo *mocks.MockScooterRepository, userRepo *mocks.MockUserRepository, locationRepo *mocks.MockLocationUpdateRepository) {
				trip := &models.Trip{
					ID:        uuid.New(),
					ScooterID: uuid.New(),
					UserID:    uuid.New(),
					Status:    models.TripStatusActive,
				}
				tripRepo.On("GetActiveByScooterID", mock.Anything, mock.Anything).Return(trip, nil)
				tripRepo.On("CancelTrip", mock.Anything, mock.Anything).Return(errors.New("database error"))
			},
			expectedError: "failed to cancel trip: database error",
		},
		{
			name:      "repository error when updating scooter status",
			scooterID: uuid.New(),
			setupMocks: func(tripRepo *mocks.MockTripRepository, scooterRepo *mocks.MockScooterRepository, userRepo *mocks.MockUserRepository, locationRepo *mocks.MockLocationUpdateRepository) {
				trip := &models.Trip{
					ID:        uuid.New(),
					ScooterID: uuid.New(),
					UserID:    uuid.New(),
					Status:    models.TripStatusActive,
				}
				tripRepo.On("GetActiveByScooterID", mock.Anything, mock.Anything).Return(trip, nil)
				tripRepo.On("CancelTrip", mock.Anything, mock.Anything).Return(nil)
				scooterRepo.On("UpdateStatusWithCheck", mock.Anything, mock.Anything, models.ScooterStatusAvailable, models.ScooterStatusOccupied).Return(errors.New("database error"))
			},
			expectedError: "failed to update scooter status: database error",
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
			trip, err := service.CancelTrip(context.Background(), tt.scooterID)

			// Assertions
			if tt.expectedError != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError)
				assert.Nil(t, trip)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, trip)
				assert.Equal(t, models.TripStatusCancelled, trip.Status)
			}

			// Verify all expectations were met
			tripRepo.AssertExpectations(t)
			scooterRepo.AssertExpectations(t)
			userRepo.AssertExpectations(t)
			locationRepo.AssertExpectations(t)
		})
	}
}

func TestTripService_GetActiveTrip(t *testing.T) {
	tests := []struct {
		name          string
		scooterID     uuid.UUID
		setupMocks    func(*mocks.MockTripRepository, *mocks.MockScooterRepository, *mocks.MockUserRepository, *mocks.MockLocationUpdateRepository)
		expectedError string
		expectedTrip  *models.Trip
	}{
		{
			name:      "successful get active trip",
			scooterID: uuid.New(),
			setupMocks: func(tripRepo *mocks.MockTripRepository, scooterRepo *mocks.MockScooterRepository, userRepo *mocks.MockUserRepository, locationRepo *mocks.MockLocationUpdateRepository) {
				trip := &models.Trip{
					ID:        uuid.New(),
					ScooterID: uuid.New(),
					UserID:    uuid.New(),
					Status:    models.TripStatusActive,
				}
				tripRepo.On("GetActiveByScooterID", mock.Anything, mock.Anything).Return(trip, nil)
			},
			expectedError: "",
			expectedTrip: &models.Trip{
				Status: models.TripStatusActive,
			},
		},
		{
			name:      "no active trip found",
			scooterID: uuid.New(),
			setupMocks: func(tripRepo *mocks.MockTripRepository, scooterRepo *mocks.MockScooterRepository, userRepo *mocks.MockUserRepository, locationRepo *mocks.MockLocationUpdateRepository) {
				tripRepo.On("GetActiveByScooterID", mock.Anything, mock.Anything).Return(nil, nil)
			},
			expectedError: "",
			expectedTrip:  nil,
		},
		{
			name:      "repository error",
			scooterID: uuid.New(),
			setupMocks: func(tripRepo *mocks.MockTripRepository, scooterRepo *mocks.MockScooterRepository, userRepo *mocks.MockUserRepository, locationRepo *mocks.MockLocationUpdateRepository) {
				tripRepo.On("GetActiveByScooterID", mock.Anything, mock.Anything).Return(nil, errors.New("database error"))
			},
			expectedError: "failed to get active trip: database error",
			expectedTrip:  nil,
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
			trip, err := service.GetActiveTrip(context.Background(), tt.scooterID)

			// Assertions
			if tt.expectedError != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError)
				assert.Nil(t, trip)
			} else {
				assert.NoError(t, err)
				if tt.expectedTrip == nil {
					assert.Nil(t, trip)
				} else {
					assert.NotNil(t, trip)
					assert.Equal(t, tt.expectedTrip.Status, trip.Status)
				}
			}

			// Verify all expectations were met
			tripRepo.AssertExpectations(t)
			scooterRepo.AssertExpectations(t)
			userRepo.AssertExpectations(t)
			locationRepo.AssertExpectations(t)
		})
	}
}

func TestTripService_GetActiveTripByUser(t *testing.T) {
	tests := []struct {
		name          string
		userID        uuid.UUID
		setupMocks    func(*mocks.MockTripRepository, *mocks.MockScooterRepository, *mocks.MockUserRepository, *mocks.MockLocationUpdateRepository)
		expectedError string
		expectedTrip  *models.Trip
	}{
		{
			name:   "successful get active trip by user",
			userID: uuid.New(),
			setupMocks: func(tripRepo *mocks.MockTripRepository, scooterRepo *mocks.MockScooterRepository, userRepo *mocks.MockUserRepository, locationRepo *mocks.MockLocationUpdateRepository) {
				trip := &models.Trip{
					ID:     uuid.New(),
					UserID: uuid.New(),
					Status: models.TripStatusActive,
				}
				tripRepo.On("GetActiveByUserID", mock.Anything, mock.Anything).Return(trip, nil)
			},
			expectedError: "",
			expectedTrip: &models.Trip{
				Status: models.TripStatusActive,
			},
		},
		{
			name:   "no active trip found for user",
			userID: uuid.New(),
			setupMocks: func(tripRepo *mocks.MockTripRepository, scooterRepo *mocks.MockScooterRepository, userRepo *mocks.MockUserRepository, locationRepo *mocks.MockLocationUpdateRepository) {
				tripRepo.On("GetActiveByUserID", mock.Anything, mock.Anything).Return(nil, nil)
			},
			expectedError: "",
			expectedTrip:  nil,
		},
		{
			name:   "repository error",
			userID: uuid.New(),
			setupMocks: func(tripRepo *mocks.MockTripRepository, scooterRepo *mocks.MockScooterRepository, userRepo *mocks.MockUserRepository, locationRepo *mocks.MockLocationUpdateRepository) {
				tripRepo.On("GetActiveByUserID", mock.Anything, mock.Anything).Return(nil, errors.New("database error"))
			},
			expectedError: "failed to get active trip: database error",
			expectedTrip:  nil,
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
			trip, err := service.GetActiveTripByUser(context.Background(), tt.userID)

			// Assertions
			if tt.expectedError != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError)
				assert.Nil(t, trip)
			} else {
				assert.NoError(t, err)
				if tt.expectedTrip == nil {
					assert.Nil(t, trip)
				} else {
					assert.NotNil(t, trip)
					assert.Equal(t, tt.expectedTrip.Status, trip.Status)
				}
			}

			// Verify all expectations were met
			tripRepo.AssertExpectations(t)
			scooterRepo.AssertExpectations(t)
			userRepo.AssertExpectations(t)
			locationRepo.AssertExpectations(t)
		})
	}
}

func TestTripService_GetTrip(t *testing.T) {
	tests := []struct {
		name          string
		tripID        uuid.UUID
		setupMocks    func(*mocks.MockTripRepository, *mocks.MockScooterRepository, *mocks.MockUserRepository, *mocks.MockLocationUpdateRepository)
		expectedError string
		expectedTrip  *models.Trip
	}{
		{
			name:   "successful get trip",
			tripID: uuid.New(),
			setupMocks: func(tripRepo *mocks.MockTripRepository, scooterRepo *mocks.MockScooterRepository, userRepo *mocks.MockUserRepository, locationRepo *mocks.MockLocationUpdateRepository) {
				trip := &models.Trip{
					ID:     uuid.New(),
					UserID: uuid.New(),
					Status: models.TripStatusCompleted,
				}
				tripRepo.On("GetByID", mock.Anything, mock.Anything).Return(trip, nil)
			},
			expectedError: "",
			expectedTrip: &models.Trip{
				Status: models.TripStatusCompleted,
			},
		},
		{
			name:   "trip not found",
			tripID: uuid.New(),
			setupMocks: func(tripRepo *mocks.MockTripRepository, scooterRepo *mocks.MockScooterRepository, userRepo *mocks.MockUserRepository, locationRepo *mocks.MockLocationUpdateRepository) {
				tripRepo.On("GetByID", mock.Anything, mock.Anything).Return(nil, nil)
			},
			expectedError: "trip not found",
			expectedTrip:  nil,
		},
		{
			name:   "repository error",
			tripID: uuid.New(),
			setupMocks: func(tripRepo *mocks.MockTripRepository, scooterRepo *mocks.MockScooterRepository, userRepo *mocks.MockUserRepository, locationRepo *mocks.MockLocationUpdateRepository) {
				tripRepo.On("GetByID", mock.Anything, mock.Anything).Return(nil, errors.New("database error"))
			},
			expectedError: "failed to get trip: database error",
			expectedTrip:  nil,
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
			trip, err := service.GetTrip(context.Background(), tt.tripID)

			// Assertions
			if tt.expectedError != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError)
				assert.Nil(t, trip)
			} else {
				assert.NoError(t, err)
				if tt.expectedTrip == nil {
					assert.Nil(t, trip)
				} else {
					assert.NotNil(t, trip)
					assert.Equal(t, tt.expectedTrip.Status, trip.Status)
				}
			}

			// Verify all expectations were met
			tripRepo.AssertExpectations(t)
			scooterRepo.AssertExpectations(t)
			userRepo.AssertExpectations(t)
			locationRepo.AssertExpectations(t)
		})
	}
}
