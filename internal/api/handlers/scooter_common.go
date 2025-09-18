package handlers

import (
	"time"

	"scootin-aboot/internal/services"

	"github.com/google/uuid"
)

// ScooterHandler handles scooter-related HTTP requests
type ScooterHandler struct {
	tripService services.TripService
}

// NewScooterHandler creates a new scooter handler
func NewScooterHandler(tripService services.TripService) *ScooterHandler {
	return &ScooterHandler{
		tripService: tripService,
	}
}

// Trip-related request/response types

// StartTripRequest represents the request body for starting a trip
type StartTripRequest struct {
	UserID         uuid.UUID `json:"user_id" binding:"required"`
	StartLatitude  float64   `json:"start_latitude" binding:"required"`
	StartLongitude float64   `json:"start_longitude" binding:"required"`
}

// StartTripResponse represents the response for starting a trip
type StartTripResponse struct {
	TripID         uuid.UUID `json:"trip_id"`
	ScooterID      uuid.UUID `json:"scooter_id"`
	UserID         uuid.UUID `json:"user_id"`
	StartTime      time.Time `json:"start_time"`
	StartLatitude  float64   `json:"start_latitude"`
	StartLongitude float64   `json:"start_longitude"`
	Status         string    `json:"status"`
}

// EndTripRequest represents the request body for ending a trip
type EndTripRequest struct {
	EndLatitude  float64 `json:"end_latitude" binding:"required"`
	EndLongitude float64 `json:"end_longitude" binding:"required"`
}

// EndTripResponse represents the response for ending a trip
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

// Location-related request/response types

// LocationUpdateRequest represents the request body for updating location
type LocationUpdateRequest struct {
	Latitude  float64   `json:"latitude" binding:"required"`
	Longitude float64   `json:"longitude" binding:"required"`
	Timestamp time.Time `json:"timestamp" binding:"required"`
}

// LocationUpdateResponse represents the response for location update
type LocationUpdateResponse struct {
	UpdateID  uuid.UUID `json:"update_id"`
	TripID    uuid.UUID `json:"trip_id"`
	Latitude  float64   `json:"latitude"`
	Longitude float64   `json:"longitude"`
	Timestamp time.Time `json:"timestamp"`
}

// Search and query types

// ScooterQueryParams represents query parameters for scooter queries
type ScooterQueryParams struct {
	Status string  `form:"status"`
	MinLat float64 `form:"min_lat"`
	MaxLat float64 `form:"max_lat"`
	MinLng float64 `form:"min_lng"`
	MaxLng float64 `form:"max_lng"`
	Limit  int     `form:"limit,default=50"`
	Offset int     `form:"offset,default=0"`
}

// ScooterListResponse represents the response for scooter list queries
type ScooterListResponse struct {
	Scooters []ScooterInfo `json:"scooters"`
	Total    int64         `json:"total"`
	Limit    int           `json:"limit"`
	Offset   int           `json:"offset"`
}

// ScooterInfo represents scooter information in list responses
type ScooterInfo struct {
	ID               uuid.UUID `json:"id"`
	Status           string    `json:"status"`
	CurrentLatitude  float64   `json:"current_latitude"`
	CurrentLongitude float64   `json:"current_longitude"`
	LastSeen         time.Time `json:"last_seen"`
	CreatedAt        time.Time `json:"created_at"`
}

// ScooterDetailsResponse represents the response for scooter details
type ScooterDetailsResponse struct {
	ID               uuid.UUID `json:"id"`
	Status           string    `json:"status"`
	CurrentLatitude  float64   `json:"current_latitude"`
	CurrentLongitude float64   `json:"current_longitude"`
	LastSeen         time.Time `json:"last_seen"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
	// TODO: Add active trip information if scooter is occupied
	ActiveTrip *TripInfo `json:"active_trip,omitempty"`
}

// TripInfo represents trip information in scooter details
type TripInfo struct {
	TripID         uuid.UUID `json:"trip_id"`
	UserID         uuid.UUID `json:"user_id"`
	StartTime      time.Time `json:"start_time"`
	StartLatitude  float64   `json:"start_latitude"`
	StartLongitude float64   `json:"start_longitude"`
}

// Closest scooters types

// ClosestScootersParams represents query parameters for closest scooters
type ClosestScootersParams struct {
	Latitude  float64 `form:"lat" binding:"required"`
	Longitude float64 `form:"lng" binding:"required"`
	Radius    float64 `form:"radius,default=1000"` // in meters
	Limit     int     `form:"limit,default=10"`
	Status    string  `form:"status"`
}

// ClosestScootersResponse represents the response for closest scooters
type ClosestScootersResponse struct {
	Scooters []ScooterWithDistance `json:"scooters"`
	Center   Location              `json:"center"`
	Radius   float64               `json:"radius_meters"`
}

// ScooterWithDistance represents a scooter with distance information
type ScooterWithDistance struct {
	ScooterInfo
	Distance float64 `json:"distance_meters"`
}

// Location represents a geographic location
type Location struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}
