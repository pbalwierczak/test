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

type ScooterService interface {
	GetScooters(ctx context.Context, params ScooterQueryParams) (*ScooterListResult, error)
	GetScooter(ctx context.Context, id uuid.UUID) (*ScooterDetailsResult, error)
	GetClosestScooters(ctx context.Context, params ClosestScootersQueryParams) (*ClosestScootersResult, error)
	UpdateLocation(ctx context.Context, scooterID uuid.UUID, lat, lng float64) error
}

type scooterService struct {
	scooterRepo  repository.ScooterRepository
	tripRepo     repository.TripRepository
	locationRepo repository.LocationUpdateRepository
	unitOfWork   repository.UnitOfWork
}

func NewScooterService(
	scooterRepo repository.ScooterRepository,
	tripRepo repository.TripRepository,
	locationRepo repository.LocationUpdateRepository,
	unitOfWork repository.UnitOfWork,
) ScooterService {
	return &scooterService{
		scooterRepo:  scooterRepo,
		tripRepo:     tripRepo,
		locationRepo: locationRepo,
		unitOfWork:   unitOfWork,
	}
}

type ScooterQueryParams struct {
	Status string
	MinLat float64
	MaxLat float64
	MinLng float64
	MaxLng float64
	Limit  int
	Offset int
}

type ScooterListResult struct {
	Scooters []*ScooterInfo
	Total    int64
	Limit    int
	Offset   int
}

type ScooterInfo struct {
	ID               uuid.UUID `json:"id"`
	Status           string    `json:"status"`
	CurrentLatitude  float64   `json:"current_latitude"`
	CurrentLongitude float64   `json:"current_longitude"`
	LastSeen         time.Time `json:"last_seen"`
	CreatedAt        time.Time `json:"created_at"`
}

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

type TripInfo struct {
	TripID         uuid.UUID `json:"trip_id"`
	UserID         uuid.UUID `json:"user_id"`
	StartTime      time.Time `json:"start_time"`
	StartLatitude  float64   `json:"start_latitude"`
	StartLongitude float64   `json:"start_longitude"`
}

type ClosestScootersQueryParams struct {
	Latitude  float64
	Longitude float64
	Radius    float64
	Limit     int
	Status    string
}

type ClosestScootersResult struct {
	Scooters []*ScooterWithDistance
	Center   Location
	Radius   float64
}

type ScooterWithDistance struct {
	*ScooterInfo
	Distance float64 `json:"distance_meters"`
}

type Location struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

