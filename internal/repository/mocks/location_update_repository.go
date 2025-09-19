package mocks

import (
	"context"

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

func (m *MockLocationUpdateRepository) GetByScooterID(ctx context.Context, scooterID uuid.UUID) ([]*models.LocationUpdate, error) {
	args := m.Called(ctx, scooterID)
	return args.Get(0).([]*models.LocationUpdate), args.Error(1)
}
