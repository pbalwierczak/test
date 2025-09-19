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

func TestNewScooterService(t *testing.T) {
	mockScooterRepo := &mocks.MockScooterRepository{}
	mockTripRepo := &mocks.MockTripRepository{}
	mockLocationRepo := &mocks.MockLocationUpdateRepository{}

	service := NewScooterService(mockScooterRepo, mockTripRepo, mockLocationRepo)

	assert.NotNil(t, service)
	assert.IsType(t, &scooterService{}, service)
}

func TestScooterService_GetScooters_NoFilters(t *testing.T) {
	mockScooterRepo := &mocks.MockScooterRepository{}
	mockTripRepo := &mocks.MockTripRepository{}
	mockLocationRepo := &mocks.MockLocationUpdateRepository{}
	service := NewScooterService(mockScooterRepo, mockTripRepo, mockLocationRepo)

	ctx := context.Background()
	params := ScooterQueryParams{
		Limit:  10,
		Offset: 0,
	}

	// Mock data
	testScooters := []*models.Scooter{
		{
			ID:               uuid.New(),
			Status:           models.ScooterStatusAvailable,
			CurrentLatitude:  45.4215,
			CurrentLongitude: -75.6972,
			LastSeen:         time.Now(),
			CreatedAt:        time.Now(),
		},
	}

	mockScooterRepo.On("List", ctx, 10, 0).Return(testScooters, nil)

	result, err := service.GetScooters(ctx, params)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result.Scooters, 1)
	assert.Equal(t, int64(1), result.Total)
	assert.Equal(t, 10, result.Limit)
	assert.Equal(t, 0, result.Offset)
	assert.Equal(t, testScooters[0].ID, result.Scooters[0].ID)
	assert.Equal(t, string(testScooters[0].Status), result.Scooters[0].Status)

	mockScooterRepo.AssertExpectations(t)
}

func TestScooterService_GetScooters_StatusFilter(t *testing.T) {
	mockScooterRepo := &mocks.MockScooterRepository{}
	mockTripRepo := &mocks.MockTripRepository{}
	mockLocationRepo := &mocks.MockLocationUpdateRepository{}
	service := NewScooterService(mockScooterRepo, mockTripRepo, mockLocationRepo)

	ctx := context.Background()
	params := ScooterQueryParams{
		Status: "available",
		Limit:  10,
		Offset: 0,
	}

	testScooters := []*models.Scooter{
		{
			ID:               uuid.New(),
			Status:           models.ScooterStatusAvailable,
			CurrentLatitude:  45.4215,
			CurrentLongitude: -75.6972,
			LastSeen:         time.Now(),
			CreatedAt:        time.Now(),
		},
	}

	mockScooterRepo.On("GetByStatus", ctx, models.ScooterStatusAvailable).Return(testScooters, nil)

	result, err := service.GetScooters(ctx, params)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result.Scooters, 1)
	mockScooterRepo.AssertExpectations(t)
}

func TestScooterService_GetScooters_GeographicBounds(t *testing.T) {
	mockScooterRepo := &mocks.MockScooterRepository{}
	mockTripRepo := &mocks.MockTripRepository{}
	mockLocationRepo := &mocks.MockLocationUpdateRepository{}
	service := NewScooterService(mockScooterRepo, mockTripRepo, mockLocationRepo)

	ctx := context.Background()
	params := ScooterQueryParams{
		MinLat: 45.0,
		MaxLat: 46.0,
		MinLng: -76.0,
		MaxLng: -75.0,
		Limit:  10,
		Offset: 0,
	}

	testScooters := []*models.Scooter{
		{
			ID:               uuid.New(),
			Status:           models.ScooterStatusAvailable,
			CurrentLatitude:  45.4215,
			CurrentLongitude: -75.6972,
			LastSeen:         time.Now(),
			CreatedAt:        time.Now(),
		},
	}

	mockScooterRepo.On("GetInBounds", ctx, 45.0, 46.0, -76.0, -75.0).Return(testScooters, nil)

	result, err := service.GetScooters(ctx, params)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result.Scooters, 1)
	mockScooterRepo.AssertExpectations(t)
}