func (s *scooterService) GetScooters(ctx context.Context, params ScooterQueryParams) (*ScooterListResult, error) {
	if err := s.validateScooterQueryParams(params); err != nil {
		return nil, fmt.Errorf("invalid query parameters: %w", err)
	}

	scooters, err := s.queryScootersByFilters(ctx, params)

	if err != nil {
		return nil, fmt.Errorf("failed to query scooters: %w", err)
	}

	scooterInfos := make([]*ScooterInfo, len(scooters))
	for i, scooter := range scooters {
		scooterInfos[i] = s.mapScooterToInfo(scooter)
	}

	total := int64(len(scooterInfos))

	if params.Status == "" && (!s.hasLocationBounds(params)) {
	} else {
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

func (s *scooterService) GetScooter(ctx context.Context, id uuid.UUID) (*ScooterDetailsResult, error) {
	scooter, err := s.scooterRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get scooter: %w", err)
	}
	if scooter == nil {
		return nil, errors.New("scooter not found")
	}

	var activeTrip *TripInfo
	if scooter.Status == models.ScooterStatusOccupied {
		trip, err := s.tripRepo.GetActiveByScooterID(ctx, id)
		if err != nil {
			// Trip fetch error is non-critical
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

func (s *scooterService) GetClosestScooters(ctx context.Context, params ClosestScootersQueryParams) (*ClosestScootersResult, error) {
	if err := s.validateClosestScootersParams(params); err != nil {
		return nil, fmt.Errorf("invalid query parameters: %w", err)
	}

	scooters, err := s.scooterRepo.GetClosestWithRadius(ctx, params.Latitude, params.Longitude, params.Radius, params.Status, params.Limit)
	if err != nil {
		return nil, fmt.Errorf("failed to query closest scooters: %w", err)
	}

	scootersWithDistance := make([]*ScooterWithDistance, len(scooters))
	for i, scooter := range scooters {
		distance := repository.HaversineDistance(params.Latitude, params.Longitude, scooter.CurrentLatitude, scooter.CurrentLongitude)

		scootersWithDistance[i] = &ScooterWithDistance{
			ScooterInfo: s.mapScooterToInfo(scooter),
			Distance:    distance * 1000,
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

func (s *scooterService) validateScooterQueryParams(params ScooterQueryParams) error {
	if params.Status != "" && params.Status != "available" && params.Status != "occupied" {
		return errors.New("status must be 'available' or 'occupied'")
	}

	if err := repository.ValidateGeographicBounds(params.MinLat, params.MaxLat, params.MinLng, params.MaxLng); err != nil {
		return err
	}

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

func (s *scooterService) validateClosestScootersParams(params ClosestScootersQueryParams) error {
	if err := validation.ValidateCoordinates(params.Latitude, params.Longitude); err != nil {
		return err
	}

	if params.Radius < 0 {
		return errors.New("radius must be non-negative")
	}
	if params.Radius > 50000 { // 50km max
		return errors.New("radius cannot exceed 50000 meters")
	}

	if params.Status != "" && params.Status != "available" && params.Status != "occupied" {
		return errors.New("status must be 'available' or 'occupied'")
	}

	if params.Limit < 0 {
		return errors.New("limit must be non-negative")
	}
	if params.Limit > 50 {
		return errors.New("limit cannot exceed 50")
	}

	return nil
}

func (s *scooterService) UpdateLocation(ctx context.Context, scooterID uuid.UUID, lat, lng float64) error {
	if err := validation.ValidateCoordinates(lat, lng); err != nil {
		return fmt.Errorf("invalid coordinates: %w", err)
	}

	tx, err := s.unitOfWork.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	var committed bool
	defer func() {
		if !committed {
			tx.Rollback()
		}
	}()

	scooterRepo := tx.ScooterRepository()
	locationRepo := tx.LocationUpdateRepository()

	scooter, err := scooterRepo.GetByID(ctx, scooterID)
	if err != nil {
		return fmt.Errorf("failed to get scooter: %w", err)
	}
	if scooter == nil {
		return errors.New("scooter not found")
	}

	locationUpdate := &models.LocationUpdate{
		ScooterID: scooterID,
		Latitude:  lat,
		Longitude: lng,
		Timestamp: time.Now(),
	}

	if err := locationRepo.Create(ctx, locationUpdate); err != nil {
		return fmt.Errorf("failed to create location update: %w", err)
	}

	if err := scooterRepo.UpdateLocation(ctx, scooterID, lat, lng); err != nil {
		return fmt.Errorf("failed to update scooter location: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	committed = true
	return nil
}

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

func (s *scooterService) queryScootersByFilters(ctx context.Context, params ScooterQueryParams) ([]*models.Scooter, error) {
	hasStatusFilter := params.Status != ""
	hasLocationFilter := s.hasLocationBounds(params)

	switch {
	case hasStatusFilter && hasLocationFilter:
		status := models.ScooterStatus(params.Status)
		return s.scooterRepo.GetByStatusInBounds(ctx, status, params.MinLat, params.MaxLat, params.MinLng, params.MaxLng)

	case hasStatusFilter:
		status := models.ScooterStatus(params.Status)
		return s.scooterRepo.GetByStatus(ctx, status)

	case hasLocationFilter:
		return s.scooterRepo.GetInBounds(ctx, params.MinLat, params.MaxLat, params.MinLng, params.MaxLng)

	default:
		return s.scooterRepo.List(ctx, params.Limit, params.Offset)
	}
}

func (s *scooterService) hasLocationBounds(params ScooterQueryParams) bool {
	return params.MinLat != 0 || params.MaxLat != 0 || params.MinLng != 0 || params.MaxLng != 0
}
