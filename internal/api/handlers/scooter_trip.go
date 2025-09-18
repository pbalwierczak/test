package handlers

import (
	"net/http"

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

	// Use trip service to start the trip
	trip, err := h.tripService.StartTrip(c.Request.Context(), scooterID, req.UserID, req.StartLatitude, req.StartLongitude)
	if err != nil {
		// Map service errors to appropriate HTTP status codes
		switch err.Error() {
		case "scooter not found":
			c.JSON(http.StatusNotFound, gin.H{"error": "Scooter not found"})
		case "user not found":
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		case "scooter is not available":
			c.JSON(http.StatusConflict, gin.H{"error": "Scooter is not available"})
		case "user already has an active trip":
			c.JSON(http.StatusConflict, gin.H{"error": "User already has an active trip"})
		case "scooter already has an active trip":
			c.JSON(http.StatusConflict, gin.H{"error": "Scooter already has an active trip"})
		default:
			if err.Error() == "invalid coordinates: invalid latitude: must be between -90 and 90" ||
				err.Error() == "invalid coordinates: invalid longitude: must be between -180 and 180" {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to start trip"})
			}
		}
		return
	}

	// Convert trip to response
	response := StartTripResponse{
		TripID:         trip.ID,
		ScooterID:      trip.ScooterID,
		UserID:         trip.UserID,
		StartTime:      trip.StartTime,
		StartLatitude:  trip.StartLatitude,
		StartLongitude: trip.StartLongitude,
		Status:         string(trip.Status),
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

	// Use trip service to end the trip
	trip, err := h.tripService.EndTrip(c.Request.Context(), scooterID, req.EndLatitude, req.EndLongitude)
	if err != nil {
		// Map service errors to appropriate HTTP status codes
		switch err.Error() {
		case "no active trip found for scooter":
			c.JSON(http.StatusNotFound, gin.H{"error": "No active trip found for scooter"})
		default:
			if err.Error() == "invalid coordinates: invalid latitude: must be between -90 and 90" ||
				err.Error() == "invalid coordinates: invalid longitude: must be between -180 and 180" {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to end trip"})
			}
		}
		return
	}

	// Calculate trip duration
	var duration int64
	if trip.EndTime != nil {
		duration = int64(trip.EndTime.Sub(trip.StartTime).Seconds())
	}

	// Convert trip to response
	response := EndTripResponse{
		TripID:         trip.ID,
		ScooterID:      trip.ScooterID,
		UserID:         trip.UserID,
		StartTime:      trip.StartTime,
		EndTime:        *trip.EndTime,
		StartLatitude:  trip.StartLatitude,
		StartLongitude: trip.StartLongitude,
		EndLatitude:    *trip.EndLatitude,
		EndLongitude:   *trip.EndLongitude,
		Status:         string(trip.Status),
		Duration:       duration,
	}

	c.JSON(http.StatusOK, response)
}