func TestScooterService_GetScooters_StatusAndBounds(t *testing.T) {
	mockScooterRepo := &mocks.MockScooterRepository{}
	mockTripRepo := &mocks.MockTripRepository{}
	mockLocationRepo := &mocks.MockLocationUpdateRepository{}
	service := NewScooterService(mockScooterRepo, mockTripRepo, mockLocationRepo)

	ctx := context.Background()
	params := ScooterQueryParams{
		Status: "available",
		MinLat: 45.0,
		MaxLat: 46.0,
		MinLng: -76.0,
		MaxLng: -75.0,
		Limit:  10,
		Offset: 0,
	}

	testScooters := []*models.Scooter{
		{
			ID:               uuid.New(),
			Status:           models.ScooterStatusAvailable,
			CurrentLatitude:  45.4215,
			CurrentLongitude: -75.6972,
			LastSeen:         time.Now(),
			CreatedAt:        time.Now(),
		},
	}

	mockScooterRepo.On("GetByStatusInBounds", ctx, models.ScooterStatusAvailable, 45.0, 46.0, -76.0, -75.0).Return(testScooters, nil)

	result, err := service.GetScooters(ctx, params)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result.Scooters, 1)
	mockScooterRepo.AssertExpectations(t)
}

func TestScooterService_GetScooters_Pagination(t *testing.T) {
	mockScooterRepo := &mocks.MockScooterRepository{}
	mockTripRepo := &mocks.MockTripRepository{}
	mockLocationRepo := &mocks.MockLocationUpdateRepository{}
	service := NewScooterService(mockScooterRepo, mockTripRepo, mockLocationRepo)

	ctx := context.Background()
	params := ScooterQueryParams{
		Status: "available",
		Limit:  5,
		Offset: 10,
	}

	// Create 15 scooters to test pagination
	testScooters := make([]*models.Scooter, 15)
	for i := 0; i < 15; i++ {
		testScooters[i] = &models.Scooter{
			ID:               uuid.New(),
			Status:           models.ScooterStatusAvailable,
			CurrentLatitude:  45.4215,
			CurrentLongitude: -75.6972,
			LastSeen:         time.Now(),
			CreatedAt:        time.Now(),
		}
	}

	mockScooterRepo.On("GetByStatus", ctx, models.ScooterStatusAvailable).Return(testScooters, nil)

	result, err := service.GetScooters(ctx, params)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result.Scooters, 5) // Should be limited to 5
	assert.Equal(t, int64(15), result.Total)
	assert.Equal(t, 5, result.Limit)
	assert.Equal(t, 10, result.Offset)
	mockScooterRepo.AssertExpectations(t)
}

func TestScooterService_GetScooters_InvalidParams(t *testing.T) {
	mockScooterRepo := &mocks.MockScooterRepository{}
	mockTripRepo := &mocks.MockTripRepository{}
	mockLocationRepo := &mocks.MockLocationUpdateRepository{}
	service := NewScooterService(mockScooterRepo, mockTripRepo, mockLocationRepo)

	ctx := context.Background()

	tests := []struct {
		name   string
		params ScooterQueryParams
	}{
		{
			name: "invalid status",
			params: ScooterQueryParams{
				Status: "invalid",
				Limit:  10,
				Offset: 0,
			},
		},
		{
			name: "negative limit",
			params: ScooterQueryParams{
				Limit:  -1,
				Offset: 0,
			},
		},
		{
			name: "negative offset",
			params: ScooterQueryParams{
				Limit:  10,
				Offset: -1,
			},
		},
		{
			name: "excessive limit",
			params: ScooterQueryParams{
				Limit:  101,
				Offset: 0,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := service.GetScooters(ctx, tt.params)
			assert.Error(t, err)
			assert.Nil(t, result)
		})
	}
}

func TestScooterService_GetScooters_RepositoryError(t *testing.T) {
	mockScooterRepo := &mocks.MockScooterRepository{}
	mockTripRepo := &mocks.MockTripRepository{}
	mockLocationRepo := &mocks.MockLocationUpdateRepository{}
	service := NewScooterService(mockScooterRepo, mockTripRepo, mockLocationRepo)

	ctx := context.Background()
	params := ScooterQueryParams{
		Limit:  10,
		Offset: 0,
	}

	expectedError := errors.New("database connection failed")
	mockScooterRepo.On("List", ctx, 10, 0).Return([]*models.Scooter(nil), expectedError)

	result, err := service.GetScooters(ctx, params)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "failed to query scooters")
	mockScooterRepo.AssertExpectations(t)
}

