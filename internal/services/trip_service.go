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
	unitOfWork   repository.UnitOfWork
}

// NewTripService creates a new trip service instance
func NewTripService(
	tripRepo repository.TripRepository,
	scooterRepo repository.ScooterRepository,
	userRepo repository.UserRepository,
	locationRepo repository.LocationUpdateRepository,
	unitOfWork repository.UnitOfWork,
) TripService {
	return &tripService{
		tripRepo:     tripRepo,
		scooterRepo:  scooterRepo,
		userRepo:     userRepo,
		locationRepo: locationRepo,
		unitOfWork:   unitOfWork,
	}
}

func (s *tripService) StartTrip(ctx context.Context, scooterID, userID uuid.UUID, lat, lng float64) (*models.Trip, error) {
	if err := validation.ValidateCoordinates(lat, lng); err != nil {
		return nil, fmt.Errorf("invalid coordinates: %w", err)
	}

	// Start a transaction using Unit of Work
	tx, err := s.unitOfWork.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}

	// Ensure transaction is rolled back if there's an error
	var committed bool
	defer func() {
		if !committed {
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				// Log rollback error but don't override the original error
				// In production, you'd want to use proper logging here
			}
		}
	}()

	// Get repositories from transaction
	userRepo := tx.UserRepository()
	tripRepo := tx.TripRepository()
	scooterRepo := tx.ScooterRepository()

	// Validate user exists
	user, err := userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	if user == nil {
		return nil, errors.New("user not found")
	}

	// Check if user already has an active trip
	activeTrip, err := tripRepo.GetActiveByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to check user's active trip: %w", err)
	}
	if activeTrip != nil {
		return nil, errors.New("user already has an active trip")
	}

	// Get scooter with lock for update
	scooter, err := scooterRepo.GetByIDForUpdate(ctx, scooterID)
	if err != nil {
		return nil, fmt.Errorf("failed to get scooter: %w", err)
	}
	if scooter == nil {
		return nil, errors.New("scooter not found")
	}
	if !scooter.IsAvailable() {
		return nil, errors.New("scooter is not available")
	}

	// Check if scooter already has an active trip
	activeScooterTrip, err := tripRepo.GetActiveByScooterID(ctx, scooterID)
	if err != nil {
		return nil, fmt.Errorf("failed to check scooter's active trip: %w", err)
	}
	if activeScooterTrip != nil {
		return nil, errors.New("scooter already has an active trip")
	}

	// Create trip record
	trip := &models.Trip{
		ScooterID:      scooterID,
		UserID:         userID,
		StartTime:      time.Now(),
		StartLatitude:  lat,
		StartLongitude: lng,
		Status:         models.TripStatusActive,
	}

	if err := tripRepo.Create(ctx, trip); err != nil {
		return nil, fmt.Errorf("failed to create trip: %w", err)
	}

	// Update scooter status to occupied
	if err := scooterRepo.UpdateStatusWithCheck(ctx, scooterID, models.ScooterStatusOccupied, models.ScooterStatusAvailable); err != nil {
		return nil, fmt.Errorf("failed to update scooter status: %w", err)
	}

	// Update scooter location
	if err := scooterRepo.UpdateLocation(ctx, scooterID, lat, lng); err != nil {
		// Log error but don't fail the trip start
		// The scooter location will be updated with the first location update
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	committed = true
	return trip, nil
}

