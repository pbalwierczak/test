package events

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"scootin-aboot/internal/config"
	"scootin-aboot/internal/logger"

	"github.com/Shopify/sarama"
)

type EventProducer interface {
	PublishTripStarted(ctx context.Context, event *TripStartedEvent) error
	PublishTripEnded(ctx context.Context, event *TripEndedEvent) error
	PublishLocationUpdated(ctx context.Context, event *LocationUpdatedEvent) error
	Close() error
}

type KafkaProducer struct {
	producer sarama.SyncProducer
	config   *config.KafkaConfig
}

func NewKafkaProducer(cfg *config.KafkaConfig) (*KafkaProducer, error) {
	if cfg == nil {
		return nil, fmt.Errorf("failed to create Kafka producer: config cannot be nil")
	}

	saramaConfig := sarama.NewConfig()
	saramaConfig.Producer.RequiredAcks = sarama.WaitForAll
	saramaConfig.Producer.Retry.Max = 3
	saramaConfig.Producer.Return.Successes = true
	saramaConfig.Producer.Timeout = 10 * time.Second
	saramaConfig.Net.DialTimeout = 5 * time.Second
	saramaConfig.Net.ReadTimeout = 5 * time.Second
	saramaConfig.Net.WriteTimeout = 5 * time.Second

	producer, err := sarama.NewSyncProducer(cfg.Brokers, saramaConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create Kafka producer: %w", err)
	}

	return &KafkaProducer{
		producer: producer,
		config:   cfg,
	}, nil
}

func (p *KafkaProducer) PublishTripStarted(ctx context.Context, event *TripStartedEvent) error {
	return p.publishEvent(ctx, p.config.Topics.TripStarted, event)
}

func (p *KafkaProducer) PublishTripEnded(ctx context.Context, event *TripEndedEvent) error {
	return p.publishEvent(ctx, p.config.Topics.TripEnded, event)
}

func (p *KafkaProducer) PublishLocationUpdated(ctx context.Context, event *LocationUpdatedEvent) error {
	return p.publishEvent(ctx, p.config.Topics.LocationUpdated, event)
}

func (p *KafkaProducer) publishEvent(ctx context.Context, topic string, event interface{}) error {
	eventJSON, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to marshal event: %w", err)
	}

	message := &sarama.ProducerMessage{
		Topic: topic,
		Key:   sarama.StringEncoder(fmt.Sprintf("%s-%d", topic, time.Now().UnixNano())),
		Value: sarama.ByteEncoder(eventJSON),
		Headers: []sarama.RecordHeader{
			{
				Key:   []byte("event-type"),
				Value: []byte(fmt.Sprintf("%T", event)),
			},
			{
				Key:   []byte("timestamp"),
				Value: []byte(time.Now().Format(time.RFC3339)),
			},
		},
	}

	partition, offset, err := p.producer.SendMessage(message)
	if err != nil {
		return fmt.Errorf("failed to send message to Kafka: %w", err)
	}

	logger.Debug("Event published to Kafka",
		logger.String("topic", topic),
		logger.String("partition", fmt.Sprintf("%d", partition)),
		logger.String("offset", fmt.Sprintf("%d", offset)),
		logger.String("event_type", fmt.Sprintf("%T", event)),
	)

	return nil
}

func (p *KafkaProducer) Close() error {
	if p.producer != nil {
		return p.producer.Close()
	}
	return nil
}

type MockProducer struct {
	Events []interface{}
}

func NewMockProducer() *MockProducer {
	return &MockProducer{
		Events: make([]interface{}, 0),
	}
}

func (m *MockProducer) PublishTripStarted(ctx context.Context, event *TripStartedEvent) error {
	m.Events = append(m.Events, event)
	logger.Debug("Mock: Trip started event published", logger.String("trip_id", event.Data.TripID))
	return nil
}

func (m *MockProducer) PublishTripEnded(ctx context.Context, event *TripEndedEvent) error {
	m.Events = append(m.Events, event)
	logger.Debug("Mock: Trip ended event published", logger.String("trip_id", event.Data.TripID))
	return nil
}

func (m *MockProducer) PublishLocationUpdated(ctx context.Context, event *LocationUpdatedEvent) error {
	m.Events = append(m.Events, event)
	logger.Debug("Mock: Location updated event published", logger.String("scooter_id", event.Data.ScooterID))
	return nil
}

func (m *MockProducer) Close() error {
	return nil
}

func (m *MockProducer) GetEvents() []interface{} {
	return m.Events
}

func (m *MockProducer) ClearEvents() {
	m.Events = make([]interface{}, 0)
}
