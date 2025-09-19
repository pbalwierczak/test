package handlers

import (
	"encoding/json"
	"net/http"
	"strings"
	"testing"
	"time"

	"scootin-aboot/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestScooterHandler_StartTrip(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("successful trip start", func(t *testing.T) {
		// Arrange
		mockScooterService, mockTripService := createMockServices()
		handler := createScooterHandler(mockScooterService, mockTripService)

		expectedTrip := createValidTrip()
		mockTripService.On("StartTrip", mock.Anything, TestData.ValidScooterID, TestData.ValidUserID, TestData.ValidLatitude, TestData.ValidLongitude).
			Return(expectedTrip, nil)

		request := createValidStartTripRequest()
		c, w := setupTestContext("POST", "/api/v1/scooters/"+TestData.ValidScooterID.String()+"/trip/start", toJSON(request))

		// Act
		handler.StartTrip(c)

		// Assert
		assert.Equal(t, http.StatusCreated, w.Code)
		assert.Equal(t, "application/json; charset=utf-8", w.Header().Get("Content-Type"))

		var response StartTripResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, TestData.ValidTripID, response.TripID)
		assert.Equal(t, TestData.ValidScooterID, response.ScooterID)
		assert.Equal(t, TestData.ValidUserID, response.UserID)
		assert.Equal(t, TestData.ValidLatitude, response.StartLatitude)
		assert.Equal(t, TestData.ValidLongitude, response.StartLongitude)
		assert.Equal(t, string(models.TripStatusCompleted), response.Status)

		mockTripService.AssertExpectations(t)
	})

	t.Run("invalid scooter ID format", func(t *testing.T) {
		// Arrange
		mockScooterService, mockTripService := createMockServices()
		handler := createScooterHandler(mockScooterService, mockTripService)

		request := createValidStartTripRequest()
		c, w := setupTestContext("POST", "/api/v1/scooters/"+TestData.InvalidUUID+"/trip/start", toJSON(request))

		// Act
		handler.StartTrip(c)

		// Assert
		assert.Equal(t, http.StatusBadRequest, w.Code)
		assertErrorResponse(t, w, http.StatusBadRequest, "Invalid scooter ID")

		mockTripService.AssertNotCalled(t, "StartTrip")
	})

	t.Run("invalid JSON request body", func(t *testing.T) {
		// Arrange
		mockScooterService, mockTripService := createMockServices()
		handler := createScooterHandler(mockScooterService, mockTripService)

		invalidJSON := `{"user_id": "invalid-uuid", "start_latitude": 52.5200, "start_longitude": 13.4050}`
		c, w := setupTestContext("POST", "/api/v1/scooters/"+TestData.ValidScooterID.String()+"/trip/start", strings.NewReader(invalidJSON))

		// Act
		handler.StartTrip(c)

		// Assert
		assert.Equal(t, http.StatusBadRequest, w.Code)
		assertErrorResponse(t, w, http.StatusBadRequest, "invalid UUID length: 12")

		mockTripService.AssertNotCalled(t, "StartTrip")
	})

	t.Run("missing required fields", func(t *testing.T) {
		// Arrange
		mockScooterService, mockTripService := createMockServices()
		handler := createScooterHandler(mockScooterService, mockTripService)

		invalidRequest := `{"user_id": "` + TestData.ValidUserID.String() + `"}`
		c, w := setupTestContext("POST", "/api/v1/scooters/"+TestData.ValidScooterID.String()+"/trip/start", strings.NewReader(invalidRequest))

		// Act
		handler.StartTrip(c)

		// Assert
		assert.Equal(t, http.StatusBadRequest, w.Code)
		assertErrorResponse(t, w, http.StatusBadRequest, "Key: 'StartTripRequest.StartLatitude' Error:Field validation for 'StartLatitude' failed on the 'required' tag\nKey: 'StartTripRequest.StartLongitude' Error:Field validation for 'StartLongitude' failed on the 'required' tag")

		mockTripService.AssertNotCalled(t, "StartTrip")
	})

	t.Run("scooter not found", func(t *testing.T) {
		// Arrange
		mockScooterService, mockTripService := createMockServices()
		handler := createScooterHandler(mockScooterService, mockTripService)

		// Create a custom error that matches the expected error message
		customError := &mockError{message: "scooter not found"}

		mockTripService.On("StartTrip", mock.Anything, TestData.ValidScooterID, TestData.ValidUserID, TestData.ValidLatitude, TestData.ValidLongitude).
			Return(nil, customError)

		request := createValidStartTripRequest()
		c, w := setupTestContext("POST", "/api/v1/scooters/"+TestData.ValidScooterID.String()+"/trip/start", toJSON(request))

		// Act
		handler.StartTrip(c)

		// Assert
		assert.Equal(t, http.StatusNotFound, w.Code)
		assertErrorResponse(t, w, http.StatusNotFound, "Scooter not found")

		mockTripService.AssertExpectations(t)
	})

	t.Run("user not found", func(t *testing.T) {
		// Arrange
		mockScooterService, mockTripService := createMockServices()
		handler := createScooterHandler(mockScooterService, mockTripService)

		// Create a custom error that matches the expected error message
		customError := &mockError{message: "user not found"}

		mockTripService.On("StartTrip", mock.Anything, TestData.ValidScooterID, TestData.ValidUserID, TestData.ValidLatitude, TestData.ValidLongitude).
			Return(nil, customError)

		request := createValidStartTripRequest()
		c, w := setupTestContext("POST", "/api/v1/scooters/"+TestData.ValidScooterID.String()+"/trip/start", toJSON(request))

		// Act
		handler.StartTrip(c)

		// Assert
		assert.Equal(t, http.StatusNotFound, w.Code)
		assertErrorResponse(t, w, http.StatusNotFound, "User not found")

		mockTripService.AssertExpectations(t)
	})

	t.Run("scooter is not available", func(t *testing.T) {
		// Arrange
		mockScooterService, mockTripService := createMockServices()
		handler := createScooterHandler(mockScooterService, mockTripService)

		// Create a custom error that matches the expected error message
		customError := &mockError{message: "scooter is not available"}

		mockTripService.On("StartTrip", mock.Anything, TestData.ValidScooterID, TestData.ValidUserID, TestData.ValidLatitude, TestData.ValidLongitude).
			Return(nil, customError)

		request := createValidStartTripRequest()
		c, w := setupTestContext("POST", "/api/v1/scooters/"+TestData.ValidScooterID.String()+"/trip/start", toJSON(request))

		// Act
		handler.StartTrip(c)

		// Assert
		assert.Equal(t, http.StatusConflict, w.Code)
		assertErrorResponse(t, w, http.StatusConflict, "Scooter is not available")

		mockTripService.AssertExpectations(t)
	})

	t.Run("user already has an active trip", func(t *testing.T) {
		// Arrange
		mockScooterService, mockTripService := createMockServices()
		handler := createScooterHandler(mockScooterService, mockTripService)

		// Create a custom error that matches the expected error message
		customError := &mockError{message: "user already has an active trip"}

		mockTripService.On("StartTrip", mock.Anything, TestData.ValidScooterID, TestData.ValidUserID, TestData.ValidLatitude, TestData.ValidLongitude).
			Return(nil, customError)

		request := createValidStartTripRequest()
		c, w := setupTestContext("POST", "/api/v1/scooters/"+TestData.ValidScooterID.String()+"/trip/start", toJSON(request))

		// Act
		handler.StartTrip(c)

		// Assert
		assert.Equal(t, http.StatusConflict, w.Code)
		assertErrorResponse(t, w, http.StatusConflict, "User already has an active trip")

		mockTripService.AssertExpectations(t)
	})

	t.Run("scooter already has an active trip", func(t *testing.T) {
		// Arrange
		mockScooterService, mockTripService := createMockServices()
		handler := createScooterHandler(mockScooterService, mockTripService)

		// Create a custom error that matches the expected error message
		customError := &mockError{message: "scooter already has an active trip"}

		mockTripService.On("StartTrip", mock.Anything, TestData.ValidScooterID, TestData.ValidUserID, TestData.ValidLatitude, TestData.ValidLongitude).
			Return(nil, customError)

		request := createValidStartTripRequest()
		c, w := setupTestContext("POST", "/api/v1/scooters/"+TestData.ValidScooterID.String()+"/trip/start", toJSON(request))

		// Act
		handler.StartTrip(c)

		// Assert
		assert.Equal(t, http.StatusConflict, w.Code)
		assertErrorResponse(t, w, http.StatusConflict, "Scooter already has an active trip")

		mockTripService.AssertExpectations(t)
	})

	t.Run("invalid coordinates - latitude", func(t *testing.T) {
		// Arrange
		mockScooterService, mockTripService := createMockServices()
		handler := createScooterHandler(mockScooterService, mockTripService)

		// Create a custom error that matches the expected error message
		customError := &mockError{message: "invalid coordinates: invalid latitude: must be between -90 and 90"}

		mockTripService.On("StartTrip", mock.Anything, TestData.ValidScooterID, TestData.ValidUserID, TestData.InvalidLatitude, TestData.ValidLongitude).
			Return(nil, customError)

		request := StartTripRequest{
			UserID:         TestData.ValidUserID,
			StartLatitude:  TestData.InvalidLatitude,
			StartLongitude: TestData.ValidLongitude,
		}
		c, w := setupTestContext("POST", "/api/v1/scooters/"+TestData.ValidScooterID.String()+"/trip/start", toJSON(request))

		// Act
		handler.StartTrip(c)

		// Assert
		assert.Equal(t, http.StatusBadRequest, w.Code)
		assertErrorResponse(t, w, http.StatusBadRequest, "invalid coordinates: invalid latitude: must be between -90 and 90")

		mockTripService.AssertExpectations(t)
	})

	t.Run("invalid coordinates - longitude", func(t *testing.T) {
		// Arrange
		mockScooterService, mockTripService := createMockServices()
		handler := createScooterHandler(mockScooterService, mockTripService)

		// Create a custom error that matches the expected error message
		customError := &mockError{message: "invalid coordinates: invalid longitude: must be between -180 and 180"}

		mockTripService.On("StartTrip", mock.Anything, TestData.ValidScooterID, TestData.ValidUserID, TestData.ValidLatitude, TestData.InvalidLongitude).
			Return(nil, customError)

		request := StartTripRequest{
			UserID:         TestData.ValidUserID,
			StartLatitude:  TestData.ValidLatitude,
			StartLongitude: TestData.InvalidLongitude,
		}
		c, w := setupTestContext("POST", "/api/v1/scooters/"+TestData.ValidScooterID.String()+"/trip/start", toJSON(request))

		// Act
		handler.StartTrip(c)

		// Assert
		assert.Equal(t, http.StatusBadRequest, w.Code)
		assertErrorResponse(t, w, http.StatusBadRequest, "invalid coordinates: invalid longitude: must be between -180 and 180")

		mockTripService.AssertExpectations(t)
	})

	t.Run("internal server error", func(t *testing.T) {
		// Arrange
		mockScooterService, mockTripService := createMockServices()
		handler := createScooterHandler(mockScooterService, mockTripService)

		mockTripService.On("StartTrip", mock.Anything, TestData.ValidScooterID, TestData.ValidUserID, TestData.ValidLatitude, TestData.ValidLongitude).
			Return(nil, assert.AnError)

		request := createValidStartTripRequest()
		c, w := setupTestContext("POST", "/api/v1/scooters/"+TestData.ValidScooterID.String()+"/trip/start", toJSON(request))

		// Act
		handler.StartTrip(c)

		// Assert
		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assertErrorResponse(t, w, http.StatusInternalServerError, "Failed to start trip")

		mockTripService.AssertExpectations(t)
	})
}

