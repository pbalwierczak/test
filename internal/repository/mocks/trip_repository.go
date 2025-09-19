package mocks

import (
	"context"

	"scootin-aboot/internal/models"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

// MockTripRepository is a mock implementation of repository.TripRepository
type MockTripRepository struct {
	mock.Mock
}

func (m *MockTripRepository) Create(ctx context.Context, trip *models.Trip) error {
	args := m.Called(ctx, trip)
	return args.Error(0)
}

func (m *MockTripRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Trip, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Trip), args.Error(1)
}

func (m *MockTripRepository) Update(ctx context.Context, trip *models.Trip) error {
	args := m.Called(ctx, trip)
	return args.Error(0)
}

func (m *MockTripRepository) Delete(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockTripRepository) List(ctx context.Context, limit, offset int) ([]*models.Trip, error) {
	args := m.Called(ctx, limit, offset)
	return args.Get(0).([]*models.Trip), args.Error(1)
}

func (m *MockTripRepository) UpdateStatus(ctx context.Context, id uuid.UUID, status models.TripStatus) error {
	args := m.Called(ctx, id, status)
	return args.Error(0)
}

func (m *MockTripRepository) EndTrip(ctx context.Context, id uuid.UUID, endLat, endLng float64) error {
	args := m.Called(ctx, id, endLat, endLng)
	return args.Error(0)
}

func (m *MockTripRepository) CancelTrip(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockTripRepository) GetActiveByUserID(ctx context.Context, userID uuid.UUID) (*models.Trip, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Trip), args.Error(1)
}

func (m *MockTripRepository) GetActiveByScooterID(ctx context.Context, scooterID uuid.UUID) (*models.Trip, error) {
	args := m.Called(ctx, scooterID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Trip), args.Error(1)
}
