package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"scootin-aboot/internal/api/middleware"
	"scootin-aboot/internal/repository"
	"scootin-aboot/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestScooterHandler_GetScooters(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("successful request with valid parameters", func(t *testing.T) {
		// Arrange
		mockScooterService, mockTripService := createMockServices()
		handler := createScooterHandler(mockScooterService, mockTripService)

		expectedResult := createValidScooterListResult()
		mockScooterService.On("GetScooters", mock.Anything, mock.AnythingOfType("services.ScooterQueryParams")).
			Return(expectedResult, nil)

		router := gin.New()
		router.Use(middleware.ErrorHandlerMiddleware())
		router.GET("/test", handler.GetScooters)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/test?status=available&limit=10&offset=0", nil)
		router.ServeHTTP(w, req)

		// Assert
		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, "application/json; charset=utf-8", w.Header().Get("Content-Type"))

		var response ScooterListResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, int64(1), response.Total)
		assert.Equal(t, 10, response.Limit)
		assert.Equal(t, 0, response.Offset)
		assert.Len(t, response.Scooters, 1)
		assert.Equal(t, TestData.ValidScooterID, response.Scooters[0].ID)

		mockScooterService.AssertExpectations(t)
	})

	t.Run("successful request with empty results", func(t *testing.T) {
		// Arrange
		mockScooterService, mockTripService := createMockServices()
		handler := createScooterHandler(mockScooterService, mockTripService)

		expectedResult := &services.ScooterListResult{
			Scooters: []*services.ScooterInfo{},
			Total:    0,
			Limit:    10,
			Offset:   0,
		}
		mockScooterService.On("GetScooters", mock.Anything, mock.AnythingOfType("services.ScooterQueryParams")).
			Return(expectedResult, nil)

		c, w := setupTestContext("GET", "/api/v1/scooters", nil)

		// Act
		handler.GetScooters(c)

		// Assert
		assert.Equal(t, http.StatusOK, w.Code)

		var response ScooterListResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, int64(0), response.Total)
		assert.Len(t, response.Scooters, 0)

		mockScooterService.AssertExpectations(t)
	})

	t.Run("service returns error", func(t *testing.T) {
		// Arrange
		mockScooterService, mockTripService := createMockServices()
		handler := createScooterHandler(mockScooterService, mockTripService)

		mockScooterService.On("GetScooters", mock.Anything, mock.AnythingOfType("services.ScooterQueryParams")).
			Return(nil, assert.AnError)

		router := gin.New()
		router.Use(middleware.ErrorHandlerMiddleware())
		router.GET("/test", handler.GetScooters)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/test", nil)
		router.ServeHTTP(w, req)

		// Assert
		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assertErrorResponse(t, w, http.StatusInternalServerError, "Internal server error")

		mockScooterService.AssertExpectations(t)
	})

	t.Run("invalid query parameters", func(t *testing.T) {
		// Arrange
		mockScooterService, mockTripService := createMockServices()
		handler := createScooterHandler(mockScooterService, mockTripService)

		router := createTestRouter(handler.GetScooters)
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/test?limit=invalid", nil)
		router.ServeHTTP(w, req)

		// Assert
		assert.Equal(t, http.StatusBadRequest, w.Code)
		assertErrorResponse(t, w, http.StatusBadRequest, "strconv.ParseInt: parsing \"invalid\": invalid syntax")

		mockScooterService.AssertNotCalled(t, "GetScooters")
	})

	t.Run("validates service parameters mapping", func(t *testing.T) {
		// Arrange
		mockScooterService, mockTripService := createMockServices()
		handler := createScooterHandler(mockScooterService, mockTripService)

		expectedResult := createValidScooterListResult()
		mockScooterService.On("GetScooters", mock.Anything, mock.MatchedBy(func(params services.ScooterQueryParams) bool {
			return params.Status == "available" &&
				params.MinLat == 52.0 &&
				params.MaxLat == 53.0 &&
				params.MinLng == 13.0 &&
				params.MaxLng == 14.0 &&
				params.Limit == 20 &&
				params.Offset == 10
		})).Return(expectedResult, nil)

		c, w := setupTestContext("GET", "/api/v1/scooters?status=available&min_lat=52.0&max_lat=53.0&min_lng=13.0&max_lng=14.0&limit=20&offset=10", nil)

		// Act
		handler.GetScooters(c)

		// Assert
		assert.Equal(t, http.StatusOK, w.Code)
		mockScooterService.AssertExpectations(t)
	})
}

