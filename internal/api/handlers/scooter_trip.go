package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// StartTrip starts a new trip for a scooter
// POST /api/v1/scooters/{id}/trip/start
func (h *ScooterHandler) StartTrip(c *gin.Context) {
	scooterIDStr := c.Param("id")
	scooterID, err := uuid.Parse(scooterIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid scooter ID"})
		return
	}

	var req StartTripRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// TODO: Implement scooter service logic
	// 1. Check if scooter exists
	// 2. Check if scooter is available
	// 3. Check if user has an active trip
	// 4. Create new trip
	// 5. Update scooter status to occupied

	// Dummy response for now
	response := StartTripResponse{
		TripID:         uuid.New(),
		ScooterID:      scooterID,
		UserID:         req.UserID,
		StartTime:      time.Now(),
		StartLatitude:  req.StartLatitude,
		StartLongitude: req.StartLongitude,
		Status:         "active",
	}

	c.JSON(http.StatusCreated, response)
}

// EndTrip ends an active trip for a scooter
// POST /api/v1/scooters/{id}/trip/end
func (h *ScooterHandler) EndTrip(c *gin.Context) {
	scooterIDStr := c.Param("id")
	scooterID, err := uuid.Parse(scooterIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid scooter ID"})
		return
	}

	var req EndTripRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// TODO: Implement scooter service logic
	// 1. Check if scooter exists
	// 2. Check if scooter has an active trip
	// 3. End the active trip
	// 4. Update scooter status to available
	// 5. Calculate trip duration

	// Dummy response for now
	startTime := time.Now().Add(-10 * time.Minute) // Simulate 10-minute trip
	endTime := time.Now()
	response := EndTripResponse{
		TripID:         uuid.New(),
		ScooterID:      scooterID,
		UserID:         uuid.New(), // TODO: Get from active trip
		StartTime:      startTime,
		EndTime:        endTime,
		StartLatitude:  45.4215, // TODO: Get from active trip
		StartLongitude: -75.6972,
		EndLatitude:    req.EndLatitude,
		EndLongitude:   req.EndLongitude,
		Status:         "completed",
		Duration:       int64(endTime.Sub(startTime).Seconds()),
	}

	c.JSON(http.StatusOK, response)
}
