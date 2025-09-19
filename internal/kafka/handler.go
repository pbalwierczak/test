package kafka

import (
	"context"

	"scootin-aboot/internal/services"
)

type EventHandler interface {
	Handle(ctx context.Context, data []byte) error
}

type HandlerDependencies struct {
	TripService    services.TripService
	ScooterService services.ScooterService
}
