package services

import (
	"errors"

	"scootin-aboot/internal/models"
	"scootin-aboot/internal/repository"
	"scootin-aboot/internal/repository/mocks"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

// ScooterTestCases contains test cases for scooter service
type ScooterTestCases struct{}

// GetScootersTestCases returns test cases for GetScooters method
func (tc *ScooterTestCases) GetScootersTestCases() []struct {
	Name          string
	Params        ScooterQueryParams
	ExpectedError string
	SetupMocks    func(*mocks.MockScooterRepository)
} {
	return []struct {
		Name          string
		Params        ScooterQueryParams
		ExpectedError string
		SetupMocks    func(*mocks.MockScooterRepository)
	}{
		{
			Name:          "successful get with no filters",
			Params:        GetValidScooterQueryParams(),
			ExpectedError: "",
			SetupMocks: func(repo *mocks.MockScooterRepository) {
				repo.On("List", mock.Anything, TestData.ValidLimit, TestData.ValidOffset).Return(GetTestScooters(1), nil)
			},
		},
		{
			Name:          "successful get with status filter",
			Params:        GetValidScooterQueryParamsWithStatus("available"),
			ExpectedError: "",
			SetupMocks: func(repo *mocks.MockScooterRepository) {
				repo.On("GetByStatus", mock.Anything, models.ScooterStatusAvailable).Return(GetTestScootersWithStatus(1, models.ScooterStatusAvailable), nil)
			},
		},
		{
			Name:          "successful get with geographic bounds",
			Params:        GetValidScooterQueryParamsWithBounds(),
			ExpectedError: "",
			SetupMocks: func(repo *mocks.MockScooterRepository) {
				repo.On("GetInBounds", mock.Anything, mock.AnythingOfType("float64"), mock.AnythingOfType("float64"), mock.AnythingOfType("float64"), mock.AnythingOfType("float64")).Return(GetTestScooters(1), nil)
			},
		},
		{
			Name: "successful get with status and bounds",
			Params: func() ScooterQueryParams {
				p := GetValidScooterQueryParamsWithBounds()
				p.Status = "available"
				return p
			}(),
			ExpectedError: "",
			SetupMocks: func(repo *mocks.MockScooterRepository) {
				repo.On("GetByStatusInBounds", mock.Anything, models.ScooterStatusAvailable, mock.AnythingOfType("float64"), mock.AnythingOfType("float64"), mock.AnythingOfType("float64"), mock.AnythingOfType("float64")).Return(GetTestScooters(1), nil)
			},
		},
		{
			Name:          "repository error",
			Params:        GetValidScooterQueryParams(),
			ExpectedError: "failed to query scooters",
			SetupMocks: func(repo *mocks.MockScooterRepository) {
				repo.On("List", mock.Anything, TestData.ValidLimit, TestData.ValidOffset).Return([]*models.Scooter(nil), errors.New("database connection failed"))
			},
		},
	}
}

// GetScooterTestCases returns test cases for GetScooter method
func (tc *ScooterTestCases) GetScooterTestCases() []struct {
	Name          string
	ScooterID     uuid.UUID
	ExpectedError string
	SetupMocks    func(*mocks.MockScooterRepository, *mocks.MockTripRepository)
} {
	return []struct {
		Name          string
		ScooterID     uuid.UUID
		ExpectedError string
		SetupMocks    func(*mocks.MockScooterRepository, *mocks.MockTripRepository)
	}{
		{
			Name:          "successful get scooter",
			ScooterID:     TestData.ValidScooterID,
			ExpectedError: "",
			SetupMocks: func(scooterRepo *mocks.MockScooterRepository, tripRepo *mocks.MockTripRepository) {
				scooter := NewTestScooterBuilder().WithID(TestData.ValidScooterID).WithStatus(models.ScooterStatusAvailable).Build()
				scooterRepo.On("GetByID", mock.Anything, TestData.ValidScooterID).Return(scooter, nil)
				// No trip query needed for available scooters
			},
		},
		{
			Name:          "scooter not found",
			ScooterID:     TestData.ValidScooterID,
			ExpectedError: "scooter not found",
			SetupMocks: func(scooterRepo *mocks.MockScooterRepository, tripRepo *mocks.MockTripRepository) {
				scooterRepo.On("GetByID", mock.Anything, TestData.ValidScooterID).Return(nil, repository.ErrScooterNotFound)
			},
		},
		{
			Name:          "scooter with active trip",
			ScooterID:     TestData.ValidScooterID,
			ExpectedError: "",
			SetupMocks: func(scooterRepo *mocks.MockScooterRepository, tripRepo *mocks.MockTripRepository) {
				scooter := NewTestScooterBuilder().WithID(TestData.ValidScooterID).WithStatus(models.ScooterStatusOccupied).Build()
				trip := NewTestTripBuilder().WithScooterID(TestData.ValidScooterID).Build()
				scooterRepo.On("GetByID", mock.Anything, TestData.ValidScooterID).Return(scooter, nil)
				tripRepo.On("GetActiveByScooterID", mock.Anything, TestData.ValidScooterID).Return(trip, nil)
			},
		},
		{
			Name:          "repository error",
			ScooterID:     TestData.ValidScooterID,
			ExpectedError: "failed to get scooter",
			SetupMocks: func(scooterRepo *mocks.MockScooterRepository, tripRepo *mocks.MockTripRepository) {
				scooterRepo.On("GetByID", mock.Anything, TestData.ValidScooterID).Return(nil, errors.New("database connection failed"))
			},
		},
		{
			Name:          "trip query error",
			ScooterID:     TestData.ValidScooterID,
			ExpectedError: "",
			SetupMocks: func(scooterRepo *mocks.MockScooterRepository, tripRepo *mocks.MockTripRepository) {
				scooter := NewTestScooterBuilder().WithID(TestData.ValidScooterID).WithStatus(models.ScooterStatusOccupied).Build()
				scooterRepo.On("GetByID", mock.Anything, TestData.ValidScooterID).Return(scooter, nil)
				tripRepo.On("GetActiveByScooterID", mock.Anything, TestData.ValidScooterID).Return(nil, errors.New("trip query failed"))
			},
		},
	}
}

// GetClosestScootersTestCases returns test cases for GetClosestScooters method
func (tc *ScooterTestCases) GetClosestScootersTestCases() []struct {
	Name          string
	Params        ClosestScootersQueryParams
	ExpectedError string
	SetupMocks    func(*mocks.MockScooterRepository)
} {
	return []struct {
		Name          string
		Params        ClosestScootersQueryParams
		ExpectedError string
		SetupMocks    func(*mocks.MockScooterRepository)
	}{
		{
			Name:          "successful get closest scooters",
			Params:        GetValidClosestScootersQueryParams(),
			ExpectedError: "",
			SetupMocks: func(repo *mocks.MockScooterRepository) {
				repo.On("GetClosestWithRadius", mock.Anything, TestData.ValidLatitude, TestData.ValidLongitude, TestData.ValidRadius, "", TestData.ValidLimit).Return(GetTestScooters(1), nil)
			},
		},
		{
			Name:          "successful get closest scooters with status",
			Params:        GetValidClosestScootersQueryParamsWithStatus("available"),
			ExpectedError: "",
			SetupMocks: func(repo *mocks.MockScooterRepository) {
				repo.On("GetClosestWithRadius", mock.Anything, TestData.ValidLatitude, TestData.ValidLongitude, TestData.ValidRadius, "available", TestData.ValidLimit).Return(GetTestScootersWithStatus(1, models.ScooterStatusAvailable), nil)
			},
		},
		{
			Name:          "repository error",
			Params:        GetValidClosestScootersQueryParams(),
			ExpectedError: "failed to query closest scooters",
			SetupMocks: func(repo *mocks.MockScooterRepository) {
				repo.On("GetClosestWithRadius", mock.Anything, TestData.ValidLatitude, TestData.ValidLongitude, TestData.ValidRadius, "", TestData.ValidLimit).Return([]*models.Scooter(nil), errors.New("database connection failed"))
			},
		},
	}
}

// UpdateLocationTestCases returns test cases for UpdateLocation method
func (tc *ScooterTestCases) UpdateLocationTestCases() []struct {
	Name          string
	ScooterID     uuid.UUID
	Latitude      float64
	Longitude     float64
	ExpectedError string
	SetupMocks    func(*mocks.MockScooterRepository, *mocks.MockLocationUpdateRepository, *mocks.MockUnitOfWork, *mocks.MockUnitOfWorkTx)
} {
	return []struct {
		Name          string
		ScooterID     uuid.UUID
		Latitude      float64
		Longitude     float64
		ExpectedError string
		SetupMocks    func(*mocks.MockScooterRepository, *mocks.MockLocationUpdateRepository, *mocks.MockUnitOfWork, *mocks.MockUnitOfWorkTx)
	}{
		{
			Name:          "successful location update",
			ScooterID:     TestData.ValidScooterID,
			Latitude:      TestData.ValidLatitude,
			Longitude:     TestData.ValidLongitude,
			ExpectedError: "",
			SetupMocks: func(scooterRepo *mocks.MockScooterRepository, locationRepo *mocks.MockLocationUpdateRepository, unitOfWork *mocks.MockUnitOfWork, mockTx *mocks.MockUnitOfWorkTx) {
				scooter := NewTestScooterBuilder().WithID(TestData.ValidScooterID).Build()
				unitOfWork.On("Begin", mock.Anything).Return(mockTx, nil)
				mockTx.On("ScooterRepository").Return(scooterRepo)
				mockTx.On("LocationUpdateRepository").Return(locationRepo)
				mockTx.On("Commit").Return(nil)
				scooterRepo.On("GetByID", mock.Anything, TestData.ValidScooterID).Return(scooter, nil)
				locationRepo.On("Create", mock.Anything, mock.AnythingOfType("*models.LocationUpdate")).Return(nil)
				scooterRepo.On("UpdateLocation", mock.Anything, TestData.ValidScooterID, TestData.ValidLatitude, TestData.ValidLongitude).Return(nil)
			},
		},
		{
			Name:          "invalid coordinates",
			ScooterID:     TestData.ValidScooterID,
			Latitude:      TestData.InvalidLatitudeHigh,
			Longitude:     TestData.ValidLongitude,
			ExpectedError: "invalid coordinates",
			SetupMocks: func(*mocks.MockScooterRepository, *mocks.MockLocationUpdateRepository, *mocks.MockUnitOfWork, *mocks.MockUnitOfWorkTx) {
			},
		},
		{
			Name:          "scooter not found",
			ScooterID:     TestData.ValidScooterID,
			Latitude:      TestData.ValidLatitude,
			Longitude:     TestData.ValidLongitude,
			ExpectedError: "scooter not found",
			SetupMocks: func(scooterRepo *mocks.MockScooterRepository, locationRepo *mocks.MockLocationUpdateRepository, unitOfWork *mocks.MockUnitOfWork, mockTx *mocks.MockUnitOfWorkTx) {
				unitOfWork.On("Begin", mock.Anything).Return(mockTx, nil)
				mockTx.On("ScooterRepository").Return(scooterRepo)
				mockTx.On("LocationUpdateRepository").Return(locationRepo)
				mockTx.On("Rollback").Return(nil)
				scooterRepo.On("GetByID", mock.Anything, TestData.ValidScooterID).Return(nil, nil)
			},
		},
		{
			Name:          "get scooter error",
			ScooterID:     TestData.ValidScooterID,
			Latitude:      TestData.ValidLatitude,
			Longitude:     TestData.ValidLongitude,
			ExpectedError: "failed to get scooter",
			SetupMocks: func(scooterRepo *mocks.MockScooterRepository, locationRepo *mocks.MockLocationUpdateRepository, unitOfWork *mocks.MockUnitOfWork, mockTx *mocks.MockUnitOfWorkTx) {
				unitOfWork.On("Begin", mock.Anything).Return(mockTx, nil)
				mockTx.On("ScooterRepository").Return(scooterRepo)
				mockTx.On("LocationUpdateRepository").Return(locationRepo)
				mockTx.On("Rollback").Return(nil)
				scooterRepo.On("GetByID", mock.Anything, TestData.ValidScooterID).Return(nil, errors.New("database error"))
			},
		},
		{
			Name:          "create location update error",
			ScooterID:     TestData.ValidScooterID,
			Latitude:      TestData.ValidLatitude,
			Longitude:     TestData.ValidLongitude,
			ExpectedError: "failed to create location update",
			SetupMocks: func(scooterRepo *mocks.MockScooterRepository, locationRepo *mocks.MockLocationUpdateRepository, unitOfWork *mocks.MockUnitOfWork, mockTx *mocks.MockUnitOfWorkTx) {
				scooter := NewTestScooterBuilder().WithID(TestData.ValidScooterID).Build()
				unitOfWork.On("Begin", mock.Anything).Return(mockTx, nil)
				mockTx.On("ScooterRepository").Return(scooterRepo)
				mockTx.On("LocationUpdateRepository").Return(locationRepo)
				mockTx.On("Rollback").Return(nil)
				scooterRepo.On("GetByID", mock.Anything, TestData.ValidScooterID).Return(scooter, nil)
				locationRepo.On("Create", mock.Anything, mock.AnythingOfType("*models.LocationUpdate")).Return(errors.New("database error"))
			},
		},
		{
			Name:          "update scooter location error",
			ScooterID:     TestData.ValidScooterID,
			Latitude:      TestData.ValidLatitude,
			Longitude:     TestData.ValidLongitude,
			ExpectedError: "failed to update scooter location",
			SetupMocks: func(scooterRepo *mocks.MockScooterRepository, locationRepo *mocks.MockLocationUpdateRepository, unitOfWork *mocks.MockUnitOfWork, mockTx *mocks.MockUnitOfWorkTx) {
				scooter := NewTestScooterBuilder().WithID(TestData.ValidScooterID).Build()
				unitOfWork.On("Begin", mock.Anything).Return(mockTx, nil)
				mockTx.On("ScooterRepository").Return(scooterRepo)
				mockTx.On("LocationUpdateRepository").Return(locationRepo)
				mockTx.On("Rollback").Return(nil)
				scooterRepo.On("GetByID", mock.Anything, TestData.ValidScooterID).Return(scooter, nil)
				locationRepo.On("Create", mock.Anything, mock.AnythingOfType("*models.LocationUpdate")).Return(nil)
				scooterRepo.On("UpdateLocation", mock.Anything, TestData.ValidScooterID, TestData.ValidLatitude, TestData.ValidLongitude).Return(errors.New("database error"))
			},
		},
	}
}
