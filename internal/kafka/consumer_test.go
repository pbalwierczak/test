package kafka

import (
	"context"
	"errors"
	"testing"
	"time"

	"scootin-aboot/internal/config"
	"scootin-aboot/internal/models"

	"github.com/Shopify/sarama"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestNewEventConsumer(t *testing.T) {
	// Skip this test since it requires actual Kafka connection
	t.Skip("Skipping NewEventConsumer test - requires Kafka connection or mocking sarama")
}

func TestEventConsumer_Setup(t *testing.T) {
	consumer := &EventConsumer{}
	session := &MockConsumerGroupSession{}

	err := consumer.Setup(session)
	assert.NoError(t, err)
}

func TestEventConsumer_Cleanup(t *testing.T) {
	consumer := &EventConsumer{}
	session := &MockConsumerGroupSession{}

	err := consumer.Cleanup(session)
	assert.NoError(t, err)
}

func TestEventConsumer_processMessage(t *testing.T) {
	tests := []struct {
		name        string
		message     *sarama.ConsumerMessage
		setupMocks  func(*MockTripService, *MockScooterService)
		expectError bool
		errorMsg    string
	}{
		{
			name: "trip started message",
			message: &sarama.ConsumerMessage{
				Topic: "trip-started",
				Value: []byte(`{"eventType":"trip.started","eventId":"test-id","timestamp":"2023-01-01T00:00:00Z","version":"1.0","data":{"tripId":"550e8400-e29b-41d4-a716-446655440000","scooterId":"550e8400-e29b-41d4-a716-446655440001","userId":"550e8400-e29b-41d4-a716-446655440002","startLatitude":45.4215,"startLongitude":-75.6972,"startTime":"2023-01-01T00:00:00Z"}}`),
			},
			setupMocks: func(tripService *MockTripService, scooterService *MockScooterService) {
				tripService.On("StartTrip", mock.Anything, mock.Anything, mock.Anything, 45.4215, -75.6972).Return(&models.Trip{}, nil)
			},
			expectError: false,
		},
		{
			name: "trip ended message",
			message: &sarama.ConsumerMessage{
				Topic: "trip-ended",
				Value: []byte(`{"eventType":"trip.ended","eventId":"test-id","timestamp":"2023-01-01T00:00:00Z","version":"1.0","data":{"tripId":"trip-123","scooterId":"550e8400-e29b-41d4-a716-446655440001","userId":"user-123","endLatitude":45.4216,"endLongitude":-75.6973,"endTime":"2023-01-01T00:30:00Z","durationSeconds":1800}}`),
			},
			setupMocks: func(tripService *MockTripService, scooterService *MockScooterService) {
				tripService.On("EndTrip", mock.Anything, mock.Anything, 45.4216, -75.6973).Return(&models.Trip{}, nil)
			},
			expectError: false,
		},
		{
			name: "location updated message",
			message: &sarama.ConsumerMessage{
				Topic: "location-updated",
				Value: []byte(`{"eventType":"location.updated","eventId":"test-id","timestamp":"2023-01-01T00:00:00Z","version":"1.0","data":{"scooterId":"550e8400-e29b-41d4-a716-446655440001","tripId":"trip-123","latitude":45.4216,"longitude":-75.6973,"heading":90.0,"speed":15.5}}`),
			},
			setupMocks: func(tripService *MockTripService, scooterService *MockScooterService) {
				tripService.On("UpdateLocation", mock.Anything, mock.Anything, 45.4216, -75.6973).Return(nil)
			},
			expectError: false,
		},
		{
			name: "unknown topic",
			message: &sarama.ConsumerMessage{
				Topic: "unknown-topic",
				Value: []byte(`{}`),
			},
			setupMocks: func(tripService *MockTripService, scooterService *MockScooterService) {
			},
			expectError: true,
			errorMsg:    "unknown topic",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tripService := &MockTripService{}
			scooterService := &MockScooterService{}
			tt.setupMocks(tripService, scooterService)

			deps := HandlerDependencies{
				TripService:    tripService,
				ScooterService: scooterService,
			}

			handlers := map[string]EventHandler{
				"trip-started":     NewTripStartedHandler(deps),
				"trip-ended":       NewTripEndedHandler(deps),
				"location-updated": NewLocationUpdatedHandler(deps),
			}

			consumer := &EventConsumer{
				config: &config.KafkaConfig{
					Topics: config.KafkaTopics{
						TripStarted:     "trip-started",
						TripEnded:       "trip-ended",
						LocationUpdated: "location-updated",
					},
				},
				handlers: handlers,
				ctx:      context.Background(),
			}

			err := consumer.processMessage(tt.message)

			if tt.expectError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errorMsg)
			} else {
				assert.NoError(t, err)
			}

			tripService.AssertExpectations(t)
			scooterService.AssertExpectations(t)
		})
	}
}

