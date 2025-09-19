package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

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

	err = h.scooterService.UpdateLocation(c.Request.Context(), scooterID, req.Latitude, req.Longitude)
	if err != nil {
		switch err.Error() {
		case "scooter not found":
			c.JSON(http.StatusNotFound, gin.H{"error": "Scooter not found"})
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

	response := LocationUpdateResponse{
		UpdateID:  uuid.New(),
		ScooterID: scooterID,
		Latitude:  req.Latitude,
		Longitude: req.Longitude,
		Timestamp: req.Timestamp,
	}

	c.JSON(http.StatusOK, response)
}
