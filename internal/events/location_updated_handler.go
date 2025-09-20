package events

import (
	"context"
	"encoding/json"
	"fmt"

	"scootin-aboot/internal/logger"

	"github.com/google/uuid"
)

type LocationUpdatedHandler struct {
	deps HandlerDependencies
}

func NewLocationUpdatedHandler(deps HandlerDependencies) *LocationUpdatedHandler {
	return &LocationUpdatedHandler{
		deps: deps,
	}
}

func (h *LocationUpdatedHandler) Handle(ctx context.Context, data []byte) error {
	var event LocationUpdatedEvent
	if err := json.Unmarshal(data, &event); err != nil {
		return fmt.Errorf("failed to unmarshal location updated event: %w", err)
	}

	logger.Debug("Processing location updated event",
		logger.String("scooter_id", event.Data.ScooterID),
		logger.String("trip_id", event.Data.TripID),
		logger.Float64("lat", event.Data.Latitude),
		logger.Float64("lng", event.Data.Longitude),
	)

	scooterID, err := uuid.Parse(event.Data.ScooterID)
	if err != nil {
		return fmt.Errorf("invalid scooter ID: %w", err)
	}

	if err := h.deps.TripService.UpdateLocation(ctx, scooterID, event.Data.Latitude, event.Data.Longitude); err != nil {
		return fmt.Errorf("failed to update location: %w", err)
	}

	logger.Debug("Location updated event processed successfully",
		logger.String("scooter_id", event.Data.ScooterID),
	)

	return nil
}