func TestScooterHandler_EndTrip(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("successful trip end", func(t *testing.T) {
		// Arrange
		mockScooterService, mockTripService := createMockServices()
		handler := createScooterHandler(mockScooterService, mockTripService)

		expectedTrip := createValidTrip()
		mockTripService.On("EndTrip", mock.Anything, TestData.ValidScooterID, TestData.ValidLatitude, TestData.ValidLongitude).
			Return(expectedTrip, nil)

		request := createValidEndTripRequest()
		c, w := setupTestContext("POST", "/api/v1/scooters/"+TestData.ValidScooterID.String()+"/trip/end", toJSON(request))

		// Act
		handler.EndTrip(c)

		// Assert
		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, "application/json; charset=utf-8", w.Header().Get("Content-Type"))

		var response EndTripResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, TestData.ValidTripID, response.TripID)
		assert.Equal(t, TestData.ValidScooterID, response.ScooterID)
		assert.Equal(t, TestData.ValidUserID, response.UserID)
		assert.Equal(t, TestData.ValidLatitude, response.StartLatitude)
		assert.Equal(t, TestData.ValidLongitude, response.StartLongitude)
		assert.Equal(t, TestData.ValidLatitude, response.EndLatitude)
		assert.Equal(t, TestData.ValidLongitude, response.EndLongitude)
		assert.Equal(t, string(models.TripStatusCompleted), response.Status)
		assert.Greater(t, response.Duration, int64(0))

		mockTripService.AssertExpectations(t)
	})

	t.Run("invalid scooter ID format", func(t *testing.T) {
		// Arrange
		mockScooterService, mockTripService := createMockServices()
		handler := createScooterHandler(mockScooterService, mockTripService)

		request := createValidEndTripRequest()
		c, w := setupTestContext("POST", "/api/v1/scooters/"+TestData.InvalidUUID+"/trip/end", toJSON(request))

		// Act
		handler.EndTrip(c)

		// Assert
		assert.Equal(t, http.StatusBadRequest, w.Code)
		assertErrorResponse(t, w, http.StatusBadRequest, "Invalid scooter ID")

		mockTripService.AssertNotCalled(t, "EndTrip")
	})

	t.Run("invalid JSON request body", func(t *testing.T) {
		// Arrange
		mockScooterService, mockTripService := createMockServices()
		handler := createScooterHandler(mockScooterService, mockTripService)

		invalidJSON := `{"end_latitude": "invalid", "end_longitude": 13.4050}`
		c, w := setupTestContext("POST", "/api/v1/scooters/"+TestData.ValidScooterID.String()+"/trip/end", strings.NewReader(invalidJSON))

		// Act
		handler.EndTrip(c)

		// Assert
		assert.Equal(t, http.StatusBadRequest, w.Code)
		assertErrorResponse(t, w, http.StatusBadRequest, "json: cannot unmarshal string into Go struct field EndTripRequest.end_latitude of type float64")

		mockTripService.AssertNotCalled(t, "EndTrip")
	})

	t.Run("missing required fields", func(t *testing.T) {
		// Arrange
		mockScooterService, mockTripService := createMockServices()
		handler := createScooterHandler(mockScooterService, mockTripService)

		invalidRequest := `{"end_latitude": 52.5200}`
		c, w := setupTestContext("POST", "/api/v1/scooters/"+TestData.ValidScooterID.String()+"/trip/end", strings.NewReader(invalidRequest))

		// Act
		handler.EndTrip(c)

		// Assert
		assert.Equal(t, http.StatusBadRequest, w.Code)
		assertErrorResponse(t, w, http.StatusBadRequest, "Key: 'EndTripRequest.EndLongitude' Error:Field validation for 'EndLongitude' failed on the 'required' tag")

		mockTripService.AssertNotCalled(t, "EndTrip")
	})

	t.Run("no active trip found for scooter", func(t *testing.T) {
		// Arrange
		mockScooterService, mockTripService := createMockServices()
		handler := createScooterHandler(mockScooterService, mockTripService)

		// Create a custom error that matches the expected error message
		customError := &mockError{message: "no active trip found for scooter"}

		mockTripService.On("EndTrip", mock.Anything, TestData.ValidScooterID, TestData.ValidLatitude, TestData.ValidLongitude).
			Return(nil, customError)

		request := createValidEndTripRequest()
		c, w := setupTestContext("POST", "/api/v1/scooters/"+TestData.ValidScooterID.String()+"/trip/end", toJSON(request))

		// Act
		handler.EndTrip(c)

		// Assert
		assert.Equal(t, http.StatusNotFound, w.Code)
		assertErrorResponse(t, w, http.StatusNotFound, "No active trip found for scooter")

		mockTripService.AssertExpectations(t)
	})

	t.Run("invalid coordinates - latitude", func(t *testing.T) {
		// Arrange
		mockScooterService, mockTripService := createMockServices()
		handler := createScooterHandler(mockScooterService, mockTripService)

		// Create a custom error that matches the expected error message
		customError := &mockError{message: "invalid coordinates: invalid latitude: must be between -90 and 90"}

		mockTripService.On("EndTrip", mock.Anything, TestData.ValidScooterID, TestData.InvalidLatitude, TestData.ValidLongitude).
			Return(nil, customError)

		request := EndTripRequest{
			EndLatitude:  TestData.InvalidLatitude,
			EndLongitude: TestData.ValidLongitude,
		}
		c, w := setupTestContext("POST", "/api/v1/scooters/"+TestData.ValidScooterID.String()+"/trip/end", toJSON(request))

		// Act
		handler.EndTrip(c)

		// Assert
		assert.Equal(t, http.StatusBadRequest, w.Code)
		assertErrorResponse(t, w, http.StatusBadRequest, "invalid coordinates: invalid latitude: must be between -90 and 90")

		mockTripService.AssertExpectations(t)
	})

	t.Run("invalid coordinates - longitude", func(t *testing.T) {
		// Arrange
		mockScooterService, mockTripService := createMockServices()
		handler := createScooterHandler(mockScooterService, mockTripService)

		// Create a custom error that matches the expected error message
		customError := &mockError{message: "invalid coordinates: invalid longitude: must be between -180 and 180"}

		mockTripService.On("EndTrip", mock.Anything, TestData.ValidScooterID, TestData.ValidLatitude, TestData.InvalidLongitude).
			Return(nil, customError)

		request := EndTripRequest{
			EndLatitude:  TestData.ValidLatitude,
			EndLongitude: TestData.InvalidLongitude,
		}
		c, w := setupTestContext("POST", "/api/v1/scooters/"+TestData.ValidScooterID.String()+"/trip/end", toJSON(request))

		// Act
		handler.EndTrip(c)

		// Assert
		assert.Equal(t, http.StatusBadRequest, w.Code)
		assertErrorResponse(t, w, http.StatusBadRequest, "invalid coordinates: invalid longitude: must be between -180 and 180")

		mockTripService.AssertExpectations(t)
	})

	t.Run("internal server error", func(t *testing.T) {
		// Arrange
		mockScooterService, mockTripService := createMockServices()
		handler := createScooterHandler(mockScooterService, mockTripService)

		mockTripService.On("EndTrip", mock.Anything, TestData.ValidScooterID, TestData.ValidLatitude, TestData.ValidLongitude).
			Return(nil, assert.AnError)

		request := createValidEndTripRequest()
		c, w := setupTestContext("POST", "/api/v1/scooters/"+TestData.ValidScooterID.String()+"/trip/end", toJSON(request))

		// Act
		handler.EndTrip(c)

		// Assert
		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assertErrorResponse(t, w, http.StatusInternalServerError, "Failed to end trip")

		mockTripService.AssertExpectations(t)
	})

	t.Run("duration calculation with different times", func(t *testing.T) {
		// Arrange
		mockScooterService, mockTripService := createMockServices()
		handler := createScooterHandler(mockScooterService, mockTripService)

		// Create a trip with specific start and end times
		startTime := time.Now().Add(-30 * time.Minute)
		endTime := time.Now()
		expectedTrip := &models.Trip{
			ID:             TestData.ValidTripID,
			ScooterID:      TestData.ValidScooterID,
			UserID:         TestData.ValidUserID,
			StartTime:      startTime,
			EndTime:        &endTime,
			StartLatitude:  TestData.ValidLatitude,
			StartLongitude: TestData.ValidLongitude,
			EndLatitude:    &TestData.ValidLatitude,
			EndLongitude:   &TestData.ValidLongitude,
			Status:         models.TripStatusCompleted,
		}

		mockTripService.On("EndTrip", mock.Anything, TestData.ValidScooterID, TestData.ValidLatitude, TestData.ValidLongitude).
			Return(expectedTrip, nil)

		request := createValidEndTripRequest()
		c, w := setupTestContext("POST", "/api/v1/scooters/"+TestData.ValidScooterID.String()+"/trip/end", toJSON(request))

		// Act
		handler.EndTrip(c)

		// Assert
		assert.Equal(t, http.StatusOK, w.Code)

		var response EndTripResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, int64(1800), response.Duration) // 30 minutes in seconds

		mockTripService.AssertExpectations(t)
	})
}