func TestScooterService_GetScooter_Success(t *testing.T) {
	mockScooterRepo := &mocks.MockScooterRepository{}
	mockTripRepo := &mocks.MockTripRepository{}
	mockLocationRepo := &mocks.MockLocationUpdateRepository{}
	service := NewScooterService(mockScooterRepo, mockTripRepo, mockLocationRepo)

	ctx := context.Background()
	scooterID := uuid.New()

	testScooter := &models.Scooter{
		ID:               scooterID,
		Status:           models.ScooterStatusAvailable,
		CurrentLatitude:  45.4215,
		CurrentLongitude: -75.6972,
		LastSeen:         time.Now(),
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}

	mockScooterRepo.On("GetByID", ctx, scooterID).Return(testScooter, nil)

	result, err := service.GetScooter(ctx, scooterID)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, scooterID, result.ID)
	assert.Equal(t, string(testScooter.Status), result.Status)
	assert.Equal(t, testScooter.CurrentLatitude, result.CurrentLatitude)
	assert.Equal(t, testScooter.CurrentLongitude, result.CurrentLongitude)
	assert.Nil(t, result.ActiveTrip) // Should be nil for available scooter
	mockScooterRepo.AssertExpectations(t)
}

func TestScooterService_GetScooter_NotFound(t *testing.T) {
	mockScooterRepo := &mocks.MockScooterRepository{}
	mockTripRepo := &mocks.MockTripRepository{}
	mockLocationRepo := &mocks.MockLocationUpdateRepository{}
	service := NewScooterService(mockScooterRepo, mockTripRepo, mockLocationRepo)

	ctx := context.Background()
	scooterID := uuid.New()

	mockScooterRepo.On("GetByID", ctx, scooterID).Return(nil, nil)

	result, err := service.GetScooter(ctx, scooterID)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "scooter not found", err.Error())
	mockScooterRepo.AssertExpectations(t)
}

func TestScooterService_GetScooter_WithActiveTrip(t *testing.T) {
	mockScooterRepo := &mocks.MockScooterRepository{}
	mockTripRepo := &mocks.MockTripRepository{}
	mockLocationRepo := &mocks.MockLocationUpdateRepository{}
	service := NewScooterService(mockScooterRepo, mockTripRepo, mockLocationRepo)

	ctx := context.Background()
	scooterID := uuid.New()
	userID := uuid.New()
	tripID := uuid.New()

	testScooter := &models.Scooter{
		ID:               scooterID,
		Status:           models.ScooterStatusOccupied,
		CurrentLatitude:  45.4215,
		CurrentLongitude: -75.6972,
		LastSeen:         time.Now(),
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}

	testTrip := &models.Trip{
		ID:             tripID,
		UserID:         userID,
		StartTime:      time.Now(),
		StartLatitude:  45.4215,
		StartLongitude: -75.6972,
	}

	mockScooterRepo.On("GetByID", ctx, scooterID).Return(testScooter, nil)
	mockTripRepo.On("GetActiveByScooterID", ctx, scooterID).Return(testTrip, nil)

	result, err := service.GetScooter(ctx, scooterID)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, scooterID, result.ID)
	assert.NotNil(t, result.ActiveTrip)
	assert.Equal(t, tripID, result.ActiveTrip.TripID)
	assert.Equal(t, userID, result.ActiveTrip.UserID)
	mockScooterRepo.AssertExpectations(t)
	mockTripRepo.AssertExpectations(t)
}

