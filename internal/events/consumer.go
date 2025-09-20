package events

import (
	"context"
	"fmt"
	"sync"
	"time"

	"scootin-aboot/internal/config"
	"scootin-aboot/internal/logger"
	"scootin-aboot/internal/services"

	"github.com/Shopify/sarama"
)

type EventConsumer struct {
	consumerGroup sarama.ConsumerGroup
	config        *config.KafkaConfig
	handlers      map[string]EventHandler
	ctx           context.Context
	cancel        context.CancelFunc
	wg            sync.WaitGroup
}

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

	deps := HandlerDependencies{
		TripService:    tripService,
		ScooterService: scooterService,
	}

	handlers := map[string]EventHandler{
		cfg.Topics.TripStarted:     NewTripStartedHandler(deps),
		cfg.Topics.TripEnded:       NewTripEndedHandler(deps),
		cfg.Topics.LocationUpdated: NewLocationUpdatedHandler(deps),
	}

	return &EventConsumer{
		consumerGroup: consumerGroup,
		config:        cfg,
		handlers:      handlers,
		ctx:           ctx,
		cancel:        cancel,
	}, nil
}

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

func (c *EventConsumer) Stop() {
	logger.Info("Stopping Kafka consumer...")
	c.cancel()
	c.wg.Wait()

	if err := c.consumerGroup.Close(); err != nil {
		logger.Error("Error closing consumer group", logger.ErrorField(err))
	}

	logger.Info("Kafka consumer stopped")
}

func (c *EventConsumer) Setup(sarama.ConsumerGroupSession) error {
	return nil
}

func (c *EventConsumer) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

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

func (c *EventConsumer) processMessage(message *sarama.ConsumerMessage) error {
	logger.Debug("Processing Kafka message",
		logger.String("topic", message.Topic),
		logger.String("partition", fmt.Sprintf("%d", message.Partition)),
		logger.String("offset", fmt.Sprintf("%d", message.Offset)),
	)

	handler, exists := c.handlers[message.Topic]
	if !exists {
		return fmt.Errorf("unknown topic: %s", message.Topic)
	}

	return handler.Handle(c.ctx, message.Value)
}
