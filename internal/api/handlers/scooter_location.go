package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// UpdateLocation updates the location of a scooter during a trip
// POST /api/v1/scooters/{id}/location
func (h *ScooterHandler) UpdateLocation(c *gin.Context) {
	scooterIDStr := c.Param("id")
	scooterID, err := uuid.Parse(scooterIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid scooter ID"})
		return
	}

	var req LocationUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Use trip service to update location
	err = h.tripService.UpdateLocation(c.Request.Context(), scooterID, req.Latitude, req.Longitude)
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
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update location"})
			}
		}
		return
	}

	// Get the active trip to return trip ID in response
	trip, err := h.tripService.GetActiveTrip(c.Request.Context(), scooterID)
	if err != nil {
		// If we can't get the trip, still return success for location update
		// but without trip ID
		response := LocationUpdateResponse{
			UpdateID:  uuid.New(),
			Latitude:  req.Latitude,
			Longitude: req.Longitude,
			Timestamp: req.Timestamp,
		}
		c.JSON(http.StatusOK, response)
		return
	}

	// Convert to response
	response := LocationUpdateResponse{
		UpdateID:  uuid.New(),
		TripID:    trip.ID,
		Latitude:  req.Latitude,
		Longitude: req.Longitude,
		Timestamp: req.Timestamp,
	}

	c.JSON(http.StatusOK, response)
}
