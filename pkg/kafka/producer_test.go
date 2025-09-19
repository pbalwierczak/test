package kafka

import (
	"context"
	"testing"
	"time"

	"scootin-aboot/internal/config"
)

func TestMockProducer(t *testing.T) {
	producer := NewMockProducer()

	// Test trip started event
	tripStartedEvent := NewTripStartedEvent("trip-123", "scooter-456", "user-789", 45.4215, -75.6972)
	err := producer.PublishTripStarted(context.Background(), tripStartedEvent)
	if err != nil {
		t.Fatalf("Failed to publish trip started event: %v", err)
	}

	// Test trip ended event
	startTime := time.Now().Add(-10 * time.Minute)
	tripEndedEvent := NewTripEndedEvent("trip-123", "scooter-456", "user-789", 45.4300, -75.6800, startTime)
	err = producer.PublishTripEnded(context.Background(), tripEndedEvent)
	if err != nil {
		t.Fatalf("Failed to publish trip ended event: %v", err)
	}

	// Test location updated event
	locationUpdatedEvent := NewLocationUpdatedEvent("scooter-456", "trip-123", 45.4250, -75.6900, 45.5, 15.2)
	err = producer.PublishLocationUpdated(context.Background(), locationUpdatedEvent)
	if err != nil {
		t.Fatalf("Failed to publish location updated event: %v", err)
	}

	// Verify events were stored
	events := producer.GetEvents()
	if len(events) != 3 {
		t.Fatalf("Expected 3 events, got %d", len(events))
	}

	// Test event types
	if _, ok := events[0].(*TripStartedEvent); !ok {
		t.Error("First event should be TripStartedEvent")
	}
	if _, ok := events[1].(*TripEndedEvent); !ok {
		t.Error("Second event should be TripEndedEvent")
	}
	if _, ok := events[2].(*LocationUpdatedEvent); !ok {
		t.Error("Third event should be LocationUpdatedEvent")
	}

	// Test clear events
	producer.ClearEvents()
	events = producer.GetEvents()
	if len(events) != 0 {
		t.Fatalf("Expected 0 events after clear, got %d", len(events))
	}
}

func TestEventCreation(t *testing.T) {
	// Test trip started event creation
	tripStarted := NewTripStartedEvent("trip-123", "scooter-456", "user-789", 45.4215, -75.6972)

	if tripStarted.EventType != "trip.started" {
		t.Errorf("Expected event type 'trip.started', got '%s'", tripStarted.EventType)
	}
	if tripStarted.Data.TripID != "trip-123" {
		t.Errorf("Expected trip ID 'trip-123', got '%s'", tripStarted.Data.TripID)
	}
	if tripStarted.Data.ScooterID != "scooter-456" {
		t.Errorf("Expected scooter ID 'scooter-456', got '%s'", tripStarted.Data.ScooterID)
	}
	if tripStarted.Data.UserID != "user-789" {
		t.Errorf("Expected user ID 'user-789', got '%s'", tripStarted.Data.UserID)
	}

	// Test trip ended event creation
	startTime := time.Now().Add(-10 * time.Minute)
	tripEnded := NewTripEndedEvent("trip-123", "scooter-456", "user-789", 45.4300, -75.6800, startTime)

	if tripEnded.EventType != "trip.ended" {
		t.Errorf("Expected event type 'trip.ended', got '%s'", tripEnded.EventType)
	}
	if tripEnded.Data.DurationSeconds < 0 {
		t.Error("Duration should be positive")
	}

	// Test location updated event creation
	locationUpdated := NewLocationUpdatedEvent("scooter-456", "trip-123", 45.4250, -75.6900, 45.5, 15.2)

	if locationUpdated.EventType != "location.updated" {
		t.Errorf("Expected event type 'location.updated', got '%s'", locationUpdated.EventType)
	}
	if locationUpdated.Data.ScooterID != "scooter-456" {
		t.Errorf("Expected scooter ID 'scooter-456', got '%s'", locationUpdated.Data.ScooterID)
	}
}

func TestKafkaConfig(t *testing.T) {
	cfg := &config.KafkaConfig{
		Brokers:          []string{"localhost:9092"},
		ClientID:         "test-client",
		SecurityProtocol: "PLAINTEXT",
		Topics: config.KafkaTopics{
			TripStarted:     "test.trip.started",
			TripEnded:       "test.trip.ended",
			LocationUpdated: "test.location.updated",
		},
	}

	if len(cfg.Brokers) != 1 {
		t.Errorf("Expected 1 broker, got %d", len(cfg.Brokers))
	}
	if cfg.Brokers[0] != "localhost:9092" {
		t.Errorf("Expected broker 'localhost:9092', got '%s'", cfg.Brokers[0])
	}
	if cfg.ClientID != "test-client" {
		t.Errorf("Expected client ID 'test-client', got '%s'", cfg.ClientID)
	}
	if cfg.Topics.TripStarted != "test.trip.started" {
		t.Errorf("Expected topic 'test.trip.started', got '%s'", cfg.Topics.TripStarted)
	}
}
