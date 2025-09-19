package services

import (
	"testing"

	"scootin-aboot/internal/models"

	"github.com/stretchr/testify/assert"
)

func TestNewTripService(t *testing.T) {
	mockSetup := &MockSetup{}
	service, _, _, _, _, _ := mockSetup.CreateTestTripService()

	assert.NotNil(t, service)
}

func TestTripService_StartTrip(t *testing.T) {
	testCases := &TripTestCases{}
	cases := testCases.StartTripTestCases()

	for _, tc := range cases {
		t.Run(tc.Name, func(t *testing.T) {
			mockSetup := &MockSetup{}
			service, tripRepo, scooterRepo, userRepo, locationRepo, unitOfWork := mockSetup.CreateTestTripService()

			tc.SetupMocks(tripRepo, scooterRepo, userRepo, locationRepo, unitOfWork)

			trip, err := service.StartTrip(TestContext(), tc.ScooterID, tc.UserID, tc.Latitude, tc.Longitude)

			if tc.ExpectedError != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tc.ExpectedError)
				assert.Nil(t, trip)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, trip)
				assert.Equal(t, tc.ScooterID, trip.ScooterID)
				assert.Equal(t, tc.UserID, trip.UserID)
				assert.Equal(t, tc.Latitude, trip.StartLatitude)
				assert.Equal(t, tc.Longitude, trip.StartLongitude)
				assert.Equal(t, models.TripStatusActive, trip.Status)
			}

			tripRepo.AssertExpectations(t)
			scooterRepo.AssertExpectations(t)
			userRepo.AssertExpectations(t)
			locationRepo.AssertExpectations(t)
			unitOfWork.AssertExpectations(t)
		})
	}
}

func TestTripService_EndTrip(t *testing.T) {
	testCases := &TripTestCases{}
	cases := testCases.EndTripTestCases()

	for _, tc := range cases {
		t.Run(tc.Name, func(t *testing.T) {
			mockSetup := &MockSetup{}
			service, tripRepo, scooterRepo, userRepo, locationRepo, unitOfWork := mockSetup.CreateTestTripService()

			tc.SetupMocks(tripRepo, scooterRepo, userRepo, locationRepo, unitOfWork)

			trip, err := service.EndTrip(TestContext(), tc.ScooterID, tc.Latitude, tc.Longitude)

			if tc.ExpectedError != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tc.ExpectedError)
				assert.Nil(t, trip)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, trip)
				assert.Equal(t, models.TripStatusCompleted, trip.Status)
				assert.NotNil(t, trip.EndTime)
				assert.NotNil(t, trip.EndLatitude)
				assert.NotNil(t, trip.EndLongitude)
			}

			tripRepo.AssertExpectations(t)
			scooterRepo.AssertExpectations(t)
			userRepo.AssertExpectations(t)
			locationRepo.AssertExpectations(t)
			unitOfWork.AssertExpectations(t)
		})
	}
}

func TestTripService_UpdateLocation(t *testing.T) {
	testCases := &TripTestCases{}
	cases := testCases.UpdateLocationTestCases()

	for _, tc := range cases {
		t.Run(tc.Name, func(t *testing.T) {
			mockSetup := &MockSetup{}
			service, tripRepo, scooterRepo, userRepo, locationRepo, unitOfWork := mockSetup.CreateTestTripService()

			tc.SetupMocks(tripRepo, scooterRepo, userRepo, locationRepo, unitOfWork)

			err := service.UpdateLocation(TestContext(), tc.ScooterID, tc.Latitude, tc.Longitude)

			if tc.ExpectedError != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tc.ExpectedError)
			} else {
				assert.NoError(t, err)
			}

			tripRepo.AssertExpectations(t)
			scooterRepo.AssertExpectations(t)
			userRepo.AssertExpectations(t)
			locationRepo.AssertExpectations(t)
			unitOfWork.AssertExpectations(t)
		})
	}
}

