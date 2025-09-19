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

type TripService interface {
	StartTrip(ctx context.Context, scooterID, userID uuid.UUID, lat, lng float64) (*models.Trip, error)
	EndTrip(ctx context.Context, scooterID uuid.UUID, lat, lng float64) (*models.Trip, error)
	CancelTrip(ctx context.Context, scooterID uuid.UUID) (*models.Trip, error)
	UpdateLocation(ctx context.Context, scooterID uuid.UUID, lat, lng float64) error
	GetActiveTrip(ctx context.Context, scooterID uuid.UUID) (*models.Trip, error)
	GetActiveTripByUser(ctx context.Context, userID uuid.UUID) (*models.Trip, error)
	GetTrip(ctx context.Context, tripID uuid.UUID) (*models.Trip, error)
}

type tripService struct {
	tripRepo     repository.TripRepository
	scooterRepo  repository.ScooterRepository
	userRepo     repository.UserRepository
	locationRepo repository.LocationUpdateRepository
	unitOfWork   repository.UnitOfWork
}

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

	tx, err := s.unitOfWork.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}

	var committed bool
	defer func() {
		if !committed {
			tx.Rollback()
		}
	}()

	userRepo := tx.UserRepository()
	tripRepo := tx.TripRepository()
	scooterRepo := tx.ScooterRepository()

	user, err := userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	if user == nil {
		return nil, errors.New("user not found")
	}

	activeTrip, err := tripRepo.GetActiveByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to check user's active trip: %w", err)
	}
	if activeTrip != nil {
		return nil, errors.New("user already has an active trip")
	}

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

	activeScooterTrip, err := tripRepo.GetActiveByScooterID(ctx, scooterID)
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

	if err := tripRepo.Create(ctx, trip); err != nil {
		return nil, fmt.Errorf("failed to create trip: %w", err)
	}

	if err := scooterRepo.UpdateStatusWithCheck(ctx, scooterID, models.ScooterStatusOccupied, models.ScooterStatusAvailable); err != nil {
		return nil, fmt.Errorf("failed to update scooter status: %w", err)
	}

	if err := scooterRepo.UpdateLocation(ctx, scooterID, lat, lng); err != nil {
		// Location will be updated with first location update
	}

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

	tx, err := s.unitOfWork.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}

	var committed bool
	defer func() {
		if !committed {
			tx.Rollback()
		}
	}()

	tripRepo := tx.TripRepository()
	scooterRepo := tx.ScooterRepository()

	trip, err := tripRepo.GetActiveByScooterID(ctx, scooterID)
	if err != nil {
		return nil, fmt.Errorf("failed to get active trip: %w", err)
	}
	if trip == nil {
		return nil, errors.New("no active trip found for scooter")
	}

	endTime := time.Now()

	if err := tripRepo.EndTrip(ctx, trip.ID, lat, lng); err != nil {
		return nil, fmt.Errorf("failed to end trip: %w", err)
	}

	if err := scooterRepo.UpdateStatusWithCheck(ctx, scooterID, models.ScooterStatusAvailable, models.ScooterStatusOccupied); err != nil {
		return nil, fmt.Errorf("failed to update scooter status: %w", err)
	}

	if err := scooterRepo.UpdateLocation(ctx, scooterID, lat, lng); err != nil {
		// Location update failure is non-critical
	}

	trip.EndTime = &endTime
	trip.EndLatitude = &lat
	trip.EndLongitude = &lng
	trip.Status = models.TripStatusCompleted

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	committed = true
	return trip, nil
}

func (s *tripService) CancelTrip(ctx context.Context, scooterID uuid.UUID) (*models.Trip, error) {
	tx, err := s.unitOfWork.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}

	var committed bool
	defer func() {
		if !committed {
			tx.Rollback()
		}
	}()

	tripRepo := tx.TripRepository()
	scooterRepo := tx.ScooterRepository()

	trip, err := tripRepo.GetActiveByScooterID(ctx, scooterID)
	if err != nil {
		return nil, fmt.Errorf("failed to get active trip: %w", err)
	}
	if trip == nil {
		return nil, errors.New("no active trip found for scooter")
	}

	if err := tripRepo.CancelTrip(ctx, trip.ID); err != nil {
		return nil, fmt.Errorf("failed to cancel trip: %w", err)
	}

	if err := scooterRepo.UpdateStatusWithCheck(ctx, scooterID, models.ScooterStatusAvailable, models.ScooterStatusOccupied); err != nil {
		return nil, fmt.Errorf("failed to update scooter status: %w", err)
	}

	trip.Status = models.TripStatusCancelled

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

	tripRepo := tx.TripRepository()
	locationRepo := tx.LocationUpdateRepository()
	scooterRepo := tx.ScooterRepository()

	trip, err := tripRepo.GetActiveByScooterID(ctx, scooterID)
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
