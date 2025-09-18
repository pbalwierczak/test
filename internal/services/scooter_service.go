package services

import (
	"context"
	"errors"
	"fmt"
	"math"
	"time"

	"scootin-aboot/internal/models"
	"scootin-aboot/internal/repository"

	"github.com/google/uuid"
)

// ScooterService defines the interface for scooter query operations
type ScooterService interface {
	// Query operations
	GetScooters(ctx context.Context, params ScooterQueryParams) (*ScooterListResult, error)
	GetScooter(ctx context.Context, id uuid.UUID) (*ScooterDetailsResult, error)
	GetClosestScooters(ctx context.Context, params ClosestScootersQueryParams) (*ClosestScootersResult, error)
}

// scooterService implements ScooterService interface
type scooterService struct {
	scooterRepo repository.ScooterRepository
	tripRepo    repository.TripRepository
}

// NewScooterService creates a new scooter service instance
func NewScooterService(
	scooterRepo repository.ScooterRepository,
	tripRepo repository.TripRepository,
) ScooterService {
	return &scooterService{
		scooterRepo: scooterRepo,
		tripRepo:    tripRepo,
	}
}

// ScooterQueryParams represents parameters for scooter queries
type ScooterQueryParams struct {
	Status string
	MinLat float64
	MaxLat float64
	MinLng float64
	MaxLng float64
	Limit  int
	Offset int
}

// ScooterListResult represents the result of a scooter list query
type ScooterListResult struct {
	Scooters []*ScooterInfo
	Total    int64
	Limit    int
	Offset   int
}

// ScooterInfo represents scooter information in query results
type ScooterInfo struct {
	ID               uuid.UUID `json:"id"`
	Status           string    `json:"status"`
	CurrentLatitude  float64   `json:"current_latitude"`
	CurrentLongitude float64   `json:"current_longitude"`
	LastSeen         time.Time `json:"last_seen"`
	CreatedAt        time.Time `json:"created_at"`
}

