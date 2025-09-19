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

// TripService defines the interface for trip management operations
type TripService interface {
	// Trip lifecycle management
	StartTrip(ctx context.Context, scooterID, userID uuid.UUID, lat, lng float64) (*models.Trip, error)
	EndTrip(ctx context.Context, scooterID uuid.UUID, lat, lng float64) (*models.Trip, error)
	CancelTrip(ctx context.Context, scooterID uuid.UUID) (*models.Trip, error)

	// Location updates
	UpdateLocation(ctx context.Context, scooterID uuid.UUID, lat, lng float64) error

	// Query operations
	GetActiveTrip(ctx context.Context, scooterID uuid.UUID) (*models.Trip, error)
	GetActiveTripByUser(ctx context.Context, userID uuid.UUID) (*models.Trip, error)
	GetTrip(ctx context.Context, tripID uuid.UUID) (*models.Trip, error)
}

// tripService implements TripService interface
type tripService struct {
	tripRepo     repository.TripRepository
	scooterRepo  repository.ScooterRepository
	userRepo     repository.UserRepository
	locationRepo repository.LocationUpdateRepository
}

// NewTripService creates a new trip service instance
func NewTripService(
	tripRepo repository.TripRepository,
	scooterRepo repository.ScooterRepository,
	userRepo repository.UserRepository,
	locationRepo repository.LocationUpdateRepository,
) TripService {
	return &tripService{
		tripRepo:     tripRepo,
		scooterRepo:  scooterRepo,
		userRepo:     userRepo,
		locationRepo: locationRepo,
	}
}

func (s *tripService) StartTrip(ctx context.Context, scooterID, userID uuid.UUID, lat, lng float64) (*models.Trip, error) {
	if err := validation.ValidateCoordinates(lat, lng); err != nil {
		return nil, fmt.Errorf("invalid coordinates: %w", err)
	}

	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	if user == nil {
		return nil, errors.New("user not found")
	}

	activeTrip, err := s.tripRepo.GetActiveByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to check user's active trip: %w", err)
	}
	if activeTrip != nil {
		return nil, errors.New("user already has an active trip")
	}

	scooter, err := s.scooterRepo.GetByIDForUpdate(ctx, scooterID)
	if err != nil {
		return nil, fmt.Errorf("failed to get scooter: %w", err)
	}
	if scooter == nil {
		return nil, errors.New("scooter not found")
	}
	if !scooter.IsAvailable() {
		return nil, errors.New("scooter is not available")
	}

	activeScooterTrip, err := s.tripRepo.GetActiveByScooterID(ctx, scooterID)
	if err != nil {
		return nil, fmt.Errorf("failed to check scooter's active trip: %w", err)
	}
	if activeScooterTrip != nil {
		return nil, errors.New("scooter already has an active trip")
	}

	trip := &models.Trip{
		ScooterID:      scooterID,
		UserID:         userID,
		StartTime:      time.Now(),
		StartLatitude:  lat,
		StartLongitude: lng,
		Status:         models.TripStatusActive,
	}

	if err := s.tripRepo.Create(ctx, trip); err != nil {
		return nil, fmt.Errorf("failed to create trip: %w", err)
	}

	if err := s.scooterRepo.UpdateStatusWithCheck(ctx, scooterID, models.ScooterStatusOccupied, models.ScooterStatusAvailable); err != nil {
		if deleteErr := s.tripRepo.Delete(ctx, trip.ID); deleteErr != nil {
			return nil, fmt.Errorf("failed to update scooter status: %w (cleanup failed: %v)", err, deleteErr)
		}
		return nil, fmt.Errorf("failed to update scooter status: %w", err)
	}

	if err := s.scooterRepo.UpdateLocation(ctx, scooterID, lat, lng); err != nil {
		// Log error but don't fail the trip start
		// The scooter location will be updated with the first location update
	}

	return trip, nil
}

