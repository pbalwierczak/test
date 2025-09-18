package handlers

import (
	"net/http"

	"scootin-aboot/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// GetScooters retrieves scooters with optional filtering
// GET /api/v1/scooters
func (h *ScooterHandler) GetScooters(c *gin.Context) {
	var params ScooterQueryParams
	if err := c.ShouldBindQuery(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Convert to service parameters
	serviceParams := services.ScooterQueryParams{
		Status: params.Status,
		MinLat: params.MinLat,
		MaxLat: params.MaxLat,
		MinLng: params.MinLng,
		MaxLng: params.MaxLng,
		Limit:  params.Limit,
		Offset: params.Offset,
	}

	// Call scooter service
	result, err := h.scooterService.GetScooters(c.Request.Context(), serviceParams)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve scooters"})
		return
	}

	// Convert service result to response format
	response := ScooterListResponse{
		Scooters: make([]ScooterInfo, len(result.Scooters)),
		Total:    result.Total,
		Limit:    result.Limit,
		Offset:   result.Offset,
	}

	for i, scooter := range result.Scooters {
		response.Scooters[i] = ScooterInfo{
			ID:               scooter.ID,
			Status:           scooter.Status,
			CurrentLatitude:  scooter.CurrentLatitude,
			CurrentLongitude: scooter.CurrentLongitude,
			LastSeen:         scooter.LastSeen,
			CreatedAt:        scooter.CreatedAt,
		}
	}

	c.JSON(http.StatusOK, response)
}

// GetScooter retrieves details for a specific scooter
// GET /api/v1/scooters/{id}
func (h *ScooterHandler) GetScooter(c *gin.Context) {
	scooterIDStr := c.Param("id")
	scooterID, err := uuid.Parse(scooterIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid scooter ID"})
		return
	}

	// Call scooter service
	result, err := h.scooterService.GetScooter(c.Request.Context(), scooterID)
	if err != nil {
		if err.Error() == "scooter not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": "Scooter not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve scooter"})
		return
	}

	// Convert service result to response format
	response := ScooterDetailsResponse{
		ID:               result.ID,
		Status:           result.Status,
		CurrentLatitude:  result.CurrentLatitude,
		CurrentLongitude: result.CurrentLongitude,
		LastSeen:         result.LastSeen,
		CreatedAt:        result.CreatedAt,
		UpdatedAt:        result.UpdatedAt,
	}

	// Add active trip information if present
	if result.ActiveTrip != nil {
		response.ActiveTrip = &TripInfo{
			TripID:         result.ActiveTrip.TripID,
			UserID:         result.ActiveTrip.UserID,
			StartTime:      result.ActiveTrip.StartTime,
			StartLatitude:  result.ActiveTrip.StartLatitude,
			StartLongitude: result.ActiveTrip.StartLongitude,
		}
	}

	c.JSON(http.StatusOK, response)
}

// GetClosestScooters finds the closest scooters to a given location
// GET /api/v1/scooters/closest
func (h *ScooterHandler) GetClosestScooters(c *gin.Context) {
	var params ClosestScootersParams
	if err := c.ShouldBindQuery(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Convert to service parameters
	serviceParams := services.ClosestScootersQueryParams{
		Latitude:  params.Latitude,
		Longitude: params.Longitude,
		Radius:    params.Radius,
		Limit:     params.Limit,
		Status:    params.Status,
	}

	// Call scooter service
	result, err := h.scooterService.GetClosestScooters(c.Request.Context(), serviceParams)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve closest scooters"})
		return
	}

	// Convert service result to response format
	response := ClosestScootersResponse{
		Scooters: make([]ScooterWithDistance, len(result.Scooters)),
		Center: Location{
			Latitude:  result.Center.Latitude,
			Longitude: result.Center.Longitude,
		},
		Radius: result.Radius,
	}

	for i, scooter := range result.Scooters {
		response.Scooters[i] = ScooterWithDistance{
			ScooterInfo: ScooterInfo{
				ID:               scooter.ID,
				Status:           scooter.Status,
				CurrentLatitude:  scooter.CurrentLatitude,
				CurrentLongitude: scooter.CurrentLongitude,
				LastSeen:         scooter.LastSeen,
				CreatedAt:        scooter.CreatedAt,
			},
			Distance: scooter.Distance,
		}
	}

	c.JSON(http.StatusOK, response)
}
