package handlers

import (
	"errors"
	"net/http"

	"scootin-aboot/internal/api/middleware"
	"scootin-aboot/internal/repository"
	"scootin-aboot/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (h *ScooterHandler) GetScooters(c *gin.Context) {
	var params ScooterQueryParams
	if err := c.ShouldBindQuery(&params); err != nil {
		c.Error(middleware.NewAPIError(http.StatusBadRequest, err.Error()))
		return
	}

	serviceParams := services.ScooterQueryParams{
		Status: params.Status,
		MinLat: params.MinLat,
		MaxLat: params.MaxLat,
		MinLng: params.MinLng,
		MaxLng: params.MaxLng,
		Limit:  params.Limit,
		Offset: params.Offset,
	}

	result, err := h.scooterService.GetScooters(c.Request.Context(), serviceParams)
	if err != nil {
		c.Error(middleware.ErrInternalServer)
		return
	}

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

func (h *ScooterHandler) GetScooter(c *gin.Context) {
	scooterIDStr := c.Param("id")
	scooterID, err := uuid.Parse(scooterIDStr)
	if err != nil {
		c.Error(middleware.NewAPIError(http.StatusBadRequest, "Invalid scooter ID"))
		return
	}

	result, err := h.scooterService.GetScooter(c.Request.Context(), scooterID)
	if err != nil {
		if errors.Is(err, repository.ErrScooterNotFound) {
			c.Error(middleware.ErrNotFound)
			return
		}
		c.Error(middleware.ErrInternalServer)
		return
	}

	response := ScooterDetailsResponse{
		ID:               result.ID,
		Status:           result.Status,
		CurrentLatitude:  result.CurrentLatitude,
		CurrentLongitude: result.CurrentLongitude,
		LastSeen:         result.LastSeen,
		CreatedAt:        result.CreatedAt,
		UpdatedAt:        result.UpdatedAt,
	}

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

func (h *ScooterHandler) GetClosestScooters(c *gin.Context) {
	var params ClosestScootersParams
	if err := c.ShouldBindQuery(&params); err != nil {
		c.Error(middleware.NewAPIError(http.StatusBadRequest, err.Error()))
		return
	}

	serviceParams := services.ClosestScootersQueryParams{
		Latitude:  params.Latitude,
		Longitude: params.Longitude,
		Radius:    params.Radius,
		Limit:     params.Limit,
		Status:    params.Status,
	}

	result, err := h.scooterService.GetClosestScooters(c.Request.Context(), serviceParams)
	if err != nil {
		c.Error(middleware.ErrInternalServer)
		return
	}

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
