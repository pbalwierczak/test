package mocks

import (
	"context"

	"scootin-aboot/internal/models"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

// MockTripService is a mock implementation of TripService
type MockTripService struct {
	mock.Mock
}

// StartTrip mocks the StartTrip method
func (m *MockTripService) StartTrip(ctx context.Context, scooterID, userID uuid.UUID, lat, lng float64) (*models.Trip, error) {
	args := m.Called(ctx, scooterID, userID, lat, lng)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Trip), args.Error(1)
}

// EndTrip mocks the EndTrip method
func (m *MockTripService) EndTrip(ctx context.Context, scooterID uuid.UUID, lat, lng float64) (*models.Trip, error) {
	args := m.Called(ctx, scooterID, lat, lng)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Trip), args.Error(1)
}

// CancelTrip mocks the CancelTrip method
func (m *MockTripService) CancelTrip(ctx context.Context, scooterID uuid.UUID) (*models.Trip, error) {
	args := m.Called(ctx, scooterID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Trip), args.Error(1)
}

// UpdateLocation mocks the UpdateLocation method
func (m *MockTripService) UpdateLocation(ctx context.Context, scooterID uuid.UUID, lat, lng float64) error {
	args := m.Called(ctx, scooterID, lat, lng)
	return args.Error(0)
}

// GetActiveTrip mocks the GetActiveTrip method
func (m *MockTripService) GetActiveTrip(ctx context.Context, scooterID uuid.UUID) (*models.Trip, error) {
	args := m.Called(ctx, scooterID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Trip), args.Error(1)
}

// GetActiveTripByUser mocks the GetActiveTripByUser method
func (m *MockTripService) GetActiveTripByUser(ctx context.Context, userID uuid.UUID) (*models.Trip, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Trip), args.Error(1)
}

// GetTrip mocks the GetTrip method
func (m *MockTripService) GetTrip(ctx context.Context, tripID uuid.UUID) (*models.Trip, error) {
	args := m.Called(ctx, tripID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Trip), args.Error(1)
}