func TestScooterService_GetScooter_WithoutActiveTrip(t *testing.T) {
	mockScooterRepo := &mocks.MockScooterRepository{}
	mockTripRepo := &mocks.MockTripRepository{}
	mockLocationRepo := &mocks.MockLocationUpdateRepository{}
	service := NewScooterService(mockScooterRepo, mockTripRepo, mockLocationRepo)

	ctx := context.Background()
	scooterID := uuid.New()

	testScooter := &models.Scooter{
		ID:               scooterID,
		Status:           models.ScooterStatusOccupied,
		CurrentLatitude:  45.4215,
		CurrentLongitude: -75.6972,
		LastSeen:         time.Now(),
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}

	mockScooterRepo.On("GetByID", ctx, scooterID).Return(testScooter, nil)
	mockTripRepo.On("GetActiveByScooterID", ctx, scooterID).Return(nil, nil)

	result, err := service.GetScooter(ctx, scooterID)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, scooterID, result.ID)
	assert.Nil(t, result.ActiveTrip) // Should be nil when no active trip found
	mockScooterRepo.AssertExpectations(t)
	mockTripRepo.AssertExpectations(t)
}

func TestScooterService_GetScooter_RepositoryError(t *testing.T) {
	mockScooterRepo := &mocks.MockScooterRepository{}
	mockTripRepo := &mocks.MockTripRepository{}
	mockLocationRepo := &mocks.MockLocationUpdateRepository{}
	service := NewScooterService(mockScooterRepo, mockTripRepo, mockLocationRepo)

	ctx := context.Background()
	scooterID := uuid.New()

	expectedError := errors.New("database connection failed")
	mockScooterRepo.On("GetByID", ctx, scooterID).Return(nil, expectedError)

	result, err := service.GetScooter(ctx, scooterID)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "failed to get scooter")
	mockScooterRepo.AssertExpectations(t)
}

func TestScooterService_GetScooter_TripQueryError(t *testing.T) {
	mockScooterRepo := &mocks.MockScooterRepository{}
	mockTripRepo := &mocks.MockTripRepository{}
	mockLocationRepo := &mocks.MockLocationUpdateRepository{}
	service := NewScooterService(mockScooterRepo, mockTripRepo, mockLocationRepo)

	ctx := context.Background()
	scooterID := uuid.New()

	testScooter := &models.Scooter{
		ID:               scooterID,
		Status:           models.ScooterStatusOccupied,
		CurrentLatitude:  45.4215,
		CurrentLongitude: -75.6972,
		LastSeen:         time.Now(),
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}

	tripError := errors.New("trip query failed")
	mockScooterRepo.On("GetByID", ctx, scooterID).Return(testScooter, nil)
	mockTripRepo.On("GetActiveByScooterID", ctx, scooterID).Return(nil, tripError)

	result, err := service.GetScooter(ctx, scooterID)

	// Should not fail the request even if trip query fails
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, scooterID, result.ID)
	assert.Nil(t, result.ActiveTrip) // Should be nil when trip query fails
	mockScooterRepo.AssertExpectations(t)
	mockTripRepo.AssertExpectations(t)
}

func TestScooterService_GetClosestScooters_Success(t *testing.T) {
	mockScooterRepo := &mocks.MockScooterRepository{}
	mockTripRepo := &mocks.MockTripRepository{}
	mockLocationRepo := &mocks.MockLocationUpdateRepository{}
	service := NewScooterService(mockScooterRepo, mockTripRepo, mockLocationRepo)

	ctx := context.Background()
	params := ClosestScootersQueryParams{
		Latitude:  45.4215,
		Longitude: -75.6972,
		Radius:    1000, // 1km
		Limit:     10,
		Status:    "available",
	}

	testScooters := []*models.Scooter{
		{
			ID:               uuid.New(),
			Status:           models.ScooterStatusAvailable,
			CurrentLatitude:  45.4216,
			CurrentLongitude: -75.6973,
			LastSeen:         time.Now(),
			CreatedAt:        time.Now(),
		},
	}

	mockScooterRepo.On("GetClosestWithRadius", ctx, 45.4215, -75.6972, 1000.0, "available", 10).Return(testScooters, nil)

	result, err := service.GetClosestScooters(ctx, params)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result.Scooters, 1)
	assert.Equal(t, 45.4215, result.Center.Latitude)
	assert.Equal(t, -75.6972, result.Center.Longitude)
	assert.Equal(t, 1000.0, result.Radius)
	assert.Greater(t, result.Scooters[0].Distance, 0.0) // Should have calculated distance
	mockScooterRepo.AssertExpectations(t)
}

