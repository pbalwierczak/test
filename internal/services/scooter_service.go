package services

import (
	"context"
	"errors"
	"fmt"
	"time"

	"scootin-aboot/internal/models"
	"scootin-aboot/internal/repository"
	"scootin-aboot/pkg/validation"

	"github.com/google/uuid"
)

// ScooterService defines the interface for scooter query operations
type ScooterService interface {
	// Query operations
	GetScooters(ctx context.Context, params ScooterQueryParams) (*ScooterListResult, error)
	GetScooter(ctx context.Context, id uuid.UUID) (*ScooterDetailsResult, error)
	GetClosestScooters(ctx context.Context, params ClosestScootersQueryParams) (*ClosestScootersResult, error)

	// Location operations
	UpdateLocation(ctx context.Context, scooterID uuid.UUID, lat, lng float64) error
}

// scooterService implements ScooterService interface
type scooterService struct {
	scooterRepo  repository.ScooterRepository
	tripRepo     repository.TripRepository
	locationRepo repository.LocationUpdateRepository
}

// NewScooterService creates a new scooter service instance
func NewScooterService(
	scooterRepo repository.ScooterRepository,
	tripRepo repository.TripRepository,
	locationRepo repository.LocationUpdateRepository,
) ScooterService {
	return &scooterService{
		scooterRepo:  scooterRepo,
		tripRepo:     tripRepo,
		locationRepo: locationRepo,
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
		distance := repository.HaversineDistance(params.Latitude, params.Longitude, scooter.CurrentLatitude, scooter.CurrentLongitude)

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
	if err := repository.ValidateGeographicBounds(params.MinLat, params.MaxLat, params.MinLng, params.MaxLng); err != nil {
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
	if err := validation.ValidateCoordinates(params.Latitude, params.Longitude); err != nil {
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

// UpdateLocation updates the location of a scooter and creates a location update record
func (s *scooterService) UpdateLocation(ctx context.Context, scooterID uuid.UUID, lat, lng float64) error {
	// Validate coordinates
	if err := validation.ValidateCoordinates(lat, lng); err != nil {
		return fmt.Errorf("invalid coordinates: %w", err)
	}

	// Check if scooter exists
	scooter, err := s.scooterRepo.GetByID(ctx, scooterID)
	if err != nil {
		return fmt.Errorf("failed to get scooter: %w", err)
	}
	if scooter == nil {
		return errors.New("scooter not found")
	}

	// Create location update record
	locationUpdate := &models.LocationUpdate{
		ScooterID: scooterID,
		Latitude:  lat,
		Longitude: lng,
		Timestamp: time.Now(),
	}

	// Save location update
	if err := s.locationRepo.Create(ctx, locationUpdate); err != nil {
		return fmt.Errorf("failed to create location update: %w", err)
	}

	// Update scooter's current location
	if err := s.scooterRepo.UpdateLocation(ctx, scooterID, lat, lng); err != nil {
		return fmt.Errorf("failed to update scooter location: %w", err)
	}

	return nil
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
