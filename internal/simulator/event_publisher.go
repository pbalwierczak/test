package simulator

import (
	"context"
	"time"

	"scootin-aboot/internal/events"
)

// EventPublisher interface for publishing events from the simulator
type EventPublisher interface {
	PublishTripStarted(ctx context.Context, tripID, scooterID, userID string, lat, lng float64) error
	PublishTripEnded(ctx context.Context, tripID, scooterID, userID string, lat, lng float64, startTime time.Time) error
	PublishLocationUpdated(ctx context.Context, scooterID, tripID string, lat, lng, heading, speed float64) error
	Close() error
}

// KafkaEventPublisher implements EventPublisher using Kafka
type KafkaEventPublisher struct {
	producer events.EventProducer
}

// NewKafkaEventPublisher creates a new Kafka event publisher
func NewKafkaEventPublisher(producer events.EventProducer) EventPublisher {
	return &KafkaEventPublisher{
		producer: producer,
	}
}

// PublishTripStarted publishes a trip started event
func (p *KafkaEventPublisher) PublishTripStarted(ctx context.Context, tripID, scooterID, userID string, lat, lng float64) error {
	event := events.NewTripStartedEvent(tripID, scooterID, userID, lat, lng)
	return p.producer.PublishTripStarted(ctx, event)
}

// PublishTripEnded publishes a trip ended event
func (p *KafkaEventPublisher) PublishTripEnded(ctx context.Context, tripID, scooterID, userID string, lat, lng float64, startTime time.Time) error {
	event := events.NewTripEndedEvent(tripID, scooterID, userID, lat, lng, startTime)
	return p.producer.PublishTripEnded(ctx, event)
}

// PublishLocationUpdated publishes a location updated event
func (p *KafkaEventPublisher) PublishLocationUpdated(ctx context.Context, scooterID, tripID string, lat, lng, heading, speed float64) error {
	event := events.NewLocationUpdatedEvent(scooterID, tripID, lat, lng, heading, speed)
	return p.producer.PublishLocationUpdated(ctx, event)
}

// Close closes the event publisher
func (p *KafkaEventPublisher) Close() error {
	return p.producer.Close()
}