func TestScooterService_GetClosestScooters_WithStatusFilter(t *testing.T) {
	mockScooterRepo := &mocks.MockScooterRepository{}
	mockTripRepo := &mocks.MockTripRepository{}
	mockLocationRepo := &mocks.MockLocationUpdateRepository{}
	service := NewScooterService(mockScooterRepo, mockTripRepo, mockLocationRepo)

	ctx := context.Background()
	params := ClosestScootersQueryParams{
		Latitude:  45.4215,
		Longitude: -75.6972,
		Radius:    1000,
		Limit:     10,
		Status:    "occupied",
	}

	testScooters := []*models.Scooter{
		{
			ID:               uuid.New(),
			Status:           models.ScooterStatusOccupied,
			CurrentLatitude:  45.4216,
			CurrentLongitude: -75.6973,
			LastSeen:         time.Now(),
			CreatedAt:        time.Now(),
		},
	}

	mockScooterRepo.On("GetClosestWithRadius", ctx, 45.4215, -75.6972, 1000.0, "occupied", 10).Return(testScooters, nil)

	result, err := service.GetClosestScooters(ctx, params)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result.Scooters, 1)
	mockScooterRepo.AssertExpectations(t)
}

func TestScooterService_GetClosestScooters_InvalidCoordinates(t *testing.T) {
	mockScooterRepo := &mocks.MockScooterRepository{}
	mockTripRepo := &mocks.MockTripRepository{}
	mockLocationRepo := &mocks.MockLocationUpdateRepository{}
	service := NewScooterService(mockScooterRepo, mockTripRepo, mockLocationRepo)

	ctx := context.Background()

	tests := []struct {
		name   string
		params ClosestScootersQueryParams
	}{
		{
			name: "invalid latitude",
			params: ClosestScootersQueryParams{
				Latitude:  91.0,
				Longitude: -75.6972,
				Radius:    1000,
				Limit:     10,
			},
		},
		{
			name: "invalid longitude",
			params: ClosestScootersQueryParams{
				Latitude:  45.4215,
				Longitude: 181.0,
				Radius:    1000,
				Limit:     10,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := service.GetClosestScooters(ctx, tt.params)
			assert.Error(t, err)
			assert.Nil(t, result)
		})
	}
}

func TestScooterService_GetClosestScooters_InvalidRadius(t *testing.T) {
	mockScooterRepo := &mocks.MockScooterRepository{}
	mockTripRepo := &mocks.MockTripRepository{}
	mockLocationRepo := &mocks.MockLocationUpdateRepository{}
	service := NewScooterService(mockScooterRepo, mockTripRepo, mockLocationRepo)

	ctx := context.Background()

	tests := []struct {
		name   string
		params ClosestScootersQueryParams
	}{
		{
			name: "negative radius",
			params: ClosestScootersQueryParams{
				Latitude:  45.4215,
				Longitude: -75.6972,
				Radius:    -100,
				Limit:     10,
			},
		},
		{
			name: "excessive radius",
			params: ClosestScootersQueryParams{
				Latitude:  45.4215,
				Longitude: -75.6972,
				Radius:    60000, // > 50km
				Limit:     10,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := service.GetClosestScooters(ctx, tt.params)
			assert.Error(t, err)
			assert.Nil(t, result)
		})
	}
}

func TestScooterService_GetClosestScooters_InvalidLimit(t *testing.T) {
	mockScooterRepo := &mocks.MockScooterRepository{}
	mockTripRepo := &mocks.MockTripRepository{}
	mockLocationRepo := &mocks.MockLocationUpdateRepository{}
	service := NewScooterService(mockScooterRepo, mockTripRepo, mockLocationRepo)

	ctx := context.Background()

	tests := []struct {
		name   string
		params ClosestScootersQueryParams
	}{
		{
			name: "negative limit",
			params: ClosestScootersQueryParams{
				Latitude:  45.4215,
				Longitude: -75.6972,
				Radius:    1000,
				Limit:     -1,
			},
		},
		{
			name: "excessive limit",
			params: ClosestScootersQueryParams{
				Latitude:  45.4215,
				Longitude: -75.6972,
				Radius:    1000,
				Limit:     51, // > 50
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := service.GetClosestScooters(ctx, tt.params)
			assert.Error(t, err)
			assert.Nil(t, result)
		})
	}
}

