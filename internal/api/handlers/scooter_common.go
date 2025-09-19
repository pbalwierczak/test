package handlers

import (
	"time"

	"scootin-aboot/internal/services"

	"github.com/google/uuid"
)

type ScooterHandler struct {
	tripService    services.TripService
	scooterService services.ScooterService
}

func NewScooterHandler(tripService services.TripService, scooterService services.ScooterService) *ScooterHandler {
	return &ScooterHandler{
		tripService:    tripService,
		scooterService: scooterService,
	}
}

type StartTripRequest struct {
	UserID         uuid.UUID `json:"user_id" binding:"required"`
	StartLatitude  float64   `json:"start_latitude" binding:"required"`
	StartLongitude float64   `json:"start_longitude" binding:"required"`
}

type StartTripResponse struct {
	TripID         uuid.UUID `json:"trip_id"`
	ScooterID      uuid.UUID `json:"scooter_id"`
	UserID         uuid.UUID `json:"user_id"`
	StartTime      time.Time `json:"start_time"`
	StartLatitude  float64   `json:"start_latitude"`
	StartLongitude float64   `json:"start_longitude"`
	Status         string    `json:"status"`
}

type EndTripRequest struct {
	EndLatitude  float64 `json:"end_latitude" binding:"required"`
	EndLongitude float64 `json:"end_longitude" binding:"required"`
}

type EndTripResponse struct {
	TripID         uuid.UUID `json:"trip_id"`
	ScooterID      uuid.UUID `json:"scooter_id"`
	UserID         uuid.UUID `json:"user_id"`
	StartTime      time.Time `json:"start_time"`
	EndTime        time.Time `json:"end_time"`
	StartLatitude  float64   `json:"start_latitude"`
	StartLongitude float64   `json:"start_longitude"`
	EndLatitude    float64   `json:"end_latitude"`
	EndLongitude   float64   `json:"end_longitude"`
	Status         string    `json:"status"`
	Duration       int64     `json:"duration_seconds"`
}

type LocationUpdateRequest struct {
	Latitude  float64   `json:"latitude" binding:"required"`
	Longitude float64   `json:"longitude" binding:"required"`
	Timestamp time.Time `json:"timestamp" binding:"required"`
}

type LocationUpdateResponse struct {
	UpdateID  uuid.UUID `json:"update_id"`
	ScooterID uuid.UUID `json:"scooter_id"`
	Latitude  float64   `json:"latitude"`
	Longitude float64   `json:"longitude"`
	Timestamp time.Time `json:"timestamp"`
}

type ScooterQueryParams struct {
	Status string  `form:"status"`
	MinLat float64 `form:"min_lat"`
	MaxLat float64 `form:"max_lat"`
	MinLng float64 `form:"min_lng"`
	MaxLng float64 `form:"max_lng"`
	Limit  int     `form:"limit,default=50"`
	Offset int     `form:"offset,default=0"`
}

type ScooterListResponse struct {
	Scooters []ScooterInfo `json:"scooters"`
	Total    int64         `json:"total"`
	Limit    int           `json:"limit"`
	Offset   int           `json:"offset"`
}

type ScooterInfo struct {
	ID               uuid.UUID `json:"id"`
	Status           string    `json:"status"`
	CurrentLatitude  float64   `json:"current_latitude"`
	CurrentLongitude float64   `json:"current_longitude"`
	LastSeen         time.Time `json:"last_seen"`
	CreatedAt        time.Time `json:"created_at"`
}

type ScooterDetailsResponse struct {
	ID               uuid.UUID `json:"id"`
	Status           string    `json:"status"`
	CurrentLatitude  float64   `json:"current_latitude"`
	CurrentLongitude float64   `json:"current_longitude"`
	LastSeen         time.Time `json:"last_seen"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
	ActiveTrip       *TripInfo `json:"active_trip,omitempty"`
}

type TripInfo struct {
	TripID         uuid.UUID `json:"trip_id"`
	UserID         uuid.UUID `json:"user_id"`
	StartTime      time.Time `json:"start_time"`
	StartLatitude  float64   `json:"start_latitude"`
	StartLongitude float64   `json:"start_longitude"`
}

type ClosestScootersParams struct {
	Latitude  float64 `form:"lat" binding:"required"`
	Longitude float64 `form:"lng" binding:"required"`
	Radius    float64 `form:"radius,default=1000"` // in meters
	Limit     int     `form:"limit,default=10"`
	Status    string  `form:"status"`
}

type ClosestScootersResponse struct {
	Scooters []ScooterWithDistance `json:"scooters"`
	Center   Location              `json:"center"`
	Radius   float64               `json:"radius_meters"`
}

type ScooterWithDistance struct {
	ScooterInfo
	Distance float64 `json:"distance_meters"`
}

type Location struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}
