package mocks

import (
	"context"
	"time"

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

func (m *MockTripRepository) GetByStatus(ctx context.Context, status models.TripStatus) ([]*models.Trip, error) {
	args := m.Called(ctx, status)
	return args.Get(0).([]*models.Trip), args.Error(1)
}

func (m *MockTripRepository) GetActive(ctx context.Context) ([]*models.Trip, error) {
	args := m.Called(ctx)
	return args.Get(0).([]*models.Trip), args.Error(1)
}

func (m *MockTripRepository) GetCompleted(ctx context.Context) ([]*models.Trip, error) {
	args := m.Called(ctx)
	return args.Get(0).([]*models.Trip), args.Error(1)
}

func (m *MockTripRepository) GetCancelled(ctx context.Context) ([]*models.Trip, error) {
	args := m.Called(ctx)
	return args.Get(0).([]*models.Trip), args.Error(1)
}

func (m *MockTripRepository) GetByUserID(ctx context.Context, userID uuid.UUID) ([]*models.Trip, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).([]*models.Trip), args.Error(1)
}

func (m *MockTripRepository) GetByScooterID(ctx context.Context, scooterID uuid.UUID) ([]*models.Trip, error) {
	args := m.Called(ctx, scooterID)
	return args.Get(0).([]*models.Trip), args.Error(1)
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

func (m *MockTripRepository) GetByDateRange(ctx context.Context, start, end time.Time) ([]*models.Trip, error) {
	args := m.Called(ctx, start, end)
	return args.Get(0).([]*models.Trip), args.Error(1)
}

func (m *MockTripRepository) GetByUserIDAndDateRange(ctx context.Context, userID uuid.UUID, start, end time.Time) ([]*models.Trip, error) {
	args := m.Called(ctx, userID, start, end)
	return args.Get(0).([]*models.Trip), args.Error(1)
}

func (m *MockTripRepository) GetByScooterIDAndDateRange(ctx context.Context, scooterID uuid.UUID, start, end time.Time) ([]*models.Trip, error) {
	args := m.Called(ctx, scooterID, start, end)
	return args.Get(0).([]*models.Trip), args.Error(1)
}

func (m *MockTripRepository) GetTripCount(ctx context.Context) (int64, error) {
	args := m.Called(ctx)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockTripRepository) GetTripCountByUser(ctx context.Context, userID uuid.UUID) (int64, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockTripRepository) GetTripCountByScooter(ctx context.Context, scooterID uuid.UUID) (int64, error) {
	args := m.Called(ctx, scooterID)
	return args.Get(0).(int64), args.Error(1)
}
