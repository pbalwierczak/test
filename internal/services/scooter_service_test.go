package services

import (
	"testing"

	"scootin-aboot/internal/models"
	"scootin-aboot/internal/repository/mocks"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestNewScooterService(t *testing.T) {
	mockSetup := &MockSetup{}
	service, _, _, _, _ := mockSetup.CreateTestScooterService()

	assert.NotNil(t, service)
}

func TestScooterService_GetScooters(t *testing.T) {
	testCases := &ScooterTestCases{}
	cases := testCases.GetScootersTestCases()

	for _, tc := range cases {
		t.Run(tc.Name, func(t *testing.T) {
			mockSetup := &MockSetup{}
			service, scooterRepo, _, _, _ := mockSetup.CreateTestScooterService()

			tc.SetupMocks(scooterRepo)

			result, err := service.GetScooters(TestContext(), tc.Params)

			if tc.ExpectedError != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tc.ExpectedError)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Len(t, result.Scooters, 1)
			}

			scooterRepo.AssertExpectations(t)
		})
	}
}

func TestScooterService_GetScooters_InvalidParams(t *testing.T) {
	invalidParams := GetInvalidScooterQueryParams()
	mockSetup := &MockSetup{}

	for _, tc := range invalidParams {
		t.Run(tc.Name, func(t *testing.T) {
			service, _, _, _, _ := mockSetup.CreateTestScooterService()

			result, err := service.GetScooters(TestContext(), tc.Params)

			assert.Error(t, err)
			assert.Nil(t, result)
		})
	}
}

func TestScooterService_GetScooters_Pagination(t *testing.T) {
	mockSetup := &MockSetup{}
	service, scooterRepo, _, _, _ := mockSetup.CreateTestScooterService()

	params := ScooterQueryParams{
		Status: "available",
		Limit:  5,
		Offset: 10,
	}

	testScooters := GetTestScooters(15)
	scooterRepo.On("GetByStatus", mock.Anything, mock.Anything).Return(testScooters, nil)

	result, err := service.GetScooters(TestContext(), params)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result.Scooters, 5)
	assert.Equal(t, int64(15), result.Total)
	assert.Equal(t, 5, result.Limit)
	assert.Equal(t, 10, result.Offset)

	scooterRepo.AssertExpectations(t)
}

func TestScooterService_GetScooter(t *testing.T) {
	testCases := &ScooterTestCases{}
	cases := testCases.GetScooterTestCases()

	for _, tc := range cases {
		t.Run(tc.Name, func(t *testing.T) {
			mockSetup := &MockSetup{}
			service, scooterRepo, tripRepo, _, _ := mockSetup.CreateTestScooterService()

			tc.SetupMocks(scooterRepo, tripRepo)

			result, err := service.GetScooter(TestContext(), tc.ScooterID)

			if tc.ExpectedError != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tc.ExpectedError)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, tc.ScooterID, result.ID)
			}

			scooterRepo.AssertExpectations(t)
			tripRepo.AssertExpectations(t)
		})
	}
}

func TestScooterService_GetClosestScooters(t *testing.T) {
	testCases := &ScooterTestCases{}
	cases := testCases.GetClosestScootersTestCases()

	for _, tc := range cases {
		t.Run(tc.Name, func(t *testing.T) {
			mockSetup := &MockSetup{}
			service, scooterRepo, _, _, _ := mockSetup.CreateTestScooterService()

			tc.SetupMocks(scooterRepo)

			result, err := service.GetClosestScooters(TestContext(), tc.Params)

			if tc.ExpectedError != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tc.ExpectedError)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Len(t, result.Scooters, 1)
				assert.Equal(t, tc.Params.Latitude, result.Center.Latitude)
				assert.Equal(t, tc.Params.Longitude, result.Center.Longitude)
				assert.Equal(t, float64(tc.Params.Radius), result.Radius)
			}

			scooterRepo.AssertExpectations(t)
		})
	}
}

func TestScooterService_GetClosestScooters_InvalidParams(t *testing.T) {
	invalidParams := GetInvalidClosestScootersQueryParams()
	mockSetup := &MockSetup{}

	for _, tc := range invalidParams {
		t.Run(tc.Name, func(t *testing.T) {
			service, _, _, _, _ := mockSetup.CreateTestScooterService()

			result, err := service.GetClosestScooters(TestContext(), tc.Params)

			assert.Error(t, err)
			assert.Nil(t, result)
		})
	}
}

func TestScooterService_UpdateLocation(t *testing.T) {
	testCases := &ScooterTestCases{}
	cases := testCases.UpdateLocationTestCases()

	for _, tc := range cases {
		t.Run(tc.Name, func(t *testing.T) {
			mockSetup := &MockSetup{}
			service, scooterRepo, _, locationRepo, unitOfWork := mockSetup.CreateTestScooterService()
			mockTx := &mocks.MockUnitOfWorkTx{}

			tc.SetupMocks(scooterRepo, locationRepo, unitOfWork, mockTx)

			err := service.UpdateLocation(TestContext(), tc.ScooterID, tc.Latitude, tc.Longitude)

			if tc.ExpectedError != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tc.ExpectedError)
			} else {
				assert.NoError(t, err)
			}

			unitOfWork.AssertExpectations(t)
			mockTx.AssertExpectations(t)
			scooterRepo.AssertExpectations(t)
			locationRepo.AssertExpectations(t)
		})
	}
}

