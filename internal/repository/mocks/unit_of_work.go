package mocks

import (
	"context"

	"scootin-aboot/internal/repository"

	"github.com/stretchr/testify/mock"
)

// MockUnitOfWork is a mock implementation of UnitOfWork
type MockUnitOfWork struct {
	mock.Mock
}

// Begin mocks the Begin method
func (m *MockUnitOfWork) Begin(ctx context.Context) (repository.UnitOfWorkTx, error) {
	args := m.Called(ctx)
	return args.Get(0).(repository.UnitOfWorkTx), args.Error(1)
}

// MockUnitOfWorkTx is a mock implementation of UnitOfWorkTx
type MockUnitOfWorkTx struct {
	mock.Mock
}

// ScooterRepository mocks the ScooterRepository method
func (m *MockUnitOfWorkTx) ScooterRepository() repository.ScooterRepository {
	args := m.Called()
	return args.Get(0).(repository.ScooterRepository)
}

// TripRepository mocks the TripRepository method
func (m *MockUnitOfWorkTx) TripRepository() repository.TripRepository {
	args := m.Called()
	return args.Get(0).(repository.TripRepository)
}

// UserRepository mocks the UserRepository method
func (m *MockUnitOfWorkTx) UserRepository() repository.UserRepository {
	args := m.Called()
	return args.Get(0).(repository.UserRepository)
}

// LocationUpdateRepository mocks the LocationUpdateRepository method
func (m *MockUnitOfWorkTx) LocationUpdateRepository() repository.LocationUpdateRepository {
	args := m.Called()
	return args.Get(0).(repository.LocationUpdateRepository)
}

// Commit mocks the Commit method
func (m *MockUnitOfWorkTx) Commit() error {
	args := m.Called()
	return args.Error(0)
}

// Rollback mocks the Rollback method
func (m *MockUnitOfWorkTx) Rollback() error {
	args := m.Called()
	return args.Error(0)
}
