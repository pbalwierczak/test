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

// StartTrip starts a new trip for a scooter
func (s *tripService) StartTrip(ctx context.Context, scooterID, userID uuid.UUID, lat, lng float64) (*models.Trip, error) {
	// Validate coordinates
	if err := validation.ValidateCoordinates(lat, lng); err != nil {
		return nil, fmt.Errorf("invalid coordinates: %w", err)
	}

	// Check if user exists
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	if user == nil {
		return nil, errors.New("user not found")
	}

	// Check if user already has an active trip
	activeTrip, err := s.tripRepo.GetActiveByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to check user's active trip: %w", err)
	}
	if activeTrip != nil {
		return nil, errors.New("user already has an active trip")
	}

	// Lock scooter for update to prevent concurrent modifications
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

	// Check if scooter already has an active trip (double-check with lock)
	activeScooterTrip, err := s.tripRepo.GetActiveByScooterID(ctx, scooterID)
	if err != nil {
		return nil, fmt.Errorf("failed to check scooter's active trip: %w", err)
	}
	if activeScooterTrip != nil {
		return nil, errors.New("scooter already has an active trip")
	}

	// Create new trip
	trip := &models.Trip{
		ScooterID:      scooterID,
		UserID:         userID,
		StartTime:      time.Now(),
		StartLatitude:  lat,
		StartLongitude: lng,
		Status:         models.TripStatusActive,
	}

	// Create trip in database
	if err := s.tripRepo.Create(ctx, trip); err != nil {
		return nil, fmt.Errorf("failed to create trip: %w", err)
	}

	// Update scooter status to occupied with status check
	if err := s.scooterRepo.UpdateStatusWithCheck(ctx, scooterID, models.ScooterStatusOccupied, models.ScooterStatusAvailable); err != nil {
		// If scooter update fails, we should clean up the trip
		// In a real system, we might want to implement compensation logic
		return nil, fmt.Errorf("failed to update scooter status: %w", err)
	}

	// Update scooter location
	if err := s.scooterRepo.UpdateLocation(ctx, scooterID, lat, lng); err != nil {
		// Log error but don't fail the trip start
		// The scooter location will be updated with the first location update
	}

	return trip, nil
}

// EndTrip ends an active trip for a scooter
func (s *tripService) EndTrip(ctx context.Context, scooterID uuid.UUID, lat, lng float64) (*models.Trip, error) {
	// Validate coordinates
	if err := validation.ValidateCoordinates(lat, lng); err != nil {
		return nil, fmt.Errorf("invalid coordinates: %w", err)
	}

	// Get active trip for scooter
	trip, err := s.tripRepo.GetActiveByScooterID(ctx, scooterID)
	if err != nil {
		return nil, fmt.Errorf("failed to get active trip: %w", err)
	}
	if trip == nil {
		return nil, errors.New("no active trip found for scooter")
	}

	// End the trip
	endTime := time.Now()
	if err := s.tripRepo.EndTrip(ctx, trip.ID, lat, lng); err != nil {
		return nil, fmt.Errorf("failed to end trip: %w", err)
	}

	// Update scooter status to available with status check
	if err := s.scooterRepo.UpdateStatusWithCheck(ctx, scooterID, models.ScooterStatusAvailable, models.ScooterStatusOccupied); err != nil {
		return nil, fmt.Errorf("failed to update scooter status: %w", err)
	}

	// Update scooter location
	if err := s.scooterRepo.UpdateLocation(ctx, scooterID, lat, lng); err != nil {
		// Log error but don't fail the trip end
	}

	// Update trip object for response
	trip.EndTime = &endTime
	trip.EndLatitude = &lat
	trip.EndLongitude = &lng
	trip.Status = models.TripStatusCompleted

	return trip, nil
}

// CancelTrip cancels an active trip for a scooter
func (s *tripService) CancelTrip(ctx context.Context, scooterID uuid.UUID) (*models.Trip, error) {
	// Get active trip for scooter
	trip, err := s.tripRepo.GetActiveByScooterID(ctx, scooterID)
	if err != nil {
		return nil, fmt.Errorf("failed to get active trip: %w", err)
	}
	if trip == nil {
		return nil, errors.New("no active trip found for scooter")
	}

	// Cancel the trip
	if err := s.tripRepo.CancelTrip(ctx, trip.ID); err != nil {
		return nil, fmt.Errorf("failed to cancel trip: %w", err)
	}

	// Update scooter status to available with status check
	if err := s.scooterRepo.UpdateStatusWithCheck(ctx, scooterID, models.ScooterStatusAvailable, models.ScooterStatusOccupied); err != nil {
		return nil, fmt.Errorf("failed to update scooter status: %w", err)
	}

	// Update trip object for response
	trip.Status = models.TripStatusCancelled

	return trip, nil
}

// UpdateLocation updates the location of a scooter during an active trip
func (s *tripService) UpdateLocation(ctx context.Context, scooterID uuid.UUID, lat, lng float64) error {
	// Validate coordinates
	if err := validation.ValidateCoordinates(lat, lng); err != nil {
		return fmt.Errorf("invalid coordinates: %w", err)
	}

	// Get active trip for scooter
	trip, err := s.tripRepo.GetActiveByScooterID(ctx, scooterID)
	if err != nil {
		return fmt.Errorf("failed to get active trip: %w", err)
	}
	if trip == nil {
		return errors.New("no active trip found for scooter")
	}

	// Create location update record
	locationUpdate := &models.LocationUpdate{
		TripID:    trip.ID,
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

// GetActiveTrip gets the active trip for a scooter
func (s *tripService) GetActiveTrip(ctx context.Context, scooterID uuid.UUID) (*models.Trip, error) {
	trip, err := s.tripRepo.GetActiveByScooterID(ctx, scooterID)
	if err != nil {
		return nil, fmt.Errorf("failed to get active trip: %w", err)
	}
	return trip, nil
}

// GetActiveTripByUser gets the active trip for a user
func (s *tripService) GetActiveTripByUser(ctx context.Context, userID uuid.UUID) (*models.Trip, error) {
	trip, err := s.tripRepo.GetActiveByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get active trip: %w", err)
	}
	return trip, nil
}

// GetTrip gets a trip by ID
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