func TestScooterHandler_GetScooter(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("successful request with valid scooter ID", func(t *testing.T) {
		// Arrange
		mockScooterService, mockTripService := createMockServices()
		handler := createScooterHandler(mockScooterService, mockTripService)

		expectedResult := createValidScooterDetailsResult()
		mockScooterService.On("GetScooter", mock.Anything, TestData.ValidScooterID).
			Return(expectedResult, nil)

		c, w := setupTestContext("GET", "/api/v1/scooters/"+TestData.ValidScooterID.String(), nil)

		// Act
		handler.GetScooter(c)

		// Assert
		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, "application/json; charset=utf-8", w.Header().Get("Content-Type"))

		var response ScooterDetailsResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, TestData.ValidScooterID, response.ID)
		assert.Equal(t, TestData.ValidStatus, response.Status)
		assert.Nil(t, response.ActiveTrip)

		mockScooterService.AssertExpectations(t)
	})

	t.Run("successful request with scooter having active trip", func(t *testing.T) {
		// Arrange
		mockScooterService, mockTripService := createMockServices()
		handler := createScooterHandler(mockScooterService, mockTripService)

		expectedResult := createValidScooterDetailsResultWithTrip()
		mockScooterService.On("GetScooter", mock.Anything, TestData.ValidScooterID).
			Return(expectedResult, nil)

		c, w := setupTestContext("GET", "/api/v1/scooters/"+TestData.ValidScooterID.String(), nil)

		// Act
		handler.GetScooter(c)

		// Assert
		assert.Equal(t, http.StatusOK, w.Code)

		var response ScooterDetailsResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, TestData.ValidScooterID, response.ID)
		assert.NotNil(t, response.ActiveTrip)
		assert.Equal(t, TestData.ValidTripID, response.ActiveTrip.TripID)
		assert.Equal(t, TestData.ValidUserID, response.ActiveTrip.UserID)

		mockScooterService.AssertExpectations(t)
	})

	t.Run("invalid scooter ID format", func(t *testing.T) {
		// Arrange
		mockScooterService, mockTripService := createMockServices()
		handler := createScooterHandler(mockScooterService, mockTripService)

		router := createTestRouterWithParam(handler.GetScooter)
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/test/"+TestData.InvalidUUID, nil)
		router.ServeHTTP(w, req)

		// Assert
		assert.Equal(t, http.StatusBadRequest, w.Code)
		assertErrorResponse(t, w, http.StatusBadRequest, "Invalid scooter ID")

		mockScooterService.AssertNotCalled(t, "GetScooter")
	})

	t.Run("scooter not found", func(t *testing.T) {
		// Arrange
		mockScooterService, mockTripService := createMockServices()
		handler := createScooterHandler(mockScooterService, mockTripService)

		mockScooterService.On("GetScooter", mock.Anything, TestData.ValidScooterID).
			Return(nil, assert.AnError)

		router := createTestRouterWithParam(handler.GetScooter)
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/test/"+TestData.ValidScooterID.String(), nil)
		router.ServeHTTP(w, req)

		// Assert
		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assertErrorResponse(t, w, http.StatusInternalServerError, "Internal server error")

		mockScooterService.AssertExpectations(t)
	})

	t.Run("scooter not found with specific error message", func(t *testing.T) {
		// Arrange
		mockScooterService, mockTripService := createMockServices()
		handler := createScooterHandler(mockScooterService, mockTripService)

		// Simulate repository.ErrScooterNotFound error
		mockScooterService.On("GetScooter", mock.Anything, TestData.ValidScooterID).
			Return(nil, repository.ErrScooterNotFound)

		router := createTestRouterWithParam(handler.GetScooter)
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/test/"+TestData.ValidScooterID.String(), nil)
		router.ServeHTTP(w, req)

		// Assert
		assert.Equal(t, http.StatusNotFound, w.Code)
		assertErrorResponse(t, w, http.StatusNotFound, "Resource not found")

		mockScooterService.AssertExpectations(t)
	})
}