// Validation tests
func TestScooterService_ValidateScooterQueryParams(t *testing.T) {
	service := &scooterService{}

	t.Run("valid params", func(t *testing.T) {
		params := GetValidScooterQueryParamsWithStatus("available")
		err := service.validateScooterQueryParams(params)
		assert.NoError(t, err)
	})

	t.Run("invalid status", func(t *testing.T) {
		params := ScooterQueryParams{Status: "invalid", Limit: 10, Offset: 0}
		err := service.validateScooterQueryParams(params)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "status must be 'available' or 'occupied'")
	})

	t.Run("negative limit", func(t *testing.T) {
		params := ScooterQueryParams{Limit: -1, Offset: 0}
		err := service.validateScooterQueryParams(params)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "limit must be non-negative")
	})

	t.Run("negative offset", func(t *testing.T) {
		params := ScooterQueryParams{Limit: 10, Offset: -1}
		err := service.validateScooterQueryParams(params)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "offset must be non-negative")
	})

	t.Run("excessive limit", func(t *testing.T) {
		params := ScooterQueryParams{Limit: 101, Offset: 0}
		err := service.validateScooterQueryParams(params)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "limit cannot exceed 100")
	})
}

func TestScooterService_ValidateClosestScootersParams(t *testing.T) {
	service := &scooterService{}

	t.Run("valid params", func(t *testing.T) {
		params := GetValidClosestScootersQueryParams()
		err := service.validateClosestScootersParams(params)
		assert.NoError(t, err)
	})

	t.Run("invalid latitude", func(t *testing.T) {
		params := ClosestScootersQueryParams{
			Latitude:  TestData.InvalidLatitudeHigh,
			Longitude: TestData.ValidLongitude,
			Radius:    TestData.ValidRadius,
			Limit:     TestData.ValidLimit,
		}
		err := service.validateClosestScootersParams(params)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid latitude")
	})

	t.Run("negative radius", func(t *testing.T) {
		params := ClosestScootersQueryParams{
			Latitude:  TestData.ValidLatitude,
			Longitude: TestData.ValidLongitude,
			Radius:    -100,
			Limit:     TestData.ValidLimit,
		}
		err := service.validateClosestScootersParams(params)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "radius must be non-negative")
	})

	t.Run("excessive radius", func(t *testing.T) {
		params := ClosestScootersQueryParams{
			Latitude:  TestData.ValidLatitude,
			Longitude: TestData.ValidLongitude,
			Radius:    TestData.ExcessiveRadius,
			Limit:     TestData.ValidLimit,
		}
		err := service.validateClosestScootersParams(params)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "radius cannot exceed 50000 meters")
	})

	t.Run("invalid status", func(t *testing.T) {
		params := ClosestScootersQueryParams{
			Latitude:  TestData.ValidLatitude,
			Longitude: TestData.ValidLongitude,
			Radius:    TestData.ValidRadius,
			Limit:     TestData.ValidLimit,
			Status:    "invalid",
		}
		err := service.validateClosestScootersParams(params)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "status must be 'available' or 'occupied'")
	})

	t.Run("excessive limit", func(t *testing.T) {
		params := ClosestScootersQueryParams{
			Latitude:  TestData.ValidLatitude,
			Longitude: TestData.ValidLongitude,
			Radius:    TestData.ValidRadius,
			Limit:     51,
		}
		err := service.validateClosestScootersParams(params)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "limit cannot exceed 50")
	})
}

func TestScooterService_MapScooterToInfo(t *testing.T) {
	service := &scooterService{}

	t.Run("complete scooter", func(t *testing.T) {
		scooter := NewTestScooterBuilder().Build()
		result := service.mapScooterToInfo(scooter)

		assert.Equal(t, scooter.ID, result.ID)
		assert.Equal(t, string(scooter.Status), result.Status)
		assert.Equal(t, scooter.CurrentLatitude, result.CurrentLatitude)
		assert.Equal(t, scooter.CurrentLongitude, result.CurrentLongitude)
		assert.Equal(t, scooter.LastSeen, result.LastSeen)
		assert.Equal(t, scooter.CreatedAt, result.CreatedAt)
	})

	t.Run("minimal scooter", func(t *testing.T) {
		scooter := NewTestScooterBuilder().
			WithStatus(models.ScooterStatusOccupied).
			WithLocation(0.0, 0.0).
			Build()
		result := service.mapScooterToInfo(scooter)

		assert.Equal(t, scooter.ID, result.ID)
		assert.Equal(t, string(scooter.Status), result.Status)
		assert.Equal(t, scooter.CurrentLatitude, result.CurrentLatitude)
		assert.Equal(t, scooter.CurrentLongitude, result.CurrentLongitude)
	})
}