func TestScooterService_GetClosestScooters_RepositoryError(t *testing.T) {
	mockScooterRepo := &mocks.MockScooterRepository{}
	mockTripRepo := &mocks.MockTripRepository{}
	mockLocationRepo := &mocks.MockLocationUpdateRepository{}
	service := NewScooterService(mockScooterRepo, mockTripRepo, mockLocationRepo)

	ctx := context.Background()
	params := ClosestScootersQueryParams{
		Latitude:  45.4215,
		Longitude: -75.6972,
		Radius:    1000,
		Limit:     10,
	}

	expectedError := errors.New("database connection failed")
	mockScooterRepo.On("GetClosestWithRadius", ctx, 45.4215, -75.6972, 1000.0, "", 10).Return([]*models.Scooter(nil), expectedError)

	result, err := service.GetClosestScooters(ctx, params)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "failed to query closest scooters")
	mockScooterRepo.AssertExpectations(t)
}

func TestValidateScooterQueryParams_ValidParams(t *testing.T) {
	service := &scooterService{}

	params := ScooterQueryParams{
		Status: "available",
		MinLat: 45.0,
		MaxLat: 46.0,
		MinLng: -76.0,
		MaxLng: -75.0,
		Limit:  10,
		Offset: 0,
	}

	err := service.validateScooterQueryParams(params)
	assert.NoError(t, err)
}

func TestValidateScooterQueryParams_InvalidStatus(t *testing.T) {
	service := &scooterService{}

	params := ScooterQueryParams{
		Status: "invalid",
		Limit:  10,
		Offset: 0,
	}

	err := service.validateScooterQueryParams(params)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "status must be 'available' or 'occupied'")
}

func TestValidateScooterQueryParams_NegativeLimit(t *testing.T) {
	service := &scooterService{}

	params := ScooterQueryParams{
		Limit:  -1,
		Offset: 0,
	}

	err := service.validateScooterQueryParams(params)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "limit must be non-negative")
}

func TestValidateScooterQueryParams_NegativeOffset(t *testing.T) {
	service := &scooterService{}

	params := ScooterQueryParams{
		Limit:  10,
		Offset: -1,
	}

	err := service.validateScooterQueryParams(params)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "offset must be non-negative")
}

func TestValidateScooterQueryParams_ExcessiveLimit(t *testing.T) {
	service := &scooterService{}

	params := ScooterQueryParams{
		Limit:  101,
		Offset: 0,
	}

	err := service.validateScooterQueryParams(params)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "limit cannot exceed 100")
}

func TestValidateClosestScootersParams_ValidParams(t *testing.T) {
	service := &scooterService{}

	params := ClosestScootersQueryParams{
		Latitude:  45.4215,
		Longitude: -75.6972,
		Radius:    1000,
		Limit:     10,
		Status:    "available",
	}

	err := service.validateClosestScootersParams(params)
	assert.NoError(t, err)
}

func TestValidateClosestScootersParams_InvalidCoordinates(t *testing.T) {
	service := &scooterService{}

	params := ClosestScootersQueryParams{
		Latitude:  91.0,
		Longitude: -75.6972,
		Radius:    1000,
		Limit:     10,
	}

	err := service.validateClosestScootersParams(params)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid latitude")
}

func TestValidateClosestScootersParams_NegativeRadius(t *testing.T) {
	service := &scooterService{}

	params := ClosestScootersQueryParams{
		Latitude:  45.4215,
		Longitude: -75.6972,
		Radius:    -100,
		Limit:     10,
	}

	err := service.validateClosestScootersParams(params)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "radius must be non-negative")
}

func TestValidateClosestScootersParams_ExcessiveRadius(t *testing.T) {
	service := &scooterService{}

	params := ClosestScootersQueryParams{
		Latitude:  45.4215,
		Longitude: -75.6972,
		Radius:    60000,
		Limit:     10,
	}

	err := service.validateClosestScootersParams(params)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "radius cannot exceed 50000 meters")
}

