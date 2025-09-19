package mocks

import (
	"context"

	"scootin-aboot/internal/models"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

// MockScooterRepository is a mock implementation of repository.ScooterRepository
type MockScooterRepository struct {
	mock.Mock
}

func (m *MockScooterRepository) Create(ctx context.Context, scooter *models.Scooter) error {
	args := m.Called(ctx, scooter)
	return args.Error(0)
}

func (m *MockScooterRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Scooter, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Scooter), args.Error(1)
}

func (m *MockScooterRepository) Update(ctx context.Context, scooter *models.Scooter) error {
	args := m.Called(ctx, scooter)
	return args.Error(0)
}

func (m *MockScooterRepository) Delete(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockScooterRepository) List(ctx context.Context, limit, offset int) ([]*models.Scooter, error) {
	args := m.Called(ctx, limit, offset)
	return args.Get(0).([]*models.Scooter), args.Error(1)
}

func (m *MockScooterRepository) UpdateStatus(ctx context.Context, id uuid.UUID, status models.ScooterStatus) error {
	args := m.Called(ctx, id, status)
	return args.Error(0)
}

func (m *MockScooterRepository) UpdateLocation(ctx context.Context, id uuid.UUID, latitude, longitude float64) error {
	args := m.Called(ctx, id, latitude, longitude)
	return args.Error(0)
}

func (m *MockScooterRepository) GetByStatus(ctx context.Context, status models.ScooterStatus) ([]*models.Scooter, error) {
	args := m.Called(ctx, status)
	return args.Get(0).([]*models.Scooter), args.Error(1)
}

func (m *MockScooterRepository) GetInBounds(ctx context.Context, minLat, maxLat, minLng, maxLng float64) ([]*models.Scooter, error) {
	args := m.Called(ctx, minLat, maxLat, minLng, maxLng)
	return args.Get(0).([]*models.Scooter), args.Error(1)
}

func (m *MockScooterRepository) GetClosest(ctx context.Context, latitude, longitude float64, limit int) ([]*models.Scooter, error) {
	args := m.Called(ctx, latitude, longitude, limit)
	return args.Get(0).([]*models.Scooter), args.Error(1)
}

func (m *MockScooterRepository) GetClosestWithRadius(ctx context.Context, latitude, longitude, radius float64, status string, limit int) ([]*models.Scooter, error) {
	args := m.Called(ctx, latitude, longitude, radius, status, limit)
	return args.Get(0).([]*models.Scooter), args.Error(1)
}

func (m *MockScooterRepository) GetByStatusInBounds(ctx context.Context, status models.ScooterStatus, minLat, maxLat, minLng, maxLng float64) ([]*models.Scooter, error) {
	args := m.Called(ctx, status, minLat, maxLat, minLng, maxLng)
	return args.Get(0).([]*models.Scooter), args.Error(1)
}

func (m *MockScooterRepository) GetByIDForUpdate(ctx context.Context, id uuid.UUID) (*models.Scooter, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Scooter), args.Error(1)
}

func (m *MockScooterRepository) UpdateStatusWithCheck(ctx context.Context, id uuid.UUID, newStatus models.ScooterStatus, expectedStatus models.ScooterStatus) error {
	args := m.Called(ctx, id, newStatus, expectedStatus)
	return args.Error(0)
}
