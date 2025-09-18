package handlers

import (
	"net/http"
	"time"

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

	// TODO: Implement scooter service logic
	// 1. Validate query parameters
	// 2. Apply filters (status, geographic bounds)
	// 3. Execute database query with pagination
	// 4. Return filtered results

	// Dummy response for now
	response := ScooterListResponse{
		Scooters: []ScooterInfo{
			{
				ID:               uuid.New(),
				Status:           "available",
				CurrentLatitude:  45.4215,
				CurrentLongitude: -75.6972,
				LastSeen:         time.Now(),
				CreatedAt:        time.Now().Add(-24 * time.Hour),
			},
		},
		Total:  1,
		Limit:  params.Limit,
		Offset: params.Offset,
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

	// TODO: Implement scooter service logic
	// 1. Check if scooter exists
	// 2. Get scooter details
	// 3. If occupied, get active trip information
	// 4. Return complete scooter information

	// Dummy response for now
	response := ScooterDetailsResponse{
		ID:               scooterID,
		Status:           "available",
		CurrentLatitude:  45.4215,
		CurrentLongitude: -75.6972,
		LastSeen:         time.Now(),
		CreatedAt:        time.Now().Add(-24 * time.Hour),
		UpdatedAt:        time.Now(),
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

	// TODO: Implement scooter service logic
	// 1. Validate location parameters
	// 2. Calculate distances using Haversine formula
	// 3. Filter by status if specified
	// 4. Filter by radius
	// 5. Sort by distance
	// 6. Apply limit
	// 7. Return closest scooters

	// Dummy response for now
	response := ClosestScootersResponse{
		Scooters: []ScooterWithDistance{
			{
				ScooterInfo: ScooterInfo{
					ID:               uuid.New(),
					Status:           "available",
					CurrentLatitude:  45.4215,
					CurrentLongitude: -75.6972,
					LastSeen:         time.Now(),
					CreatedAt:        time.Now().Add(-24 * time.Hour),
				},
				Distance: 150.5,
			},
		},
		Center: Location{
			Latitude:  params.Latitude,
			Longitude: params.Longitude,
		},
		Radius: params.Radius,
	}

	c.JSON(http.StatusOK, response)
}
