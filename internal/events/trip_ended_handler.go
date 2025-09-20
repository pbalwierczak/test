package events

import (
	"context"
	"encoding/json"
	"fmt"

	"scootin-aboot/internal/logger"

	"github.com/google/uuid"
)

type TripEndedHandler struct {
	deps HandlerDependencies
}

func NewTripEndedHandler(deps HandlerDependencies) *TripEndedHandler {
	return &TripEndedHandler{
		deps: deps,
	}
}

func (h *TripEndedHandler) Handle(ctx context.Context, data []byte) error {
	var event TripEndedEvent
	if err := json.Unmarshal(data, &event); err != nil {
		return fmt.Errorf("failed to unmarshal trip ended event: %w", err)
	}

	logger.Info("Processing trip ended event",
		logger.String("trip_id", event.Data.TripID),
		logger.String("scooter_id", event.Data.ScooterID),
		logger.String("user_id", event.Data.UserID),
		logger.Int("duration_seconds", event.Data.DurationSeconds),
	)

	scooterID, err := uuid.Parse(event.Data.ScooterID)
	if err != nil {
		return fmt.Errorf("invalid scooter ID: %w", err)
	}

	trip, err := h.deps.TripService.EndTrip(ctx, scooterID, event.Data.EndLatitude, event.Data.EndLongitude)
	if err != nil {
		return fmt.Errorf("failed to end trip: %w", err)
	}

	logger.Info("Trip ended event processed successfully",
		logger.String("trip_id", trip.ID.String()),
		logger.String("scooter_id", event.Data.ScooterID),
		logger.String("user_id", event.Data.UserID),
	)

	return nil
}