func TestScooterHandler_GetClosestScooters(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("successful request with valid parameters", func(t *testing.T) {
		// Arrange
		mockScooterService, mockTripService := createMockServices()
		handler := createScooterHandler(mockScooterService, mockTripService)

		expectedResult := createValidClosestScootersResult()
		mockScooterService.On("GetClosestScooters", mock.Anything, mock.AnythingOfType("services.ClosestScootersQueryParams")).
			Return(expectedResult, nil)

		c, w := setupTestContext("GET", "/api/v1/scooters/closest?lat=52.5200&lng=13.4050&radius=1000&limit=10", nil)

		// Act
		handler.GetClosestScooters(c)

		// Assert
		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, "application/json; charset=utf-8", w.Header().Get("Content-Type"))

		var response ClosestScootersResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, TestData.ValidLatitude, response.Center.Latitude)
		assert.Equal(t, TestData.ValidLongitude, response.Center.Longitude)
		assert.Equal(t, TestData.ValidRadius, response.Radius)
		assert.Len(t, response.Scooters, 1)
		assert.Equal(t, TestData.ValidScooterID, response.Scooters[0].ID)
		assert.Equal(t, 500.0, response.Scooters[0].Distance)

		mockScooterService.AssertExpectations(t)
	})

	t.Run("missing required parameters", func(t *testing.T) {
		// Arrange
		mockScooterService, mockTripService := createMockServices()
		handler := createScooterHandler(mockScooterService, mockTripService)

		router := createTestRouter(handler.GetClosestScooters)
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/test", nil)
		router.ServeHTTP(w, req)

		// Assert
		assert.Equal(t, http.StatusBadRequest, w.Code)
		// The validation error includes both missing fields
		assertErrorResponse(t, w, http.StatusBadRequest, "Key: 'ClosestScootersParams.Latitude' Error:Field validation for 'Latitude' failed on the 'required' tag\nKey: 'ClosestScootersParams.Longitude' Error:Field validation for 'Longitude' failed on the 'required' tag")

		mockScooterService.AssertNotCalled(t, "GetClosestScooters")
	})

	t.Run("service returns error", func(t *testing.T) {
		// Arrange
		mockScooterService, mockTripService := createMockServices()
		handler := createScooterHandler(mockScooterService, mockTripService)

		mockScooterService.On("GetClosestScooters", mock.Anything, mock.AnythingOfType("services.ClosestScootersQueryParams")).
			Return(nil, assert.AnError)

		router := createTestRouter(handler.GetClosestScooters)
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/test?lat=52.5200&lng=13.4050", nil)
		router.ServeHTTP(w, req)

		// Assert
		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assertErrorResponse(t, w, http.StatusInternalServerError, "Internal server error")

		mockScooterService.AssertExpectations(t)
	})

	t.Run("validates service parameters mapping", func(t *testing.T) {
		// Arrange
		mockScooterService, mockTripService := createMockServices()
		handler := createScooterHandler(mockScooterService, mockTripService)

		expectedResult := createValidClosestScootersResult()
		mockScooterService.On("GetClosestScooters", mock.Anything, mock.MatchedBy(func(params services.ClosestScootersQueryParams) bool {
			return params.Latitude == 52.5200 &&
				params.Longitude == 13.4050 &&
				params.Radius == 2000.0 &&
				params.Limit == 5 &&
				params.Status == "available"
		})).Return(expectedResult, nil)

		c, w := setupTestContext("GET", "/api/v1/scooters/closest?lat=52.5200&lng=13.4050&radius=2000&limit=5&status=available", nil)

		// Act
		handler.GetClosestScooters(c)

		// Assert
		assert.Equal(t, http.StatusOK, w.Code)
		mockScooterService.AssertExpectations(t)
	})

	t.Run("empty results", func(t *testing.T) {
		// Arrange
		mockScooterService, mockTripService := createMockServices()
		handler := createScooterHandler(mockScooterService, mockTripService)

		expectedResult := &services.ClosestScootersResult{
			Scooters: []*services.ScooterWithDistance{},
			Center: services.Location{
				Latitude:  TestData.ValidLatitude,
				Longitude: TestData.ValidLongitude,
			},
			Radius: TestData.ValidRadius,
		}
		mockScooterService.On("GetClosestScooters", mock.Anything, mock.AnythingOfType("services.ClosestScootersQueryParams")).
			Return(expectedResult, nil)

		c, w := setupTestContext("GET", "/api/v1/scooters/closest?lat=52.5200&lng=13.4050", nil)

		// Act
		handler.GetClosestScooters(c)

		// Assert
		assert.Equal(t, http.StatusOK, w.Code)

		var response ClosestScootersResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Len(t, response.Scooters, 0)

		mockScooterService.AssertExpectations(t)
	})
}
