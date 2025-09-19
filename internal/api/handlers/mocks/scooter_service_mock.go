package mocks

import (
	"context"

	"scootin-aboot/internal/services"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

// MockScooterService is a mock implementation of ScooterService
type MockScooterService struct {
	mock.Mock
}

// GetScooters mocks the GetScooters method
func (m *MockScooterService) GetScooters(ctx context.Context, params services.ScooterQueryParams) (*services.ScooterListResult, error) {
	args := m.Called(ctx, params)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*services.ScooterListResult), args.Error(1)
}

// GetScooter mocks the GetScooter method
func (m *MockScooterService) GetScooter(ctx context.Context, id uuid.UUID) (*services.ScooterDetailsResult, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*services.ScooterDetailsResult), args.Error(1)
}

// GetClosestScooters mocks the GetClosestScooters method
func (m *MockScooterService) GetClosestScooters(ctx context.Context, params services.ClosestScootersQueryParams) (*services.ClosestScootersResult, error) {
	args := m.Called(ctx, params)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*services.ClosestScootersResult), args.Error(1)
}

// UpdateLocation mocks the UpdateLocation method
func (m *MockScooterService) UpdateLocation(ctx context.Context, scooterID uuid.UUID, lat, lng float64) error {
	args := m.Called(ctx, scooterID, lat, lng)
	return args.Error(0)
}
