package mocks

import (
	"context"

	"scootin-aboot/internal/repository"

	"github.com/stretchr/testify/mock"
)

type MockUnitOfWork struct {
	mock.Mock
}

func (m *MockUnitOfWork) Begin(ctx context.Context) (repository.UnitOfWorkTx, error) {
	args := m.Called(ctx)
	return args.Get(0).(repository.UnitOfWorkTx), args.Error(1)
}

type MockUnitOfWorkTx struct {
	mock.Mock
}

func (m *MockUnitOfWorkTx) ScooterRepository() repository.ScooterRepository {
	args := m.Called()
	return args.Get(0).(repository.ScooterRepository)
}

func (m *MockUnitOfWorkTx) TripRepository() repository.TripRepository {
	args := m.Called()
	return args.Get(0).(repository.TripRepository)
}

func (m *MockUnitOfWorkTx) UserRepository() repository.UserRepository {
	args := m.Called()
	return args.Get(0).(repository.UserRepository)
}

func (m *MockUnitOfWorkTx) LocationUpdateRepository() repository.LocationUpdateRepository {
	args := m.Called()
	return args.Get(0).(repository.LocationUpdateRepository)
}

func (m *MockUnitOfWorkTx) Commit() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockUnitOfWorkTx) Rollback() error {
	args := m.Called()
	return args.Error(0)
}
