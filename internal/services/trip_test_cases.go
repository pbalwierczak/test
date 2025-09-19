package services

import (
	"errors"
	"time"

	"scootin-aboot/internal/models"
	"scootin-aboot/internal/repository/mocks"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

// TripTestCases contains test cases for trip service
type TripTestCases struct{}

// StartTripTestCases returns test cases for StartTrip method
func (tc *TripTestCases) StartTripTestCases() []struct {
	Name          string
	ScooterID     uuid.UUID
	UserID        uuid.UUID
	Latitude      float64
	Longitude     float64
	ExpectedError string
	SetupMocks    func(*mocks.MockTripRepository, *mocks.MockScooterRepository, *mocks.MockUserRepository, *mocks.MockLocationUpdateRepository, *mocks.MockUnitOfWork)
} {
	return []struct {
		Name          string
		ScooterID     uuid.UUID
		UserID        uuid.UUID
		Latitude      float64
		Longitude     float64
		ExpectedError string
		SetupMocks    func(*mocks.MockTripRepository, *mocks.MockScooterRepository, *mocks.MockUserRepository, *mocks.MockLocationUpdateRepository, *mocks.MockUnitOfWork)
	}{
		{
			Name:      "successful trip start",
			ScooterID: TestData.ValidScooterID,
			UserID:    TestData.ValidUserID,
			Latitude:  TestData.ValidLatitude,
			Longitude: TestData.ValidLongitude,
			ExpectedError: "",
			SetupMocks: func(tripRepo *mocks.MockTripRepository, scooterRepo *mocks.MockScooterRepository, userRepo *mocks.MockUserRepository, locationRepo *mocks.MockLocationUpdateRepository, unitOfWork *mocks.MockUnitOfWork) {
				mockSetup := &MockSetup{}
				mockSetup.SetupBasicUnitOfWork(unitOfWork, tripRepo, scooterRepo, userRepo, locationRepo)
				
				user := NewTestUserBuilder().WithID(TestData.ValidUserID).Build()
				userRepo.On("GetByID", mock.Anything, TestData.ValidUserID).Return(user, nil)
				tripRepo.On("GetActiveByUserID", mock.Anything, TestData.ValidUserID).Return(nil, nil)

				scooter := NewTestScooterBuilder().
					WithID(TestData.ValidScooterID).
					WithStatus(models.ScooterStatusAvailable).
					WithLocation(TestData.ValidLatitude, TestData.ValidLongitude).
					Build()
				scooterRepo.On("GetByIDForUpdate", mock.Anything, TestData.ValidScooterID).Return(scooter, nil)
				tripRepo.On("GetActiveByScooterID", mock.Anything, TestData.ValidScooterID).Return(nil, nil)
				tripRepo.On("Create", mock.Anything, mock.AnythingOfType("*models.Trip")).Return(nil)
				scooterRepo.On("UpdateStatusWithCheck", mock.Anything, TestData.ValidScooterID, models.ScooterStatusOccupied, models.ScooterStatusAvailable).Return(nil)
				scooterRepo.On("UpdateLocation", mock.Anything, TestData.ValidScooterID, TestData.ValidLatitude, TestData.ValidLongitude).Return(nil)
			},
		},
		{
			Name:      "invalid coordinates - latitude too high",
			ScooterID: TestData.ValidScooterID,
			UserID:    TestData.ValidUserID,
			Latitude:  TestData.InvalidLatitudeHigh,
			Longitude: TestData.ValidLongitude,
			ExpectedError: "invalid coordinates: invalid latitude: must be between -90 and 90",
			SetupMocks: func(tripRepo *mocks.MockTripRepository, scooterRepo *mocks.MockScooterRepository, userRepo *mocks.MockUserRepository, locationRepo *mocks.MockLocationUpdateRepository, unitOfWork *mocks.MockUnitOfWork) {
				mockSetup := &MockSetup{}
				mockSetup.SetupBasicUnitOfWork(unitOfWork, tripRepo, scooterRepo, userRepo, locationRepo)
			},
		},
		{
			Name:      "user not found",
			ScooterID: TestData.ValidScooterID,
			UserID:    TestData.ValidUserID,
			Latitude:  TestData.ValidLatitude,
			Longitude: TestData.ValidLongitude,
			ExpectedError: "user not found",
			SetupMocks: func(tripRepo *mocks.MockTripRepository, scooterRepo *mocks.MockScooterRepository, userRepo *mocks.MockUserRepository, locationRepo *mocks.MockLocationUpdateRepository, unitOfWork *mocks.MockUnitOfWork) {
				mockSetup := &MockSetup{}
				mockSetup.SetupBasicUnitOfWork(unitOfWork, tripRepo, scooterRepo, userRepo, locationRepo)
				userRepo.On("GetByID", mock.Anything, TestData.ValidUserID).Return(nil, nil)
			},
		},
		{
			Name:      "user already has active trip",
			ScooterID: TestData.ValidScooterID,
			UserID:    TestData.ValidUserID,
			Latitude:  TestData.ValidLatitude,
			Longitude: TestData.ValidLongitude,
			ExpectedError: "user already has an active trip",
			SetupMocks: func(tripRepo *mocks.MockTripRepository, scooterRepo *mocks.MockScooterRepository, userRepo *mocks.MockUserRepository, locationRepo *mocks.MockLocationUpdateRepository, unitOfWork *mocks.MockUnitOfWork) {
				mockSetup := &MockSetup{}
				mockSetup.SetupBasicUnitOfWork(unitOfWork, tripRepo, scooterRepo, userRepo, locationRepo)
				user := NewTestUserBuilder().WithID(TestData.ValidUserID).Build()
				userRepo.On("GetByID", mock.Anything, TestData.ValidUserID).Return(user, nil)
				activeTrip := NewTestTripBuilder().WithUserID(TestData.ValidUserID).Build()
				tripRepo.On("GetActiveByUserID", mock.Anything, TestData.ValidUserID).Return(activeTrip, nil)
			},
		},
		{
			Name:      "scooter not found",
			ScooterID: TestData.ValidScooterID,
			UserID:    TestData.ValidUserID,
			Latitude:  TestData.ValidLatitude,
			Longitude: TestData.ValidLongitude,
			ExpectedError: "scooter not found",
			SetupMocks: func(tripRepo *mocks.MockTripRepository, scooterRepo *mocks.MockScooterRepository, userRepo *mocks.MockUserRepository, locationRepo *mocks.MockLocationUpdateRepository, unitOfWork *mocks.MockUnitOfWork) {
				mockSetup := &MockSetup{}
				mockSetup.SetupBasicUnitOfWork(unitOfWork, tripRepo, scooterRepo, userRepo, locationRepo)
				user := NewTestUserBuilder().WithID(TestData.ValidUserID).Build()
				userRepo.On("GetByID", mock.Anything, TestData.ValidUserID).Return(user, nil)
				tripRepo.On("GetActiveByUserID", mock.Anything, TestData.ValidUserID).Return(nil, nil)
				scooterRepo.On("GetByIDForUpdate", mock.Anything, TestData.ValidScooterID).Return(nil, nil)
			},
		},
		{
			Name:      "scooter not available",
			ScooterID: TestData.ValidScooterID,
			UserID:    TestData.ValidUserID,
			Latitude:  TestData.ValidLatitude,
			Longitude: TestData.ValidLongitude,
			ExpectedError: "scooter is not available",
			SetupMocks: func(tripRepo *mocks.MockTripRepository, scooterRepo *mocks.MockScooterRepository, userRepo *mocks.MockUserRepository, locationRepo *mocks.MockLocationUpdateRepository, unitOfWork *mocks.MockUnitOfWork) {
				mockSetup := &MockSetup{}
				mockSetup.SetupBasicUnitOfWork(unitOfWork, tripRepo, scooterRepo, userRepo, locationRepo)
				user := NewTestUserBuilder().WithID(TestData.ValidUserID).Build()
				userRepo.On("GetByID", mock.Anything, TestData.ValidUserID).Return(user, nil)
				tripRepo.On("GetActiveByUserID", mock.Anything, TestData.ValidUserID).Return(nil, nil)

				scooter := NewTestScooterBuilder().
					WithID(TestData.ValidScooterID).
					WithStatus(models.ScooterStatusOccupied).
					WithLocation(TestData.ValidLatitude, TestData.ValidLongitude).
					Build()
				scooterRepo.On("GetByIDForUpdate", mock.Anything, TestData.ValidScooterID).Return(scooter, nil)
			},
		},
	}
}

// EndTripTestCases returns test cases for EndTrip method
func (tc *TripTestCases) EndTripTestCases() []struct {
	Name          string
	ScooterID     uuid.UUID
	Latitude      float64
	Longitude     float64
	ExpectedError string
	SetupMocks    func(*mocks.MockTripRepository, *mocks.MockScooterRepository, *mocks.MockUserRepository, *mocks.MockLocationUpdateRepository, *mocks.MockUnitOfWork)
} {
	return []struct {
		Name          string
		ScooterID     uuid.UUID
		Latitude      float64
		Longitude     float64
		ExpectedError string
		SetupMocks    func(*mocks.MockTripRepository, *mocks.MockScooterRepository, *mocks.MockUserRepository, *mocks.MockLocationUpdateRepository, *mocks.MockUnitOfWork)
	}{
		{
			Name:      "successful trip end",
			ScooterID: TestData.ValidScooterID,
			Latitude:  TestData.ValidLatitude + 0.001,
			Longitude: TestData.ValidLongitude + 0.001,
			ExpectedError: "",
			SetupMocks: func(tripRepo *mocks.MockTripRepository, scooterRepo *mocks.MockScooterRepository, userRepo *mocks.MockUserRepository, locationRepo *mocks.MockLocationUpdateRepository, unitOfWork *mocks.MockUnitOfWork) {
				mockSetup := &MockSetup{}
				mockSetup.SetupBasicUnitOfWork(unitOfWork, tripRepo, scooterRepo, userRepo, locationRepo)
				
				trip := NewTestTripBuilder().
					WithScooterID(TestData.ValidScooterID).
					WithUserID(TestData.ValidUserID).
					WithStartTime(time.Now().Add(-10 * time.Minute)).
					WithStatus(models.TripStatusActive).
					WithStartLocation(TestData.ValidLatitude, TestData.ValidLongitude).
					Build()
				tripRepo.On("GetActiveByScooterID", mock.Anything, TestData.ValidScooterID).Return(trip, nil)
				tripRepo.On("EndTrip", mock.Anything, TestData.ValidScooterID, TestData.ValidLatitude+0.001, TestData.ValidLongitude+0.001).Return(nil)
				scooterRepo.On("UpdateStatusWithCheck", mock.Anything, TestData.ValidScooterID, models.ScooterStatusAvailable, models.ScooterStatusOccupied).Return(nil)
				scooterRepo.On("UpdateLocation", mock.Anything, TestData.ValidScooterID, TestData.ValidLatitude+0.001, TestData.ValidLongitude+0.001).Return(nil)
			},
		},
		{
			Name:      "no active trip found",
			ScooterID: TestData.ValidScooterID,
			Latitude:  TestData.ValidLatitude,
			Longitude: TestData.ValidLongitude,
			ExpectedError: "no active trip found for scooter",
			SetupMocks: func(tripRepo *mocks.MockTripRepository, scooterRepo *mocks.MockScooterRepository, userRepo *mocks.MockUserRepository, locationRepo *mocks.MockLocationUpdateRepository, unitOfWork *mocks.MockUnitOfWork) {
				mockSetup := &MockSetup{}
				mockSetup.SetupBasicUnitOfWork(unitOfWork, tripRepo, scooterRepo, userRepo, locationRepo)
				tripRepo.On("GetActiveByScooterID", mock.Anything, TestData.ValidScooterID).Return(nil, nil)
			},
		},
		{
			Name:      "invalid coordinates",
			ScooterID: TestData.ValidScooterID,
			Latitude:  TestData.InvalidLatitudeHigh,
			Longitude: TestData.ValidLongitude,
			ExpectedError: "invalid coordinates: invalid latitude: must be between -90 and 90",
			SetupMocks: func(tripRepo *mocks.MockTripRepository, scooterRepo *mocks.MockScooterRepository, userRepo *mocks.MockUserRepository, locationRepo *mocks.MockLocationUpdateRepository, unitOfWork *mocks.MockUnitOfWork) {
				mockSetup := &MockSetup{}
				mockSetup.SetupBasicUnitOfWork(unitOfWork, tripRepo, scooterRepo, userRepo, locationRepo)
			},
		},
	}
}

// UpdateLocationTestCases returns test cases for UpdateLocation method
func (tc *TripTestCases) UpdateLocationTestCases() []struct {
	Name          string
	ScooterID     uuid.UUID
	Latitude      float64
	Longitude     float64
	ExpectedError string
	SetupMocks    func(*mocks.MockTripRepository, *mocks.MockScooterRepository, *mocks.MockUserRepository, *mocks.MockLocationUpdateRepository, *mocks.MockUnitOfWork)
} {
	return []struct {
		Name          string
		ScooterID     uuid.UUID
		Latitude      float64
		Longitude     float64
		ExpectedError string
		SetupMocks    func(*mocks.MockTripRepository, *mocks.MockScooterRepository, *mocks.MockUserRepository, *mocks.MockLocationUpdateRepository, *mocks.MockUnitOfWork)
	}{
		{
			Name:      "successful location update",
			ScooterID: TestData.ValidScooterID,
			Latitude:  TestData.ValidLatitude + 0.001,
			Longitude: TestData.ValidLongitude + 0.001,
			ExpectedError: "",
			SetupMocks: func(tripRepo *mocks.MockTripRepository, scooterRepo *mocks.MockScooterRepository, userRepo *mocks.MockUserRepository, locationRepo *mocks.MockLocationUpdateRepository, unitOfWork *mocks.MockUnitOfWork) {
				mockSetup := &MockSetup{}
				mockSetup.SetupBasicUnitOfWork(unitOfWork, tripRepo, scooterRepo, userRepo, locationRepo)
				
				trip := NewTestTripBuilder().
					WithScooterID(TestData.ValidScooterID).
					WithUserID(TestData.ValidUserID).
					WithStatus(models.TripStatusActive).
					Build()
				tripRepo.On("GetActiveByScooterID", mock.Anything, TestData.ValidScooterID).Return(trip, nil)
				locationRepo.On("Create", mock.Anything, mock.AnythingOfType("*models.LocationUpdate")).Return(nil)
				scooterRepo.On("UpdateLocation", mock.Anything, TestData.ValidScooterID, TestData.ValidLatitude+0.001, TestData.ValidLongitude+0.001).Return(nil)
			},
		},
		{
			Name:      "no active trip found",
			ScooterID: TestData.ValidScooterID,
			Latitude:  TestData.ValidLatitude,
			Longitude: TestData.ValidLongitude,
			ExpectedError: "no active trip found for scooter",
			SetupMocks: func(tripRepo *mocks.MockTripRepository, scooterRepo *mocks.MockScooterRepository, userRepo *mocks.MockUserRepository, locationRepo *mocks.MockLocationUpdateRepository, unitOfWork *mocks.MockUnitOfWork) {
				mockSetup := &MockSetup{}
				mockSetup.SetupBasicUnitOfWork(unitOfWork, tripRepo, scooterRepo, userRepo, locationRepo)
				tripRepo.On("GetActiveByScooterID", mock.Anything, TestData.ValidScooterID).Return(nil, nil)
			},
		},
		{
			Name:      "invalid coordinates",
			ScooterID: TestData.ValidScooterID,
			Latitude:  TestData.InvalidLatitudeHigh,
			Longitude: TestData.ValidLongitude,
			ExpectedError: "invalid coordinates: invalid latitude: must be between -90 and 90",
			SetupMocks: func(tripRepo *mocks.MockTripRepository, scooterRepo *mocks.MockScooterRepository, userRepo *mocks.MockUserRepository, locationRepo *mocks.MockLocationUpdateRepository, unitOfWork *mocks.MockUnitOfWork) {
				mockSetup := &MockSetup{}
				mockSetup.SetupBasicUnitOfWork(unitOfWork, tripRepo, scooterRepo, userRepo, locationRepo)
			},
		},
	}
}

// CancelTripTestCases returns test cases for CancelTrip method
func (tc *TripTestCases) CancelTripTestCases() []struct {
	Name          string
	ScooterID     uuid.UUID
	ExpectedError string
	SetupMocks    func(*mocks.MockTripRepository, *mocks.MockScooterRepository, *mocks.MockUserRepository, *mocks.MockLocationUpdateRepository, *mocks.MockUnitOfWork)
} {
	return []struct {
		Name          string
		ScooterID     uuid.UUID
		ExpectedError string
		SetupMocks    func(*mocks.MockTripRepository, *mocks.MockScooterRepository, *mocks.MockUserRepository, *mocks.MockLocationUpdateRepository, *mocks.MockUnitOfWork)
	}{
		{
			Name:      "successful trip cancellation",
			ScooterID: TestData.ValidScooterID,
			ExpectedError: "",
			SetupMocks: func(tripRepo *mocks.MockTripRepository, scooterRepo *mocks.MockScooterRepository, userRepo *mocks.MockUserRepository, locationRepo *mocks.MockLocationUpdateRepository, unitOfWork *mocks.MockUnitOfWork) {
				mockSetup := &MockSetup{}
				mockSetup.SetupBasicUnitOfWork(unitOfWork, tripRepo, scooterRepo, userRepo, locationRepo)
				
				trip := NewTestTripBuilder().
					WithScooterID(TestData.ValidScooterID).
					WithUserID(TestData.ValidUserID).
					WithStatus(models.TripStatusActive).
					Build()
				tripRepo.On("GetActiveByScooterID", mock.Anything, TestData.ValidScooterID).Return(trip, nil)
				tripRepo.On("CancelTrip", mock.Anything, TestData.ValidScooterID).Return(nil)
				scooterRepo.On("UpdateStatusWithCheck", mock.Anything, TestData.ValidScooterID, models.ScooterStatusAvailable, models.ScooterStatusOccupied).Return(nil)
			},
		},
		{
			Name:      "no active trip found",
			ScooterID: TestData.ValidScooterID,
			ExpectedError: "no active trip found for scooter",
			SetupMocks: func(tripRepo *mocks.MockTripRepository, scooterRepo *mocks.MockScooterRepository, userRepo *mocks.MockUserRepository, locationRepo *mocks.MockLocationUpdateRepository, unitOfWork *mocks.MockUnitOfWork) {
				mockSetup := &MockSetup{}
				mockSetup.SetupBasicUnitOfWork(unitOfWork, tripRepo, scooterRepo, userRepo, locationRepo)
				tripRepo.On("GetActiveByScooterID", mock.Anything, TestData.ValidScooterID).Return(nil, nil)
			},
		},
		{
			Name:      "repository error when getting active trip",
			ScooterID: TestData.ValidScooterID,
			ExpectedError: "failed to get active trip: database error",
			SetupMocks: func(tripRepo *mocks.MockTripRepository, scooterRepo *mocks.MockScooterRepository, userRepo *mocks.MockUserRepository, locationRepo *mocks.MockLocationUpdateRepository, unitOfWork *mocks.MockUnitOfWork) {
				mockSetup := &MockSetup{}
				mockSetup.SetupBasicUnitOfWork(unitOfWork, tripRepo, scooterRepo, userRepo, locationRepo)
				tripRepo.On("GetActiveByScooterID", mock.Anything, TestData.ValidScooterID).Return(nil, errors.New("database error"))
			},
		},
		{
			Name:      "repository error when cancelling trip",
			ScooterID: TestData.ValidScooterID,
			ExpectedError: "failed to cancel trip: database error",
			SetupMocks: func(tripRepo *mocks.MockTripRepository, scooterRepo *mocks.MockScooterRepository, userRepo *mocks.MockUserRepository, locationRepo *mocks.MockLocationUpdateRepository, unitOfWork *mocks.MockUnitOfWork) {
				mockSetup := &MockSetup{}
				mockSetup.SetupBasicUnitOfWork(unitOfWork, tripRepo, scooterRepo, userRepo, locationRepo)
				
				trip := NewTestTripBuilder().
					WithScooterID(TestData.ValidScooterID).
					WithUserID(TestData.ValidUserID).
					WithStatus(models.TripStatusActive).
					Build()
				tripRepo.On("GetActiveByScooterID", mock.Anything, TestData.ValidScooterID).Return(trip, nil)
				tripRepo.On("CancelTrip", mock.Anything, TestData.ValidScooterID).Return(errors.New("database error"))
			},
		},
		{
			Name:      "repository error when updating scooter status",
			ScooterID: TestData.ValidScooterID,
			ExpectedError: "failed to update scooter status: database error",
			SetupMocks: func(tripRepo *mocks.MockTripRepository, scooterRepo *mocks.MockScooterRepository, userRepo *mocks.MockUserRepository, locationRepo *mocks.MockLocationUpdateRepository, unitOfWork *mocks.MockUnitOfWork) {
				mockSetup := &MockSetup{}
				mockSetup.SetupBasicUnitOfWork(unitOfWork, tripRepo, scooterRepo, userRepo, locationRepo)
				
				trip := NewTestTripBuilder().
					WithScooterID(TestData.ValidScooterID).
					WithUserID(TestData.ValidUserID).
					WithStatus(models.TripStatusActive).
					Build()
				tripRepo.On("GetActiveByScooterID", mock.Anything, TestData.ValidScooterID).Return(trip, nil)
				tripRepo.On("CancelTrip", mock.Anything, TestData.ValidScooterID).Return(nil)
				scooterRepo.On("UpdateStatusWithCheck", mock.Anything, TestData.ValidScooterID, models.ScooterStatusAvailable, models.ScooterStatusOccupied).Return(errors.New("database error"))
			},
		},
	}
}

// GetActiveTripTestCases returns test cases for GetActiveTrip method
func (tc *TripTestCases) GetActiveTripTestCases() []struct {
	Name          string
	ScooterID     uuid.UUID
	ExpectedError string
	ExpectedTrip  *models.Trip
	SetupMocks    func(*mocks.MockTripRepository, *mocks.MockScooterRepository, *mocks.MockUserRepository, *mocks.MockLocationUpdateRepository, *mocks.MockUnitOfWork)
} {
	return []struct {
		Name          string
		ScooterID     uuid.UUID
		ExpectedError string
		ExpectedTrip  *models.Trip
		SetupMocks    func(*mocks.MockTripRepository, *mocks.MockScooterRepository, *mocks.MockUserRepository, *mocks.MockLocationUpdateRepository, *mocks.MockUnitOfWork)
	}{
		{
			Name:      "successful get active trip",
			ScooterID: TestData.ValidScooterID,
			ExpectedError: "",
			ExpectedTrip: &models.Trip{
				Status: models.TripStatusActive,
			},
			SetupMocks: func(tripRepo *mocks.MockTripRepository, scooterRepo *mocks.MockScooterRepository, userRepo *mocks.MockUserRepository, locationRepo *mocks.MockLocationUpdateRepository, unitOfWork *mocks.MockUnitOfWork) {
				mockSetup := &MockSetup{}
				mockSetup.SetupBasicUnitOfWork(unitOfWork, tripRepo, scooterRepo, userRepo, locationRepo)
				
				trip := NewTestTripBuilder().
					WithScooterID(TestData.ValidScooterID).
					WithUserID(TestData.ValidUserID).
					WithStatus(models.TripStatusActive).
					Build()
				tripRepo.On("GetActiveByScooterID", mock.Anything, TestData.ValidScooterID).Return(trip, nil)
			},
		},
		{
			Name:      "no active trip found",
			ScooterID: TestData.ValidScooterID,
			ExpectedError: "",
			ExpectedTrip:  nil,
			SetupMocks: func(tripRepo *mocks.MockTripRepository, scooterRepo *mocks.MockScooterRepository, userRepo *mocks.MockUserRepository, locationRepo *mocks.MockLocationUpdateRepository, unitOfWork *mocks.MockUnitOfWork) {
				mockSetup := &MockSetup{}
				mockSetup.SetupBasicUnitOfWork(unitOfWork, tripRepo, scooterRepo, userRepo, locationRepo)
				tripRepo.On("GetActiveByScooterID", mock.Anything, TestData.ValidScooterID).Return(nil, nil)
			},
		},
		{
			Name:      "repository error",
			ScooterID: TestData.ValidScooterID,
			ExpectedError: "failed to get active trip: database error",
			ExpectedTrip:  nil,
			SetupMocks: func(tripRepo *mocks.MockTripRepository, scooterRepo *mocks.MockScooterRepository, userRepo *mocks.MockUserRepository, locationRepo *mocks.MockLocationUpdateRepository, unitOfWork *mocks.MockUnitOfWork) {
				mockSetup := &MockSetup{}
				mockSetup.SetupBasicUnitOfWork(unitOfWork, tripRepo, scooterRepo, userRepo, locationRepo)
				tripRepo.On("GetActiveByScooterID", mock.Anything, TestData.ValidScooterID).Return(nil, errors.New("database error"))
			},
		},
	}
}

// GetActiveTripByUserTestCases returns test cases for GetActiveTripByUser method
func (tc *TripTestCases) GetActiveTripByUserTestCases() []struct {
	Name          string
	UserID        uuid.UUID
	ExpectedError string
	ExpectedTrip  *models.Trip
	SetupMocks    func(*mocks.MockTripRepository, *mocks.MockScooterRepository, *mocks.MockUserRepository, *mocks.MockLocationUpdateRepository, *mocks.MockUnitOfWork)
} {
	return []struct {
		Name          string
		UserID        uuid.UUID
		ExpectedError string
		ExpectedTrip  *models.Trip
		SetupMocks    func(*mocks.MockTripRepository, *mocks.MockScooterRepository, *mocks.MockUserRepository, *mocks.MockLocationUpdateRepository, *mocks.MockUnitOfWork)
	}{
		{
			Name:      "successful get active trip by user",
			UserID:    TestData.ValidUserID,
			ExpectedError: "",
			ExpectedTrip: &models.Trip{
				Status: models.TripStatusActive,
			},
			SetupMocks: func(tripRepo *mocks.MockTripRepository, scooterRepo *mocks.MockScooterRepository, userRepo *mocks.MockUserRepository, locationRepo *mocks.MockLocationUpdateRepository, unitOfWork *mocks.MockUnitOfWork) {
				mockSetup := &MockSetup{}
				mockSetup.SetupBasicUnitOfWork(unitOfWork, tripRepo, scooterRepo, userRepo, locationRepo)
				
				trip := NewTestTripBuilder().
					WithUserID(TestData.ValidUserID).
					WithStatus(models.TripStatusActive).
					Build()
				tripRepo.On("GetActiveByUserID", mock.Anything, TestData.ValidUserID).Return(trip, nil)
			},
		},
		{
			Name:      "no active trip found for user",
			UserID:    TestData.ValidUserID,
			ExpectedError: "",
			ExpectedTrip:  nil,
			SetupMocks: func(tripRepo *mocks.MockTripRepository, scooterRepo *mocks.MockScooterRepository, userRepo *mocks.MockUserRepository, locationRepo *mocks.MockLocationUpdateRepository, unitOfWork *mocks.MockUnitOfWork) {
				mockSetup := &MockSetup{}
				mockSetup.SetupBasicUnitOfWork(unitOfWork, tripRepo, scooterRepo, userRepo, locationRepo)
				tripRepo.On("GetActiveByUserID", mock.Anything, TestData.ValidUserID).Return(nil, nil)
			},
		},
		{
			Name:      "repository error",
			UserID:    TestData.ValidUserID,
			ExpectedError: "failed to get active trip: database error",
			ExpectedTrip:  nil,
			SetupMocks: func(tripRepo *mocks.MockTripRepository, scooterRepo *mocks.MockScooterRepository, userRepo *mocks.MockUserRepository, locationRepo *mocks.MockLocationUpdateRepository, unitOfWork *mocks.MockUnitOfWork) {
				mockSetup := &MockSetup{}
				mockSetup.SetupBasicUnitOfWork(unitOfWork, tripRepo, scooterRepo, userRepo, locationRepo)
				tripRepo.On("GetActiveByUserID", mock.Anything, TestData.ValidUserID).Return(nil, errors.New("database error"))
			},
		},
	}
}

// GetTripTestCases returns test cases for GetTrip method
func (tc *TripTestCases) GetTripTestCases() []struct {
	Name          string
	TripID        uuid.UUID
	ExpectedError string
	ExpectedTrip  *models.Trip
	SetupMocks    func(*mocks.MockTripRepository, *mocks.MockScooterRepository, *mocks.MockUserRepository, *mocks.MockLocationUpdateRepository, *mocks.MockUnitOfWork)
} {
	return []struct {
		Name          string
		TripID        uuid.UUID
		ExpectedError string
		ExpectedTrip  *models.Trip
		SetupMocks    func(*mocks.MockTripRepository, *mocks.MockScooterRepository, *mocks.MockUserRepository, *mocks.MockLocationUpdateRepository, *mocks.MockUnitOfWork)
	}{
		{
			Name:      "successful get trip",
			TripID:    TestData.ValidTripID,
			ExpectedError: "",
			ExpectedTrip: &models.Trip{
				Status: models.TripStatusCompleted,
			},
			SetupMocks: func(tripRepo *mocks.MockTripRepository, scooterRepo *mocks.MockScooterRepository, userRepo *mocks.MockUserRepository, locationRepo *mocks.MockLocationUpdateRepository, unitOfWork *mocks.MockUnitOfWork) {
				mockSetup := &MockSetup{}
				mockSetup.SetupBasicUnitOfWork(unitOfWork, tripRepo, scooterRepo, userRepo, locationRepo)
				
				trip := NewTestTripBuilder().
					WithID(TestData.ValidTripID).
					WithUserID(TestData.ValidUserID).
					WithStatus(models.TripStatusCompleted).
					Build()
				tripRepo.On("GetByID", mock.Anything, TestData.ValidTripID).Return(trip, nil)
			},
		},
		{
			Name:      "trip not found",
			TripID:    TestData.ValidTripID,
			ExpectedError: "trip not found",
			ExpectedTrip:  nil,
			SetupMocks: func(tripRepo *mocks.MockTripRepository, scooterRepo *mocks.MockScooterRepository, userRepo *mocks.MockUserRepository, locationRepo *mocks.MockLocationUpdateRepository, unitOfWork *mocks.MockUnitOfWork) {
				mockSetup := &MockSetup{}
				mockSetup.SetupBasicUnitOfWork(unitOfWork, tripRepo, scooterRepo, userRepo, locationRepo)
				tripRepo.On("GetByID", mock.Anything, TestData.ValidTripID).Return(nil, nil)
			},
		},
		{
			Name:      "repository error",
			TripID:    TestData.ValidTripID,
			ExpectedError: "failed to get trip: database error",
			ExpectedTrip:  nil,
			SetupMocks: func(tripRepo *mocks.MockTripRepository, scooterRepo *mocks.MockScooterRepository, userRepo *mocks.MockUserRepository, locationRepo *mocks.MockLocationUpdateRepository, unitOfWork *mocks.MockUnitOfWork) {
				mockSetup := &MockSetup{}
				mockSetup.SetupBasicUnitOfWork(unitOfWork, tripRepo, scooterRepo, userRepo, locationRepo)
				tripRepo.On("GetByID", mock.Anything, TestData.ValidTripID).Return(nil, errors.New("database error"))
			},
		},
	}
}
