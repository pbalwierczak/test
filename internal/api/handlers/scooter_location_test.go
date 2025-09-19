package handlers

import (
	"encoding/json"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestScooterHandler_UpdateLocation(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("successful location update", func(t *testing.T) {
		// Arrange
		mockScooterService, mockTripService := createMockServices()
		handler := createScooterHandler(mockScooterService, mockTripService)

		mockScooterService.On("UpdateLocation", mock.Anything, TestData.ValidScooterID, TestData.ValidLatitude, TestData.ValidLongitude).
			Return(nil)

		request := createValidLocationUpdateRequest()
		c, w := setupTestContext("POST", "/api/v1/scooters/"+TestData.ValidScooterID.String()+"/location", toJSON(request))

		// Act
		handler.UpdateLocation(c)

		// Assert
		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, "application/json; charset=utf-8", w.Header().Get("Content-Type"))

		var response LocationUpdateResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, TestData.ValidScooterID, response.ScooterID)
		assert.Equal(t, TestData.ValidLatitude, response.Latitude)
		assert.Equal(t, TestData.ValidLongitude, response.Longitude)
		assert.WithinDuration(t, request.Timestamp, response.Timestamp, time.Second)

		mockScooterService.AssertExpectations(t)
	})

	t.Run("invalid scooter ID format", func(t *testing.T) {
		// Arrange
		mockScooterService, mockTripService := createMockServices()
		handler := createScooterHandler(mockScooterService, mockTripService)

		request := createValidLocationUpdateRequest()
		c, w := setupTestContext("POST", "/api/v1/scooters/"+TestData.InvalidUUID+"/location", toJSON(request))

		// Act
		handler.UpdateLocation(c)

		// Assert
		assert.Equal(t, http.StatusBadRequest, w.Code)
		assertErrorResponse(t, w, http.StatusBadRequest, "Invalid scooter ID")

		mockScooterService.AssertNotCalled(t, "UpdateLocation")
	})

	t.Run("invalid JSON request body", func(t *testing.T) {
		// Arrange
		mockScooterService, mockTripService := createMockServices()
		handler := createScooterHandler(mockScooterService, mockTripService)

		invalidJSON := `{"latitude": 52.5200, "longitude": 13.4050, "timestamp": "invalid-timestamp"}`
		c, w := setupTestContext("POST", "/api/v1/scooters/"+TestData.ValidScooterID.String()+"/location", strings.NewReader(invalidJSON))

		// Act
		handler.UpdateLocation(c)

		// Assert
		assert.Equal(t, http.StatusBadRequest, w.Code)
		assertErrorResponse(t, w, http.StatusBadRequest, "parsing time \"invalid-timestamp\" as \"2006-01-02T15:04:05Z07:00\": cannot parse \"invalid-timestamp\" as \"2006\"")

		mockScooterService.AssertNotCalled(t, "UpdateLocation")
	})

	t.Run("missing required fields", func(t *testing.T) {
		// Arrange
		mockScooterService, mockTripService := createMockServices()
		handler := createScooterHandler(mockScooterService, mockTripService)

		invalidRequest := `{"latitude": 52.5200}`
		c, w := setupTestContext("POST", "/api/v1/scooters/"+TestData.ValidScooterID.String()+"/location", strings.NewReader(invalidRequest))

		// Act
		handler.UpdateLocation(c)

		// Assert
		assert.Equal(t, http.StatusBadRequest, w.Code)
		assertErrorResponse(t, w, http.StatusBadRequest, "Key: 'LocationUpdateRequest.Longitude' Error:Field validation for 'Longitude' failed on the 'required' tag\nKey: 'LocationUpdateRequest.Timestamp' Error:Field validation for 'Timestamp' failed on the 'required' tag")

		mockScooterService.AssertNotCalled(t, "UpdateLocation")
	})

	t.Run("scooter not found", func(t *testing.T) {
		// Arrange
		mockScooterService, mockTripService := createMockServices()
		handler := createScooterHandler(mockScooterService, mockTripService)

		mockScooterService.On("UpdateLocation", mock.Anything, TestData.ValidScooterID, TestData.ValidLatitude, TestData.ValidLongitude).
			Return(assert.AnError)

		request := createValidLocationUpdateRequest()
		c, w := setupTestContext("POST", "/api/v1/scooters/"+TestData.ValidScooterID.String()+"/location", toJSON(request))

		// Act
		handler.UpdateLocation(c)

		// Assert
		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assertErrorResponse(t, w, http.StatusInternalServerError, "Failed to update location")

		mockScooterService.AssertExpectations(t)
	})

	t.Run("scooter not found with specific error message", func(t *testing.T) {
		// Arrange
		mockScooterService, mockTripService := createMockServices()
		handler := createScooterHandler(mockScooterService, mockTripService)

		customError := &mockError{message: "scooter not found"}

		mockScooterService.On("UpdateLocation", mock.Anything, TestData.ValidScooterID, TestData.ValidLatitude, TestData.ValidLongitude).
			Return(customError)

		request := createValidLocationUpdateRequest()
		c, w := setupTestContext("POST", "/api/v1/scooters/"+TestData.ValidScooterID.String()+"/location", toJSON(request))

		// Act
		handler.UpdateLocation(c)

		// Assert
		assert.Equal(t, http.StatusNotFound, w.Code)
		assertErrorResponse(t, w, http.StatusNotFound, "Scooter not found")

		mockScooterService.AssertExpectations(t)
	})

	t.Run("invalid latitude - too high", func(t *testing.T) {
		// Arrange
		mockScooterService, mockTripService := createMockServices()
		handler := createScooterHandler(mockScooterService, mockTripService)

		customError := &mockError{message: "invalid coordinates: invalid latitude: must be between -90 and 90"}

		mockScooterService.On("UpdateLocation", mock.Anything, TestData.ValidScooterID, TestData.InvalidLatitude, TestData.ValidLongitude).
			Return(customError)

		request := LocationUpdateRequest{
			Latitude:  TestData.InvalidLatitude,
			Longitude: TestData.ValidLongitude,
			Timestamp: createValidLocationUpdateRequest().Timestamp,
		}
		c, w := setupTestContext("POST", "/api/v1/scooters/"+TestData.ValidScooterID.String()+"/location", toJSON(request))

		// Act
		handler.UpdateLocation(c)

		// Assert
		assert.Equal(t, http.StatusBadRequest, w.Code)
		assertErrorResponse(t, w, http.StatusBadRequest, "invalid coordinates: invalid latitude: must be between -90 and 90")

		mockScooterService.AssertExpectations(t)
	})

	t.Run("invalid longitude - too high", func(t *testing.T) {
		// Arrange
		mockScooterService, mockTripService := createMockServices()
		handler := createScooterHandler(mockScooterService, mockTripService)

		customError := &mockError{message: "invalid coordinates: invalid longitude: must be between -180 and 180"}

		mockScooterService.On("UpdateLocation", mock.Anything, TestData.ValidScooterID, TestData.ValidLatitude, TestData.InvalidLongitude).
			Return(customError)

		request := LocationUpdateRequest{
			Latitude:  TestData.ValidLatitude,
			Longitude: TestData.InvalidLongitude,
			Timestamp: createValidLocationUpdateRequest().Timestamp,
		}
		c, w := setupTestContext("POST", "/api/v1/scooters/"+TestData.ValidScooterID.String()+"/location", toJSON(request))

		// Act
		handler.UpdateLocation(c)

		// Assert
		assert.Equal(t, http.StatusBadRequest, w.Code)
		assertErrorResponse(t, w, http.StatusBadRequest, "invalid coordinates: invalid longitude: must be between -180 and 180")

		mockScooterService.AssertExpectations(t)
	})

	t.Run("boundary values - valid coordinates", func(t *testing.T) {
		// Arrange
		mockScooterService, mockTripService := createMockServices()
		handler := createScooterHandler(mockScooterService, mockTripService)

		boundaryLat := 90.0
		boundaryLng := 180.0

		mockScooterService.On("UpdateLocation", mock.Anything, TestData.ValidScooterID, boundaryLat, boundaryLng).
			Return(nil)

		request := LocationUpdateRequest{
			Latitude:  boundaryLat,
			Longitude: boundaryLng,
			Timestamp: createValidLocationUpdateRequest().Timestamp,
		}
		c, w := setupTestContext("POST", "/api/v1/scooters/"+TestData.ValidScooterID.String()+"/location", toJSON(request))

		// Act
		handler.UpdateLocation(c)

		// Assert
		assert.Equal(t, http.StatusOK, w.Code)

		mockScooterService.AssertExpectations(t)
	})

	t.Run("negative boundary values - valid coordinates", func(t *testing.T) {
		// Arrange
		mockScooterService, mockTripService := createMockServices()
		handler := createScooterHandler(mockScooterService, mockTripService)

		boundaryLat := -90.0
		boundaryLng := -180.0

		mockScooterService.On("UpdateLocation", mock.Anything, TestData.ValidScooterID, boundaryLat, boundaryLng).
			Return(nil)

		request := LocationUpdateRequest{
			Latitude:  boundaryLat,
			Longitude: boundaryLng,
			Timestamp: createValidLocationUpdateRequest().Timestamp,
		}
		c, w := setupTestContext("POST", "/api/v1/scooters/"+TestData.ValidScooterID.String()+"/location", toJSON(request))

		// Act
		handler.UpdateLocation(c)

		// Assert
		assert.Equal(t, http.StatusOK, w.Code)

		mockScooterService.AssertExpectations(t)
	})
}
