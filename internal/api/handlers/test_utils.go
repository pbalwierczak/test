package handlers

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"scootin-aboot/internal/api/handlers/mocks"
	"scootin-aboot/internal/models"
	"scootin-aboot/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

var TestData = struct {
	ValidScooterID   uuid.UUID
	ValidUserID      uuid.UUID
	ValidTripID      uuid.UUID
	ValidCoordinates struct{ Lat, Lng float64 }
	InvalidUUID      string
	InvalidLatitude  float64
	InvalidLongitude float64
	ValidLatitude    float64
	ValidLongitude   float64
	ValidRadius      float64
	ValidStatus      string
	ValidLimit       int
	ValidOffset      int
}{
	ValidScooterID:   uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"),
	ValidUserID:      uuid.MustParse("123e4567-e89b-12d3-a456-426614174001"),
	ValidTripID:      uuid.MustParse("123e4567-e89b-12d3-a456-426614174002"),
	ValidCoordinates: struct{ Lat, Lng float64 }{52.5200, 13.4050}, // Berlin
	InvalidUUID:      "invalid-uuid",
	InvalidLatitude:  91.0,
	InvalidLongitude: 181.0,
	ValidLatitude:    52.5200,
	ValidLongitude:   13.4050,
	ValidRadius:      1000.0,
	ValidStatus:      "available",
	ValidLimit:       10,
	ValidOffset:      0,
}

func setupTestContext(method, url string, body io.Reader) (*gin.Context, *httptest.ResponseRecorder) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	req, _ := http.NewRequest(method, url, body)
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	c.Request = req

	if id := extractIDFromURL(url); id != "" {
		c.Params = gin.Params{
			{Key: "id", Value: id},
		}
	}

	return c, w
}

func extractIDFromURL(url string) string {
	parts := strings.Split(url, "/")
	for i, part := range parts {
		if part == "scooters" && i+1 < len(parts) {
			return parts[i+1]
		}
	}
	return ""
}

func createMockServices() (*mocks.MockScooterService, *mocks.MockTripService) {
	return &mocks.MockScooterService{}, &mocks.MockTripService{}
}

func createScooterHandler(mockScooterService *mocks.MockScooterService, mockTripService *mocks.MockTripService) *ScooterHandler {
	return NewScooterHandler(mockTripService, mockScooterService)
}

func assertJSONResponse(t *testing.T, expected interface{}, actual string) {
	expectedJSON, err := json.Marshal(expected)
	assert.NoError(t, err)
	assert.JSONEq(t, string(expectedJSON), actual)
}

func assertErrorResponse(t *testing.T, w *httptest.ResponseRecorder, expectedStatus int, expectedMessage string) {
	assert.Equal(t, expectedStatus, w.Code)

	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, expectedMessage, response["error"])
}

func createValidScooterListResult() *services.ScooterListResult {
	return &services.ScooterListResult{
		Scooters: []*services.ScooterInfo{
			{
				ID:               TestData.ValidScooterID,
				Status:           TestData.ValidStatus,
				CurrentLatitude:  TestData.ValidLatitude,
				CurrentLongitude: TestData.ValidLongitude,
				LastSeen:         time.Now(),
				CreatedAt:        time.Now().Add(-time.Hour),
			},
		},
		Total:  1,
		Limit:  TestData.ValidLimit,
		Offset: TestData.ValidOffset,
	}
}

func createValidScooterDetailsResult() *services.ScooterDetailsResult {
	return &services.ScooterDetailsResult{
		ID:               TestData.ValidScooterID,
		Status:           TestData.ValidStatus,
		CurrentLatitude:  TestData.ValidLatitude,
		CurrentLongitude: TestData.ValidLongitude,
		LastSeen:         time.Now(),
		CreatedAt:        time.Now().Add(-time.Hour),
		UpdatedAt:        time.Now(),
		ActiveTrip:       nil,
	}
}

func createValidScooterDetailsResultWithTrip() *services.ScooterDetailsResult {
	now := time.Now()
	return &services.ScooterDetailsResult{
		ID:               TestData.ValidScooterID,
		Status:           "occupied",
		CurrentLatitude:  TestData.ValidLatitude,
		CurrentLongitude: TestData.ValidLongitude,
		LastSeen:         now,
		CreatedAt:        now.Add(-time.Hour),
		UpdatedAt:        now,
		ActiveTrip: &services.TripInfo{
			TripID:         TestData.ValidTripID,
			UserID:         TestData.ValidUserID,
			StartTime:      now.Add(-30 * time.Minute),
			StartLatitude:  TestData.ValidLatitude,
			StartLongitude: TestData.ValidLongitude,
		},
	}
}

func createValidClosestScootersResult() *services.ClosestScootersResult {
	return &services.ClosestScootersResult{
		Scooters: []*services.ScooterWithDistance{
			{
				ScooterInfo: &services.ScooterInfo{
					ID:               TestData.ValidScooterID,
					Status:           TestData.ValidStatus,
					CurrentLatitude:  TestData.ValidLatitude,
					CurrentLongitude: TestData.ValidLongitude,
					LastSeen:         time.Now(),
					CreatedAt:        time.Now().Add(-time.Hour),
				},
				Distance: 500.0,
			},
		},
		Center: services.Location{
			Latitude:  TestData.ValidLatitude,
			Longitude: TestData.ValidLongitude,
		},
		Radius: TestData.ValidRadius,
	}
}

func createValidTrip() *models.Trip {
	now := time.Now()
	return &models.Trip{
		ID:             TestData.ValidTripID,
		ScooterID:      TestData.ValidScooterID,
		UserID:         TestData.ValidUserID,
		StartTime:      now.Add(-30 * time.Minute),
		EndTime:        &now,
		StartLatitude:  TestData.ValidLatitude,
		StartLongitude: TestData.ValidLongitude,
		EndLatitude:    &TestData.ValidLatitude,
		EndLongitude:   &TestData.ValidLongitude,
		Status:         models.TripStatusCompleted,
	}
}

func toJSON(obj interface{}) *bytes.Buffer {
	jsonData, _ := json.Marshal(obj)
	return bytes.NewBuffer(jsonData)
}

type mockError struct {
	message string
}

func (e *mockError) Error() string {
	return e.message
}