func (s *tripService) EndTrip(ctx context.Context, scooterID uuid.UUID, lat, lng float64) (*models.Trip, error) {
	if err := validation.ValidateCoordinates(lat, lng); err != nil {
		return nil, fmt.Errorf("invalid coordinates: %w", err)
	}

	// Start a transaction using Unit of Work
	tx, err := s.unitOfWork.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}

	// Ensure transaction is rolled back if there's an error
	var committed bool
	defer func() {
		if !committed {
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				// Log rollback error but don't override the original error
				// In production, you'd want to use proper logging here
			}
		}
	}()

	// Get repositories from transaction
	tripRepo := tx.TripRepository()
	scooterRepo := tx.ScooterRepository()

	// Get active trip
	trip, err := tripRepo.GetActiveByScooterID(ctx, scooterID)
	if err != nil {
		return nil, fmt.Errorf("failed to get active trip: %w", err)
	}
	if trip == nil {
		return nil, errors.New("no active trip found for scooter")
	}

	endTime := time.Now()

	// End the trip
	if err := tripRepo.EndTrip(ctx, trip.ID, lat, lng); err != nil {
		return nil, fmt.Errorf("failed to end trip: %w", err)
	}

	// Update scooter status back to available
	if err := scooterRepo.UpdateStatusWithCheck(ctx, scooterID, models.ScooterStatusAvailable, models.ScooterStatusOccupied); err != nil {
		return nil, fmt.Errorf("failed to update scooter status: %w", err)
	}

	// Update scooter location
	if err := scooterRepo.UpdateLocation(ctx, scooterID, lat, lng); err != nil {
		// Log error but don't fail the trip end
	}

	// Update trip object for return
	trip.EndTime = &endTime
	trip.EndLatitude = &lat
	trip.EndLongitude = &lng
	trip.Status = models.TripStatusCompleted

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	committed = true
	return trip, nil
}

func (s *tripService) CancelTrip(ctx context.Context, scooterID uuid.UUID) (*models.Trip, error) {
	// Start a transaction using Unit of Work
	tx, err := s.unitOfWork.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}

	// Ensure transaction is rolled back if there's an error
	var committed bool
	defer func() {
		if !committed {
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				// Log rollback error but don't override the original error
				// In production, you'd want to use proper logging here
			}
		}
	}()

	// Get repositories from transaction
	tripRepo := tx.TripRepository()
	scooterRepo := tx.ScooterRepository()

	// Get active trip
	trip, err := tripRepo.GetActiveByScooterID(ctx, scooterID)
	if err != nil {
		return nil, fmt.Errorf("failed to get active trip: %w", err)
	}
	if trip == nil {
		return nil, errors.New("no active trip found for scooter")
	}

	// Cancel the trip
	if err := tripRepo.CancelTrip(ctx, trip.ID); err != nil {
		return nil, fmt.Errorf("failed to cancel trip: %w", err)
	}

	// Update scooter status back to available
	if err := scooterRepo.UpdateStatusWithCheck(ctx, scooterID, models.ScooterStatusAvailable, models.ScooterStatusOccupied); err != nil {
		return nil, fmt.Errorf("failed to update scooter status: %w", err)
	}

	// Update trip object for return
	trip.Status = models.TripStatusCancelled

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	committed = true
	return trip, nil
}

func (s *tripService) UpdateLocation(ctx context.Context, scooterID uuid.UUID, lat, lng float64) error {
	if err := validation.ValidateCoordinates(lat, lng); err != nil {
		return fmt.Errorf("invalid coordinates: %w", err)
	}

	// Start a transaction using Unit of Work
	tx, err := s.unitOfWork.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	// Ensure transaction is rolled back if there's an error
	var committed bool
	defer func() {
		if !committed {
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				// Log rollback error but don't override the original error
				// In production, you'd want to use proper logging here
			}
		}
	}()

	// Get repositories from transaction
	tripRepo := tx.TripRepository()
	locationRepo := tx.LocationUpdateRepository()
	scooterRepo := tx.ScooterRepository()

	// Check if there's an active trip
	trip, err := tripRepo.GetActiveByScooterID(ctx, scooterID)
	if err != nil {
		return fmt.Errorf("failed to get active trip: %w", err)
	}
	if trip == nil {
		return errors.New("no active trip found for scooter")
	}

	// Create location update record
	locationUpdate := &models.LocationUpdate{
		ScooterID: scooterID,
		Latitude:  lat,
		Longitude: lng,
		Timestamp: time.Now(),
	}

	if err := locationRepo.Create(ctx, locationUpdate); err != nil {
		return fmt.Errorf("failed to create location update: %w", err)
	}

	// Update scooter location
	if err := scooterRepo.UpdateLocation(ctx, scooterID, lat, lng); err != nil {
		return fmt.Errorf("failed to update scooter location: %w", err)
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	committed = true
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