func TestEventConsumer_ConsumeClaim(t *testing.T) {
	tests := []struct {
		name          string
		setupMocks    func(*MockConsumerGroupSession, *MockConsumerGroupClaim)
		expectError   bool
		expectedCalls int
	}{
		{
			name: "successful message processing",
			setupMocks: func(session *MockConsumerGroupSession, claim *MockConsumerGroupClaim) {
				ctx := context.Background()
				session.On("Context").Return(ctx)
				claim.On("Messages").Return(nil) // Return nil to use the internal channel
				session.On("MarkMessage", mock.Anything, "").Return()

				// Send a test message
				go func() {
					time.Sleep(10 * time.Millisecond)
					claim.SendMessage(&sarama.ConsumerMessage{
						Topic:     "trip-started",
						Partition: 0,
						Offset:    1,
						Value:     []byte(`{"eventType":"trip.started","eventId":"test-id","timestamp":"2023-01-01T00:00:00Z","version":"1.0","data":{"tripId":"550e8400-e29b-41d4-a716-446655440000","scooterId":"550e8400-e29b-41d4-a716-446655440001","userId":"550e8400-e29b-41d4-a716-446655440002","startLatitude":45.4215,"startLongitude":-75.6972,"startTime":"2023-01-01T00:00:00Z"}}`),
					})
					time.Sleep(10 * time.Millisecond)
					claim.Close()
				}()
			},
			expectError:   false,
			expectedCalls: 1,
		},
		{
			name: "context cancellation",
			setupMocks: func(session *MockConsumerGroupSession, claim *MockConsumerGroupClaim) {
				ctx, cancel := context.WithCancel(context.Background())
				session.On("Context").Return(ctx)
				claim.On("Messages").Return(nil) // Return nil to use the internal channel

				// Cancel context after a short delay
				go func() {
					time.Sleep(10 * time.Millisecond)
					cancel()
				}()
			},
			expectError:   false,
			expectedCalls: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tripService := &MockTripService{}
			scooterService := &MockScooterService{}
			session := &MockConsumerGroupSession{}
			claim := NewMockConsumerGroupClaim()

			if tt.name == "successful message processing" {
				tripService.On("StartTrip", mock.Anything, mock.Anything, mock.Anything, 45.4215, -75.6972).Return(&models.Trip{}, nil)
			}

			tt.setupMocks(session, claim)

			deps := HandlerDependencies{
				TripService:    tripService,
				ScooterService: scooterService,
			}

			handlers := map[string]EventHandler{
				"trip-started":     NewTripStartedHandler(deps),
				"trip-ended":       NewTripEndedHandler(deps),
				"location-updated": NewLocationUpdatedHandler(deps),
			}

			consumer := &EventConsumer{
				config: &config.KafkaConfig{
					Topics: config.KafkaTopics{
						TripStarted:     "trip-started",
						TripEnded:       "trip-ended",
						LocationUpdated: "location-updated",
					},
				},
				handlers: handlers,
				ctx:      context.Background(),
			}

			err := consumer.ConsumeClaim(session, claim)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			session.AssertExpectations(t)
			claim.AssertExpectations(t)
			tripService.AssertExpectations(t)
		})
	}
}

func TestEventConsumer_Start(t *testing.T) {
	tests := []struct {
		name          string
		consumerGroup *MockConsumerGroup
		expectError   bool
		expectedCalls int
	}{
		{
			name: "successful start",
			consumerGroup: func() *MockConsumerGroup {
				mockGroup := &MockConsumerGroup{}
				mockGroup.On("Consume", mock.Anything, mock.Anything, mock.Anything).Return(nil)
				mockGroup.On("Close").Return(nil)
				return mockGroup
			}(),
			expectError:   false,
			expectedCalls: 1,
		},
		{
			name: "consume error",
			consumerGroup: func() *MockConsumerGroup {
				mockGroup := &MockConsumerGroup{}
				mockGroup.On("Consume", mock.Anything, mock.Anything, mock.Anything).Return(errors.New("consume error"))
				mockGroup.On("Close").Return(nil)
				return mockGroup
			}(),
			expectError:   false, // Start doesn't return error, it logs and retries
			expectedCalls: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithCancel(context.Background())
			consumer := &EventConsumer{
				consumerGroup: tt.consumerGroup,
				config: &config.KafkaConfig{
					Topics: config.KafkaTopics{
						TripStarted:     "trip-started",
						TripEnded:       "trip-ended",
						LocationUpdated: "location-updated",
					},
				},
				ctx:    ctx,
				cancel: cancel,
			}

			err := consumer.Start()
			assert.NoError(t, err)

			// Wait a bit for the goroutine to start
			time.Sleep(50 * time.Millisecond)

			// Stop the consumer to clean up
			consumer.Stop()

			tt.consumerGroup.AssertExpectations(t)
		})
	}
}

func TestEventConsumer_Stop(t *testing.T) {
	tests := []struct {
		name          string
		consumerGroup *MockConsumerGroup
		expectError   bool
	}{
		{
			name: "successful stop",
			consumerGroup: func() *MockConsumerGroup {
				mockGroup := &MockConsumerGroup{}
				mockGroup.On("Close").Return(nil)
				return mockGroup
			}(),
			expectError: false,
		},
		{
			name: "close error",
			consumerGroup: func() *MockConsumerGroup {
				mockGroup := &MockConsumerGroup{}
				mockGroup.On("Close").Return(errors.New("close error"))
				return mockGroup
			}(),
			expectError: false, // Stop doesn't return error, it logs the error
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithCancel(context.Background())
			consumer := &EventConsumer{
				consumerGroup: tt.consumerGroup,
				ctx:           ctx,
				cancel:        cancel,
			}

			consumer.Stop()

			tt.consumerGroup.AssertExpectations(t)
		})
	}
}