func TestValidateClosestScootersParams_InvalidStatus(t *testing.T) {
	service := &scooterService{}

	params := ClosestScootersQueryParams{
		Latitude:  45.4215,
		Longitude: -75.6972,
		Radius:    1000,
		Limit:     10,
		Status:    "invalid",
	}

	err := service.validateClosestScootersParams(params)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "status must be 'available' or 'occupied'")
}

func TestValidateClosestScootersParams_ExcessiveLimit(t *testing.T) {
	service := &scooterService{}

	params := ClosestScootersQueryParams{
		Latitude:  45.4215,
		Longitude: -75.6972,
		Radius:    1000,
		Limit:     51,
	}

	err := service.validateClosestScootersParams(params)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "limit cannot exceed 50")
}

func TestMapScooterToInfo_CompleteScooter(t *testing.T) {
	service := &scooterService{}

	now := time.Now()
	scooter := &models.Scooter{
		ID:               uuid.New(),
		Status:           models.ScooterStatusAvailable,
		CurrentLatitude:  45.4215,
		CurrentLongitude: -75.6972,
		LastSeen:         now,
		CreatedAt:        now,
	}

	result := service.mapScooterToInfo(scooter)

	assert.Equal(t, scooter.ID, result.ID)
	assert.Equal(t, string(scooter.Status), result.Status)
	assert.Equal(t, scooter.CurrentLatitude, result.CurrentLatitude)
	assert.Equal(t, scooter.CurrentLongitude, result.CurrentLongitude)
	assert.Equal(t, scooter.LastSeen, result.LastSeen)
	assert.Equal(t, scooter.CreatedAt, result.CreatedAt)
}

func TestMapScooterToInfo_MinimalScooter(t *testing.T) {
	service := &scooterService{}

	scooter := &models.Scooter{
		ID:               uuid.New(),
		Status:           models.ScooterStatusOccupied,
		CurrentLatitude:  0.0,
		CurrentLongitude: 0.0,
	}

	result := service.mapScooterToInfo(scooter)

	assert.Equal(t, scooter.ID, result.ID)
	assert.Equal(t, string(scooter.Status), result.Status)
	assert.Equal(t, scooter.CurrentLatitude, result.CurrentLatitude)
	assert.Equal(t, scooter.CurrentLongitude, result.CurrentLongitude)
}

func TestScooterService_UpdateLocation_Success(t *testing.T) {
	mockScooterRepo := &mocks.MockScooterRepository{}
	mockTripRepo := &mocks.MockTripRepository{}
	mockLocationRepo := &mocks.MockLocationUpdateRepository{}
	service := NewScooterService(mockScooterRepo, mockTripRepo, mockLocationRepo)

	ctx := context.Background()
	scooterID := uuid.New()
	latitude := 45.4215
	longitude := -75.6972

	// Mock scooter exists
	scooter := &models.Scooter{
		ID:               scooterID,
		Status:           models.ScooterStatusAvailable,
		CurrentLatitude:  45.0,
		CurrentLongitude: -75.0,
		LastSeen:         time.Now(),
		CreatedAt:        time.Now(),
	}

	mockScooterRepo.On("GetByID", ctx, scooterID).Return(scooter, nil)
	mockLocationRepo.On("Create", ctx, mock.AnythingOfType("*models.LocationUpdate")).Return(nil)
	mockScooterRepo.On("UpdateLocation", ctx, scooterID, latitude, longitude).Return(nil)

	err := service.UpdateLocation(ctx, scooterID, latitude, longitude)

	assert.NoError(t, err)
	mockScooterRepo.AssertExpectations(t)
	mockLocationRepo.AssertExpectations(t)
}

func TestScooterService_UpdateLocation_InvalidCoordinates(t *testing.T) {
	mockScooterRepo := &mocks.MockScooterRepository{}
	mockTripRepo := &mocks.MockTripRepository{}
	mockLocationRepo := &mocks.MockLocationUpdateRepository{}
	service := NewScooterService(mockScooterRepo, mockTripRepo, mockLocationRepo)

	ctx := context.Background()
	scooterID := uuid.New()
	latitude := 91.0 // Invalid latitude
	longitude := -75.6972

	err := service.UpdateLocation(ctx, scooterID, latitude, longitude)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid coordinates")
}

