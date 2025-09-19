package kafka

import (
	"context"
	"errors"
	"testing"
	"time"

	"scootin-aboot/internal/config"

	"github.com/Shopify/sarama"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockSyncProducer is a mock implementation of sarama.SyncProducer
type MockSyncProducer struct {
	mock.Mock
}

func (m *MockSyncProducer) SendMessage(msg *sarama.ProducerMessage) (partition int32, offset int64, err error) {
	args := m.Called(msg)
	return args.Get(0).(int32), args.Get(1).(int64), args.Error(2)
}

func (m *MockSyncProducer) SendMessages(msgs []*sarama.ProducerMessage) error {
	args := m.Called(msgs)
	return args.Error(0)
}

func (m *MockSyncProducer) Close() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockSyncProducer) TxnStatus() sarama.ProducerTxnStatusFlag {
	args := m.Called()
	return args.Get(0).(sarama.ProducerTxnStatusFlag)
}

func (m *MockSyncProducer) IsTransactional() bool {
	args := m.Called()
	return args.Bool(0)
}

func (m *MockSyncProducer) BeginTxn() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockSyncProducer) CommitTxn() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockSyncProducer) AbortTxn() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockSyncProducer) AddOffsetsToTxn(offsets map[string][]*sarama.PartitionOffsetMetadata, groupId string) error {
	args := m.Called(offsets, groupId)
	return args.Error(0)
}

func (m *MockSyncProducer) AddMessageToTxn(msg *sarama.ConsumerMessage, groupId string, metadata *string) error {
	args := m.Called(msg, groupId, metadata)
	return args.Error(0)
}

func TestNewKafkaProducer(t *testing.T) {
	tests := []struct {
		name        string
		config      *config.KafkaConfig
		expectError bool
		errorMsg    string
	}{
		{
			name:        "nil config",
			config:      nil,
			expectError: true,
			errorMsg:    "failed to create Kafka producer",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			producer, err := NewKafkaProducer(tt.config)

			if tt.expectError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errorMsg)
				assert.Nil(t, producer)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, producer)
				assert.Equal(t, tt.config, producer.config)
			}
		})
	}
}
func TestKafkaProducer_PublishTripStarted(t *testing.T) {
	tests := []struct {
		name        string
		event       *TripStartedEvent
		setupMock   func(*MockSyncProducer)
		expectError bool
		errorMsg    string
	}{
		{
			name: "successful publish",
			event: &TripStartedEvent{
				BaseEvent: BaseEvent{
					EventType: "trip.started",
					EventID:   uuid.New().String(),
					Timestamp: time.Now(),
					Version:   "1.0",
				},
				Data: TripStartedData{
					TripID:         "trip-123",
					ScooterID:      "scooter-123",
					UserID:         "user-123",
					StartLatitude:  45.4215,
					StartLongitude: -75.6972,
					StartTime:      time.Now().Format(time.RFC3339),
				},
			},
			setupMock: func(mockProducer *MockSyncProducer) {
				mockProducer.On("SendMessage", mock.AnythingOfType("*sarama.ProducerMessage")).Return(int32(0), int64(1), nil)
			},
			expectError: false,
		},
		{
			name: "publish error",
			event: &TripStartedEvent{
				BaseEvent: BaseEvent{
					EventType: "trip.started",
					EventID:   uuid.New().String(),
					Timestamp: time.Now(),
					Version:   "1.0",
				},
				Data: TripStartedData{
					TripID:         "trip-123",
					ScooterID:      "scooter-123",
					UserID:         "user-123",
					StartLatitude:  45.4215,
					StartLongitude: -75.6972,
					StartTime:      time.Now().Format(time.RFC3339),
				},
			},
			setupMock: func(mockProducer *MockSyncProducer) {
				mockProducer.On("SendMessage", mock.AnythingOfType("*sarama.ProducerMessage")).Return(int32(0), int64(0), errors.New("publish error"))
			},
			expectError: true,
			errorMsg:    "failed to send message to Kafka: publish error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockProducer := &MockSyncProducer{}
			tt.setupMock(mockProducer)

			producer := &KafkaProducer{
				producer: mockProducer,
				config: &config.KafkaConfig{
					Topics: config.KafkaTopics{
						TripStarted: "trip-started",
					},
				},
			}

			err := producer.PublishTripStarted(context.Background(), tt.event)

			if tt.expectError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errorMsg)
			} else {
				assert.NoError(t, err)
			}

			mockProducer.AssertExpectations(t)
		})
	}
}

func TestKafkaProducer_Close(t *testing.T) {
	tests := []struct {
		name        string
		setupMock   func(*MockSyncProducer)
		expectError bool
		errorMsg    string
	}{
		{
			name: "successful close",
			setupMock: func(mockProducer *MockSyncProducer) {
				mockProducer.On("Close").Return(nil)
			},
			expectError: false,
		},
		{
			name: "close error",
			setupMock: func(mockProducer *MockSyncProducer) {
				mockProducer.On("Close").Return(errors.New("close error"))
			},
			expectError: true,
			errorMsg:    "close error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockProducer := &MockSyncProducer{}
			tt.setupMock(mockProducer)

			producer := &KafkaProducer{
				producer: mockProducer,
			}

			err := producer.Close()

			if tt.expectError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errorMsg)
			} else {
				assert.NoError(t, err)
			}

			mockProducer.AssertExpectations(t)
		})
	}
}

func TestNewTripStartedEvent(t *testing.T) {
	tripID := "trip-123"
	scooterID := "scooter-123"
	userID := "user-123"
	startLat := 45.4215
	startLng := -75.6972

	event := NewTripStartedEvent(tripID, scooterID, userID, startLat, startLng)

	assert.Equal(t, "trip.started", event.EventType)
	assert.NotEmpty(t, event.EventID)
	assert.NotZero(t, event.Timestamp)
	assert.Equal(t, "1.0", event.Version)
	assert.Equal(t, tripID, event.Data.TripID)
	assert.Equal(t, scooterID, event.Data.ScooterID)
	assert.Equal(t, userID, event.Data.UserID)
	assert.Equal(t, startLat, event.Data.StartLatitude)
	assert.Equal(t, startLng, event.Data.StartLongitude)
	assert.NotEmpty(t, event.Data.StartTime)
}
