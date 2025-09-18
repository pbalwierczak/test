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
	_, err := uuid.Parse(scooterIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid scooter ID"})
		return
	}

	var req LocationUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// TODO: Implement scooter service logic
	// 1. Check if scooter exists
	// 2. Check if scooter has an active trip
	// 3. Create location update record
	// 4. Update scooter's current location
	// 5. Update scooter's last_seen timestamp

	// Dummy response for now
	response := LocationUpdateResponse{
		UpdateID:  uuid.New(),
		TripID:    uuid.New(), // TODO: Get from active trip
		Latitude:  req.Latitude,
		Longitude: req.Longitude,
		Timestamp: req.Timestamp,
	}

	c.JSON(http.StatusOK, response)
}