func (s *tripService) EndTrip(ctx context.Context, scooterID uuid.UUID, lat, lng float64) (*models.Trip, error) {
	if err := validation.ValidateCoordinates(lat, lng); err != nil {
		return nil, fmt.Errorf("invalid coordinates: %w", err)
	}

	trip, err := s.tripRepo.GetActiveByScooterID(ctx, scooterID)
	if err != nil {
		return nil, fmt.Errorf("failed to get active trip: %w", err)
	}
	if trip == nil {
		return nil, errors.New("no active trip found for scooter")
	}

	endTime := time.Now()
	if err := s.tripRepo.EndTrip(ctx, trip.ID, lat, lng); err != nil {
		return nil, fmt.Errorf("failed to end trip: %w", err)
	}

	if err := s.scooterRepo.UpdateStatusWithCheck(ctx, scooterID, models.ScooterStatusAvailable, models.ScooterStatusOccupied); err != nil {
		return nil, fmt.Errorf("failed to update scooter status: %w", err)
	}

	if err := s.scooterRepo.UpdateLocation(ctx, scooterID, lat, lng); err != nil {
		// Log error but don't fail the trip end
	}

	trip.EndTime = &endTime
	trip.EndLatitude = &lat
	trip.EndLongitude = &lng
	trip.Status = models.TripStatusCompleted

	return trip, nil
}

func (s *tripService) CancelTrip(ctx context.Context, scooterID uuid.UUID) (*models.Trip, error) {
	trip, err := s.tripRepo.GetActiveByScooterID(ctx, scooterID)
	if err != nil {
		return nil, fmt.Errorf("failed to get active trip: %w", err)
	}
	if trip == nil {
		return nil, errors.New("no active trip found for scooter")
	}

	if err := s.tripRepo.CancelTrip(ctx, trip.ID); err != nil {
		return nil, fmt.Errorf("failed to cancel trip: %w", err)
	}

	if err := s.scooterRepo.UpdateStatusWithCheck(ctx, scooterID, models.ScooterStatusAvailable, models.ScooterStatusOccupied); err != nil {
		return nil, fmt.Errorf("failed to update scooter status: %w", err)
	}

	trip.Status = models.TripStatusCancelled

	return trip, nil
}

func (s *tripService) UpdateLocation(ctx context.Context, scooterID uuid.UUID, lat, lng float64) error {
	if err := validation.ValidateCoordinates(lat, lng); err != nil {
		return fmt.Errorf("invalid coordinates: %w", err)
	}

	trip, err := s.tripRepo.GetActiveByScooterID(ctx, scooterID)
	if err != nil {
		return fmt.Errorf("failed to get active trip: %w", err)
	}
	if trip == nil {
		return errors.New("no active trip found for scooter")
	}

	locationUpdate := &models.LocationUpdate{
		ScooterID: scooterID,
		Latitude:  lat,
		Longitude: lng,
		Timestamp: time.Now(),
	}

	if err := s.locationRepo.Create(ctx, locationUpdate); err != nil {
		return fmt.Errorf("failed to create location update: %w", err)
	}

	if err := s.scooterRepo.UpdateLocation(ctx, scooterID, lat, lng); err != nil {
		return fmt.Errorf("failed to update scooter location: %w", err)
	}

	return nil
}

func (s *tripService) GetActiveTrip(ctx context.Context, scooterID uuid.UUID) (*models.Trip, error) {
	trip, err := s.tripRepo.GetActiveByScooterID(ctx, scooterID)
	if err != nil {
		return nil, fmt.Errorf("failed to get active trip: %w", err)
	}
	return trip, nil
}

func (s *tripService) GetActiveTripByUser(ctx context.Context, userID uuid.UUID) (*models.Trip, error) {
	trip, err := s.tripRepo.GetActiveByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get active trip: %w", err)
	}
	return trip, nil
}

func (s *tripService) GetTrip(ctx context.Context, tripID uuid.UUID) (*models.Trip, error) {
	trip, err := s.tripRepo.GetByID(ctx, tripID)
	if err != nil {
		return nil, fmt.Errorf("failed to get trip: %w", err)
	}
	if trip == nil {
		return nil, errors.New("trip not found")
	}
	return trip, nil
}
