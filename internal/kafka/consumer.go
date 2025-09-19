package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"scootin-aboot/internal/config"
	"scootin-aboot/internal/services"
	"scootin-aboot/pkg/kafka"
	"scootin-aboot/pkg/logger"

	"github.com/Shopify/sarama"
	"github.com/google/uuid"
)

// EventConsumer handles consuming events from Kafka
type EventConsumer struct {
	consumerGroup  sarama.ConsumerGroup
	config         *config.KafkaConfig
	tripService    services.TripService
	scooterService services.ScooterService
	ctx            context.Context
	cancel         context.CancelFunc
	wg             sync.WaitGroup
}

// NewEventConsumer creates a new event consumer
func NewEventConsumer(cfg *config.KafkaConfig, tripService services.TripService, scooterService services.ScooterService) (*EventConsumer, error) {
	saramaConfig := sarama.NewConfig()
	saramaConfig.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategyRoundRobin
	saramaConfig.Consumer.Offsets.Initial = sarama.OffsetNewest
	saramaConfig.Consumer.Group.Session.Timeout = 10 * time.Second
	saramaConfig.Consumer.Group.Heartbeat.Interval = 3 * time.Second
	saramaConfig.Consumer.MaxProcessingTime = 500 * time.Millisecond

	consumerGroup, err := sarama.NewConsumerGroup(cfg.Brokers, "scooter-service", saramaConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create consumer group: %w", err)
	}

	ctx, cancel := context.WithCancel(context.Background())

	return &EventConsumer{
		consumerGroup:  consumerGroup,
		config:         cfg,
		tripService:    tripService,
		scooterService: scooterService,
		ctx:            ctx,
		cancel:         cancel,
	}, nil
}

// Start starts consuming events
func (c *EventConsumer) Start() error {
	topics := []string{
		c.config.Topics.TripStarted,
		c.config.Topics.TripEnded,
		c.config.Topics.LocationUpdated,
	}

	c.wg.Add(1)
	go func() {
		defer c.wg.Done()
		for {
			select {
			case <-c.ctx.Done():
				return
			default:
				if err := c.consumerGroup.Consume(c.ctx, topics, c); err != nil {
					logger.Error("Error consuming from Kafka", logger.ErrorField(err))
					time.Sleep(1 * time.Second)
				}
			}
		}
	}()

	logger.Info("Kafka consumer started", logger.Strings("topics", topics))
	return nil
}

// Stop stops consuming events
func (c *EventConsumer) Stop() {
	logger.Info("Stopping Kafka consumer...")
	c.cancel()
	c.wg.Wait()

	if err := c.consumerGroup.Close(); err != nil {
		logger.Error("Error closing consumer group", logger.ErrorField(err))
	}

	logger.Info("Kafka consumer stopped")
}

// Setup is run at the beginning of a new session, before ConsumeClaim
func (c *EventConsumer) Setup(sarama.ConsumerGroupSession) error {
	return nil
}

// Cleanup is run at the end of a session, once all ConsumeClaim goroutines have exited
func (c *EventConsumer) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

// ConsumeClaim must start a consumer loop of ConsumerGroupClaim's Messages()
func (c *EventConsumer) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for {
		select {
		case message := <-claim.Messages():
			if message == nil {
				return nil
			}

			if err := c.processMessage(message); err != nil {
				logger.Error("Error processing message",
					logger.String("topic", message.Topic),
					logger.String("partition", fmt.Sprintf("%d", message.Partition)),
					logger.String("offset", fmt.Sprintf("%d", message.Offset)),
					logger.ErrorField(err),
				)
			}

			session.MarkMessage(message, "")

		case <-session.Context().Done():
			return nil
		}
	}
}

// processMessage processes a single Kafka message
func (c *EventConsumer) processMessage(message *sarama.ConsumerMessage) error {
	logger.Debug("Processing Kafka message",
		logger.String("topic", message.Topic),
		logger.String("partition", fmt.Sprintf("%d", message.Partition)),
		logger.String("offset", fmt.Sprintf("%d", message.Offset)),
	)

	switch message.Topic {
	case c.config.Topics.TripStarted:
		return c.handleTripStarted(message.Value)
	case c.config.Topics.TripEnded:
		return c.handleTripEnded(message.Value)
	case c.config.Topics.LocationUpdated:
		return c.handleLocationUpdated(message.Value)
	default:
		return fmt.Errorf("unknown topic: %s", message.Topic)
	}
}

// handleTripStarted handles trip started events
func (c *EventConsumer) handleTripStarted(data []byte) error {
	var event kafka.TripStartedEvent
	if err := json.Unmarshal(data, &event); err != nil {
		return fmt.Errorf("failed to unmarshal trip started event: %w", err)
	}

	logger.Info("Processing trip started event",
		logger.String("trip_id", event.Data.TripID),
		logger.String("scooter_id", event.Data.ScooterID),
		logger.String("user_id", event.Data.UserID),
	)

	// Convert string IDs to UUIDs
	scooterID, err := uuid.Parse(event.Data.ScooterID)
	if err != nil {
		return fmt.Errorf("invalid scooter ID: %w", err)
	}

	userID, err := uuid.Parse(event.Data.UserID)
	if err != nil {
		return fmt.Errorf("invalid user ID: %w", err)
	}

	// Call trip service to start the trip
	trip, err := c.tripService.StartTrip(c.ctx, scooterID, userID, event.Data.StartLatitude, event.Data.StartLongitude)
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

// handleTripEnded handles trip ended events
func (c *EventConsumer) handleTripEnded(data []byte) error {
	var event kafka.TripEndedEvent
	if err := json.Unmarshal(data, &event); err != nil {
		return fmt.Errorf("failed to unmarshal trip ended event: %w", err)
	}

	logger.Info("Processing trip ended event",
		logger.String("trip_id", event.Data.TripID),
		logger.String("scooter_id", event.Data.ScooterID),
		logger.String("user_id", event.Data.UserID),
		logger.Int("duration_seconds", event.Data.DurationSeconds),
	)

	// Convert string ID to UUID
	scooterID, err := uuid.Parse(event.Data.ScooterID)
	if err != nil {
		return fmt.Errorf("invalid scooter ID: %w", err)
	}

	// Call trip service to end the trip
	trip, err := c.tripService.EndTrip(c.ctx, scooterID, event.Data.EndLatitude, event.Data.EndLongitude)
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

// handleLocationUpdated handles location updated events
func (c *EventConsumer) handleLocationUpdated(data []byte) error {
	var event kafka.LocationUpdatedEvent
	if err := json.Unmarshal(data, &event); err != nil {
		return fmt.Errorf("failed to unmarshal location updated event: %w", err)
	}

	logger.Debug("Processing location updated event",
		logger.String("scooter_id", event.Data.ScooterID),
		logger.String("trip_id", event.Data.TripID),
		logger.Float64("lat", event.Data.Latitude),
		logger.Float64("lng", event.Data.Longitude),
	)

	// Convert string ID to UUID
	scooterID, err := uuid.Parse(event.Data.ScooterID)
	if err != nil {
		return fmt.Errorf("invalid scooter ID: %w", err)
	}

	// Call trip service to update location
	if err := c.tripService.UpdateLocation(c.ctx, scooterID, event.Data.Latitude, event.Data.Longitude); err != nil {
		return fmt.Errorf("failed to update location: %w", err)
	}

	logger.Debug("Location updated event processed successfully",
		logger.String("scooter_id", event.Data.ScooterID),
	)

	return nil
}
