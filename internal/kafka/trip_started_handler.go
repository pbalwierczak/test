package kafka

import (
	"context"
	"encoding/json"
	"fmt"

	"scootin-aboot/pkg/kafka"
	"scootin-aboot/pkg/logger"

	"github.com/google/uuid"
)

type TripStartedHandler struct {
	deps HandlerDependencies
}

func NewTripStartedHandler(deps HandlerDependencies) *TripStartedHandler {
	return &TripStartedHandler{
		deps: deps,
	}
}

func (h *TripStartedHandler) Handle(ctx context.Context, data []byte) error {
	var event kafka.TripStartedEvent
	if err := json.Unmarshal(data, &event); err != nil {
		return fmt.Errorf("failed to unmarshal trip started event: %w", err)
	}

	logger.Info("Processing trip started event",
		logger.String("trip_id", event.Data.TripID),
		logger.String("scooter_id", event.Data.ScooterID),
		logger.String("user_id", event.Data.UserID),
	)

	scooterID, err := uuid.Parse(event.Data.ScooterID)
	if err != nil {
		return fmt.Errorf("invalid scooter ID: %w", err)
	}

	userID, err := uuid.Parse(event.Data.UserID)
	if err != nil {
		return fmt.Errorf("invalid user ID: %w", err)
	}

	trip, err := h.deps.TripService.StartTrip(ctx, scooterID, userID, event.Data.StartLatitude, event.Data.StartLongitude)
	if err != nil {
		return fmt.Errorf("failed to start trip: %w", err)
	}

	logger.Info("Trip started event processed successfully",
		logger.String("trip_id", trip.ID.String()),
		logger.String("scooter_id", event.Data.ScooterID),
		logger.String("user_id", event.Data.UserID),
	)

	return nil
}