// ScooterDetailsResult represents the result of a single scooter query
type ScooterDetailsResult struct {
	ID               uuid.UUID `json:"id"`
	Status           string    `json:"status"`
	CurrentLatitude  float64   `json:"current_latitude"`
	CurrentLongitude float64   `json:"current_longitude"`
	LastSeen         time.Time `json:"last_seen"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
	ActiveTrip       *TripInfo `json:"active_trip,omitempty"`
}

// TripInfo represents trip information in scooter details
type TripInfo struct {
	TripID         uuid.UUID `json:"trip_id"`
	UserID         uuid.UUID `json:"user_id"`
	StartTime      time.Time `json:"start_time"`
	StartLatitude  float64   `json:"start_latitude"`
	StartLongitude float64   `json:"start_longitude"`
}

// ClosestScootersQueryParams represents parameters for closest scooters queries
type ClosestScootersQueryParams struct {
	Latitude  float64
	Longitude float64
	Radius    float64
	Limit     int
	Status    string
}

// ClosestScootersResult represents the result of a closest scooters query
type ClosestScootersResult struct {
	Scooters []*ScooterWithDistance
	Center   Location
	Radius   float64
}

// ScooterWithDistance represents a scooter with distance information
type ScooterWithDistance struct {
	*ScooterInfo
	Distance float64 `json:"distance_meters"`
}

// Location represents a geographic location
type Location struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

// GetScooters retrieves scooters with optional filtering
func (s *scooterService) GetScooters(ctx context.Context, params ScooterQueryParams) (*ScooterListResult, error) {
	// Validate parameters
	if err := s.validateScooterQueryParams(params); err != nil {
		return nil, fmt.Errorf("invalid query parameters: %w", err)
	}

	var scooters []*models.Scooter
	var err error

	// Determine which repository method to use based on filters
	if params.Status != "" && (params.MinLat != 0 || params.MaxLat != 0 || params.MinLng != 0 || params.MaxLng != 0) {
		// Both status and geographic bounds
		status := models.ScooterStatus(params.Status)
		scooters, err = s.scooterRepo.GetByStatusInBounds(ctx, status, params.MinLat, params.MaxLat, params.MinLng, params.MaxLng)
	} else if params.Status != "" {
		// Only status filter
		status := models.ScooterStatus(params.Status)
		scooters, err = s.scooterRepo.GetByStatus(ctx, status)
	} else if params.MinLat != 0 || params.MaxLat != 0 || params.MinLng != 0 || params.MaxLng != 0 {
		// Only geographic bounds
		scooters, err = s.scooterRepo.GetInBounds(ctx, params.MinLat, params.MaxLat, params.MinLng, params.MaxLng)
	} else {
		// No filters, get all with pagination
		scooters, err = s.scooterRepo.List(ctx, params.Limit, params.Offset)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to query scooters: %w", err)
	}

	// Convert to response format
	scooterInfos := make([]*ScooterInfo, len(scooters))
	for i, scooter := range scooters {
		scooterInfos[i] = s.mapScooterToInfo(scooter)
	}

	// For now, we'll use the length as total count
	// In a production system, you might want to implement a separate count query
	total := int64(len(scooterInfos))

	// Apply pagination if not already applied by repository
	if params.Status == "" && (params.MinLat == 0 && params.MaxLat == 0 && params.MinLng == 0 && params.MaxLng == 0) {
		// Pagination was already applied by List method
	} else {
		// Apply pagination to results
		start := params.Offset
		end := start + params.Limit
		if end > len(scooterInfos) {
			end = len(scooterInfos)
		}
		if start > len(scooterInfos) {
			scooterInfos = []*ScooterInfo{}
		} else {
			scooterInfos = scooterInfos[start:end]
		}
	}

	return &ScooterListResult{
		Scooters: scooterInfos,
		Total:    total,
		Limit:    params.Limit,
		Offset:   params.Offset,
	}, nil
}

// GetScooter retrieves details for a specific scooter
func (s *scooterService) GetScooter(ctx context.Context, id uuid.UUID) (*ScooterDetailsResult, error) {
	scooter, err := s.scooterRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get scooter: %w", err)
	}
	if scooter == nil {
		return nil, errors.New("scooter not found")
	}

	// Get active trip if scooter is occupied
	var activeTrip *TripInfo
	if scooter.Status == models.ScooterStatusOccupied {
		trip, err := s.tripRepo.GetActiveByScooterID(ctx, id)
		if err != nil {
			// Log error but don't fail the request
			// In production, you might want to log this
		} else if trip != nil {
			activeTrip = &TripInfo{
				TripID:         trip.ID,
				UserID:         trip.UserID,
				StartTime:      trip.StartTime,
				StartLatitude:  trip.StartLatitude,
				StartLongitude: trip.StartLongitude,
			}
		}
	}

	return &ScooterDetailsResult{
		ID:               scooter.ID,
		Status:           string(scooter.Status),
		CurrentLatitude:  scooter.CurrentLatitude,
		CurrentLongitude: scooter.CurrentLongitude,
		LastSeen:         scooter.LastSeen,
		CreatedAt:        scooter.CreatedAt,
		UpdatedAt:        scooter.UpdatedAt,
		ActiveTrip:       activeTrip,
	}, nil
}

// GetClosestScooters finds the closest scooters to a given location
func (s *scooterService) GetClosestScooters(ctx context.Context, params ClosestScootersQueryParams) (*ClosestScootersResult, error) {
	// Validate parameters
	if err := s.validateClosestScootersParams(params); err != nil {
		return nil, fmt.Errorf("invalid query parameters: %w", err)
	}

	// Query scooters using repository
	scooters, err := s.scooterRepo.GetClosestWithRadius(ctx, params.Latitude, params.Longitude, params.Radius, params.Status, params.Limit)
	if err != nil {
		return nil, fmt.Errorf("failed to query closest scooters: %w", err)
	}

	// Convert to response format with distances
	scootersWithDistance := make([]*ScooterWithDistance, len(scooters))
	for i, scooter := range scooters {
		// Calculate distance using Haversine formula
		distance := s.calculateDistance(params.Latitude, params.Longitude, scooter.CurrentLatitude, scooter.CurrentLongitude)

		scootersWithDistance[i] = &ScooterWithDistance{
			ScooterInfo: s.mapScooterToInfo(scooter),
			Distance:    distance * 1000, // Convert km to meters
		}
	}

	return &ClosestScootersResult{
		Scooters: scootersWithDistance,
		Center: Location{
			Latitude:  params.Latitude,
			Longitude: params.Longitude,
		},
		Radius: params.Radius,
	}, nil
}

// validateScooterQueryParams validates scooter query parameters
func (s *scooterService) validateScooterQueryParams(params ScooterQueryParams) error {
	// Validate status
	if params.Status != "" && params.Status != "available" && params.Status != "occupied" {
		return errors.New("status must be 'available' or 'occupied'")
	}

	// Validate geographic bounds
	if err := s.validateGeographicBounds(params.MinLat, params.MaxLat, params.MinLng, params.MaxLng); err != nil {
		return err
	}

	// Validate pagination
	if params.Limit < 0 {
		return errors.New("limit must be non-negative")
	}
	if params.Offset < 0 {
		return errors.New("offset must be non-negative")
	}
	if params.Limit > 100 {
		return errors.New("limit cannot exceed 100")
	}

	return nil
}

// validateClosestScootersParams validates closest scooters query parameters
func (s *scooterService) validateClosestScootersParams(params ClosestScootersQueryParams) error {
	// Validate coordinates
	if err := s.validateCoordinates(params.Latitude, params.Longitude); err != nil {
		return err
	}

	// Validate radius
	if params.Radius < 0 {
		return errors.New("radius must be non-negative")
	}
	if params.Radius > 50000 { // 50km max
		return errors.New("radius cannot exceed 50000 meters")
	}

	// Validate status
	if params.Status != "" && params.Status != "available" && params.Status != "occupied" {
		return errors.New("status must be 'available' or 'occupied'")
	}

	// Validate limit
	if params.Limit < 0 {
		return errors.New("limit must be non-negative")
	}
	if params.Limit > 50 {
		return errors.New("limit cannot exceed 50")
	}

	return nil
}

// validateGeographicBounds validates geographic bounding box parameters
func (s *scooterService) validateGeographicBounds(minLat, maxLat, minLng, maxLng float64) error {
	// Check if any bounds are provided
	if minLat == 0 && maxLat == 0 && minLng == 0 && maxLng == 0 {
		return nil // No bounds provided, that's okay
	}

	// Validate individual coordinates
	if err := s.validateCoordinates(minLat, minLng); err != nil {
		return fmt.Errorf("invalid min bounds: %w", err)
	}
	if err := s.validateCoordinates(maxLat, maxLng); err != nil {
		return fmt.Errorf("invalid max bounds: %w", err)
	}

	// Validate bounds consistency
	if minLat >= maxLat {
		return errors.New("min_lat must be less than max_lat")
	}
	if minLng >= maxLng {
		return errors.New("min_lng must be less than max_lng")
	}

	// Check for reasonable bounds (not too large)
	latDiff := maxLat - minLat
	lngDiff := maxLng - minLng
	if latDiff > 10 || lngDiff > 10 {
		return errors.New("geographic bounds are too large (max 10 degrees)")
	}

	return nil
}

// validateCoordinates validates latitude and longitude values
func (s *scooterService) validateCoordinates(lat, lng float64) error {
	if lat < -90 || lat > 90 {
		return errors.New("latitude must be between -90 and 90")
	}
	if lng < -180 || lng > 180 {
		return errors.New("longitude must be between -180 and 180")
	}
	return nil
}

// calculateDistance calculates distance between two points using Haversine formula
func (s *scooterService) calculateDistance(lat1, lon1, lat2, lon2 float64) float64 {
	const earthRadius = 6371 // Earth's radius in kilometers

	// Convert degrees to radians
	lat1Rad := lat1 * math.Pi / 180
	lon1Rad := lon1 * math.Pi / 180
	lat2Rad := lat2 * math.Pi / 180
	lon2Rad := lon2 * math.Pi / 180

	// Haversine formula
	dlat := lat2Rad - lat1Rad
	dlon := lon2Rad - lon1Rad

	a := math.Sin(dlat/2)*math.Sin(dlat/2) +
		math.Cos(lat1Rad)*math.Cos(lat2Rad)*
			math.Sin(dlon/2)*math.Sin(dlon/2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	return earthRadius * c
}

// mapScooterToInfo maps a scooter model to scooter info
func (s *scooterService) mapScooterToInfo(scooter *models.Scooter) *ScooterInfo {
	return &ScooterInfo{
		ID:               scooter.ID,
		Status:           string(scooter.Status),
		CurrentLatitude:  scooter.CurrentLatitude,
		CurrentLongitude: scooter.CurrentLongitude,
		LastSeen:         scooter.LastSeen,
		CreatedAt:        scooter.CreatedAt,
	}
}