func TestScooterService_UpdateLocation_ScooterNotFound(t *testing.T) {
	mockScooterRepo := &mocks.MockScooterRepository{}
	mockTripRepo := &mocks.MockTripRepository{}
	mockLocationRepo := &mocks.MockLocationUpdateRepository{}
	service := NewScooterService(mockScooterRepo, mockTripRepo, mockLocationRepo)

	ctx := context.Background()
	scooterID := uuid.New()
	latitude := 45.4215
	longitude := -75.6972

	mockScooterRepo.On("GetByID", ctx, scooterID).Return(nil, nil)

	err := service.UpdateLocation(ctx, scooterID, latitude, longitude)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "scooter not found")
	mockScooterRepo.AssertExpectations(t)
}

func TestScooterService_UpdateLocation_GetScooterError(t *testing.T) {
	mockScooterRepo := &mocks.MockScooterRepository{}
	mockTripRepo := &mocks.MockTripRepository{}
	mockLocationRepo := &mocks.MockLocationUpdateRepository{}
	service := NewScooterService(mockScooterRepo, mockTripRepo, mockLocationRepo)

	ctx := context.Background()
	scooterID := uuid.New()
	latitude := 45.4215
	longitude := -75.6972

	expectedError := errors.New("database error")
	mockScooterRepo.On("GetByID", ctx, scooterID).Return(nil, expectedError)

	err := service.UpdateLocation(ctx, scooterID, latitude, longitude)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to get scooter")
	mockScooterRepo.AssertExpectations(t)
}

func TestScooterService_UpdateLocation_CreateLocationUpdateError(t *testing.T) {
	mockScooterRepo := &mocks.MockScooterRepository{}
	mockTripRepo := &mocks.MockTripRepository{}
	mockLocationRepo := &mocks.MockLocationUpdateRepository{}
	service := NewScooterService(mockScooterRepo, mockTripRepo, mockLocationRepo)

	ctx := context.Background()
	scooterID := uuid.New()
	latitude := 45.4215
	longitude := -75.6972

	scooter := &models.Scooter{
		ID:               scooterID,
		Status:           models.ScooterStatusAvailable,
		CurrentLatitude:  45.0,
		CurrentLongitude: -75.0,
		LastSeen:         time.Now(),
		CreatedAt:        time.Now(),
	}

	expectedError := errors.New("database error")
	mockScooterRepo.On("GetByID", ctx, scooterID).Return(scooter, nil)
	mockLocationRepo.On("Create", ctx, mock.AnythingOfType("*models.LocationUpdate")).Return(expectedError)

	err := service.UpdateLocation(ctx, scooterID, latitude, longitude)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to create location update")
	mockScooterRepo.AssertExpectations(t)
	mockLocationRepo.AssertExpectations(t)
}

func TestScooterService_UpdateLocation_UpdateScooterLocationError(t *testing.T) {
	mockScooterRepo := &mocks.MockScooterRepository{}
	mockTripRepo := &mocks.MockTripRepository{}
	mockLocationRepo := &mocks.MockLocationUpdateRepository{}
	service := NewScooterService(mockScooterRepo, mockTripRepo, mockLocationRepo)

	ctx := context.Background()
	scooterID := uuid.New()
	latitude := 45.4215
	longitude := -75.6972

	scooter := &models.Scooter{
		ID:               scooterID,
		Status:           models.ScooterStatusAvailable,
		CurrentLatitude:  45.0,
		CurrentLongitude: -75.0,
		LastSeen:         time.Now(),
		CreatedAt:        time.Now(),
	}

	expectedError := errors.New("database error")
	mockScooterRepo.On("GetByID", ctx, scooterID).Return(scooter, nil)
	mockLocationRepo.On("Create", ctx, mock.AnythingOfType("*models.LocationUpdate")).Return(nil)
	mockScooterRepo.On("UpdateLocation", ctx, scooterID, latitude, longitude).Return(expectedError)

	err := service.UpdateLocation(ctx, scooterID, latitude, longitude)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to update scooter location")
	mockScooterRepo.AssertExpectations(t)
	mockLocationRepo.AssertExpectations(t)
}