func TestTripService_CancelTrip(t *testing.T) {
	testCases := &TripTestCases{}
	cases := testCases.CancelTripTestCases()

	for _, tc := range cases {
		t.Run(tc.Name, func(t *testing.T) {
			mockSetup := &MockSetup{}
			service, tripRepo, scooterRepo, userRepo, locationRepo, unitOfWork := mockSetup.CreateTestTripService()

			tc.SetupMocks(tripRepo, scooterRepo, userRepo, locationRepo, unitOfWork)

			trip, err := service.CancelTrip(TestContext(), tc.ScooterID)

			if tc.ExpectedError != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tc.ExpectedError)
				assert.Nil(t, trip)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, trip)
				assert.Equal(t, models.TripStatusCancelled, trip.Status)
			}

			tripRepo.AssertExpectations(t)
			scooterRepo.AssertExpectations(t)
			userRepo.AssertExpectations(t)
			locationRepo.AssertExpectations(t)
			unitOfWork.AssertExpectations(t)
		})
	}
}

func TestTripService_GetActiveTrip(t *testing.T) {
	testCases := &TripTestCases{}
	cases := testCases.GetActiveTripTestCases()

	for _, tc := range cases {
		t.Run(tc.Name, func(t *testing.T) {
			mockSetup := &MockSetup{}
			service, tripRepo, scooterRepo, userRepo, locationRepo, unitOfWork := mockSetup.CreateTestTripService()

			tc.SetupMocks(tripRepo, scooterRepo, userRepo, locationRepo, unitOfWork)

			trip, err := service.GetActiveTrip(TestContext(), tc.ScooterID)

			if tc.ExpectedError != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tc.ExpectedError)
				assert.Nil(t, trip)
			} else {
				assert.NoError(t, err)
				if tc.ExpectedTrip == nil {
					assert.Nil(t, trip)
				} else {
					assert.NotNil(t, trip)
					assert.Equal(t, tc.ExpectedTrip.Status, trip.Status)
				}
			}

			tripRepo.AssertExpectations(t)
			scooterRepo.AssertExpectations(t)
			userRepo.AssertExpectations(t)
			locationRepo.AssertExpectations(t)
			unitOfWork.AssertExpectations(t)
		})
	}
}

func TestTripService_GetActiveTripByUser(t *testing.T) {
	testCases := &TripTestCases{}
	cases := testCases.GetActiveTripByUserTestCases()

	for _, tc := range cases {
		t.Run(tc.Name, func(t *testing.T) {
			mockSetup := &MockSetup{}
			service, tripRepo, scooterRepo, userRepo, locationRepo, unitOfWork := mockSetup.CreateTestTripService()

			tc.SetupMocks(tripRepo, scooterRepo, userRepo, locationRepo, unitOfWork)

			trip, err := service.GetActiveTripByUser(TestContext(), tc.UserID)

			if tc.ExpectedError != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tc.ExpectedError)
				assert.Nil(t, trip)
			} else {
				assert.NoError(t, err)
				if tc.ExpectedTrip == nil {
					assert.Nil(t, trip)
				} else {
					assert.NotNil(t, trip)
					assert.Equal(t, tc.ExpectedTrip.Status, trip.Status)
				}
			}

			tripRepo.AssertExpectations(t)
			scooterRepo.AssertExpectations(t)
			userRepo.AssertExpectations(t)
			locationRepo.AssertExpectations(t)
			unitOfWork.AssertExpectations(t)
		})
	}
}

func TestTripService_GetTrip(t *testing.T) {
	testCases := &TripTestCases{}
	cases := testCases.GetTripTestCases()

	for _, tc := range cases {
		t.Run(tc.Name, func(t *testing.T) {
			mockSetup := &MockSetup{}
			service, tripRepo, scooterRepo, userRepo, locationRepo, unitOfWork := mockSetup.CreateTestTripService()

			tc.SetupMocks(tripRepo, scooterRepo, userRepo, locationRepo, unitOfWork)

			trip, err := service.GetTrip(TestContext(), tc.TripID)

			if tc.ExpectedError != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tc.ExpectedError)
				assert.Nil(t, trip)
			} else {
				assert.NoError(t, err)
				if tc.ExpectedTrip == nil {
					assert.Nil(t, trip)
				} else {
					assert.NotNil(t, trip)
					assert.Equal(t, tc.ExpectedTrip.Status, trip.Status)
				}
			}

			tripRepo.AssertExpectations(t)
			scooterRepo.AssertExpectations(t)
			userRepo.AssertExpectations(t)
			locationRepo.AssertExpectations(t)
			unitOfWork.AssertExpectations(t)
		})
	}
}
