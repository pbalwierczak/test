package mocks

import (
	"context"
	"time"

	"scootin-aboot/internal/models"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

// MockLocationUpdateRepository is a mock implementation of repository.LocationUpdateRepository
type MockLocationUpdateRepository struct {
	mock.Mock
}

func (m *MockLocationUpdateRepository) Create(ctx context.Context, update *models.LocationUpdate) error {
	args := m.Called(ctx, update)
	return args.Error(0)
}

func (m *MockLocationUpdateRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.LocationUpdate, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.LocationUpdate), args.Error(1)
}

func (m *MockLocationUpdateRepository) Update(ctx context.Context, update *models.LocationUpdate) error {
	args := m.Called(ctx, update)
	return args.Error(0)
}

func (m *MockLocationUpdateRepository) Delete(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockLocationUpdateRepository) List(ctx context.Context, limit, offset int) ([]*models.LocationUpdate, error) {
	args := m.Called(ctx, limit, offset)
	return args.Get(0).([]*models.LocationUpdate), args.Error(1)
}

func (m *MockLocationUpdateRepository) GetByTripID(ctx context.Context, tripID uuid.UUID) ([]*models.LocationUpdate, error) {
	args := m.Called(ctx, tripID)
	return args.Get(0).([]*models.LocationUpdate), args.Error(1)
}

func (m *MockLocationUpdateRepository) GetByTripIDOrdered(ctx context.Context, tripID uuid.UUID) ([]*models.LocationUpdate, error) {
	args := m.Called(ctx, tripID)
	return args.Get(0).([]*models.LocationUpdate), args.Error(1)
}

func (m *MockLocationUpdateRepository) GetLatestByTripID(ctx context.Context, tripID uuid.UUID) (*models.LocationUpdate, error) {
	args := m.Called(ctx, tripID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.LocationUpdate), args.Error(1)
}

func (m *MockLocationUpdateRepository) GetByDateRange(ctx context.Context, start, end time.Time) ([]*models.LocationUpdate, error) {
	args := m.Called(ctx, start, end)
	return args.Get(0).([]*models.LocationUpdate), args.Error(1)
}

func (m *MockLocationUpdateRepository) GetByTripIDAndDateRange(ctx context.Context, tripID uuid.UUID, start, end time.Time) ([]*models.LocationUpdate, error) {
	args := m.Called(ctx, tripID, start, end)
	return args.Get(0).([]*models.LocationUpdate), args.Error(1)
}

func (m *MockLocationUpdateRepository) GetInBounds(ctx context.Context, minLat, maxLat, minLng, maxLng float64) ([]*models.LocationUpdate, error) {
	args := m.Called(ctx, minLat, maxLat, minLng, maxLng)
	return args.Get(0).([]*models.LocationUpdate), args.Error(1)
}

func (m *MockLocationUpdateRepository) GetInRadius(ctx context.Context, latitude, longitude, radiusKm float64) ([]*models.LocationUpdate, error) {
	args := m.Called(ctx, latitude, longitude, radiusKm)
	return args.Get(0).([]*models.LocationUpdate), args.Error(1)
}

func (m *MockLocationUpdateRepository) GetUpdateCount(ctx context.Context) (int64, error) {
	args := m.Called(ctx)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockLocationUpdateRepository) GetUpdateCountByTrip(ctx context.Context, tripID uuid.UUID) (int64, error) {
	args := m.Called(ctx, tripID)
	return args.Get(0).(int64), args.Error(1)
}
